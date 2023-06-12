package afas

import (
	"embed"
	"fmt"
	"os"
	"path"

	"github.com/RemcoE33/go-afas-appconnector-cli/models"
)

var Folder embed.FS

func MakeMetaRequest(s models.StoredApps, generate bool) error {
	var app models.AppConnector
	if err := AfasRequest(&app, "metainfo", s.Member, s.EnvironmentType, s.Token); err != nil {
		return err
	}

	if !generate {
		return nil
	}

	path := path.Join(s.Location, "appconnector")
	os.Mkdir(path, 6110)

	if err := OrgistrateGetConnector(app.GetConnectors, s); err != nil {
		return fmt.Errorf("OrgistrateGetConnector: %v", err)
	}

	if err := OrgistrateUpdateConnector(app.UpdateConnectors, s); err != nil {
		return fmt.Errorf("OrgistrateUpdateConnector: %v", err)
	}

	return nil
}
