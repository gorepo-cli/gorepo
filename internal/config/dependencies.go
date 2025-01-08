package config

import (
	"gorepo-cli/pkg"
)

type Dependencies struct {
	Effects *pkg.Effects
	Config  *Config
}

func NewDependencies(e *pkg.Effects, c *Config) *Dependencies {
	return &Dependencies{
		Effects: e,
		Config:  c,
	}
}
