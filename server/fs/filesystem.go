package fs

import (
	"context"
	"errors"
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

	buckets, err := minioClient.ListBuckets(context.Background())
	log.Println(buckets)
	log.Println(err)
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

func (SharedWebDavFS) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	log.Println(ctx)
	log.Println(name, perm)
	return nil
}
func (swdfs SharedWebDavFS) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	log.Println(ctx, name, flag, perm)
	log.Println(ctx, name)
	path := strings.Split(name, "/")
	log.Println(path, len(path))
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
			log.Println("L:", l, "U:", ctx.Value("userID").(int32))
			submissions, _ := swdfs.submissionRepository.FindSubmissionsBySubmitterIDAndLectureID(int(ctx.Value("userID").(int32)), int(l))
			log.Println(submissions)
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
			log.Println("s:", s, "U:", ctx.Value("userID").(int32))
			//FIXME(henrik): Permission check
			submissionsFiles, _ := swdfs.submissionRepository.FindSubmissionsFilesBySubmissionID(s, swdfs.minioClient)
			log.Println("ssf", submissionsFiles)
			return admin.Submission{SubmissionsFiles: submissionsFiles}, nil
		}
	}
	//Else return file
	if len(path) == 5 && path[3] != "" && path[4] != "" {
		//overview over one submission
		if strings.LastIndex(path[3], "$") != -1 {
			submissionId := strings.Split(path[3], "$")[len(strings.Split(path[3], "$"))-1]
			s, _ := uuid.Parse(submissionId)
			log.Println("s:", s, "U:", ctx.Value("userID").(int32))
			//FIXME(henrik): Permission check
			submissionsFiles, _ := swdfs.submissionRepository.FindSubmissionsFilesBySubmissionID(s, swdfs.minioClient)
			for _, v := range submissionsFiles {
				if v.Name_ == path[4] {
					return v, nil
				}
			}
			//NOT FOUND --> Create file
			log.Println("P:", perm)
			log.Println("F:", flag)
			return swdfs.submissionRepository.CreateSubmissionsFile(s, path[4], int(ctx.Value("userID").(int32)), swdfs.minioClient)
		}
	}
	log.Println(ctx.Value("userID"))
	return File{}, nil
}
func (swdfs SharedWebDavFS) RemoveAll(ctx context.Context, name string) error {
	log.Println("remove", ctx, name)
	path := strings.Split(name, "/")
	//Else return file
	if len(path) == 5 && path[3] != "" && path[4] != "" {
		//overview over one submission
		if strings.LastIndex(path[3], "$") != -1 {
			submissionId := strings.Split(path[3], "$")[len(strings.Split(path[3], "$"))-1]
			s, _ := uuid.Parse(submissionId)
			log.Println("s:", s, "U:", ctx.Value("userID").(int32))
			//FIXME(henrik): Permission check
			submissionsFiles, _ := swdfs.submissionRepository.FindSubmissionsFilesBySubmissionID(s, swdfs.minioClient)
			for _, v := range submissionsFiles {
				if v.Name_ == path[4] {
					//delete v
					swdfs.minioClient.RemoveObject(ctx, bucketName, v.ID.String(), minio.RemoveObjectOptions{})
					swdfs.submissionRepository.DeleteSubmissionsFile(v.ID)
					return nil
				}
			}
			//NOT FOUND --> Create file
			return errors.New("file not found")
		}
	}
	return nil
}
func (SharedWebDavFS) Rename(ctx context.Context, oldName, newName string) error {
	log.Println(ctx, oldName, newName)
	return nil
}
func (swdfs SharedWebDavFS) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	log.Println(ctx, name)
	path := strings.Split(name, "/")
	log.Println("Stat", path, len(path))
	if len(path) == 2 && path[1] == "" {
		return admin.DirInfo{}, nil
	}
	log.Println(ctx.Value("userID"))
	if len(path) == 2 && path[1] != "" {
		if strings.LastIndex(path[1], "-") != -1 {
			termId := strings.Split(path[1], "-")[len(strings.Split(path[1], "-"))-1]
			t, _ := strconv.ParseInt(termId, 10, 64)
			lectures, _ := swdfs.mampfLectureServiceClient.GetLecturesForUser(ctx, &pb.LecturesByUserAndTermRequest{TermId: int32(t), UserId: ctx.Value("userID").(int32)})
			return admin.LecturesOverview{Lectures: lectures.GetLectures()}.Stat()
		}
	}
	return FileInfo{}, nil
}
