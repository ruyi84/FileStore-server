package main

import (
	"fmt"
	. "github.com/filestore-server/handler"
	"net/http"
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.HandleFunc("/file/upload", UploadHandler)
	http.HandleFunc("/file/upload/suc", UploadSucHandler)
	http.HandleFunc("/file/meta", GetFileMetaHandler)
	http.HandleFunc("/file/query", FileQueryHandler)
	http.HandleFunc("/file/download", DownloadHandler)
	http.HandleFunc("/file/update", FileUpdateMetaHandler)
	http.HandleFunc("/file/delete", FileDelHandler)
	http.HandleFunc("/file/fastupload", HTTPInterceptor(TryFastUploadHandler))

	http.HandleFunc("/user/signup", SignupHander)
	http.HandleFunc("/user/signin", SigninHandler)
	http.HandleFunc("/user/info", HTTPInterceptor(UserInfoHandler))

	http.HandleFunc("/file/mpupload/init", HTTPInterceptor(InitialMultipartUploadHandle))
	http.HandleFunc("/file/mpupload/uppart", HTTPInterceptor(UploadPartHandler))
	http.HandleFunc("/file/mpupload/complete", HTTPInterceptor(CompleteUploadHandler))
	//http.HandlerFunc("/file/mpupload/cancle",)
	//http.HandlerFunc("/file/mpupload/status",)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Failed to start server,err%s\n", err.Error())
	}
}
