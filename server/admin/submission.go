package admin

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/henrixapp/webdav-submission/server/db"
)

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
}

func (s Submission) NameWithId() string {
	return fmt.Sprint(s.Assignment.Title, "-", s.ID)
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
