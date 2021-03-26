package admin

import (
	"fmt"
	"io/fs"
	"time"

	pb "github.com/henrixapp/mampf-rpc/grpc"
)

//TermsOverview is a directory, implements FileInfo
type TermsOverview struct {
	Terms []*pb.Term
}

//TermsOverview is a directory, implements FileInfo
type LecturesOverview struct {
	Lectures []*pb.Lecture
}

type DirInfo struct {
	Name_ string
}
type Entry interface {
	NameWithId() string
}
type Overview struct {
	Entries []Entry
}

func (termsO TermsOverview) Readdir(count int) ([]fs.FileInfo, error) {
	res := make([]fs.FileInfo, len(termsO.Terms))
	for i := range termsO.Terms {
		res[i] = DirInfo{Name_: termsO.Terms[i].GetSeason() + fmt.Sprint(termsO.Terms[i].GetYear()) + "-" + fmt.Sprint(termsO.Terms[i].GetId())}
	}
	return res, nil
}
func (termsO TermsOverview) Stat() (fs.FileInfo, error) {
	return DirInfo{}, nil
}
func (termsO TermsOverview) Close() error {
	return nil
}

func (termsO TermsOverview) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (termsO TermsOverview) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}
func (termsO TermsOverview) Write(p []byte) (n int, err error) {
	return 0, nil
}

func (t DirInfo) IsDir() bool {
	return true
}

func (t DirInfo) ModTime() time.Time {
	return time.Now()
}

func (t DirInfo) Name() string {
	return t.Name_
}

func (t DirInfo) Mode() fs.FileMode {
	return fs.ModeDir
}
func (t DirInfo) Size() int64 {
	return 4096
}

func (t DirInfo) Sys() interface{} {
	return nil
}

func (lecturesO LecturesOverview) Readdir(count int) ([]fs.FileInfo, error) {
	res := make([]fs.FileInfo, len(lecturesO.Lectures))
	for i := range lecturesO.Lectures {
		res[i] = DirInfo{Name_: lecturesO.Lectures[i].GetCourse().GetTitle() + "-" + fmt.Sprint(lecturesO.Lectures[i].GetId())}
	}
	return res, nil
}
func (lecturesO LecturesOverview) Stat() (fs.FileInfo, error) {
	return DirInfo{}, nil
}
func (lecturesO LecturesOverview) Close() error {
	return nil
}

func (lecturesO LecturesOverview) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (lecturesO LecturesOverview) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}
func (lecturesO LecturesOverview) Write(p []byte) (n int, err error) {
	return 0, nil
}

func (o Overview) Readdir(count int) ([]fs.FileInfo, error) {
	res := make([]fs.FileInfo, len(o.Entries))
	for i := range o.Entries {
		res[i] = DirInfo{Name_: o.Entries[i].NameWithId()}
	}
	return res, nil
}
func (o Overview) Stat() (fs.FileInfo, error) {
	return DirInfo{}, nil
}
func (o Overview) Close() error {
	return nil
}

func (o Overview) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (o Overview) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}
func (o Overview) Write(p []byte) (n int, err error) {
	return 0, nil
}
