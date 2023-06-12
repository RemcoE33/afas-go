package models

// Single getconnector information
type GetConnector struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Fields      []GetFields `json:"fields"`
}

// Fields information on every property inside the getconnector
type GetFields struct {
	ID              string `json:"id"`
	FieldID         string `json:"fieldId"`
	DataType        string `json:"dataType"`
	Label           string `json:"label"`
	Length          int    `json:"length"`
	ControlType     int    `json:"controlType"`
	Decimals        int    `json:"decimals"`
	DecimalsFieldID string `json:"decimalsFieldId"`
}
