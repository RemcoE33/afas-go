package models

// Get all Get/Update connector within the app connector.
type AppConnector struct {
	UpdateConnectors []UpdateConnectors `json:"updateConnectors"`
	GetConnectors    []GetConnectors    `json:"getConnectors"`
	Info             Info               `json:"info"`
}

// Array of updateconnectors
type UpdateConnectors struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

// Array of getconnectors
type GetConnectors struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

// Appconnector information
type Info struct {
	Envid       string `json:"envid"`
	AppName     string `json:"appName"`
	Group       string `json:"group"`
	TokenExpiry string `json:"tokenExpiry"`
}
