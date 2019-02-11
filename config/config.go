package config

type Config struct {
	Addr string

	MongoDB struct {
		URI string
	}
}
