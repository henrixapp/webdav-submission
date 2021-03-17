package fs

import (
	"io/fs"
	"log"
	"time"
)

//This is a low level File
type File struct{}
type FileInfo struct{}

func (file File) Readdir(count int) ([]fs.FileInfo, error) {
	res := make([]fs.FileInfo, 2)
	res[0] = FileInfo{}
	res[1] = FileInfo{}
	log.Println("FOlder returned")
	return res, nil
}
func (file File) Stat() (fs.FileInfo, error) {
	return FileInfo{}, nil
}
func (file File) Close() error {
	return nil
}

func (file File) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (file File) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}
func (file File) Write(p []byte) (n int, err error) {
	return 0, nil
}
func (f FileInfo) IsDir() bool {
	return true
}

func (f FileInfo) ModTime() time.Time {
	return time.Now()
}

func (f FileInfo) Name() string {
	return "test"
}

func (f FileInfo) Mode() fs.FileMode {
	return fs.ModeDir
}
func (fileInfo FileInfo) Size() int64 {
	return 4096
}

func (fileInfo FileInfo) Sys() interface{} {
	return nil
}
