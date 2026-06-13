package property

import "time"

type Schema struct {
	Server  ServerProperty
	Swagger SwaggerProperty
}

type ServerProperty struct {
	ServiceName        string
	ServiceDescription string
	ServiceVersion     string
	Port               string
	AppEnv             string
}

type SwaggerProperty struct {
	ApiDocs        bool
	ApiDocsVersion string
	ApiDocsHost    string
	ApiDocsSchema  string
}

type TokenProperty struct {
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}
