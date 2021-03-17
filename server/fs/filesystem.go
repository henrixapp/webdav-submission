package fs

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"

	pb "github.com/henrixapp/mampf-rpc/grpc"
	"github.com/henrixapp/webdav-submission/server/admin"
	"github.com/henrixapp/webdav-submission/server/auth"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"golang.org/x/net/webdav"
	"google.golang.org/grpc"
)

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
}

func NewSharedWebDavFS(minioParams MinioParams, mampfParams auth.MampfParams, conn *grpc.ClientConn) SharedWebDavFS {
	minioClient, err := minio.New(minioParams.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioParams.AccessKeyID, minioParams.SecretAccessKey, ""),
		Secure: minioParams.UseSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}
	return SharedWebDavFS{minioClient: minioClient, mampfClient: auth.MaMpfClientImpl{}, mampfTermsClient: pb.NewMaMpfTermServiceClient(conn),
		mampfLectureServiceClient: pb.NewMaMpfLectureServiceClient(conn)}
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
	log.Println(ctx.Value("userID"))
	return File{}, nil
}
func (SharedWebDavFS) RemoveAll(ctx context.Context, name string) error {
	log.Println(ctx, name)
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
