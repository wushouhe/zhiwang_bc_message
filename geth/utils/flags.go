package utils

import (
	"gopkg.in/urfave/cli.v1"
	"io/ioutil"
	. "zhiwang_bc_message/geth/config"
	"gopkg.in/yaml.v2"
	"errors"
	"path/filepath"
	"os"
	"flag"
	"fmt"
)


func NewApp(usage string) *cli.App {
	app := cli.NewApp()
	app.Name = filepath.Base(os.Args[0])
	app.Author = "gcc2ge"
	app.Email = "1414786333@qq.com"
	app.Version = Version
	app.Usage = usage
	return app
}

var (
	ConfigFileFlag = cli.StringFlag{
		Name:"config",
		Usage:"为程序提供启动配置文件",
		Value:"./cfg/config.yaml",
	}
)

func ReadConfig(ctx *cli.Context) error {
	configFile := ctx.GlobalString(ConfigFileFlag.Name)
	if configFile == "" {
		return errors.New("config file not set")
	}
	source, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal([]byte(source), &Cfg)
	if err != nil {
		return err
	}
	return nil
}


//glog
func glogFlagShim(fakeVals map[string]string) {
	flag.VisitAll(func(fl *flag.Flag) {
		if val, ok := fakeVals[fl.Name]; ok {
			fl.Value.Set(val)
		}
	})
}

func GlogGangstaShim(c *cli.Context) {
	_ = flag.CommandLine.Parse([]string{})
	glogFlagShim(map[string]string{
		"v":                fmt.Sprint(c.Int("v")),
		"logtostderr":      fmt.Sprint(c.Bool("logtostderr")),
		"stderrthreshold":  fmt.Sprint(c.Int("stderrthreshold")),
		"alsologtostderr":  fmt.Sprint(c.Bool("alsologtostderr")),
		"vmodule":          c.String("vmodule"),
		"log_dir":          c.String("log_dir"),
		"log_backtrace_at": c.String("log_backtrace_at"),
	})
}

var GlogGangstaFlags = []cli.Flag{
	//cli.IntFlag{
	//	Name: "v", Value: 0, Usage: "log level for V logs",
	//},
	cli.BoolFlag{
		Name: "logtostderr", Usage: "log to standard error instead of files",
	},
	cli.IntFlag{
		Name:  "stderrthreshold",
		Usage: "logs at or above this threshold go to stderr",
	},
	cli.BoolFlag{
		Name: "alsologtostderr", Usage: "log to standard error as well as files",
	},
	cli.StringFlag{
		Name:  "vmodule",
		Usage: "comma-separated list of pattern=N settings for file-filtered logging",
	},
	cli.StringFlag{
		Name: "log_dir", Usage: "If non-empty, write log files in this directory",
	},
	cli.StringFlag{
		Name:  "log_backtrace_at",
		Usage: "when logging hits line file:N, emit a stack trace",
		Value: ":0",
	},
}