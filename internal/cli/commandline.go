package cli

/*

type configGetter interface {
	GetExternalURLKitFilename() string
	GetConfigFilename() string
	GetDebug() (bool, bool)
	GetHome() string
	GetLogFile() string
	GetUpdaterConfigFilename() string
	GetVDebugSetting() string
}

type CommandLine interface {
	configGetter

	// Lower-level functions
	GetGString(string) string
	GetString(string) string
	GetBool(string, bool) (bool, bool)
}


type LogContext interface {
	GetLog() logger.Logger
}

type VDebugLog struct {
}

type VLogContext interface {
	LogContext
	GetVDebugLog() *VDebugLog
}


// APIContext defines methods for accessing API server
type APIContext interface {
	GetAPI() API
	GetExternalAPI() ExternalAPI
	GetServerURI() (string, error)
}


type API interface {
	Get(MetaContext, APIArg) (*APIRes, error)
	GetDecode(MetaContext, APIArg, APIResponseWrapper) error
	GetDecodeCtx(context.Context, APIArg, APIResponseWrapper) error
	GetResp(MetaContext, APIArg) (*http.Response, func(), error)
	Post(MetaContext, APIArg) (*APIRes, error)
	PostJSON(MetaContext, APIArg) (*APIRes, error)
	PostDecode(MetaContext, APIArg, APIResponseWrapper) error
	PostDecodeCtx(context.Context, APIArg, APIResponseWrapper) error
	PostRaw(MetaContext, APIArg, string, io.Reader) (*APIRes, error)
	Delete(MetaContext, APIArg) (*APIRes, error)
}

type ExternalAPI interface {
	Get(MetaContext, APIArg) (*ExternalAPIRes, error)
	Post(MetaContext, APIArg) (*ExternalAPIRes, error)
	GetHTML(MetaContext, APIArg) (*ExternalHTMLRes, error)
	GetText(MetaContext, APIArg) (*ExternalTextRes, error)
	PostHTML(MetaContext, APIArg) (*ExternalHTMLRes, error)
}
*/