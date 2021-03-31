/*
 * Submissions
 *
 * This API specifies the submissions service, as accessed by the web admin UI used by students, lecturers and tutors.
 *
 * API version: 1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/henrixapp/webdav-submission/server/admin"
)

type SubmissionsWebAPIProvider struct {
	submissionRepository admin.SubmissionRepository
}

// AssignmentsAssignmentIDDelete -
func (api *SubmissionsWebAPIProvider) AssignmentsAssignmentIDDelete(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// AssignmentsAssignmentIDGet -
func (api *SubmissionsWebAPIProvider) AssignmentsAssignmentIDGet(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// AssignmentsAssignmentIDPut -
func (api *SubmissionsWebAPIProvider) AssignmentsAssignmentIDPut(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// InvitationsGet -
func (api *SubmissionsWebAPIProvider) InvitationsGet(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// InvitationsInvitationIDModePost -
func (api *SubmissionsWebAPIProvider) InvitationsInvitationIDModePost(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// LectureLectureIDAssignmentsGet -
func (api *SubmissionsWebAPIProvider) LectureLectureIDAssignmentsGet(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// LectureLectureIDAssignmentsPost -
func (api *SubmissionsWebAPIProvider) LectureLectureIDAssignmentsPost(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// LectureLectureIDSubmissionsGet -
func (api *SubmissionsWebAPIProvider) LectureLectureIDSubmissionsGet(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// LectureLectureIDSubmissionsPost -
func (api *SubmissionsWebAPIProvider) LectureLectureIDSubmissionsPost(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// LectureLectureIDTutorialsGet -
func (api *SubmissionsWebAPIProvider) LectureLectureIDTutorialsGet(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// LectureLectureIDTutorialsPost -
func (api *SubmissionsWebAPIProvider) LectureLectureIDTutorialsPost(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// LecturesLectureIDInvitationsGet -
func (api *SubmissionsWebAPIProvider) LecturesLectureIDInvitationsGet(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// SubmissionsSubmissionIDDelete -
func (api *SubmissionsWebAPIProvider) SubmissionsSubmissionIDDelete(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// SubmissionsSubmissionIDGet -
func (api *SubmissionsWebAPIProvider) SubmissionsSubmissionIDGet(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// SubmissionsSubmissionIDInvitationsGet -
func (api *SubmissionsWebAPIProvider) SubmissionsSubmissionIDInvitationsGet(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// SubmissionsSubmissionIDInvitationsPost -
func (api *SubmissionsWebAPIProvider) SubmissionsSubmissionIDInvitationsPost(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// SubmissionsSubmissionIDPut -
func (api *SubmissionsWebAPIProvider) SubmissionsSubmissionIDPut(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// SubmissionsTokenJoinPost -
func (api *SubmissionsWebAPIProvider) SubmissionsTokenJoinPost(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// TutorialsTutorialIDTutorsGet -
func (api *SubmissionsWebAPIProvider) TutorialsTutorialIDTutorsGet(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// TutorialsTutorialIDTutorsUserIDDelete -
func (api *SubmissionsWebAPIProvider) TutorialsTutorialIDTutorsUserIDDelete(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// TutorialsTutorialIDTutorsUserIDGet -
func (api *SubmissionsWebAPIProvider) TutorialsTutorialIDTutorsUserIDGet(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// TutorialsTutorialIDTutorsUserIDPost -
func (api *SubmissionsWebAPIProvider) TutorialsTutorialIDTutorsUserIDPost(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}
