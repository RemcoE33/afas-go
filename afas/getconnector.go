package afas

import (
	"bytes"
	"fmt"
	"go/format"
	"html/template"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/RemcoE33/go-afas-appconnector-cli/models"
)

// The entrypoint function to get all the GetConnectoren and make the template.
func OrgistrateGetConnector(get []models.GetConnectors, s models.StoredApps) error {

	l := len(get)
	errs := make(chan error, l)
	var wg sync.WaitGroup
	wg.Add(l)

	// For each GetConnector in the AppConnector get the meta data.
	for _, conn := range get {
		go func(conn models.GetConnectors, errs chan<- error, wg *sync.WaitGroup) {

			endpoint := "metainfo/get/" + conn.ID
			var g models.GetConnector
			var td models.TemplateGetData

			if err := AfasRequest(&g, endpoint, s.Member, s.EnvironmentType, s.Token); err != nil {
				errs <- fmt.Errorf("MetaInfo.Get from %v. %v", conn.ID, err)
				wg.Done()
			}

			//Filling the template data struct
			td.Connector = g

			//Check if one of the fields is a date property so we can  inport time.Time
			if hasGetDateField(g.Fields) {
				td.ImportTime = true
			} else {
				td.ImportTime = false
			}

			if err := createGetConnectorTemplate(td, s); err != nil {
				errs <- fmt.Errorf("createGetConnectorTemplate: %v", err)
				wg.Done()
			}

			fmt.Printf("GetConnector:\t\t%s | %s | %s \n", s.Member, s.EnvironmentType, Title(g.Name))
			errs <- nil
			wg.Done()
		}(conn, errs, &wg)
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}

// For each GetConnector execute the template function
func createGetConnectorTemplate(d models.TemplateGetData, s models.StoredApps) error {

	tplPath := path.Join("tpl", "GetConnector.go.tmpl")

	tpl, err := template.New("GetConnector.go.tmpl").Funcs(template.FuncMap{
		//CamelCase the struct name and properties
		"Title": func(s string) string {
			return Title(s)
		},
		// Returns the go type from models.mapping
		"Type": func(t string) string {
			return GetType(t)
		},
	}).ParseFS(Folder, tplPath)

	if err != nil {
		return fmt.Errorf("template.New: %v", err)
	}

	// Init writer for the template.execute so that the code can be formatted before writing to the file
	var buf bytes.Buffer

	if err := tpl.Execute(&buf, d); err != nil {
		return fmt.Errorf("tpl.Excecute: %v", err)
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Println("error formatting go code")
	}

	//Naming is: afas.get.test.SalesLines.go
	outPath := filepath.Join(s.Location, "afas.get"+"."+s.EnvironmentType+"."+Title(d.Connector.Name)+".go")
	err = os.WriteFile(outPath, formatted, 0644)

	if err != nil {
		return fmt.Errorf("os.Writefile: %v", err)
	}

	return nil
}
