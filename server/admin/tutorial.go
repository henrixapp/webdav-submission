package admin

import (
	"github.com/google/uuid"
	"github.com/henrixapp/webdav-submission/server/db"
)

//Tutorial is the group in which submissions are corrected --> User
type Tutorial struct {
	db.BaseObject
	Title     string `json:"title,omitempty"`
	LectureID int    `json:"lectureID,omitempty"`
}

//Tutor is the permission to access certain submissions
type Tutor struct {
	UserID     int       `json:"userID,omitempty"`
	TutorialID uuid.UUID `gorm:"type:uuid" json:"tutorialID,omitempty"`
}
