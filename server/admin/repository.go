package admin

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

type SubmissionRepository interface {
	CreateTutorial(tutorial Tutorial, userID int) (uuid.UUID, error)
	CreateTutor(tutor Tutor, userID int) error
	UpdateTutorial(tutorial Tutorial, userID int) error
	//DeleteTutorial deletes the tutorial, returns err!=nil if dependencies exist.
	DeleteTutorial(tutorialID uuid.UUID, userID int) error
	FindTutorialByLectureID(lectureID int, userID int) ([]Tutorial, error)
	FindTutorialByID(uuid uuid.UUID, userID int) (Tutorial, error)

	CreateAssignment(assignment Assignment, userID int) (uuid.UUID, error)
	UpdateAssignment(assignment Assignment, userID int) error
	DeleteAssignment(assignmentID uuid.UUID, userID int) error
	FindAssignmentByLectureID(lectureID int, userID int) ([]Assignment, error)
	FindAssignmentByID(uuid uuid.UUID, userID int) (Assignment, error)
	//Also generates correct submission right
	CreateSubmission(submission Submission, userID int) error
	UpdateSubmission(submission Submission, userID int) error
	DeleteSubmission(submissionID uuid.UUID, userID int) error
	FindSubmissionsBySubmitterIDAndLectureID(submitterID int, lectureID int, userID int) ([]Submission, error)
	FindSubmissionByAssignmentIDAndTutorialID(assignmentID uuid.UUID, tutorialID uuid.UUID, userID int) (Submission, error)
	FindSubmissionsByLectureIDAndTutorialID(lectureID int, tutorialID uuid.UUID, userID int) ([]Submission, error)
	FindSubmissionByID(uuid uuid.UUID, userID int) (Submission, error)
	FindSubmissionBySubmitterIDAndAssignmentID(submitterID int, assignmentID uuid.UUID, userID int) (Submission, error)
	SaveInviteToSubmission(invite Invitation, userID int) error
	AcceptInvitation(invite Invitation, userID int) error

	//Returns DIR
	FindSubmissionsFilesBySubmissionID(uuid uuid.UUID, userID int) (map[string]SubmissionsFile, error)
	CountSubmissionsFilesBySubmissionID(uuid uuid.UUID, userID int) (int, error)
	FindSubmissionsSubFilesBySubmissionID(parent uuid.UUID, userID int) (map[string]SubmissionsFile, error)
	CreateSubmissionsFile(submissionUUID uuid.UUID, parent uuid.UUID, isDir bool, name string, userID int) (SubmissionsFile, error)
	//FIXME(henrik): DeleteSubmissionsFile does not delete the bucket file
	DeleteSubmissionsFile(id uuid.UUID, userID int) error
	//bool isParent
	TraverseToFile(root SubmissionsFile, path []string, userID int) (SubmissionsFile, bool, error)
}

type SubmissionRepositoryGorm struct {
	db *gorm.DB
	//TODO(henrik): remove policies on deletion
	enforcer *casbin.Enforcer
	//the actual location for the files
	minioClient *minio.Client
}

func NewSubmissionRepositoryGorm(db *gorm.DB, client *minio.Client) SubmissionRepository {
	db.AutoMigrate(&Tutor{})
	db.AutoMigrate(&Tutorial{})
	db.AutoMigrate(&Assignment{})
	db.AutoMigrate(&Submission{})
	db.AutoMigrate(&Submitter{})
	db.AutoMigrate(&Invitation{})
	db.AutoMigrate(&SubmissionsFile{})
	a, _ := gormadapter.NewAdapterByDB(db)
	e, _ := casbin.NewEnforcer("rbac.conf", a)
	return SubmissionRepositoryGorm{db: db, enforcer: e, minioClient: client}
}

