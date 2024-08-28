package main

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log/slog"
	"math/rand"
	"net/http"
	"strings"
	"path/filepath"
)

//go:embed templates static
var embedFS embed.FS

type Server struct {
	mediaLib      *MediaLibrary
	tmpl          *template.Template
	staticVersion string
}

func httpError(r *http.Request, w http.ResponseWriter, err error, code int) {
	http.Error(w, err.Error(), code)
	slog.Error("failed request",
		err,
		slog.String("url", r.URL.String()),
		slog.Int("code", code),
	)
}

// ValidatePath provides a basic protection from the path traversal vulnerability.
func ValidatePath(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "./") || strings.Contains(r.URL.Path, ".\\") {
			httpError(r, w, errors.New("invalid path"), http.StatusBadRequest)
			return
		}
		h(w, r)
	}
}

// NormalizePath normalizes the request URL by removing the delimeter suffix.
func NormalizePath(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimRight(r.URL.Path, Delimiter)
		h(w, r)
	}
}

// DisableFileListing disables file listing under directories. It can be used with the built-in http.FileServer.
func DisableFileListing(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		h.ServeHTTP(w, r)
	})
}

type TemplateData struct {
	StaticVersion string
	*MediaListing
}

func (s *Server) ListingHandler(w http.ResponseWriter, r *http.Request) {
	listing, err := s.mediaLib.List(r.URL.Path)
	if err != nil {
		httpError(r, w, err, http.StatusInternalServerError)
		return
	}
	tmplData := TemplateData{
		StaticVersion: s.staticVersion,
		MediaListing:  listing,
	}
	if err := s.tmpl.ExecuteTemplate(w, "listing.gohtml", tmplData); err != nil {
		httpError(r, w, err, http.StatusInternalServerError)
		return
	}
}

func (s *Server) StreamHandler(w http.ResponseWriter, r *http.Request) {
	url, err := s.mediaLib.ContentURL(r.URL.Path)
	if err != nil {
		httpError(r, w, err, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, url, http.StatusFound)
}

func (s *Server) AudioHandler(w http.ResponseWriter, r *http.Request) {
	// Open Audio File
	audioFile, err := s.mediaLib.OpenFile(r.URL.Path)
	if err != nil {
		http.Error(w, "Unable to open audio file", http.StatusInternalServerError)
		return
	}


	// Get File Info
	fileInfo, err := audioFile.Stat()
	if err != nil {
		http.Error(w, "Unable to get file info", http.StatusInternalServerError)
		return
	}
	// Set Content-Type
	w.Header().Set("Content-Type", "audio/mpeg")

	// Set Content-Length
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	// Set Content-Disposition (Optional, using for downloading)
	w.Header().Set("Content-Disposition", "inline; filename="+filepath.Base(audioFile.Name()))

	// Copy audio stream to response
	_, err = io.Copy(w, audioFile)
	defer audioFile.Close()
	if err != nil {
		http.Error(w, "Error streaming audio", http.StatusInternalServerError)
		return
	}
}


// Don't include sprig just for one function.
var templateFunctions = map[string]any{
	"defaultString": func(s string, def string) string {
		if s == "" {
			return def
		}
		return s
	},
}

// StartServer starts HTTP server.
func StartServer(mediaLib *MediaLibrary, addr string) error {
	tmpl, err := template.New("").Funcs(templateFunctions).ParseFS(embedFS, "templates/*.gohtml")
	if err != nil {
		return err
	}

	mux := http.NewServeMux()

	mux.Handle("/", http.RedirectHandler("/library/", http.StatusMovedPermanently))

	staticVersion := fmt.Sprintf("%x", rand.Uint64())
	staticFS, err := fs.Sub(embedFS, "static")
	if err != nil {
		return err
	}
	staticPath := fmt.Sprintf("/static/%s/", staticVersion)
	mux.Handle(staticPath, DisableFileListing(http.StripPrefix(staticPath, http.FileServer(http.FS(staticFS)))))

	s := Server{
		mediaLib:      mediaLib,
		tmpl:          tmpl,
		staticVersion: staticVersion,
	}
	mux.Handle("/library/", http.StripPrefix("/library/", ValidatePath(NormalizePath(s.ListingHandler))))
	mux.Handle("/stream/", http.StripPrefix("/stream/", ValidatePath(NormalizePath(s.StreamHandler))))
	mux.Handle("/audio/", http.StripPrefix("/audio/", ValidatePath(NormalizePath(s.AudioHandler))))

	return http.ListenAndServe(addr, mux)
}
