package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/PureParadise56b/fileupload/files"
	"github.com/PureParadise56b/fileupload/handlers"
	"github.com/PureParadise56b/fileupload/middlewares"
)

func main() {

	addr := flag.String("addr", "3000", "server listen address.")
	baseDir := flag.String("dir", "./data", "file storage base directory path.")
	maxSize := flag.Int64("filelimit", 5, "maximum file size that can be uploaded by server in mb.")
	flag.Parse()

	mux := http.NewServeMux()

	store := files.NewLocalStorage(*baseDir, 1024*1024*(*maxSize))

	fh := handlers.NewFileHandler(store)
	fileserver := http.FileServer(http.Dir(*baseDir))

	// REST Apis
	mux.Handle("POST /files/{id}/{filename}", middlewares.ValidatePath(fh.UploadREST))
	mux.Handle("GET /files/{id}/{filename}", middlewares.ValidatePath(http.StripPrefix("/files/", fileserver).ServeHTTP))

	// Multipart form routes
	mux.HandleFunc("GET /upload", fh.ShowFormPage)
	mux.HandleFunc("POST /files/upload", middlewares.ValidateMultipartForm(fh.UplaodMultipart))

	server := http.Server{
		Addr:    ":" + *addr,
		Handler: mux,
	}

	log.Println("Server started at port:", *addr)
	log.Fatal(server.ListenAndServe())
}
