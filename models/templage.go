package models

// To pass along the tempalate Execute function.
type TemplateGetData struct {
	Connector  GetConnector
	ImportTime bool
}

// To pass along the tempalate Execute function.
type TemplateUpdateData struct {
	Connector   UpdateConnector
	ImportTime  bool
	FlatObjects map[string]UpdateObjects
}
