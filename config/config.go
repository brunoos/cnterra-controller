package config

import "os"

var (
	BaseDir = "/opt/cnterra-loader/"
	TmpDir  = "/opt/cnterra-loader/tmp/"

	DbAddress  = "localhost"
	DbPort     = "27017"
	DbName     = "cnterra"
	DbUser     = "cnterra"
	DbPassword = "cnterra"
)

func Initialize() {
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
