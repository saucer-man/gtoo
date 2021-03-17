package main

import (
	"gtoo/cmd"
	"os"
	"runtime"

	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
)

func main() {
	cmd.Execute()
}

func init() {
	// 也许不需要，暂时还不需要
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 全局日志设置
	log.SetFormatter(&nested.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

}
