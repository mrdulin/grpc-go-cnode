package configs

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

const (
	PORT                        string = "PORT"
	BASE_URL                    string = "BASE_URL"
	ACCESS_TOKEN                string = "ACCESS_TOKEN"
	GRPC_GO_LOG_SEVERITY_LEVEL  string = "GRPC_GO_LOG_SEVERITY_LEVEL"
	GRPC_GO_LOG_VERBOSITY_LEVEL string = "GRPC_GO_LOG_VERBOSITY_LEVEL"
)

func Read() *viper.Viper {
	fmt.Println("load env vars from local config.yaml file")
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)
	viper.AddConfigPath(basePath)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	viper.SetDefault(PORT, "3000")
	return viper.GetViper()
}