func (srg SubmissionRepositoryGorm) CreateTutorial(tutorial Tutorial, userID int) (uuid.UUID, error) {
	srg.db.Create(&tutorial) //TODO(henrik): error handleing
	return tutorial.ID, nil
}
func (srg SubmissionRepositoryGorm) UpdateTutorial(tutorial Tutorial, userID int) error {
	srg.db.Save(&tutorial)
	return nil
}
func (srg SubmissionRepositoryGorm) DeleteTutorial(tutorialID uuid.UUID, userID int) error {
	srg.db.Where("id = ?", tutorialID).Delete(&Tutorial{})
	return nil
}
func (srg SubmissionRepositoryGorm) FindTutorialByLectureID(lectureID int, userID int) ([]Tutorial, error) {
	var tuts []Tutorial
	srg.db.Where("lecture_id = ?", lectureID).Find(&tuts)
	return tuts, nil
}
func (srg SubmissionRepositoryGorm) FindTutorialByID(uuid uuid.UUID, userID int) (Tutorial, error) {
	var tutorial Tutorial
	srg.db.First(&tutorial, "id = ?", uuid)
	return tutorial, nil
}
func (srg SubmissionRepositoryGorm) CreateAssignment(assignment Assignment, userID int) (uuid.UUID, error) {
	srg.db.Create(&assignment)
	return assignment.ID, nil
}
func (srg SubmissionRepositoryGorm) UpdateAssignment(assignment Assignment, userID int) error {
	srg.db.Save(assignment)
	return nil
}
func (srg SubmissionRepositoryGorm) DeleteAssignment(assignmentID uuid.UUID, userID int) error {
	srg.db.Where("id = ?", assignmentID).Delete(&Assignment{})
	return nil
}
func (srg SubmissionRepositoryGorm) FindAssignmentByLectureID(lectureID int, userID int) ([]Assignment, error) {
	var assignments []Assignment
	srg.db.Where("lecture_id = ?", lectureID).Find(&assignments)
	return assignments, nil
}
func (srg SubmissionRepositoryGorm) FindAssignmentByID(uuid uuid.UUID, userID int) (Assignment, error) {
	var assignment Assignment
	srg.db.Where("id = ?", uuid).First(&assignment)
	return assignment, nil
}
func (srg SubmissionRepositoryGorm) CreateSubmission(submission Submission, userID int) error {
	//TODO(henrik): Check also time here
	srg.db.Create(&submission)
	srg.db.Create(&Submitter{UserID: userID, SubmissionID: submission.ID})
	return nil
}
func (srg SubmissionRepositoryGorm) UpdateSubmission(submission Submission, userID int) error {
	if submission.ID == uuid.Nil {
		return errors.New("Without id")
	}
	srg.db.Save(&submission)
	return nil
}
func (srg SubmissionRepositoryGorm) DeleteSubmission(submissionID uuid.UUID, userID int) error {
	if submissionID == uuid.Nil {
		return errors.New("Without id")
	}
	srg.db.Where("id = ?", submissionID).Delete(&Submission{})
	return nil
}
func (srg SubmissionRepositoryGorm) FindSubmissionByAssignmentIDAndTutorialID(assignmentID uuid.UUID, tutorialID uuid.UUID, userID int) (Submission, error) {
	var submission Submission
	srg.db.Where("assignment_id = ? AND tutorial_id = ?", assignmentID, tutorialID).Find(&submission)
	return submission, nil
}
func (srg SubmissionRepositoryGorm) FindSubmissionByID(uuid uuid.UUID, userID int) (Submission, error) {
	var sub []Submission
	srg.db.Preload("Assignment").Joins("JOIN assignments on submissions.assignment_id = assignments.id").Where("submissions.id = ?", uuid).Find(&sub)
	if len(sub) == 0 {
		return Submission{}, errors.New("not found")
	}
	return sub[0], nil
}
func (srg SubmissionRepositoryGorm) FindSubmissionBySubmitterIDAndAssignmentID(submitterID int, assignmentID uuid.UUID, userID int) (Submission, error) {
	var submission Submission
	srg.db.Joins("JOIN submitters ON submitters.submission_id = submissions.id AND submitters.user_id = ?", submitterID).Where("assignment_id = ?", assignmentID).Find(&submission)
	return submission, nil
}

func (srg SubmissionRepositoryGorm) FindSubmissionsBySubmitterIDAndLectureID(submitterID int, lectureId int, userID int) ([]Submission, error) {
	var submissions []Submission
	srg.db.Preload("Assignment").Joins("JOIN submitters ON submitters.submission_id = submissions.id AND submitters.user_id = ? JOIN assignments ON submissions.assignment_id = assignments.id", submitterID).Where("lecture_id = ?", lectureId).Find(&submissions)
	return submissions, nil
}

func (srg SubmissionRepositoryGorm) SaveInviteToSubmission(invite Invitation, userID int) error {
	//TODO(henrik): check if user is already a submitter
	srg.db.Save(&invite)
	return nil
}
func (srg SubmissionRepositoryGorm) AcceptInvitation(invite Invitation, userID int) error {
	//TODO(henrik): Check access rights
	srg.db.Create(&Submitter{SubmissionID: invite.SubmissionID, UserID: invite.InvitedUserID})
	srg.db.Delete(&invite)
	return nil
}

