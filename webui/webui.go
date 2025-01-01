package webui

import (
	"cloud-commis/config"
	"cloud-commis/logger"
	"cloud-commis/storage"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"text/template"
)

//go:embed statics
var staticFiles embed.FS

//go:embed templates
var templateFiles embed.FS

type RootPageData struct {
	UserEmail string
}

func getVM(writer http.ResponseWriter, request *http.Request) {

	writer.WriteHeader(http.StatusOK)

	storage.Configure()

	jsonData, err := storage.Data.Read()
	if err != nil {
		logger.Log.Error("Read storage failure")
	}

	var templatesFS = fs.FS(templateFiles)

	tmpl := template.Must(template.ParseFS(templatesFS, "templates/vmlist.tmpl"))
	tmplErr := tmpl.Execute(writer, jsonData.Data)
	if tmplErr != nil {
		logger.Log.Error("Cannot render template vmlist.tmpl" + tmplErr.Error())
	}
}

func getConfig(writer http.ResponseWriter, request *http.Request) {
	var templatesFS = fs.FS(templateFiles)
	tmpl := template.Must(template.ParseFS(templatesFS, "templates/config.tmpl"))
	err := tmpl.Execute(writer, config.ParsedData.All())
	if err != nil {
		logger.Log.Error("Cannot render template vmlist.tmpl" + err.Error())
	}
}

func Start(router *http.ServeMux) {

	var staticFS = fs.FS(staticFiles)
	htmlContent, err := fs.Sub(staticFS, "statics")
	if err != nil {
		log.Fatal(err)
	}

	router.Handle("GET /", http.FileServer(http.FS(htmlContent)))

	router.HandleFunc("GET /config", getConfig)
	router.HandleFunc("GET /vmlist", getVM)

}
