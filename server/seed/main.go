package main

import (
	"github.com/minio/minio-go/v7"
	"log"
	"time"

	"github.com/henrixapp/webdav-submission/server/admin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func initializeDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatalln("Fatal error while connecting to DB:", err)
	}
	return db.Debug()
}
func main() {
	db := initializeDB()
	submissionRep := admin.NewSubmissionRepositoryGorm(db, &minio.Client{})
	assId, _ := submissionRep.CreateAssignment(admin.Assignment{LectureID: 1, MediumID: 51, Title: "Ters Übung", Deadline: time.Now().Add(time.Hour * 5), AcceptedFileType: ".pdf", MaxFileCount: 5})
	tutid, _ := submissionRep.CreateTutorial(admin.Tutorial{Title: "Tutorial 1", LectureID: 1})
	submissionRep.CreateTutor(admin.Tutor{TutorialID: tutid, UserID: 5})
	submissionRep.CreateSubmission(admin.Submission{AssignmentID: assId, TutorialID: tutid, Token: "blabla", LastModificationByUserAt: time.Now()}, 1)
	log.Println("Successfully seeded")
}
