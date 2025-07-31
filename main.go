package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/karnerfly/fileupload/files"
	"github.com/karnerfly/fileupload/handlers"
	"github.com/karnerfly/fileupload/middlewares"
)

func main() {
	addr := flag.String("addr", "3000", "server listen address.")
	baseDir := flag.String("dir", "uploads", "file storage base directory path.")
	maxSize := flag.Int64("filelimit", 5, "maximum file size that can be uploaded by server in mb.")
	flag.Parse()

	mux := http.NewServeMux()

	store := files.NewLocalStorage(*baseDir, (*maxSize)<<20)

	fh := handlers.NewFileHandler(store)
	fileserver := http.FileServer(http.Dir(*baseDir))

	// default handler to get files
	mux.Handle("GET /files/{id}/{filename}", middlewares.ValidatePath(http.StripPrefix("/files/", fileserver)))

	// REST Apis
	mux.Handle("POST /api/files/{id}/{filename}", middlewares.ValidatePath(http.HandlerFunc(fh.UploadREST)))

	// Multipart form routes
	mux.Handle("GET /upload", http.HandlerFunc(fh.ShowFormPage))
	mux.Handle("POST /files/upload", middlewares.ValidateMultipartForm(http.HandlerFunc(fh.UplaodMultipart)))

	server := http.Server{
		Addr:    ":" + *addr,
		Handler: mux,
	}

	log.Println("Server started at port:", *addr)
	log.Fatal(server.ListenAndServe())
}
