package main

import (
	"context"
	"log"
	"net/http"
	"time"

	pb "github.com/henrixapp/mampf-rpc/grpc"
	"github.com/henrixapp/webdav-submission/server/admin"
	"github.com/henrixapp/webdav-submission/server/auth"
	"github.com/henrixapp/webdav-submission/server/fs"
	"golang.org/x/net/webdav"
	"google.golang.org/grpc"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var webd webdav.Handler

func Log(r *http.Request, err error) {
	//	log.Println(err, r)
	//	log.Println(r.BasicAuth())
}
func BasicAuth(handler http.HandlerFunc, mampfAuthServiceClient pb.MaMpfAuthServiceClient, realm string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user, pass, ok := r.BasicAuth()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		res, err := mampfAuthServiceClient.Login(ctx, &pb.LoginInformation{Email: user, Password: pass})
		if err != nil {
			log.Println(err)
		}
		if !ok || !res.Success {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			w.WriteHeader(401)
			w.Write([]byte("Unauthorised.\n"))
			return
		}
		req := r.WithContext(context.WithValue(context.Background(), "userID", res.User.Id))
		handler(w, req)
	}
}
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
	authService := pb.NewMaMpfAuthServiceClient(conn)
	db := initializeDB()
	webd = webdav.Handler{Logger: Log, FileSystem: fs.NewSharedWebDavFS(fs.MinioParams{Endpoint: "127.0.0.1:9000", AccessKeyID: "apfel", SecretAccessKey: "kuchensahne"}, auth.MampfParams{}, conn, admin.NewSubmissionRepositoryGorm(db)), LockSystem: webdav.NewMemLS()}

	log.Panicln(http.ListenAndServe(":3002", BasicAuth(webd.ServeHTTP, authService, "MaMpf")))
}
