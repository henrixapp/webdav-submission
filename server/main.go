package main

import (
	"log"
	"net/http"

	"github.com/henrixapp/webdav-submission/server/auth"
	"github.com/henrixapp/webdav-submission/server/fs"
	"golang.org/x/net/webdav"
)

func Log(r *http.Request, err error) {
	log.Println(err, r)
	log.Println(r.BasicAuth())
}
func main() {
	webd := webdav.Handler{Logger: Log, FileSystem: fs.NewSharedWebDavFS(fs.MinioParams{Endpoint: "localhost", AccessKeyID: "apfel", SecretAccessKey: "kuchensahne"}, auth.MampfParams{}), LockSystem: webdav.NewMemLS()}
	log.Panicln(http.ListenAndServe(":3000", &webd))
}
