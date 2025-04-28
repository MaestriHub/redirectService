package configs

import (
	"os"
)

type Config struct {
	*DbConfig
	*AppStoreLinksConfig
}

type DbConfig struct {
	Port     string
	Name     string
	Password string
	Host     string
	Username string
}

type AppStoreLinksConfig struct {
	AndroidLink string
	IOSLink     string
}

func NewConfig() *Config {
	return &Config{
		DbConfig:            NewDBConfig(),
		AppStoreLinksConfig: NewAppStoreLinksConfig(),
	}
}

func NewDBConfig() *DbConfig {
	return &DbConfig{
		Port:     os.Getenv("DATABASE_PORT"),
		Name:     os.Getenv("DATABASE_NAME"),
		Password: os.Getenv("DATABASE_PASSWORD"),
		Host:     os.Getenv("DATABASE_HOST"),
		Username: os.Getenv("DATABASE_USERNAME"),
	}
}

func NewAppStoreLinksConfig() *AppStoreLinksConfig {
	return &AppStoreLinksConfig{
		AndroidLink: os.Getenv("ANDROID_APP_STORE_LINK"),
		IOSLink:     os.Getenv("IOS_APP_STORE_LINK"),
	}
}

// TODO: допилить config для cors
// type CORSConfig struct {
// 	AllowOrigins     []string,
// 	AllowMethods     []string,
// 	AllowHeaders     []string,
// 	ExposeHeaders    []string,
// 	AllowCredentials bool,
// 	AllowOriginFunc  func(origin string) bool {
// 	return origin == "https://github.com"
// },
// 	MaxAge: 12 * time.Hour,
// }
