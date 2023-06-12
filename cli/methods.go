package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/RemcoE33/go-afas-appconnector-cli/afas"
	"github.com/RemcoE33/go-afas-appconnector-cli/lib"
	"github.com/RemcoE33/go-afas-appconnector-cli/models"
)

type Apps []models.StoredApps

// Add app to apps and then write to the config file
func (a *Apps) Add(member, envtype, token, location string) error {

	//Do input checks
	checkCoreInput(&member, token)
	checkAndSetLocationString(&location)

	//Format for in the API header
	base64 := Base64(token)

	//Getting the app name and makes sure the inputed data is correct, otherwise a 401 / 403 / 500
	var metainfo models.AppConnector
	if err := afas.AfasRequest(&metainfo, "metainfo", member, envtype, base64); err != nil {
		return errors.New("error connecting to AFAS metadata endpoint: " + err.Error())
	}

	appName := metainfo.Info.AppName
	app := models.StoredApps{
		Member:          member,
		EnvironmentType: envtype,
		Desciption:      appName,
		Token:           lib.EncryptAES(base64),
		Location:        location,
	}

	*a = append(*a, app)

	if err := a.StoreConfig(); err != nil {
		return err
	}

	fmt.Println()
	fmt.Println(lib.Blue, "Succefully added: ", appName, lib.Reset)
	return nil
}

// Deletes an app from config file
func (a *Apps) Delete(n int) error {
	// Is number in range
	checkNumber(&n)

	ls := *a
	index := n - 1

	if index < 0 {
		return fmt.Errorf("number %d to delete is not found", n)
	}

	*a = append(ls[:index], ls[index+1:]...)

	if err := a.StoreConfig(); err != nil {
		return err
	}

	fmt.Println()
	fmt.Println(lib.Blue, "Succefully deleted:", n, ls[index].Desciption, lib.Reset)

	return nil
}

// Load config file into memory from UserConfigDir
func (a *Apps) LoadConfig() error {
	file, err := os.ReadFile(ConfigPath)
	if err != nil {
		return errors.New("app.LoadConfig() " + err.Error())
	}

	err = json.Unmarshal(file, a)
	if err != nil {
		return errors.New("app.LoadConfig(): invalid config file")
	}

	return nil
}

// Store the config in UserConfigDir
func (a *Apps) StoreConfig() error {
	data, err := json.MarshalIndent(a, "", "\t")
	if err != nil {
		return errors.New("app.StoreConfig(): error writing config file")
	}
	return os.WriteFile(ConfigPath, data, 6110)
}

// Run one of the saved apps
func (a *Apps) Run(n int) (string, error) {
	ls := *a
	index := n - 1
	app := ls[index]
	app.Token = lib.DecryptAES(app.Token)

	return app.Location, afas.MakeMetaRequest(app, true)
}

// Direcly run the program without adding to the config file.
func (a *Apps) Manually(member, env, token string) (string, error) {
	checkCoreInput(&member, token)
	base64 := Base64(token)
	var loc string
	checkAndSetLocationString(&loc)

	app := models.StoredApps{
		Member:          member,
		EnvironmentType: env,
		Token:           base64,
	}

	return loc, afas.MakeMetaRequest(app, true)
}

// Table view of all the apps stored in config
func (a *Apps) List() {
	ls := *a
	w := new(tabwriter.Writer)

	w.Init(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug|tabwriter.AlignRight)
	fmt.Fprintln(w, lib.Blue, "#", "\t", "Member", "\t", "Environment type", "\t", "Description", "\t", "Location", lib.Reset)
	fmt.Fprintln(w, lib.White, "", "\t", "", "\t", "", "\t", "", "\t", "", lib.Reset)

	for i, v := range ls {
		fmt.Fprintln(w, lib.White, i+1, "\t", v.Member, "\t", v.EnvironmentType, "\t", v.Desciption, "\t", v.Location, lib.Reset)
	}

	w.Flush()
	fmt.Println()
}
