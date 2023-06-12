package afas

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"regexp"

	"github.com/RemcoE33/go-afas-appconnector-cli/models"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Generic function to call AFAS endpoints.
func AfasRequest[T any](s *T, endpoint, member, envtype, token string) error {
	et := models.EntTypeMapping[envtype]
	url := fmt.Sprintf("https://%v.%v.afas.online/ProfitRestServices/%v", member, et, endpoint)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return fmt.Errorf("http.NewRequest: %w", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "AfasToken "+token)

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("client.Do: %w", err)
	}
	defer res.Body.Close()

	if os.Getenv("AS-VERBOSE") == "TRUE" {
		httpDumpLog(res)
	}

	if res.StatusCode >= 400 {
		fmt.Println()
		return fmt.Errorf("return status: %d", res.StatusCode)
	}

	if err := json.NewDecoder(res.Body).Decode(&s); err != nil {
		return fmt.Errorf("json.Decode: %w", err)
	}

	return nil
}

// Log the response output
func httpDumpLog(res *http.Response) {
	respDump, err := httputil.DumpResponse(res, true)

	if err != nil {
		log.Fatal("httputil.DumpResponse: ", err)
	}

	fmt.Printf("RESPONSE:\n%s", string(respDump))
}

// To CamelCase the struct title and the struct fields
func Title(title string) string {
	c := cases.Title(language.English)
	replacedAllSpecialChars := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	var str string

	for _, s := range replacedAllSpecialChars.Split(title, -1) {
		str += c.String(s)
	}

	return str
}

// To map the type from AFAS to the Go type
func GetType(t string) string {
	return models.TypesMap[t]
}

// Check if a field need time.Time for the import statement in the template
func hasUpdateDateField(fields []models.UpdateFields) bool {
	for _, f := range fields {
		if f.DataType == "date" {
			return true
		}
	}
	return false
}

// Check if a field need time.Time for the import statement in the template
func hasGetDateField(fields []models.GetFields) bool {
	for _, f := range fields {
		if f.DataType == "date" {
			return true
		}
	}
	return false
}
