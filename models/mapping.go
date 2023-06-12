package models

// Mapping between the types from the AFAS metadata API and Go types.
var TypesMap = map[string]string{
	"string":  "string",
	"blob":    "string",
	"boolean": "bool",
	"int":     "int",
	"decimal": "float64",
	"date":    "time.Time",
}

// Mapping between CLI environment type and the AFAS url.
var EntTypeMapping = map[string]string{
	"test":       "resttest",
	"accept":     "restaccept",
	"production": "rest",
}
