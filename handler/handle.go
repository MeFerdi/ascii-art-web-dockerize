package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"ascii-art-web-dockerize/ascii"
)

var templates *template.Template

// init initializes the template by parsing the index.html file.

func init() {
	var err error
	templates, err = template.ParseFiles(filepath.Join("templates", "index.html"))
	if err != nil {
		fmt.Println("Unable to parseFile: index.html missing")
		os.Exit(0)
	}
}

// HomeHandler handles requests to the home page.
// It only allows GET requests and serves the index.html template.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is GET
	if r.Method != http.MethodGet {
		http.ServeFile(w, r, "templates/405.html")
		return
	}

	err := templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error executing template index.html: %v", err)
	}
}

// AsciiArtHandler handles requests for generating ASCII art.
// It only allows POST requests and processes the input data to generate art.
func AsciiArtHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST
	if r.Method != http.MethodPost {
		http.ServeFile(w, r, "templates/405.html")
		return
	}

	str := r.FormValue("textData")
	bannerStyle := r.FormValue("banner")

	if len(str) == 0 {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}

	art, err := ascii.PrintAscii(str, bannerStyle)

	for _, cha := range str {
		if (cha < 32 || cha > 126) && err != nil {
			http.ServeFile(w, r, "templates/400.html")
			return
		}
	}
	if err != nil {
		http.ServeFile(w, r, "templates/500.html")
		return
	}

	// Prepare the data to pass to the template
	data := struct {
		Art string // Field to hold the generated ASCII art
	}{
		Art: art,
	}

	renderTemplate(w, "index", data)
}

// renderTemplate renders the specified template with the provided data.
func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	// Check if templates are initialized
	if templates == nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Execute the specified template with the provided data
	err := templates.ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error executing template %s: %v", tmpl, err)
	}
}
