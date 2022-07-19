package declare

type ConfigRequest struct {
	Cookies []string
	UA string
	Host string
	Referer string
	TimeOut int
}

type ConfigDatabase struct {
	Url string
	Name string
}
type Config struct {
	Request ConfigRequest
	Database ConfigDatabase
}

