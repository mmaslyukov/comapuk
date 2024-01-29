package util 

type Config struct {
	Login string `toml:"login"`
	Password string `toml:"password"`
	Certificate string
}

