package main

import (
	"os"

	"github.com/climbcomp/climbcomp-api/cmd"
	"github.com/climbcomp/climbcomp-api/conf"
	log "github.com/sirupsen/logrus"
)

func main() {
	configureLogger()
	app := cmd.NewApp()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func configureLogger() {
	log.SetOutput(os.Stdout)

	config := conf.Instance()
	level, err := log.ParseLevel(config.LogLevel)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(level)

	if config.LogFormat == "json" {
		log.SetFormatter(&log.JSONFormatter{
			// fluentd/stackdriver fieldnames
			FieldMap: log.FieldMap{
				log.FieldKeyTime:  "time",
				log.FieldKeyLevel: "severity",
				log.FieldKeyMsg:   "message",
			},
		})
	} else {
		log.SetFormatter(&log.TextFormatter{
			ForceColors: true,
		})
	}
}
