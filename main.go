package main

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/RemcoE33/go-afas-appconnector-cli/afas"
	"github.com/RemcoE33/go-afas-appconnector-cli/cli"
	"github.com/RemcoE33/go-afas-appconnector-cli/lib"
)

var Version = "v0.1.0"
var ConfigPath = ""
var Filename = "go-afas-struct.json"

//go:embed tpl/*
var folder embed.FS

// check if config file excists, otherwise create the config in comnfig dir or in specified path when building the program
func init() {
	if ConfigPath != "" {
		return
	}

	confDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Print("error getting config directory with os.UserConfigDir()")
		os.Exit(1)
	}

	// Creat the config folder
	confDirPath := path.Join(confDir, "GoAfas")
	os.Mkdir(confDirPath, 6110)

	filePath := path.Join(confDirPath, Filename)
	ConfigPath = filePath

	emptySlice := []string{}

	// create an empty array in the config.json file.
	if _, err := os.Stat(filePath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println(lib.Blue, "Welcom! initializing config file in: ", filePath, lib.Reset)

			data, err := json.MarshalIndent(emptySlice, "", "\t")
			if err != nil {
				fmt.Println("erorr initializing empty json file")
			}

			err = os.WriteFile(filePath, data, 6110)
			if err != nil {
				fmt.Print("error writing config file")
			}
		}
	}

}

func main() {
	afas.Folder = folder
	cli.Cli(Version, ConfigPath)
}
