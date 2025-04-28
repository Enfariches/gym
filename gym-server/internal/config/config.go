package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env     string     `yaml:"env" env-default:"dev"`
	Storage string     `yaml:"storage" env-required:"true"`
	TokenTTL time.Duration `yaml:"token_ttl" env-required:"true"`
	AuthTokenTTL time.Duration `yaml:"auth_token_ttl" env-required:"true"`
	SMTP SMTPConfig `yaml:"smtp" env-required:"true"`
	GRPC    GRPCConfig `yaml:"grpc" env-required:"true"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type SMTPConfig struct {
	EmailSender string `yaml:"email_sender"`
	EmailPassword string `yaml:"email_password"`
	Username string `yaml:"username"`
	Host string `yaml:"host"`
	Port int `yaml:"port"`
}

var cfg Config

func MustLoad() *Config{
	path := fetchConfigPath()
	if path == ""{
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err){
		panic("config file does not exist: " + path)
	}

	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil{
		panic(err)
	}

	return &cfg
}

func fetchConfigPath() string{
	var path string

	flag.StringVar(&path, "config", "", "path to config file")
	flag.Parse()

	if path == ""{
		path = os.Getenv("CONFIG_PATH")
	}

	return path
}
