package log

import "github.com/kelseyhightower/envconfig"

type Specs struct {
	Debug      bool   `envconfig:"DEBUG" default:"false"`
	LogFormat  string `envconfig:"LOG_FORMAT" default:"console" example:"console, json"`
	TimeFormat string `envconfig:"TIME_FORMAT" default:"2006-01-02 15:04:05"`
}

var specs Specs

func loadSpecs() error {
	return envconfig.Process("", &specs)
}
