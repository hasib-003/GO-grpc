package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

const DEFAULT_MAX_OPEN_CONNS = 5

type EnumDimensionType string

type (

	// Config -.
	Config struct {
		Auth      `mapstructure:",squash"`
		JwtSecret string `env:"JWT_SECRET"`
		HTTP      `mapstructure:",squash"`
		GRPC      `mapstructure:",squash"`
		Database  `mapstructure:",squash"`
		Logger    `mapstructure:",squash"`
		MapApiKey string `env:"MAP_API_KEY"`
		Image     `mapstructure:",squash"`
	}
	Auth struct {
	}
	// App -.
	App struct {
		Name    string `mapstructure:"NAME" json:"name"`
		Version string `mapstructure:"VERSION" json:"version"`
	}

	// HTTP -.
	HTTP struct {
		HTTPAddress string `env:"HTTP_ADDRESS"`
	}
	GRPC struct {
		GrpcPort string `env:"GRPC_PORT"`
	}

	// DB -.
	Database struct {
		DBHost          string `env:"DBHOST"`
		DbUser          string `env:"DBUSER"`
		DbPass          string `env:"DBPASS"`
		DbPort          string `env:"DBPORT"`
		DbName          string `env:"DBNAME"`
		DbSchema        string `env:"DBSCHEMA"`
		SetMaxOpenConns int    `env:"SETMAXOPENCONNS" env-default:"5"`
		Debug           bool   `env:"DEBUG" env-default:"false"`
	}

	// Logger config
	Logger struct {
		Development       bool
		DisableCaller     bool
		DisableStacktrace bool
		Encoding          string
		Level             string
	}

	// Image config
	Image struct {
		ThumbnailWidth  int   `env:"THUMBNAIL_WIDTH"`
		ThumbnailHeight int   `env:"THUMBNAIL_HEIGHT"`
		ImazeMaxSize    int   `env:"IMAGE_MAX_SIZE"`
		VideoMaxSize    int64 `env:"VIDEO_MAX_SIZE"`
	}
)

func NewConfig(configFile string) *Config {
	config := Config{}
	godotenv.Load(configFile)
	err := cleanenv.ReadEnv(&config)
	if err != nil {
		log.Fatalln(err)
	}
	return &config
}
