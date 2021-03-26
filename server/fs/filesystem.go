package fs

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
	pb "github.com/henrixapp/mampf-rpc/grpc"
	"github.com/henrixapp/webdav-submission/server/admin"
	"github.com/henrixapp/webdav-submission/server/auth"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"golang.org/x/net/webdav"
	"google.golang.org/grpc"
)

const bucketName = "mybucket"

type MinioParams struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
}
type SharedWebDavFS struct {
	mampfClient               auth.MaMpfClient
	mampfTermsClient          pb.MaMpfTermServiceClient
	mampfLectureServiceClient pb.MaMpfLectureServiceClient
	minioClient               *minio.Client
	submissionRepository      admin.SubmissionRepository
}

func NewSharedWebDavFS(minioParams MinioParams, mampfParams auth.MampfParams, conn *grpc.ClientConn, submissionRepository admin.SubmissionRepository) SharedWebDavFS {
	minioClient, err := minio.New(minioParams.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioParams.AccessKeyID, minioParams.SecretAccessKey, ""),
		Secure: minioParams.UseSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{Region: "default"})
	if err != nil {
		log.Println(err)
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(context.Background(), bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln("NOPE", err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}
	return SharedWebDavFS{minioClient: minioClient, mampfClient: auth.MaMpfClientImpl{}, mampfTermsClient: pb.NewMaMpfTermServiceClient(conn),
		mampfLectureServiceClient: pb.NewMaMpfLectureServiceClient(conn),
		submissionRepository:      submissionRepository}
}

func (swdfs SharedWebDavFS) Mkdir(ctx context.Context, name string, perm os.FileMode) error {

	path := strings.Split(name, "/")
	//Else return file
	if len(path) >= 5 && path[3] != "" && path[4] != "" {
		//overview over one submission
		if strings.LastIndex(path[3], "$") != -1 {
			submissionId := strings.Split(path[3], "$")[len(strings.Split(path[3], "$"))-1]
			s, _ := uuid.Parse(submissionId)
			//FIXME(henrik): Permission check

			submissionsFiles, _ := swdfs.submissionRepository.FindSubmissionsFilesBySubmissionID(s, swdfs.minioClient)

			sf, ok := submissionsFiles[path[4]]
			if ok {
				var err error
				sf, _, err = swdfs.submissionRepository.TraverseToFile(sf, path[5:len(path)-1], swdfs.minioClient)
				if err != nil {
					return err
				}
			}
			//implicitly ID is null, if root entry
			parent := sf.ID
			swdfs.submissionRepository.CreateSubmissionsFile(s, parent, true, path[len(path)-1], int(ctx.Value("userID").(int32)), swdfs.minioClient)

		}
	}
	return nil
}
func (swdfs SharedWebDavFS) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	path := strings.Split(name, "/")
	if len(path) == 2 && path[1] == "" {
		terms, _ := swdfs.mampfTermsClient.GetTerms(ctx, &pb.TermsRequest{UserId: ctx.Value("userID").(int32)})
		return admin.TermsOverview{Terms: terms.GetTerms()}, nil
	}
	if len(path) == 3 && path[1] != "" {
		if strings.LastIndex(path[1], "-") != -1 {
			termId := strings.Split(path[1], "-")[len(strings.Split(path[1], "-"))-1]
			t, _ := strconv.ParseInt(termId, 10, 64)
			lectures, _ := swdfs.mampfLectureServiceClient.GetLecturesForUser(ctx, &pb.LecturesByUserAndTermRequest{TermId: int32(t), UserId: ctx.Value("userID").(int32)})
			return admin.LecturesOverview{Lectures: lectures.GetLectures()}, nil
		}
	}
	//Find submissions for certain lecture
	if len(path) == 4 && path[2] != "" && path[3] == "" {
		if strings.LastIndex(path[2], "-") != -1 {
			lectureId := strings.Split(path[2], "-")[len(strings.Split(path[2], "-"))-1]
			l, _ := strconv.ParseInt(lectureId, 10, 64)
			submissions, _ := swdfs.submissionRepository.FindSubmissionsBySubmitterIDAndLectureID(int(ctx.Value("userID").(int32)), int(l))
			entries := make([]admin.Entry, len(submissions))
			for i, v := range submissions {
				entries[i] = v
			}
			return admin.Overview{Entries: entries}, nil
		}
	}
	if (len(path) == 5 && path[3] != "" && path[4] == "") || (len(path) == 4 && path[3] != "") {
		//overview over one submission
		if strings.LastIndex(path[3], "$") != -1 {
			submissionId := strings.Split(path[3], "$")[len(strings.Split(path[3], "$"))-1]
			s, _ := uuid.Parse(submissionId)
			//FIXME(henrik): Permission check
			submissionsFiles, _ := swdfs.submissionRepository.FindSubmissionsFilesBySubmissionID(s, swdfs.minioClient)
			return admin.Submission{SubmissionsFiles: submissionsFiles}, nil
		}
	}
	//Else return file
	if len(path) >= 5 && path[3] != "" && path[4] != "" {
		//overview over one submission
		if strings.LastIndex(path[3], "$") != -1 {
			submissionId := strings.Split(path[3], "$")[len(strings.Split(path[3], "$"))-1]
			s, _ := uuid.Parse(submissionId)
			//FIXME(henrik): Permission check
			submissionsFiles, _ := swdfs.submissionRepository.FindSubmissionsFilesBySubmissionID(s, swdfs.minioClient)

			sf, ok := submissionsFiles[path[4]]

			var isParent bool
			if ok {
				sf, isParent, _ = swdfs.submissionRepository.TraverseToFile(sf, path[5:], swdfs.minioClient)
				if isParent {
					if flag&os.O_CREATE != 0 {
						return swdfs.submissionRepository.CreateSubmissionsFile(s, sf.ID, false, path[len(path)-1], int(ctx.Value("userID").(int32)), swdfs.minioClient)
					}
				}
			}
			if sf.ID != uuid.Nil && !isParent {
				return sf, nil
			}
			//implicitly ID is null, if root entry
			parent := sf.ID
			//TODO(henrik) what if not all exist?
			//NOT FOUND --> Create file
			if flag&os.O_CREATE != 0 {
				if flag&os.O_EXCL != 0 && sf.ID != uuid.Nil {
					return nil, os.ErrExist
				}
				if sf.ID == uuid.Nil {
					return swdfs.submissionRepository.CreateSubmissionsFile(s, parent, false, path[len(path)-1], int(ctx.Value("userID").(int32)), swdfs.minioClient)
				}
			}

		}
	}
	return File{}, os.ErrNotExist
}
func (swdfs SharedWebDavFS) RemoveAll(ctx context.Context, name string) error {
	path := strings.Split(name, "/")
	//Else return file
	if len(path) >= 5 && path[3] != "" && path[4] != "" {
		//overview over one submission
		if strings.LastIndex(path[3], "$") != -1 {
			submissionId := strings.Split(path[3], "$")[len(strings.Split(path[3], "$"))-1]
			s, _ := uuid.Parse(submissionId)
			//FIXME(henrik): Permission check
			submissionsFiles, _ := swdfs.submissionRepository.FindSubmissionsFilesBySubmissionID(s, swdfs.minioClient)

			sf, ok := submissionsFiles[path[4]]
			if ok {
				var err error
				sf, _, err = swdfs.submissionRepository.TraverseToFile(sf, path[5:], swdfs.minioClient)
				if err != nil {
					return err
				}
				if !sf.IsDir() {
					swdfs.minioClient.RemoveObject(ctx, bucketName, sf.ID.String(), minio.RemoveObjectOptions{})
				}
				//FIXME(henrik): recursive
				swdfs.submissionRepository.DeleteSubmissionsFile(sf.ID)
				return nil

			}
			//NOT FOUND --> Create file
			return os.ErrNotExist
		}
	}
	return nil
}
func (SharedWebDavFS) Rename(ctx context.Context, oldName, newName string) error {
	//TODO(henrik): Implement it
	return nil
}
func (swdfs SharedWebDavFS) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	path := strings.Split(name, "/")
	if len(path) == 2 && path[1] == "" {
		return admin.DirInfo{}, nil
	}
	if len(path) == 2 && path[1] != "" {
		if strings.LastIndex(path[1], "-") != -1 {
			termId := strings.Split(path[1], "-")[len(strings.Split(path[1], "-"))-1]
			t, _ := strconv.ParseInt(termId, 10, 64)
			lectures, _ := swdfs.mampfLectureServiceClient.GetLecturesForUser(ctx, &pb.LecturesByUserAndTermRequest{TermId: int32(t), UserId: ctx.Value("userID").(int32)})
			return admin.LecturesOverview{Lectures: lectures.GetLectures()}.Stat()
		}
	}
	//Else return file
	if len(path) >= 5 && path[3] != "" && path[4] != "" {
		//overview over one submission
		if strings.LastIndex(path[3], "$") != -1 {
			submissionId := strings.Split(path[3], "$")[len(strings.Split(path[3], "$"))-1]
			s, _ := uuid.Parse(submissionId)
			//FIXME(henrik): Permission check
			submissionsFiles, _ := swdfs.submissionRepository.FindSubmissionsFilesBySubmissionID(s, swdfs.minioClient)

			sf, ok := submissionsFiles[path[4]]
			if ok {
				sf, _, _ = swdfs.submissionRepository.TraverseToFile(sf, path[5:], swdfs.minioClient)
			} else {
				return FileInfo{}, os.ErrNotExist
			}
			if sf.ID != uuid.Nil {
				return sf, nil
			}
			//implicitly ID is null, if root entry
		}
	}
	return FileInfo{}, nil
}
