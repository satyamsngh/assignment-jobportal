package config

import (
	"log"

	env "github.com/Netflix/go-env"
)

var cfg Config

type Config struct {
	AppConfig          AppConfig
	DataConfig         DataConfig
	RedisConfig        RedisConfig
	KeyConfig          KeyConfig
	OtpGeneratorConfig OtpGeneratorConfig
}

type AppConfig struct {
	Port string `env:"APP_PORT,required=true"`
}

type DataConfig struct {
	Host     string `env:"POSTGRES_HOST,required=true"`
	UserName string `env:"POSTGRES_USER,required=true"`
	Password string `env:"POSTGRES_PASSWORD,required=true"`
	DBName   string `env:"POSTGRES_DB,required=true"`
	Port     string `env:"POSTGRES_PORT,default=5432"`
	SSLMode  string `env:"POSTGRES_SSL_MODE,default=false"`
	Time     string `env:"POSTGRES_TIME,default=Asia/Shanghai"`
}
type RedisConfig struct {
	Host     string `env:"REDIS_HOST,default=localhost"`
	Port     string `env:"REDIS_PORT,default=6379"`
	Password string `env:"REDIS_PASSWORD,default=false"`
	DB       int    `env:"REDIS_DB,default=false"`
}
type KeyConfig struct {
	PublicKeyPath  string `env:"PUBLIC_KEY_PATH"`
	PrivateKeyPath string `env:"PRIVATE_KEY_PATH"`
}
type OtpGeneratorConfig struct {
	Port int `env:"APP_PORT_SMTP,required=true"`
}

func init() {
	_, err := env.UnmarshalFromEnviron(&cfg)
	if err != nil {
		log.Panic(err)
	}
}

func GetConfig() Config {
	return cfg
}
