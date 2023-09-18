package config

var (
	ApiUrl = "https://probe.fbrq.cloud/v1/send/"
	JWT    = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjYzMjIxMzAsImlzcyI6ImZhYnJpcXVlIiwibmFtZSI6Imh0dHBzOi8vdC5tZS90dWxsanVzIn0.iBkqElFmI0gisF2Ctuomb12IQoMYSX2yaztnth-_E8c"
)

type Config struct {
	Server   Server   `yaml:"Server"`
	Postgres Postgres `yaml:"Postgres"`
}

type Server struct {
	Port string `yaml:"port"`
}

type Postgres struct {
	DB             string `yaml:"db"`
	ConnectionToDB string `yaml:"connectionToDB"`
}
