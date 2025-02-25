package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/karnerfly/fileupload/files"
	"github.com/karnerfly/fileupload/handlers"
	"github.com/karnerfly/fileupload/middlewares"
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
	mux.Handle("POST /files/{id}/{filename}", middlewares.ValidatePath(http.HandlerFunc(fh.UploadREST)))
	mux.Handle("GET /files/{id}/{filename}", middlewares.ValidatePath(http.StripPrefix("/files/", fileserver)))

	// Multipart form routes
	mux.Handle("GET /upload", http.HandlerFunc(fh.ShowFormPage))
	mux.Handle("POST /files/upload", middlewares.ValidateMultipartForm(http.HandlerFunc(fh.UplaodMultipart)))

	// test the gzip middleware
	mux.Handle("GET /test", middlewares.GzipMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile("./demo.txt")
		if err != nil {
			log.Println(err)
			return
		}

		w.Write(data)
	})))

	server := http.Server{
		Addr:    ":" + *addr,
		Handler: mux,
	}

	log.Println("Server started at port:", *addr)
	log.Fatal(server.ListenAndServe())
}
