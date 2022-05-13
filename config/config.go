package config

import "os"

var (
	BaseDir = "/opt/cnterra-loader/"
	TmpDir  = "/opt/cnterra-loader/tmp/"

	DbAddress  = "localhost"
	DbPort     = "5432"
	DbName     = "cnterra"
	DbUser     = "cnterra"
	DbPassword = "cnterra"
	DbSslMode  = "disable"
)

func Initialize() {
	if str, found := os.LookupEnv("POSTGRES_ADDRESS"); found {
		DbAddress = str
	}
	if str, found := os.LookupEnv("POSTGRES_PORT"); found {
		DbPort = str
	}
	if str, found := os.LookupEnv("POSTGRES_DATABASE"); found {
		DbName = str
	}
	if str, found := os.LookupEnv("POSTGRES_USER"); found {
		DbUser = str
	}
	if str, found := os.LookupEnv("POSTGRES_PASSWORD"); found {
		DbPassword = str
	}
	if str, found := os.LookupEnv("POSTGRES_SSLMODE"); found {
		DbSslMode = str
	}
}
