package config

import (
	"os"
	"strconv"
	"sync"
)

type WebServer struct {
	SERVER_PORT int
	TIMEOUT     int
}

type Ollama struct {
	Endpoint string
	Model    string
}

type Config struct {
	WebServer WebServer
	Ollama    Ollama
}

var once sync.Once

var config *Config

func GetInstance() *Config {

	if config == nil {
		once.Do(
			func() {
				webserver_config := WebServer{
					SERVER_PORT: parseEnvToInt("SERVER_PORT", "5000"),
					TIMEOUT:     parseEnvToInt("TIMEOUT", "10"),
				}

				ollama_config := Ollama{
					Endpoint: getEnv("OLLAMA_ENDPOINT", "http://localhost:11434"),
					Model:    getEnv("OLLAMA_MODEL", "llama3.2:3b"),
				}

				config = &Config{
					WebServer: webserver_config,
					Ollama:    ollama_config,
				}
			},
		)
	}

	return config

}

func parseEnvToInt(envName, defaultValue string) int {
	num, err := strconv.Atoi(getEnv(envName, defaultValue))
	if err != nil {
		return 0
	}
	return num
}

func parseEnvToBool(envName string, defaultValue string) bool {
	boolValue, err := strconv.ParseBool(getEnv(envName, defaultValue))
	if err != nil {
		return boolValue
	}
	return boolValue
}

func getEnv(env, defaultValue string) string {
	enviroment := os.Getenv(env)
	if enviroment == "" {
		return defaultValue
	}

	return enviroment
}
