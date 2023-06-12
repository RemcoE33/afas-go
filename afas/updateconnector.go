package afas

import (
	"bytes"
	"fmt"
	"go/format"
	"html/template"
	"os"
	"path/filepath"
	"sync"

	"github.com/RemcoE33/go-afas-appconnector-cli/models"
)

// The entrypoint function to get all the GetConnectoren and make the template.
func OrgistrateUpdateConnector(up []models.UpdateConnectors, s models.StoredApps) error {

	l := len(up)
	errs := make(chan error, l)
	var wg sync.WaitGroup
	wg.Add(l)

	// For each GetConnector in the AppConnector get the meta data.
	for _, conn := range up {
		go func(conn models.UpdateConnectors, wg *sync.WaitGroup) {
			var data models.TemplateUpdateData
			endpoint := "metainfo/update/" + conn.ID

			if err := AfasRequest(&data.Connector, endpoint, s.Member, s.EnvironmentType, s.Token); err != nil {
				errs <- fmt.Errorf("MetaInfo.Get from %v: %v", conn.ID, err)
				wg.Done()
			}

			var objects []models.UpdateObjects
			hasDate := false

			//flatten all nested objects
			extractObjectFields(data.Connector.Objects, &objects)

			// Contains all the nested objects, we use map to make sure there are no nested structs
			flatObject := make(map[string]models.UpdateObjects)

			for _, obj := range objects {
				flatObject[obj.Name] = obj

				// Check if nested fields has a date field
				if !hasDate {
					r := hasUpdateDateField(obj.Fields)
					if r {
						hasDate = r
					}
				}
			}

			//Check if top level fields has a date field
			if !hasDate {
				r := hasUpdateDateField(data.Connector.Fields)
				if r {
					hasDate = r
				}
			}

			// used for the import statement if there is a date field
			data.ImportTime = hasDate
			data.FlatObjects = flatObject

			if err := createUpdateConnectorTemplate(data, s); err != nil {
				errs <- fmt.Errorf("createUpdateConnectorTemplate: %v", err)
				wg.Done()
			}

			fmt.Printf("UpdateConnector:\t%s | %s | %s \n", s.Member, s.EnvironmentType, data.Connector.Name)
			errs <- nil
			wg.Done()
		}(conn, &wg)
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

// Recursive function to flat all the nested objects. This for unnested structs in the template
func extractObjectFields(obj []models.UpdateObjects, s *[]models.UpdateObjects) {
	for _, o := range obj {
		*s = append(*s, o)

		if len(o.Objects) > 0 {
			extractObjectFields(o.Objects, s)
		}
	}
}

// Creat the updateconnector structs from template
func createUpdateConnectorTemplate(d models.TemplateUpdateData, s models.StoredApps) error {

	tplPath := filepath.Join("tpl", "UpdateConnector.go.tmpl")

	tpl, err := template.New("UpdateConnector.go.tmpl").Funcs(template.FuncMap{
		// To CamelCase the struct names and properties
		"Title": func(s string) string {
			return Title(s)
		},
		// Returns the mapped type from models.mapping
		"Type": func(t string) string {
			return GetType(t)
		},
	}).ParseFS(Folder, tplPath)

	if err != nil {
		return fmt.Errorf("template.New: %v", err)
	}

	// Init a writer for the template.execute so that the code can be formatted before writing to the file
	var buf bytes.Buffer

	if err := tpl.Execute(&buf, d); err != nil {
		return fmt.Errorf("tpl.Excecute: %v", err)
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Println("format.Source(): ", err)
	}

	err = os.MkdirAll(filepath.Join(s.Location, d.Connector.ID), os.ModePerm)
	if err != nil {
		return fmt.Errorf("os.MkdirAll.update: %v", err)
	}

	//Naming is: afas.update.test.FbSales.go
	outPath := filepath.Join(s.Location, d.Connector.ID, "afas.update."+s.EnvironmentType+"."+d.Connector.ID+".go")
	err = os.WriteFile(outPath, formatted, 0644)

	if err != nil {
		return fmt.Errorf("os.WriteFile.Update: %v", err)
	}

	return nil
}
