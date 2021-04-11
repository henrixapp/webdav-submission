package admin

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/google/uuid"
	pb "github.com/henrixapp/mampf-rpc/grpc"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

type SubmissionRepository interface {
	CreateTutorial(tutorial Tutorial, userID int32) (uuid.UUID, error)
	CreateTutor(tutor Tutor, userID int32) error
	DeleteTutor(tutor Tutor, userID int32) error
	FindTutorsByTutorialID(tutorialID uuid.UUID) ([]Tutor, error)
	UpdateTutorial(tutorial Tutorial, userID int32) error
	//DeleteTutorial deletes the tutorial, returns err!=nil if dependencies exist.
	DeleteTutorial(tutorialID uuid.UUID, userID int32) error
	FindTutorialByLectureID(lectureID int, userID int32) ([]Tutorial, error)
	FindTutorialByID(uuid uuid.UUID, userID int32) (Tutorial, error)

	CreateAssignment(assignment Assignment, userID int32) (uuid.UUID, error)
	UpdateAssignment(assignment Assignment, userID int32) error
	DeleteAssignment(assignmentID uuid.UUID, userID int32) error
	FindAssignmentByLectureID(lectureID int, userID int32) ([]Assignment, error)
	FindAssignmentByID(uuid uuid.UUID, userID int32) (Assignment, error)
	//Also generates correct submission right
	CreateSubmission(submission Submission, userID int32) (uuid.UUID, error)
	UpdateSubmission(submission Submission, userID int32) error
	DeleteSubmission(submissionID uuid.UUID, userID int32) error
	FindSubmissionsBySubmitterIDAndLectureID(submitterID int, lectureID int, userID int32) ([]Submission, error)
	FindSubmissionByAssignmentIDAndTutorialID(assignmentID uuid.UUID, tutorialID uuid.UUID, userID int32) (Submission, error)
	FindSubmissionsByLectureIDAndTutorialID(lectureID int, tutorialID uuid.UUID, userID int32) ([]Submission, error)
	FindSubmissionByID(uuid uuid.UUID, userID int32) (Submission, error)
	FindSubmissionBySubmitterIDAndAssignmentID(submitterID int, assignmentID uuid.UUID, userID int32) (Submission, error)

	SubmissionJoinByToken(token string, userID int32) (Submission, error)
	SaveInviteToSubmission(invite Invitation, userID int32) error
	AcceptInvitation(invite Invitation, userID int32) error
	FindInvitation(invitationID uuid.UUID) (Invitation, error)
	FindInvitationByUserID(userID int32) ([]Invitation, error)
	FindInvitationBySubmissionsID(submissionID uuid.UUID, userID int32) ([]Invitation, error)
	DeleteInvitation(invitation Invitation, userID int32) error
	//Returns DIR
	FindSubmissionsFilesBySubmissionID(uuid uuid.UUID, userID int32) (map[string]SubmissionsFile, error)
	CountSubmissionsFilesBySubmissionID(uuid uuid.UUID, userID int32) (int, error)
	FindSubmissionsSubFilesBySubmissionID(parent uuid.UUID, userID int32) (map[string]SubmissionsFile, error)
	CreateSubmissionsFile(submissionUUID uuid.UUID, parent uuid.UUID, isDir bool, name string, userID int32) (SubmissionsFile, error)
	//FIXME(henrik): DeleteSubmissionsFile does not delete the bucket file
	DeleteSubmissionsFile(id uuid.UUID, userID int32) error
	//bool isParent
	TraverseToFile(root SubmissionsFile, path []string, userID int32) (SubmissionsFile, bool, error)
}

type SubmissionRepositoryGorm struct {
	db *gorm.DB
	//TODO(henrik): remove policies on deletion
	enforcer *casbin.Enforcer
	//the actual location for the files
	minioClient               *minio.Client
	mampfLectureServiceClient pb.MaMpfLectureServiceClient
}

