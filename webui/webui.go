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

	jsonData, err := storage.Data.Read()
	if err != nil {
		logger.Log.Error("Read storage failure")
	}

	var templatesFS = fs.FS(templateFiles)

	tmpl := template.Must(template.ParseFS(templatesFS, "templates/vmlist.tmpl"))
	tmplErr := tmpl.Execute(writer, jsonData.AwsAccounts)
	if tmplErr != nil {
		logger.Log.Error("Cannot render template vmlist.tmpl" + tmplErr.Error())
	}
}

func getHome(writer http.ResponseWriter, request *http.Request) {

	type stats struct {
		TotalVMs      int
		SumVMByRegion map[string]int
	}

	writer.WriteHeader(http.StatusOK)

	jsonData, err := storage.Data.Read()
	if err != nil {
		logger.Log.Error("Read storage failure")
	}

	homeStats := stats{}
	homeStats.SumVMByRegion = make(map[string]int)
	for _, account := range jsonData.AwsAccounts {
		for _, region := range account.AwsRegions {

			homeStats.TotalVMs += len(region.VirtualMachines)
			if len(region.VirtualMachines) > 0 {
				homeStats.SumVMByRegion[region.RegionName] += len(region.VirtualMachines)
			}
		}
	}

	var templatesFS = fs.FS(templateFiles)

	tmpl := template.Must(template.ParseFS(templatesFS, "templates/home.tmpl"))
	tmplErr := tmpl.Execute(writer, homeStats)
	if tmplErr != nil {
		logger.Log.Error("Cannot render template home.tmpl" + tmplErr.Error())
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
	router.HandleFunc("GET /home", getHome)
	router.HandleFunc("GET /vmlist", getVM)

}
