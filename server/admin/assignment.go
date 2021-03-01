package admin

import (
	"time"

	"github.com/henrixapp/webdav-submission/server/db"
)

//courses are defined in mampf. so not gone define them
//here

type Assignment struct {
	LectureID        int
	MediumID         int
	Title            string
	Deadline         time.Time
	AcceptedFileType string
	db.BaseObjectBaseObject
}
