package cli

import (
	b64 "encoding/base64"
	"fmt"
	"os"
	"regexp"
)

// Check if app number is in range.
func checkNumber(n *int) {
	if *n == 0 || *n > len(app) {
		fmt.Println("invalid number or out of bound")
		os.Exit(1)
	}
}

// Makes sure the mandetory values are set
func checkCoreInput(member *string, token string) {
	if *member == "" || token == "" {
		fmt.Println("missing member or token")
		os.Exit(1)
	}

	// Checking if token is valid
	match, err := regexp.MatchString(`<token><version>\d</version><data>[A-Z0-9]{64}</data></token>`, token)
	if err != nil || !match {
		fmt.Println("token does not match regex pattern, make sure you provide full token")
		os.Exit(1)
	}

	// Strips not neader member characrters. Member: O11111AA becomes: 11111
	r, _ := regexp.Compile("[0-9]{5}")
	*member = r.FindString(*member)
}

// Returns the Environment type based on the given flag
func envType(t, a, p bool, envType *string) {
	if t {
		*envType = "test"
		return
	}

	if a {
		*envType = "accept"
		return
	}

	if p {
		*envType = "production"
		return
	}

	fmt.Println("missing enviroment type flag")
	os.Exit(1)
}

// Sets the shell variable to trigger the httpdump function
func setVerbose(verbose bool) {
	if verbose {
		os.Setenv("AS-VERBOSE", "TRUE")
	}
}

// Sets the shell variable to trigger the httpdump function
func unSetVerbose() {
	os.Unsetenv("AS-VERBOSE")
}

// If location is empty get the working directory
func checkAndSetLocationString(location *string) {
	if *location == "" {
		loc, err := os.Getwd()
		if err != nil {
			fmt.Print("could not get current working directory for setting output location")
			os.Exit(1)
		}

		*location = loc
	}
}

// Converting token to Base64 for in the Auth header
func Base64(t string) string {
	return b64.StdEncoding.EncodeToString([]byte(t))
}
