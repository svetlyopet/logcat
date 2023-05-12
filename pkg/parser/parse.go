package parser

// RequestLogs stores the request log entries which are read from file
type RequestLogs struct {
	timestamp string
	ip        string
	user      string
	method    string
	path      string
	status    string
	size      string
}

// BillingLogs stores the billing log entry we create
type BillingLogs struct {
	Timestamp       string `json:"billing_timestamp"`
	ServerName      string `json:"server_name"`
	Service         string `json:"service"`
	Action          string `json:"action"`
	RemoteIP        string `json:"ip"`
	Repository      string `json:"repository"`
	Project         string `json:"project"`
	ArtifactoryPath string `json:"artifactory_path"`
	User            string `json:"user_name"`
	ConsumptionUnit string `json:"consumption_unit"`
	Quantity        int64  `json:"quantity"`
}
