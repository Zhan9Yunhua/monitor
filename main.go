package main

import (
	"flag"
	"log"
	"os"
	"runtime"
	"strings"
)

var (
	cfg     *Config
	lang    string
	script  string
	curPath string
	output  string
	cfgFile string
	cmdArgs string
	exit    chan bool
)

func init() {
	flag.StringVar(&output, "o", "", "go build output")
	flag.StringVar(&cfgFile, "f", "", "config file")
	flag.StringVar(&cmdArgs, "args", "", "app run args. like: -args='-host=:8080,-name=demo'")
	flag.StringVar(&lang, "lang", "", "language")
	flag.StringVar(&script, "s", "", "run script")
}

func main() {
	flag.Parse()

	curPath, _ = os.Getwd()
	curPath += "/test"

	cfgHandler()

	flagHandler()
	appHandler()
}

func flagHandler() {
	if lang != "" {
		cfg.Lang = lang
	}

	if script != "" {
		cfg.Script = script
	}

	if output != "" {
		cfg.Output = output
	}

	if cmdArgs != "" {
		cfg.CmdArgs = strings.Split(cmdArgs, " ")
	}

	if cfg.Lang != "go" && cfg.Script == "" {
		log.Fatalln("Script can not empty !")
	}
}

func appHandler() {
	paths := []string{}
	parseDir(curPath, &paths)

	mon(paths)
	run()

	for {
		select {
		case <-exit:
			runtime.Goexit()
		}
	}
}