func NewSubmissionRepositoryGorm(db *gorm.DB, client *minio.Client, mampfLectureServiceClient pb.MaMpfLectureServiceClient) SubmissionRepository {
	db.AutoMigrate(&Tutor{})
	db.AutoMigrate(&Tutorial{})
	db.AutoMigrate(&Assignment{})
	db.AutoMigrate(&Submission{})
	db.AutoMigrate(&Submitter{})
	db.AutoMigrate(&Invitation{})
	db.AutoMigrate(&SubmissionsFile{})
	a, _ := gormadapter.NewAdapterByDB(db)
	e, _ := casbin.NewEnforcer("rbac.conf", a)

	return SubmissionRepositoryGorm{db: db, enforcer: e, minioClient: client, mampfLectureServiceClient: mampfLectureServiceClient}
}

func (srg SubmissionRepositoryGorm) CreateTutorial(tutorial Tutorial, userID int32) (uuid.UUID, error) {
	res, err := srg.mampfLectureServiceClient.GetIsEditor(context.TODO(), &pb.IsEditorRequest{User: userID, Lecture: int32(tutorial.LectureID)})
	if err != nil || !res.IsEditor {
		return uuid.Nil, errors.New("permission denied. not an editor")
	}
	srg.db.Create(&tutorial) //TODO(henrik): error handleing
	srg.enforcer.AddPolicy(tutorialCasbin(tutorial.ID), "", "tut")
	//TODO(henrik): check permissions

	return tutorial.ID, nil
}
func (srg SubmissionRepositoryGorm) UpdateTutorial(tutorial Tutorial, userID int32) error {
	res, err := srg.mampfLectureServiceClient.GetIsEditor(context.TODO(), &pb.IsEditorRequest{User: userID, Lecture: int32(tutorial.LectureID)})
	if err != nil || !res.IsEditor {
		return errors.New("permission denied. not an editor")
	}
	srg.db.Save(&tutorial)
	return nil
}
func (srg SubmissionRepositoryGorm) DeleteTutorial(tutorialID uuid.UUID, userID int32) error {
	tutorial, err := srg.FindTutorialByID(tutorialID, userID)
	if err != nil {
		return err
	}
	res, err := srg.mampfLectureServiceClient.GetIsEditor(context.TODO(), &pb.IsEditorRequest{User: userID, Lecture: int32(tutorial.LectureID)})
	if err != nil || !res.IsEditor {
		return errors.New("permission denied. not an editor")
	}
	srg.db.Where("id = ?", tutorialID).Delete(&Tutorial{})
	return nil
}
func (srg SubmissionRepositoryGorm) FindTutorialByLectureID(lectureID int, userID int32) ([]Tutorial, error) {
	res, err := srg.mampfLectureServiceClient.GetIsParticipantInLecture(context.TODO(), &pb.IsParticipantRequest{User: userID, Lecture: int32(lectureID)})
	if err != nil || !res.IsParticipant {
		return make([]Tutorial, 0), errors.New("permission denied. not a member of course")
	}
	var tuts []Tutorial
	srg.db.Where("lecture_id = ?", lectureID).Find(&tuts)
	//no permission check here
	return tuts, nil
}
func (srg SubmissionRepositoryGorm) FindTutorialByID(uuid uuid.UUID, userID int32) (Tutorial, error) {
	var tutorial Tutorial
	srg.db.First(&tutorial, "id = ?", uuid)
	//TODO(henrik): Discuss whether an check for permission must be done here
	return tutorial, nil
}
func (srg SubmissionRepositoryGorm) CreateAssignment(assignment Assignment, userID int32) (uuid.UUID, error) {
	res, err := srg.mampfLectureServiceClient.GetIsEditor(context.TODO(), &pb.IsEditorRequest{User: userID, Lecture: int32(assignment.LectureID)})
	if err != nil || !res.IsEditor {
		return uuid.Nil, errors.New("permission denied. not an editor")
	}
	srg.db.Create(&assignment)
	return assignment.ID, nil
}
func (srg SubmissionRepositoryGorm) UpdateAssignment(assignment Assignment, userID int32) error {
	res, err := srg.mampfLectureServiceClient.GetIsEditor(context.TODO(), &pb.IsEditorRequest{User: userID, Lecture: int32(assignment.LectureID)})
	if err != nil || !res.IsEditor {
		return errors.New("permission denied. not an editor")
	}
	srg.db.Save(assignment)
	return nil
}
func (srg SubmissionRepositoryGorm) DeleteAssignment(assignmentID uuid.UUID, userID int32) error {
	assignment, err := srg.FindAssignmentByID(assignmentID, userID)
	res, err := srg.mampfLectureServiceClient.GetIsEditor(context.TODO(), &pb.IsEditorRequest{User: userID, Lecture: int32(assignment.LectureID)})
	if err != nil || !res.IsEditor {
		return errors.New("permission denied. not an editor")
	}
	srg.db.Where("id = ?", assignmentID).Delete(&Assignment{})
	return nil
}
func (srg SubmissionRepositoryGorm) FindAssignmentByLectureID(lectureID int, userID int32) ([]Assignment, error) {
	res, err := srg.mampfLectureServiceClient.GetIsParticipantInLecture(context.TODO(), &pb.IsParticipantRequest{User: userID, Lecture: int32(lectureID)})
	resEditor, err := srg.mampfLectureServiceClient.GetIsEditor(context.TODO(), &pb.IsEditorRequest{User: userID, Lecture: int32(lectureID)})
	if err != nil || (!res.IsParticipant && !resEditor.IsEditor) {
		return make([]Assignment, 0), errors.New("permission denied. not a member of course")
	}
	var assignments []Assignment
	srg.db.Where("lecture_id = ?", lectureID).Find(&assignments)
	return assignments, nil
}
func (srg SubmissionRepositoryGorm) FindAssignmentByID(uuid uuid.UUID, userID int32) (Assignment, error) {
	var assignment Assignment
	srg.db.Where("id = ?", uuid).First(&assignment)
	//TODO(henrik): Discuss whether an check for permission must be done here
	return assignment, nil
}
func (srg SubmissionRepositoryGorm) CreateSubmission(submission Submission, userID int32) (uuid.UUID, error) {
	//TODO(henrik): Check also time here
	srg.db.Create(&submission)
	srg.db.Create(&Submitter{UserID: int(userID), SubmissionID: submission.ID})
	//We manage tutorials through a policy
	//srg.enforcer.AddPolicy(tutorialCasbin(submission.TutorialID), "tut")
	srg.enforcer.AddPolicy(fmt.Sprint(userID), submissionCasbin(submission.ID), "submit")
	return submission.ID, nil
}
func (srg SubmissionRepositoryGorm) UpdateSubmission(submission Submission, userID int32) error {

	if submission.ID == uuid.Nil {
		return errors.New("Without id")
	}

	ok, err := srg.enforcer.Enforce(fmt.Sprint(userID), submissionCasbin(submission.ID), "submit")

	tutOk, err := srg.enforcer.Enforce(fmt.Sprint(userID), tutorialCasbin(submission.TutorialID), "tut")
	if (!ok && !tutOk) || err != nil {
		return fmt.Errorf("%d has not the permission to edit", userID)
	}
	srg.db.Save(&submission)
	return nil
}
func submissionCasbin(submissionID uuid.UUID) string {
	return fmt.Sprintf("sub-%s", submissionID.String())
}
func tutorialCasbin(tutorialID uuid.UUID) string {
	return fmt.Sprintf("tut-%s", tutorialID.String())
}
func (srg SubmissionRepositoryGorm) DeleteSubmission(submissionID uuid.UUID, userID int32) error {
	if submissionID == uuid.Nil {
		return errors.New("Without id")
	}
	ok, err := srg.enforcer.Enforce(fmt.Sprint(userID), submissionCasbin(submissionID), "submit")
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("%d has not the permission to delete", userID)
	}
	srg.db.Where("id = ?", submissionID).Delete(&Submission{})
	return nil
}
func (srg SubmissionRepositoryGorm) FindSubmissionByAssignmentIDAndTutorialID(assignmentID uuid.UUID, tutorialID uuid.UUID, userID int32) (Submission, error) {
	var submission Submission
	srg.db.Where("assignment_id = ? AND tutorial_id = ?", assignmentID, tutorialID).Find(&submission)
	return submission, nil
}
func (srg SubmissionRepositoryGorm) FindSubmissionByID(uuid uuid.UUID, userID int32) (Submission, error) {
	var sub []Submission
	srg.db.Preload("Assignment").Joins("JOIN assignments on submissions.assignment_id = assignments.id").Where("submissions.id = ?", uuid).Find(&sub)
	if len(sub) == 0 {
		return Submission{}, errors.New("not found")
	}
	return sub[0], nil
}
func (srg SubmissionRepositoryGorm) FindSubmissionBySubmitterIDAndAssignmentID(submitterID int, assignmentID uuid.UUID, userID int32) (Submission, error) {
	var submission Submission
	srg.db.Joins("JOIN submitters ON submitters.submission_id = submissions.id AND submitters.user_id = ?", submitterID).Where("assignment_id = ?", assignmentID).Find(&submission)
	return submission, nil
}

