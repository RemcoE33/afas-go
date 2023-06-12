package models

// Single updateconnector information
type UpdateConnector struct {
	ID          string          `json:"id"`
	Description string          `json:"description"`
	Name        string          `json:"name"`
	Fields      []UpdateFields  `json:"fields"`
	Objects     []UpdateObjects `json:"objects"`
}

// Fields information on every property inside the updateconnector
type UpdateFields struct {
	FieldID        string   `json:"fieldId"`
	PrimaryKey     bool     `json:"primaryKey"`
	DataType       string   `json:"dataType"`
	Label          string   `json:"label"`
	Mandatory      bool     `json:"mandatory"`
	Length         int      `json:"length"`
	Decimals       int      `json:"decimals"`
	DecimalFieldID string   `json:"decimalFieldId"`
	Notzero        bool     `json:"notzero"`
	ControlType    int      `json:"controlType"`
	Values         []Values `json:"values,omitempty"`
}

// Enum values for an specific field
type Values struct {
	Id          string `json:"id"`
	Description string `json:"description"`
}

// Nested objects in the updateconnector fields
type UpdateObjects struct {
	Name    string          `json:"name"`
	Fields  []UpdateFields  `json:"fields"`
	Objects []UpdateObjects `json:"objects"`
}
