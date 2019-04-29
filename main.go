package main

import (
	"flag"
	"os"
	"runtime"
	"strings"

	"github.com/Zhan9Yunhua/logger"
)

var (
	cfg       *Config
	lang      string
	script    string
	curPath   string
	output    string
	cfgFile   string
	cmdArgs   string
	buildPath string
	exit      chan bool
)

func init() {
	flag.StringVar(&output, "o", "", "go build output")
	flag.StringVar(&cfgFile, "f", "", "config file")
	flag.StringVar(&buildPath, "b", "", "build package path")
	flag.StringVar(&cmdArgs, "args", "", "app run args. like: -args='-host=:8080,-name=demo'")
	flag.StringVar(&lang, "lang", "", "language")
	flag.StringVar(&script, "s", "", "run script")
}

func main() {
	logger.Infoln("monitor starting...")
	flag.Parse()

	curPath, _ = os.Getwd()
	//curPath += "/test"

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

	if buildPath != "" {
		cfg.BuildPath = buildPath
	}

	if cfg.Lang != "go" && cfg.Script == "" {
		logger.Fatalln("Script can not empty !")
	}

	if len(cfg.Exts) == 0 || !isIn(cfg.Exts, "."+cfg.Lang) {
		cfg.Exts = append(cfg.Exts, "."+cfg.Lang)
	}

	if cfg.BuildPath != "" {
		cfg.Files = strings.Split(cfg.BuildPath, ",")
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