func (srg SubmissionRepositoryGorm) FindSubmissionsBySubmitterIDAndLectureID(submitterID int, lectureId int, userID int32) ([]Submission, error) {
	var submissions []Submission
	srg.db.Preload("Assignment").Joins("JOIN submitters ON submitters.submission_id = submissions.id AND submitters.user_id = ? JOIN assignments ON submissions.assignment_id = assignments.id", submitterID).Where("lecture_id = ?", lectureId).Find(&submissions)
	return submissions, nil
}

func (srg SubmissionRepositoryGorm) SaveInviteToSubmission(invite Invitation, userID int32) error {
	ok, err := srg.enforcer.Enforce(fmt.Sprint(userID), submissionCasbin(invite.SubmissionID), "submit")
	if !ok || err != nil {
		return errors.New("error occured during enforcement")
	}
	srg.db.Save(&invite)
	return nil
}
func (srg SubmissionRepositoryGorm) AcceptInvitation(invite Invitation, userID int32) error {
	//TODO(henrik): Check time?
	if invite.InvitedUserID != int(userID) {
		return errors.New("only the user can accept an invite")
	}
	srg.db.Create(&Submitter{UserID: int(userID), SubmissionID: invite.SubmissionID})
	//TODO(henrik) remove policy & casbin as they are only double information
	srg.enforcer.AddPolicy(fmt.Sprint(invite.InvitedUserID), submissionCasbin(invite.SubmissionID), "submit")

	srg.db.Delete(&invite)
	return nil
}

