package main

import (
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type Page struct {
	Title string
	Body []byte
}

const PUBLIC_PATH = "../../public/"
const HTML_PATH = PUBLIC_PATH + "html/"

const RESOURCE_PATH = "../"
const TEXT_PATH = RESOURCE_PATH + "text/"

const EXPEND_STRING = ".txt"

const TOP_PATH_LENGTH = len("/top/")
const VIEW_PATH_LENGTH = len("/view/")
const EDIT_PATH_LENGTH = len("/edit/")
const SAVE_PATH_LENGTH = len("/save/")

var titleValidator = regexp.MustCompile("^[a-zA-Z0-9]+$")

var templates = make(map[string]*template.Template)

func init() {
	for _, fileName := range []string{"top", "view", "edit"} {
		template := template.Must(template.ParseFiles(HTML_PATH + fileName + ".html"))
		templates[fileName] = template
	}
}

// Model func
func (page *Page)save() error {
	fileName := page.Title + ".txt"
	// 0600はread + writeのアクセス権限設定(自分のみ)
	return ioutil.WriteFile(TEXT_PATH + fileName, page.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	fileName := title + ".txt"
	body, error := ioutil.ReadFile(TEXT_PATH + fileName)
	if error != nil {
		return nil, error
	}
	return &Page{Title: title, Body: body}, nil
}

// Helper or Common func or Base Class func
func renderTemplate(writer http.ResponseWriter, fileName string, page *Page) {
	error := templates[fileName].Execute(writer, page)
	if error != nil {
		http.Error(writer, error.Error(), http.StatusInternalServerError)
	}
}

// Controller
func makeHandler(function func(http.ResponseWriter, *http.Request, string), PathLength int) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		title := request.URL.Path[PathLength:]
		if !titleValidator.MatchString(title) {
			http.NotFound(writer, request)
			error := errors.New("Invalid Page Title: " + title)
			log.Print(error)
			return
		}
		function(writer, request, title)
	}
}

func topHandler(writer http.ResponseWriter, request *http.Request) {
	files, error := ioutil.ReadDir(TEXT_PATH)
	if error != nil {
		error := errors.New(".txt file not found. line of 67")
		log.Print(error)
		return
	}

	var paths []string
	var fileName []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), EXPEND_STRING) {
			fileName = strings.Split(string(file.Name()), EXPEND_STRING)
			paths = append(paths, fileName[0])
		}
	}

	if paths == nil {
		error := errors.New(".txt file not found. line of 82")
		log.Print(error)
	}

	template := template.Must(template.ParseFiles(HTML_PATH + "top.html"))
	error = template.Execute(writer, paths)
	if error != nil {
		http.Error(writer, error.Error(), http.StatusInternalServerError)
		return
	}
}

func viewHandler(writer http.ResponseWriter, request *http.Request, title string) {
	page, error := loadPage(title)
	if error != nil {
		http.Redirect(writer, request, "/edit/" + title, http.StatusFound)
		return
	}
	renderTemplate(writer, "view", page)
}

func editHandler(writer http.ResponseWriter, request *http.Request, title string) {
	page, error := loadPage(title)
	if error != nil {
		page = &Page{Title: title}
	}
	renderTemplate(writer, "edit", page)
}

func saveHandler(writer http.ResponseWriter, request *http.Request, title string) {
	body := request.FormValue("body")
	page := &Page{Title: title, Body: []byte(body)}
	error := page.save()
	if error != nil {
		http.Error(writer, error.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(writer, request, "/view/" + title, http.StatusFound)
}

// rooting and listen serve.
func main() {
	http.HandleFunc("/top/", topHandler)
	http.HandleFunc("/view/", makeHandler(viewHandler, VIEW_PATH_LENGTH))
	http.HandleFunc("/edit/", makeHandler(editHandler, EDIT_PATH_LENGTH))
	http.HandleFunc("/save/", makeHandler(saveHandler, SAVE_PATH_LENGTH))
	http.ListenAndServe(":8080", nil)
}
