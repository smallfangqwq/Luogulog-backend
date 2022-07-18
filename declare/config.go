package declare

type ConfigRequest struct {
	Cookies []string
	UA string
	Host string
	Referer string
	TimeOut int
}

type Config struct {
	Request ConfigRequest
}