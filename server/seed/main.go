package main

import (
	"log"
	"time"

	pb "github.com/henrixapp/mampf-rpc/grpc"
	"github.com/minio/minio-go/v7"
	"google.golang.org/grpc"

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
	conn, err := grpc.Dial("localhost:9001", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	lectureService := pb.NewMaMpfLectureServiceClient(conn)
	db := initializeDB()
	submissionRep := admin.NewSubmissionRepositoryGorm(db, &minio.Client{}, lectureService)
	assId, err := submissionRep.CreateAssignment(admin.Assignment{LectureID: 27, MediumID: 51, Title: "Ters Übung", Deadline: time.Now().Add(time.Hour * 5), AcceptedFileType: ".pdf", MaxFileCount: 5}, 1)
	if err != nil {
		log.Fatalln(err)
	}
	tutid, err := submissionRep.CreateTutorial(admin.Tutorial{Title: "Tutorial 1", LectureID: 27}, 1)
	if err != nil {
		log.Fatalln(err)
	}
	submissionRep.CreateTutor(admin.Tutor{TutorialID: tutid, UserID: 5}, 1)
	id, err := submissionRep.CreateSubmission(admin.Submission{AssignmentID: assId, TutorialID: tutid, Token: "blabla", LastModificationByUserAt: time.Now()}, 1)
	files, err := submissionRep.FindSubmissionsFilesBySubmissionID(id, 5)
	log.Println(files)
	log.Println(err)
	log.Println("Successfully seeded")
}
