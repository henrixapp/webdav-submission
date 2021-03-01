package fs

import (
	"context"
	"log"
	"os"

	"github.com/henrixapp/webdav-submission/server/auth"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"golang.org/x/net/webdav"
)

type MinioParams struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
}
type SharedWebDavFS struct {
	mampfClient auth.MaMpfClient
	minioClient *minio.Client
}

func NewSharedWebDavFS(minioParams MinioParams, mampfParams auth.MampfParams) SharedWebDavFS {
	minioClient, err := minio.New(minioParams.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioParams.AccessKeyID, minioParams.SecretAccessKey, ""),
		Secure: minioParams.UseSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}
	return SharedWebDavFS{minioClient: minioClient, mampfClient: auth.MaMpfClientImpl{}}
}

func (SharedWebDavFS) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	log.Println(ctx)
	log.Println(name, perm)
	return nil
}
func (SharedWebDavFS) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	log.Println(ctx, name, flag, perm)
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
func (SharedWebDavFS) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	log.Println(ctx, name)
	return FileInfo{}, nil
}
