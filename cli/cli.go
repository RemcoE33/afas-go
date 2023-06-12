package cli

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"

	"github.com/RemcoE33/go-afas-appconnector-cli/lib"
)

// Global variables that are used in this package
var app Apps
var ConfigPath string

// Functional entrypoint
func Cli(version, configPath string) {
	ConfigPath = configPath
	if err := app.LoadConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var test bool
	var testUsage = "Environment type test (resttest)"

	var accept bool
	var acceptUsage = "Environment type accept (restaccept)"

	var production bool
	var productionUsage = "Environment type production (rest)"

	var member string
	var memberUsage = "AFAS member number (example: 11111) (required)"

	var token string
	var tokenUsage = "Full AppConnector token (<token>XXXX<token>) (required)"

	var location string
	var locationUsage = "Full location path (defaults to current directory) (cli creates subdirectoy)"

	var verbose bool
	var verboseUsage = "Verbose output of the metadata request"

	var number int
	var numberUsage = "Number of the list to run (required)"

	// Read functions
	flag.Bool("list", true, "List the stored app connectors from config file")
	flag.Bool("version", true, "Print version number")
	flag.Bool("v", true, "Print version number"+" (shorthand)")

	// To add a new app to the config file
	add := flag.NewFlagSet("add", flag.ExitOnError)
	add.BoolVar(&test, "test", false, testUsage)
	add.BoolVar(&test, "et", false, testUsage+" (shorthand)")

	add.BoolVar(&accept, "accept", false, acceptUsage)
	add.BoolVar(&accept, "ea", false, acceptUsage+" (shorthand)")

	add.BoolVar(&production, "production", false, productionUsage)
	add.BoolVar(&production, "ep", false, productionUsage+" (shorthand)")

	add.StringVar(&member, "member", "", memberUsage)
	add.StringVar(&member, "m", "", memberUsage+" (shorthand)")

	add.StringVar(&token, "token", "", tokenUsage)
	add.StringVar(&token, "t", "", tokenUsage+" (shorthand)")

	add.StringVar(&location, "location", "", locationUsage)
	add.StringVar(&location, "l", "", locationUsage+" (shorthand)")

	add.BoolVar(&verbose, "verbose", false, verboseUsage)
	add.BoolVar(&verbose, "v", false, verboseUsage+" (shorthand)")

	// Run an app from the config file based on listed number
	run := flag.NewFlagSet("run", flag.ExitOnError)
	run.IntVar(&number, "number", 0, numberUsage)
	run.IntVar(&number, "n", 0, numberUsage+" (shorthand)")
	run.BoolVar(&verbose, "verbose", false, verboseUsage)
	run.BoolVar(&verbose, "v", false, verboseUsage+" (shorthand)")

	// Run manually without storing to the config file
	man := flag.NewFlagSet("man", flag.ExitOnError)
	man.BoolVar(&test, "test", false, testUsage)
	man.BoolVar(&test, "et", false, testUsage+" (shorthand)")

	man.BoolVar(&accept, "accept", false, acceptUsage)
	man.BoolVar(&accept, "ea", false, acceptUsage+" (shorthand)")

	man.BoolVar(&production, "production", false, productionUsage)
	man.BoolVar(&production, "ep", false, productionUsage+" (shorthand)")

	man.StringVar(&member, "member", "", memberUsage)
	man.StringVar(&member, "m", "", memberUsage+" (shorthand)")

	man.StringVar(&token, "token", "", tokenUsage)
	man.StringVar(&token, "t", "", tokenUsage+" (shorthand)")

	man.BoolVar(&verbose, "verbose", false, verboseUsage)
	man.BoolVar(&verbose, "v", false, verboseUsage+" (shorthand)")

	// Delete an app
	del := flag.NewFlagSet("del", flag.ExitOnError)
	del.IntVar(&number, "number", 0, numberUsage)
	del.IntVar(&number, "n", 0, numberUsage+" (shorthand)")

	if len(os.Args) < 2 {
		fmt.Println("expected list, add, run, man or del subcommand")
		os.Exit(1)
	}

	flag.Parse()

	// If verbose is true then set it in env so the httpdump function is ran for API output
	setVerbose(verbose)

	switch os.Args[1] {
	case "version":
		fmt.Printf("Version: %s", version)
	case "v":
		fmt.Printf("Version: %s", version)

	case "list":
		app.List()

	case "add":
		add.Parse(os.Args[2:])
		var env string
		envType(test, accept, production, &env)
		if err := app.Add(member, env, token, location); err != nil {
			fmt.Print(err)
		}

	case "run":
		run.Parse(os.Args[2:])
		path, err := app.Run(number)
		if err != nil {
			fmt.Print(err)
		}
		OpenFolder(path)

	case "man":
		man.Parse(os.Args[2:])
		var env string
		envType(test, accept, production, &env)

		path, err := app.Manually(member, env, token)
		if err != nil {
			fmt.Print(err)
		}
		OpenFolder(path)

	case "del":
		del.Parse(os.Args[2:])
		if err := app.Delete(number); err != nil {
			fmt.Print(err)
		}

	default:
		app.List()
	}

}

// User can open the targeted folder or open de code editor
func OpenFolder(path string) {
	fmt.Printf("\n%vStruct files are waiting for you in: \n%s\n\nYou want to open the folder (code/y/n)?  %v", lib.Blue, path, lib.Reset)

	var ync string
	fmt.Scan(&ync)

	match, err := regexp.MatchString("n|y|code", ync)
	if err != nil {
		fmt.Println("error matching your input")
	}

	if ync == "n" || !match {
		return
	}

	var c string
	s := runtime.GOOS

	if ync == "y" {
		switch s {
		case "windowns":
			c = "start"
		case "darwin":
			c = "open"
		default:
			c = "xdg-open"
		}
	} else {
		c = ync
	}

	cmd := exec.Command(c, path)
	cmd.Run()

	unSetVerbose()
}
