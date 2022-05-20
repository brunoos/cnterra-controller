package config

import "os"

var (
	BaseDir = "/opt/cnterra-loader/"
	TmpDir  = "/opt/cnterra-loader/tmp/"

	Address = "0.0.0.0"
	Port    = "8080"

	DbAddress  = "cnterra-mongo"
	DbPort     = "27017"
	DbName     = "cnterra"
	DbUser     = "cnterra"
	DbPassword = "cnterra"
)

func Initialize() {
	if str, found := os.LookupEnv("CNTERRA_ADDRESS"); found {
		Address = str
	}
	if str, found := os.LookupEnv("CNTERRA_PORT"); found {
		Port = str
	}

	if str, found := os.LookupEnv("MONGO_ADDRESS"); found {
		DbAddress = str
	}
	if str, found := os.LookupEnv("MONGO_PORT"); found {
		DbPort = str
	}
	if str, found := os.LookupEnv("MONGO_DATABASE"); found {
		DbName = str
	}
	if str, found := os.LookupEnv("MONGO_USER"); found {
		DbUser = str
	}
	if str, found := os.LookupEnv("MONGO_PASSWORD"); found {
		DbPassword = str
	}
}