func (srg SubmissionRepositoryGorm) FindSubmissionsByLectureIDAndTutorialID(lectureID int, tutorialID uuid.UUID, userID int) ([]Submission, error) {
	var submissions []Submission
	srg.db.Joins("JOIN assigments ON assignments.id = submissions.assignment_id AND submissions.tutorial_id = ?", tutorialID).Where("lecture_id = ?", lectureID).Find(&submissions)
	return submissions, nil
}

func (srg SubmissionRepositoryGorm) CreateTutor(tutor Tutor, userID int) error {
	srg.db.Create(&tutor)
	return nil
}

func (srg SubmissionRepositoryGorm) FindSubmissionsFilesBySubmissionID(submissionID uuid.UUID, userID int) (map[string]SubmissionsFile, error) {
	var submissionsFiles []SubmissionsFile
	res := make(map[string]SubmissionsFile)
	srg.db.Where("submission_id = ? AND parent = ?", submissionID, uuid.Nil).Find(&submissionsFiles)
	for s, sf := range submissionsFiles {
		submissionsFiles[s].minioClient = srg.minioClient
		submissionsFiles[s].get()
		res[sf.Name()] = submissionsFiles[s]
	}
	return res, nil
}

func (srg SubmissionRepositoryGorm) CountSubmissionsFilesBySubmissionID(submissionID uuid.UUID, userID int) (int, error) {
	var res int64
	srg.db.Model(&SubmissionsFile{}).Where("submission_id = ?", submissionID).Count(&res)
	log.Println("Res:", res)
	return int(res), nil
}

//TODO(henrik): enforce max recursion limit
func (srg SubmissionRepositoryGorm) FindSubmissionsSubFilesBySubmissionID(parent uuid.UUID, userID int) (map[string]SubmissionsFile, error) {
	var submissionsFiles []SubmissionsFile
	res := make(map[string]SubmissionsFile)
	srg.db.Where("parent = ?", parent).Find(&submissionsFiles)
	for _, sf := range submissionsFiles {
		sf.minioClient = srg.minioClient
		if sf.IsDir() {
			// var err error
			// sf.Children, err = srg.FindSubmissionsSubFilesBySubmissionID(sf.ID, mC)
			// if err != nil {
			// 	return res, err
			// }
		} else {
			sf.get()
		}
		res[sf.Name_] = sf
	}
	return res, nil
}
func (srg SubmissionRepositoryGorm) CreateSubmissionsFile(submissionUUID uuid.UUID, parent uuid.UUID, isDir bool, name string, user int) (SubmissionsFile, error) {
	subFI := SubmissionsFile{}
	subFI.buffer = &fileBuffer{data: make([]byte, 0)}
	subFI.minioClient = srg.minioClient
	subFI.SubmissionID = submissionUUID
	subFI.LastEditedBy = user
	subFI.Name_ = name
	subFI.IsSolution = false
	subFI.IsVisible = true //TODO(henrik): Fixme
	subFI.Parent = parent
	subFI.Dir = isDir
	srg.db.Create(&subFI)
	if !subFI.IsDir() {
		subFI.buffer.pos = 0
		subFI.minioClient.PutObject(context.Background(), bucketName, subFI.ID.String(), subFI.buffer, int64(len(subFI.buffer.data)), minio.PutObjectOptions{ContentType: "application/octet-stream"})
	}
	return subFI, nil
}

func (srg SubmissionRepositoryGorm) DeleteSubmissionsFile(subID uuid.UUID, userID int) error {
	srg.db.Where("id = ?", subID).Delete(&SubmissionsFile{})
	return nil
}

//TraverseToFile Traverses To a specific File, does not open other files
func (srg SubmissionRepositoryGorm) TraverseToFile(root SubmissionsFile, path []string, userID int) (SubmissionsFile, bool, error) {
	var sf SubmissionsFile = root
	for _, p := range path {

		if p != "" {
			var submissionsFile SubmissionsFile
			srg.db.Where("parent = ? AND name_ = ? ", sf.ID, p).Find(&submissionsFile)
			if submissionsFile.ID == uuid.Nil {
				return sf, true, os.ErrNotExist //return last folder existing
			}
			sf = submissionsFile
			if !submissionsFile.IsDir() {
				//reached the file
				break
			}
		}
	}
	if !sf.IsDir() {
		sf.minioClient = srg.minioClient
		sf.get()
	} else {
		//load subfiles
		sf.Children, _ = srg.FindSubmissionsSubFilesBySubmissionID(sf.ID, userID)
	}
	return sf, false, nil
}
