package middlewares

import (
	"context"
	"log"
	"net/http"
	"regexp"
)

type ContextKey string
type ContextValue struct {
	Id       string
	Filename string
}

func ValidatePath(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		filename := r.PathValue("filename")

		if !isValidImageId(id) || !isValidFileName(filename) {
			log.Printf("invalid id or filename for id: %v, filename:%v\n", id, filename)
			http.Error(w, "400 Bad Request", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func ValidateMultipartForm(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(5 * 1024 * 1024)
		if err != nil {
			log.Println(err)
			http.Error(w, "400 Bad Request", http.StatusBadRequest)
			return
		}

		id := r.FormValue("id")
		f, mhf, err := r.FormFile("file")
		if err != nil {
			log.Println(err)
			http.Error(w, "400 Bad Request", http.StatusBadRequest)
			return
		}

		filename := mhf.Filename

		if !isValidImageId(id) || !isValidFileName(filename) {
			log.Printf("invalid id or filename for id: %v, filename:%v\n", id, filename)
			http.Error(w, "400 Bad Request", http.StatusBadRequest)
			return
		}

		r.Body = f
		ctxKey := ContextKey("f_metadata")
		ctxValue := ContextValue{
			Id:       id,
			Filename: filename,
		}
		ctx := context.WithValue(context.Background(), ctxKey, ctxValue)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}
}

func isValidImageId(id string) bool {
	if id == "" {
		return false
	}

	regx, err := regexp.Compile("[0-9]+")
	if err != nil {
		log.Println(err)
		return false
	}
	return regx.MatchString(id)
}

func isValidFileName(name string) bool {
	if name == "" {
		return false
	}

	regx, err := regexp.Compile(`^\w+.(jpg|png|jpeg|mp4|mkv|mp3|avi)$`)
	if err != nil {
		log.Println(err)
		return false
	}

	return regx.MatchString(name)
}
