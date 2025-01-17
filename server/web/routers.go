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

// Route is the information for every URI.
type Route struct {
	// Name is the name of this Route.
	Name string
	// Method is the string for the HTTP method. ex) GET, POST etc..
	Method string
	// Pattern is the pattern of the URI.
	Pattern string
	// HandlerFunc is the handler function of this route.
	HandlerFunc gin.HandlerFunc
}

// Routes is the list of the generated Route.
type Routes []Route

// NewRouter returns a new router.
func NewRouter(submissionsRep admin.SubmissionRepository) *gin.Engine {
	router := gin.Default()
	api := SubmissionsWebAPIProvider{submissionRepository: submissionsRep}
	var routes = Routes{
		{
			"Index",
			http.MethodGet,
			"/v1/submissions/",
			Index,
		},

		{
			"AssignmentsAssignmentIDDelete",
			http.MethodDelete,
			"/v1/submissions/assignments/:assignmentID",
			api.AssignmentsAssignmentIDDelete,
		},

		{
			"AssignmentsAssignmentIDGet",
			http.MethodGet,
			"/v1/submissions/assignments/:assignmentID",
			api.AssignmentsAssignmentIDGet,
		},

		{
			"AssignmentsAssignmentIDPut",
			http.MethodPut,
			"/v1/submissions/assignments/:assignmentID",
			api.AssignmentsAssignmentIDPut,
		},

		{
			"InvitationsGet",
			http.MethodGet,
			"/v1/submissions/invitations",
			api.InvitationsGet,
		},

		{
			"InvitationsInvitationIDModePost",
			http.MethodPost,
			"/v1/submissions/invitations/:invitationID/:mode",
			api.InvitationsInvitationIDModePost,
		},

		{
			"LectureLectureIDAssignmentsGet",
			http.MethodGet,
			"/v1/submissions/lecture/:lectureID/assignments",
			api.LectureLectureIDAssignmentsGet,
		},

		{
			"LectureLectureIDAssignmentsPost",
			http.MethodPost,
			"/v1/submissions/lecture/:lectureID/assignments",
			api.LectureLectureIDAssignmentsPost,
		},

		{
			"LectureLectureIDSubmissionsGet",
			http.MethodGet,
			"/v1/submissions/lecture/:lectureID/submissions",
			api.LectureLectureIDSubmissionsGet,
		},

		{
			"LectureLectureIDSubmissionsPost",
			http.MethodPost,
			"/v1/submissions/lecture/:lectureID/submissions",
			api.LectureLectureIDSubmissionsPost,
		},

		{
			"LectureLectureIDTutorialsGet",
			http.MethodGet,
			"/v1/submissions/lecture/:lectureID/tutorials",
			api.LectureLectureIDTutorialsGet,
		},

		{
			"LectureLectureIDTutorialsPost",
			http.MethodPost,
			"/v1/submissions/lecture/:lectureID/tutorials",
			api.LectureLectureIDTutorialsPost,
		},

		{
			"LecturesLectureIDInvitationsGet",
			http.MethodGet,
			"/v1/submissions/lectures/:lectureID/invitations",
			api.LecturesLectureIDInvitationsGet,
		},

		{
			"SubmissionsSubmissionIDDelete",
			http.MethodDelete,
			"/v1/submissions/submissions/:submissionID",
			api.SubmissionsSubmissionIDDelete,
		},

		{
			"SubmissionsSubmissionIDGet",
			http.MethodGet,
			"/v1/submissions/submissions/:submissionID",
			api.SubmissionsSubmissionIDGet,
		},

		{
			"SubmissionsSubmissionIDInvitationsGet",
			http.MethodGet,
			"/v1/submissions/submissions/:submissionID/invitations",
			api.SubmissionsSubmissionIDInvitationsGet,
		},

		{
			"SubmissionsSubmissionIDInvitationsPost",
			http.MethodPost,
			"/v1/submissions/submissions/:submissionID/invitations",
			api.SubmissionsSubmissionIDInvitationsPost,
		},

		{
			"SubmissionsSubmissionIDPut",
			http.MethodPut,
			"/v1/submissions/submissions/:submissionID",
			api.SubmissionsSubmissionIDPut,
		},

		{
			"SubmissionsTokenJoinPost",
			http.MethodPost,
			"/v1/submissions/token/:token/join",
			api.SubmissionsTokenJoinPost,
		},

		{
			"TutorialsTutorialIDTutorsGet",
			http.MethodGet,
			"/v1/submissions/tutorials/:tutorialID/tutors",
			api.TutorialsTutorialIDTutorsGet,
		},

		{
			"TutorialsTutorialIDTutorsUserIDDelete",
			http.MethodDelete,
			"/v1/submissions/tutorials/:tutorialID/tutors/:userID",
			api.TutorialsTutorialIDTutorsUserIDDelete,
		},

		{
			"TutorialsTutorialIDTutorsUserIDGet",
			http.MethodGet,
			"/v1/submissions/tutorials/:tutorialID/tutors/:userID",
			api.TutorialsTutorialIDTutorsUserIDGet,
		},

		{
			"TutorialsTutorialIDTutorsUserIDPost",
			http.MethodPost,
			"/v1/submissions/tutorials/:tutorialID/tutors/:userID",
			api.TutorialsTutorialIDTutorsUserIDPost,
		},
	}

	for _, route := range routes {
		switch route.Method {
		case http.MethodGet:
			router.GET(route.Pattern, route.HandlerFunc)
		case http.MethodPost:
			router.POST(route.Pattern, route.HandlerFunc)
		case http.MethodPut:
			router.PUT(route.Pattern, route.HandlerFunc)
		case http.MethodDelete:
			router.DELETE(route.Pattern, route.HandlerFunc)
		}
	}

	return router
}

// Index is the index handler.
func Index(c *gin.Context) {
	c.String(http.StatusOK, "Hello World!")
}
