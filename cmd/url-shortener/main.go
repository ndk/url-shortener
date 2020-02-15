package main

import (
	"github.com/joeshaw/envdecode"

	"url-shortener/internal/logger"
	"url-shortener/internal/logger/log"
	"url-shortener/internal/service"
)

var (
	gitCommit = "undefined"
	gitBranch = "undefined"
)

func main() {
	cfg := &service.Config{}
	if err := envdecode.StrictDecode(cfg); err != nil {
		log.Fatal().Err(err).Str("git_commit", gitCommit).Str("git_branch", gitBranch).Msg("Cannot decode the envs to the config")
	}

	l := logger.NewLogger(&cfg.Logger)
	l.Info().Str("git_commit", gitCommit).Str("git_branch", gitBranch).Interface("config", cfg).Msg("The gathered config")

	l.Info().Msg("Running the service...")
	if err := service.Run(cfg, l); err != nil {
		l.Fatal().Err(err).Msg("The service has been stopped with the error")
	}
	l.Info().Msg("The service has been stopped")
}
