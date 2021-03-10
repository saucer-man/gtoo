package main

import (
	"gtoo/cmd"
	"os"

	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
)

func main() {
	cmd.Execute()
}

func init() {
	log.SetFormatter(&nested.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

}
