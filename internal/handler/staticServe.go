package handler

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
)

type servableFileSystem struct {
	fs http.FileSystem
}

func (sfs servableFileSystem) Open(path string) (http.File, error) {
	file, err := sfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	fileStat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if fileStat.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := sfs.fs.Open(index); err != nil {
			closeErr := file.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return file, nil
}

func StaticServe(r chi.Router, path string, root string) {
	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		servableFiles := servableFileSystem{http.Dir(root)}
		fs := http.StripPrefix(pathPrefix, http.FileServer(servableFiles))
		fs.ServeHTTP(w, r)
	})
}
