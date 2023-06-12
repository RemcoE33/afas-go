package models

// Struct for new added apps in the config.json file
type StoredApps struct {
	Member          string
	EnvironmentType string
	Desciption      string
	Token           string
	Location        string
}
