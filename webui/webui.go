package webui

import (
	"cloud-commis/config"
	"cloud-commis/logger"
	"cloud-commis/storage"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

//go:embed statics
var staticFiles embed.FS

//go:embed templates
var templateFiles embed.FS

type RootPageData struct {
	UserEmail string
}

func getVersion(writer http.ResponseWriter, request *http.Request) {

	writer.WriteHeader(http.StatusOK)

	writer.Write([]byte(config.Version))
}

func getVM(writer http.ResponseWriter, request *http.Request) {

	writer.WriteHeader(http.StatusOK)

	jsonData, err := storage.Data.Read()
	if err != nil {
		logger.Log.Error("Read storage failure")
	}

	var templatesFS = fs.FS(templateFiles)

	tmpl := template.Must(template.ParseFS(templatesFS, "templates/vmlist.tmpl"))
	tmplErr := tmpl.Execute(writer, jsonData)
	if tmplErr != nil {
		logger.Log.Error("Cannot render template vmlist.tmpl" + tmplErr.Error())
	}
}

func getVMDetails(writer http.ResponseWriter, request *http.Request) {

	type vmData struct {
		VmDetails  storage.VirtualMachine
		AmiDetails storage.AwsImage
	}

	awsAccount, _ := strconv.Atoi(request.PathValue("awsaccount"))
	region := request.PathValue("region")
	vmId := request.PathValue("vmid")

	writer.WriteHeader(http.StatusOK)

	storedData, err := storage.Data.Read()
	if err != nil {
		logger.Log.Error("Read storage failure")
	}

	var data vmData
	amiId := storedData.AwsAccounts[awsAccount].AwsRegions[region].VirtualMachines[vmId].ImageId
	data.VmDetails = storedData.AwsAccounts[awsAccount].AwsRegions[region].VirtualMachines[vmId]
	data.AmiDetails = storedData.AwsImages[amiId]

	var templatesFS = fs.FS(templateFiles)

	tmpl := template.Must(template.ParseFS(templatesFS, "templates/vmdetails.tmpl"))

	tmplErr := tmpl.Execute(writer, data)
	if tmplErr != nil {
		logger.Log.Error("Cannot render template vmdetails.tmpl" + tmplErr.Error())
	}
}

func getHome(writer http.ResponseWriter, request *http.Request) {

	type stats struct {
		TotalVMs       int
		SumVmByRegion  map[string]int
		SumVmByType    map[string]int
		SumVmByOS      map[string]int
		SumVmByAccount map[int]int
	}

	writer.WriteHeader(http.StatusOK)

	jsonData, err := storage.Data.Read()
	if err != nil {
		logger.Log.Error("Read storage failure")
	}

	homeStats := stats{}
	homeStats.SumVmByRegion = make(map[string]int)
	homeStats.SumVmByType = make(map[string]int)
	homeStats.SumVmByOS = make(map[string]int)
	homeStats.SumVmByAccount = make(map[int]int)

	for accountId, accountData := range jsonData.AwsAccounts {
		for region, regionData := range accountData.AwsRegions {

			homeStats.TotalVMs += len(regionData.VirtualMachines)
			if len(regionData.VirtualMachines) > 0 {
				homeStats.SumVmByRegion[region] += len(regionData.VirtualMachines)
			}
			for _, vmData := range regionData.VirtualMachines {
				homeStats.SumVmByType[vmData.InstanceType] += 1
				homeStats.SumVmByOS[vmData.BootImage] += 1
				homeStats.SumVmByAccount[accountId] += 1
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
	router.HandleFunc("GET /version", getVersion)
	router.HandleFunc("GET /vmdetails/{awsaccount}/{region}/{vmid}", getVMDetails)

}
