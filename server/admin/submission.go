package admin

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/henrixapp/webdav-submission/server/db"
	"github.com/minio/minio-go/v7"
)

//TODO(henrik): Size limits!
//only one bucket for all TODO(henrik): Where to split?
const bucketName = "mybucket"

//Submission is the handin of a group of Students, visible in the Filesystem
type Submission struct {
	AssignmentID             uuid.UUID
	Assignment               Assignment `gorm:"foreignKey:AssignmentID"`
	TutorialID               uuid.UUID
	Tutorial                 Tutorial `gorm:"foreignKey:TutorialID"`
	Token                    string
	ManuscriptData           string
	LastModificationByUserAt time.Time
	Accepted                 bool

	db.BaseObject

	SubmissionsFiles []SubmissionsFile
}

func (s Submission) NameWithId() string {
	return fmt.Sprint(s.Assignment.Title, "$", s.ID)
}

//Submitter is the permission to upload to a submission
type Submitter struct {
	UserID       int
	SubmissionID uuid.UUID
}

//Invitation is the permission to become a submitter.
type Invitation struct {
	db.BaseObject
	InvitedUserID  int
	InvitingUserID int
	SubmissionID   uuid.UUID
}

//SubmissionsFileInfo mimics a solution/submission part
//Implements a FileInfo
type SubmissionsFileInfo struct {
	db.BaseObject
	CreatedAt time.Time
	UpdatedAt time.Time
	Name_     string
	//IsSolution marks that an object is editable by tutor
	IsSolution bool
	//IsVisible determines whether file is readable
	IsVisible    bool
	LastEditedBy int
	minioClient  *minio.Client
	buffer       *fileBuffer
}
type SubmissionsFile struct {
	SubmissionID uuid.UUID
	SubmissionsFileInfo
}

func (submissionsFile SubmissionsFile) Readdir(count int) ([]fs.FileInfo, error) {
	res := make([]fs.FileInfo, 0)
	return res, nil
}
func (submissionsFile SubmissionsFile) Stat() (fs.FileInfo, error) {
	return submissionsFile.SubmissionsFileInfo, nil
}
func (submissionsFile SubmissionsFile) Close() error {
	//TODO(henrik): Clean up?
	return nil
}

//TODO(henrik): init filemust open file
func (submissionsFile SubmissionsFile) Read(p []byte) (n int, err error) {
	if submissionsFile.buffer != nil {
		return submissionsFile.buffer.Read(p)
	} else {
		submissionsFile.get()
	}

	return 0, nil
}

func (submissionsFile SubmissionsFile) Seek(offset int64, whence int) (int64, error) {
	return submissionsFile.buffer.Seek(offset, whence)
}
func (submissionsFile SubmissionsFile) Write(p []byte) (int, error) {
	//TODO(henrik): IF dir is supported later on... restrict it here
	n, err := submissionsFile.buffer.Write(p)
	if submissionsFile.buffer != nil {
		submissionsFile.buffer.pos = 0
		submissionsFile.minioClient.PutObject(context.Background(), bucketName, submissionsFile.ID.String(), submissionsFile.buffer, int64(len(submissionsFile.buffer.data)), minio.PutObjectOptions{ContentType: "application/octet-stream"})
	}
	return n, err
}

func (t SubmissionsFileInfo) IsDir() bool {
	return false
}

func (t SubmissionsFileInfo) ModTime() time.Time {
	return t.UpdatedAt
}

func (t SubmissionsFileInfo) Name() string {
	return t.Name_
}

func (t SubmissionsFileInfo) Mode() fs.FileMode {
	//TODO(henrik): Implement write protection here
	return 777
}
func (t SubmissionsFileInfo) Size() int64 {
	objInfo, err := t.minioClient.StatObject(context.Background(), bucketName, t.ID.String(), minio.StatObjectOptions{})
	if err != nil {
		log.Println(err)
		return 0
	}
	return objInfo.Size
}

func (t SubmissionsFileInfo) Sys() interface{} {
	return nil
}

type fileBuffer struct {
	pos  int
	data []byte
	mu   sync.Mutex
}

func (fb *fileBuffer) Write(p []byte) (n int, err error) {
	lenp := len(p)
	fb.mu.Lock()
	defer fb.mu.Unlock()

	if fb.pos < len(fb.data) {
		n := copy(fb.data[fb.pos:], p)
		fb.pos += n
		p = p[n:]
	} else if fb.pos > len(fb.data) {
		if fb.pos <= cap(fb.data) {
			oldLen := len(fb.data)
			fb.data = fb.data[:fb.pos]
			hole := fb.data[oldLen:]
			for i := range hole {
				hole[i] = 0
			}
		} else {
			d := make([]byte, fb.pos, fb.pos+len(p))
			copy(d, fb.data)
			fb.data = d
		}
	}

	if len(p) > 0 {
		fb.data = append(fb.data, p...)
		fb.pos = len(fb.data)
	}
	return lenp, nil
}

func (f *fileBuffer) Read(p []byte) (int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.pos >= len(f.data) {
		return 0, io.EOF
	}
	n := copy(p, f.data[f.pos:])
	f.pos += n
	return n, nil
}

func (f *fileBuffer) Seek(offset int64, whence int) (int64, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	npos := f.pos
	switch whence {
	case os.SEEK_SET:
		npos = int(offset)
	case os.SEEK_CUR:
		npos += int(offset)
	case os.SEEK_END:
		npos = len(f.data) + int(offset)
	default:
		npos = -1
	}
	if npos < 0 {
		return 0, os.ErrInvalid
	}
	f.pos = npos
	return int64(f.pos), nil
}

func (o Submission) Readdir(count int) ([]fs.FileInfo, error) {
	res := make([]fs.FileInfo, len(o.SubmissionsFiles))
	for i := range o.SubmissionsFiles {
		res[i] = o.SubmissionsFiles[i]
	}
	return res, nil
}
func (o Submission) Stat() (fs.FileInfo, error) {
	return DirInfo{}, nil
}
func (o Submission) Close() error {
	return nil
}

func (o Submission) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (o Submission) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}
func (o Submission) Write(p []byte) (n int, err error) {
	return 0, nil
}

func (sf SubmissionsFile) get() {
	sf.buffer = &fileBuffer{data: make([]byte, 0)}
	object, err := sf.minioClient.GetObject(context.Background(), bucketName, sf.ID.String(), minio.GetObjectOptions{})
	if err != nil {
		log.Println(err)
		return
	}
	if _, err = io.Copy(sf.buffer, object); err != nil {
		fmt.Println(err)
		return
	}
}
