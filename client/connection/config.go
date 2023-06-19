package connection

import "time"

type Config struct {
	Username string
	Password string
	Host     string
	Port     int
	Insecure bool
	Cacert   string
	Timeout  time.Duration
}
