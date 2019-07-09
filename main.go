package main

import (
	"fmt"
	. "github.com/filestore-server/handler"
	"net/http"
)

func main() {
	http.HandleFunc("/file/upload", UploadHandler)
	http.HandleFunc("/file/upload/suc", UploadSucHandler)
	http.HandleFunc("/file/meta", GetFileMetaHandler)
	http.HandleFunc("/file/download", DownloadHandler)
	http.HandleFunc("/file/update", FileUpdateMetaHandler)
	http.HandleFunc("/file/delete", FileDelHandler)

	http.HandleFunc("/user/signup", SignupHander)
	http.HandleFunc("/user/signin", SigninHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Failed to start server,err%s\n", err.Error())
	}
}
