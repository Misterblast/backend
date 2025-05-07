package config

import (
	"flag"
	"log"
	"os"

	"github.com/ghulammuzz/misterblast/pkg/middleware"
	"github.com/joho/godotenv"
)

func Init() {
	env := flag.String("env", "prod", "Environment for (stg/prod)")
	flag.Parse()

	switch *env {
	case "stg":
		err := godotenv.Load("./stg.env")
		if err != nil {
			log.Println("Error loading stg.env file")
		} else {
			log.Println("Environment: staging (stg.env loaded)")
		}
	default:
		log.Println("Environment: production (using system environment variables)")
	}

	middleware.Log(*env, true, os.Getenv("LOKI_URL"))
	middleware.Validator()
}
