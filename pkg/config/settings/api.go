package settings

// Client is the interface for interacting with the runtime settings API
type Client interface {
	Get(key string) (interface{}, error)
	Set(key string, value string) (bool, error)
	List() (map[string]RuntimeSettingResponse, error)
	FullConfig() (string, error)
}

// ClientFunc represents a function returning a runtime settings API client
type ClientFunc func() (Client, error)
