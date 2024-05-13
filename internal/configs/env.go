package configs

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type env struct {
	DBSTRING  string `env:"DBSTRING, required"`
	Port      string `env:"PORT, required"`
	JWTSECRET string `env:"JWT_SECRET,required"`
}

// initiating because godotenv doesn't play nice with docker.
var Envs = &env{
	DBSTRING:  os.Getenv("DBSTRING"),
	Port:      os.Getenv("PORT"),
	JWTSECRET: os.Getenv("JWT_SECRET"),
}

func InitEnv(ctx context.Context) error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("error on loading gotdotenv: %v", err.Error())
	}
	if err := envconfig.Process(ctx, Envs); err != nil {
		return fmt.Errorf("error on loading godotenv process: %v", err.Error())
	}
	return nil
}