func (srg SubmissionRepositoryGorm) FindSubmissionsByLectureIDAndTutorialID(lectureID int, tutorialID uuid.UUID, userID int32) ([]Submission, error) {
	var submissions []Submission
	srg.db.Joins("JOIN assigments ON assignments.id = submissions.assignment_id AND submissions.tutorial_id = ?", tutorialID).Where("lecture_id = ?", lectureID).Find(&submissions)
	return submissions, nil
}

func (srg SubmissionRepositoryGorm) CreateTutor(tutor Tutor, userID int32) error {
	tutorial, err := srg.FindTutorialByID(tutor.TutorialID, userID)
	if err != nil {
		return err
	}
	res, err := srg.mampfLectureServiceClient.GetIsEditor(context.TODO(), &pb.IsEditorRequest{User: userID, Lecture: int32(tutorial.LectureID)})
	if err != nil || !res.IsEditor {
		return errors.New("permission denied. not an editor")
	}
	//srg.db.Create(&tutor)
	srg.enforcer.AddPolicy(fmt.Sprint(tutor.UserID), tutorialCasbin(tutor.TutorialID), "tut")
	return nil
}

func (srg SubmissionRepositoryGorm) DeleteTutor(tutor Tutor, userID int32) error {
	tutorial, err := srg.FindTutorialByID(tutor.TutorialID, userID)
	if err != nil {
		return err
	}
	res, err := srg.mampfLectureServiceClient.GetIsEditor(context.TODO(), &pb.IsEditorRequest{User: userID, Lecture: int32(tutorial.LectureID)})
	if err != nil || !res.IsEditor {
		return errors.New("permission denied. not an editor")
	}
	//srg.db.Create(&tutor)
	srg.enforcer.RemovePolicy(fmt.Sprint(tutor.UserID), tutorialCasbin(tutor.TutorialID), "tut")
	return nil
}

