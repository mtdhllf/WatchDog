package main

import (
	"fmt"
	"github.com/go-vgo/robotgo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os/exec"
	"strings"
	"time"
)

var failCount = 0

func initConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Errorf("Fatal error config file: %s", err)
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

//wmic process where name="adb.exe" get ProcessId,ExecutablePath
func main() {
	initConfig()
	for {
		run()
		time.Sleep(time.Duration(viper.GetInt64("global.interval")) * time.Millisecond)
	}
}

type program struct {
	Name string
	Path string
}

func run() {
	pp := getProgram()
	if len(pp) > 0 {
		for _, p := range pp {
			fpid, _ := robotgo.FindIds(p.Name)
			fmt.Println("fpid---", fpid)
			if len(fpid) == 0 {
				failCount++
				if failCount > 2 {
					c := exec.Command("cmd", "/C", "start", p.Path)
					if err := c.Run(); err != nil {
						logrus.Error("Error: ", err)
					} else {
						logrus.Infof("start %s ", p.Path)
						failCount = 0
					}
				}
			}
		}
	}
}

func getProgram() []program {
	all := viper.AllKeys()
	pp := make([]program, 0)
	for _, item := range all {
		if strings.HasSuffix(item, ".executablepath") {
			p := program{
				Name: strings.ReplaceAll(item, ".executablepath", ""),
				Path: viper.GetString(item),
			}
			pp = append(pp, p)
		}
	}
	logrus.Info(pp)
	return pp
}
