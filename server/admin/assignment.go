package admin

import (
	"time"

	"github.com/henrixapp/webdav-submission/server/db"
)

//courses are defined in mampf. so not gone define them
//here
//Assignment is one homework
type Assignment struct {
	//LectureID is the ID of the lecture for this assignment --> from mampf-rpc
	LectureID int `json:"lectureID,omitempty"`
	//MediumID is the ID of the media the assigment is designed for
	MediumID int `json:"mediumID,omitempty"`
	// Title is  a descriptive title of an exercise
	Title string `json:"deadline,omitempty"`
	//Deadline is the end timestamp for submissions, editing afterwards might be forbidden
	Deadline time.Time `json:"title,omitempty"`
	//AcceptedFileType is the filetype accepted TODO(henrik): what's about directories
	AcceptedFileType string `json:"acceptedFileType,omitempty"`
	//MaxFileCount gives the maximum number of Files
	MaxFileCount int `json:"maxFileCount,omitempty"`
	db.BaseObject
}