func (srg SubmissionRepositoryGorm) FindSubmissionsFilesBySubmissionID(submissionID uuid.UUID, userID int32) (map[string]SubmissionsFile, error) {
	submission, err := srg.FindSubmissionByID(submissionID, userID)
	ok, err := srg.enforcer.Enforce(fmt.Sprint(userID), submissionCasbin(submissionID), "submit")
	tutOk, err := srg.enforcer.Enforce(fmt.Sprint(userID), tutorialCasbin(submission.TutorialID), "tut")
	if err != nil {
		return make(map[string]SubmissionsFile), err
	}
	if !ok && !tutOk {
		return make(map[string]SubmissionsFile), fmt.Errorf("%d has not the permission to view files", userID)
	}
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

func (srg SubmissionRepositoryGorm) CountSubmissionsFilesBySubmissionID(submissionID uuid.UUID, userID int32) (int, error) {
	var res int64
	srg.db.Model(&SubmissionsFile{}).Where("submission_id = ?", submissionID).Count(&res)
	log.Println("Res:", res)
	return int(res), nil
}

//TODO(henrik): enforce max recursion limit
func (srg SubmissionRepositoryGorm) FindSubmissionsSubFilesBySubmissionID(parent uuid.UUID, userID int32) (map[string]SubmissionsFile, error) {
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
func (srg SubmissionRepositoryGorm) CreateSubmissionsFile(submissionUUID uuid.UUID, parent uuid.UUID, isDir bool, name string, user int32) (SubmissionsFile, error) {
	subFI := SubmissionsFile{}
	subFI.buffer = &fileBuffer{data: make([]byte, 0)}
	subFI.minioClient = srg.minioClient
	subFI.SubmissionID = submissionUUID
	subFI.LastEditedBy = int(user)
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

func (srg SubmissionRepositoryGorm) DeleteSubmissionsFile(subID uuid.UUID, userID int32) error {
	srg.db.Where("id = ?", subID).Delete(&SubmissionsFile{})
	return nil
}

//TraverseToFile Traverses To a specific File, does not open other files
func (srg SubmissionRepositoryGorm) TraverseToFile(root SubmissionsFile, path []string, userID int32) (SubmissionsFile, bool, error) {
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

func (srg SubmissionRepositoryGorm) FindInvitation(invitationID uuid.UUID) (Invitation, error) {
	var inv Invitation
	srg.db.First(&inv, "id = ?", invitationID)
	if inv.ID != invitationID {
		return inv, errors.New("did not find the correct invitation")
	}
	return inv, nil
}

func (srg SubmissionRepositoryGorm) DeleteInvitation(invitation Invitation, userID int32) error {
	if invitation.InvitedUserID != int(userID) && invitation.InvitingUserID != int(userID) {
		return fmt.Errorf("%d is not allowed to delete/withdraw from inv %s", userID, invitation.ID)
	}
	if invitation.ID == uuid.Nil {
		return errors.New("nil id is not allowed")
	}
	srg.db.Where("id = ? ", invitation.ID).Delete(&Invitation{})
	return nil
}
func (srg SubmissionRepositoryGorm) FindTutorsByTutorialID(tutorialID uuid.UUID) ([]Tutor, error) {
	tutors, err := srg.enforcer.GetImplicitUsersForPermission(tutorialCasbin(tutorialID), "tut")
	if err != nil {
		return make([]Tutor, 0), err
	}
	res := make([]Tutor, len(tutors))
	for i := range res {
		id, _ := strconv.Atoi(tutors[i])
		res[i] = Tutor{TutorialID: tutorialID, UserID: id}
	}
	return res, nil
}

func (srg SubmissionRepositoryGorm) FindInvitationBySubmissionsID(submissionID uuid.UUID, userID int32) ([]Invitation, error) {

	ok, err := srg.enforcer.Enforce(fmt.Sprint(userID), submissionCasbin(submissionID), "submit")
	if !ok || err != nil {
		return make([]Invitation, 0), errors.New("error occured during enforcement")
	}
	var invites []Invitation
	srg.db.Where("submissionID = ?", submissionID).Find(&invites)
	return invites, nil
}
func (srg SubmissionRepositoryGorm) SubmissionJoinByToken(token string, userID int32) (Submission, error) {
	var submission Submission
	srg.db.Where("token = ?", token).First(&submission)
	srg.db.Create(&Submitter{UserID: int(userID), SubmissionID: submission.ID})
	//TODO(henrik) remove policy & casbin as they are only double information
	srg.enforcer.AddPolicy(fmt.Sprint(userID), submissionCasbin(submission.ID), "submit")
	return submission, nil
}
func (srg SubmissionRepositoryGorm) FindInvitationByUserID(userID int32) ([]Invitation, error) {
	var invitations []Invitation
	srg.db.Where("inviteduserID = ?", userID).Find(&invitations)
	return invitations, nil
}
