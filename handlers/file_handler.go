package handlers

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"

	"github.com/Pureparadise56b/fileupload/files"
	"github.com/Pureparadise56b/fileupload/middlewares"
)

type File struct {
	store files.Storage
}

func NewFileHandler(s files.Storage) *File {
	return &File{
		store: s,
	}
}

// file handling in restful architecture
func (f *File) UploadREST(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	id := r.PathValue("id")
	filename := r.PathValue("filename")

	err := f.saveFile(id, filename, r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Println("file uploaded:", filename)
	w.Write([]byte("OK"))
}

func (f *File) UplaodMultipart(w http.ResponseWriter, r *http.Request) {
	metadata := r.Context().Value(middlewares.ContextKey("f_metadata")).(middlewares.ContextValue)

	err := f.saveFile(metadata.Id, metadata.Filename, r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Println("file uploaded:", metadata.Filename)
	w.Write([]byte("OK"))
}

func (f *File) ShowFormPage(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./pages/index.html")
	if err != nil {
		log.Println(err)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
	}
}

func (f *File) saveFile(id, filename string, r io.Reader) error {
	fp := filepath.Join(id, filename)
	return f.store.Save(fp, r)
}
