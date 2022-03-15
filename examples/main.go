package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/geoffomen/go-app/pkg/config"
	"github.com/geoffomen/go-app/pkg/config/viperimp"
	"github.com/geoffomen/go-app/pkg/cronjob"
	"github.com/geoffomen/go-app/pkg/cronjob/robfigimp"
	"github.com/geoffomen/go-app/pkg/database"
	"github.com/geoffomen/go-app/pkg/database/gormimp"
	"github.com/geoffomen/go-app/pkg/mylog"
	"github.com/geoffomen/go-app/pkg/mylog/zapimp"
	"github.com/geoffomen/go-app/pkg/webfw"
	"github.com/geoffomen/go-app/pkg/webfw/ginimp"

	"github.com/geoffomen/go-app/examples/hello/controller"
	"github.com/geoffomen/go-app/examples/user/userctl"
	"github.com/geoffomen/go-app/examples/useracc/useraccctl"
)

var (
	branchName string
	commitId   string
	buildTime  string

	showVer = flag.Bool("v", false, "show version")
)

func main() {
	if *showVer {
		fmt.Printf("%s: %s\t%s\n", branchName, commitId, buildTime)
		os.Exit(0)
	}

	profile := flag.String("profile", "example", "Environment profile, something similar to spring profiles")
	flag.Parse()
	vp, err := viperimp.New(*profile)
	if err != nil {
		panic(fmt.Sprintf("failed to initrialize config component, err: %v", err))
	}
	config.SetInstance(vp)
	cf := config.GetInstance()

	zp, err := zapimp.New(zapimp.Configuration{
		EnableConsole:     cf.GetBoolOrDefault("log.enableConsole", true),
		ConsoleJSONFormat: cf.GetBoolOrDefault("log.consoleJSONFormat", true),
		ConsoleLevel:      cf.GetStringOrDefault("log.consoleLevel", "debug"),
		EnableFile:        cf.GetBoolOrDefault("log.enableFile", true),
		FileJSONFormat:    cf.GetBoolOrDefault("log.fileJSONFormat", true),
		FileLevel:         cf.GetStringOrDefault("log.fileLevel", "info"),
		FileLocation:      cf.GetStringOrDefault("log.fileLocation", "/tmp/miis/back/info.log"),
		ErrFileLevel:      cf.GetStringOrDefault("log.errFileLevel", "error"),
		ErrFileLocation:   cf.GetStringOrDefault("log.errFileLocation", "/tmp/miis/back/err.log"),
	})
	if err != nil {
		panic(fmt.Sprintf("failed to initrialize logger component, err: %s", err))
	}
	mylog.SetInstance(zp)

	rp, err := robfigimp.New(mylog.GetInstance())
	if err != nil {
		panic(fmt.Sprintf("failed to initrialize cronjob component, err: %v", err))
	}
	cronjob.SetInstance(rp)
	cronjob.GetInstance().Start()

	db, err := gormimp.NewGorm(gormimp.GormConfig{
		Dialect:     cf.GetStringOrDefault("database.dialect", ""),
		UserName:    cf.GetStringOrDefault("database.userName", ""),
		Password:    cf.GetStringOrDefault("database.password", ""),
		Host:        cf.GetStringOrDefault("database.host", "localhost"),
		Port:        cf.GetIntOrDefault("database.port", 3306),
		Db:          cf.GetStringOrDefault("database.db", "test"),
		OtherParams: cf.GetStringOrDefault("database.otherParams", ""),
	}, mylog.GetInstance())
	if err != nil {
		panic(fmt.Sprintf("failed to initrialize database component, err: %v", err))
	}
	database.SetInstance(db)

	gp, err := ginimp.New(cf, mylog.GetInstance())
	if err != nil {
		panic(fmt.Sprintf("failed to initrialize webfw component, err: %v", err))
	}
	webfw.SetInstance(gp)
	webfw.GetInstance().RegisterHandler(controller.Controller())
	webfw.GetInstance().RegisterHandler(useraccctl.Controller())
	webfw.GetInstance().RegisterHandler(userctl.Controller())
	webfw.GetInstance().Start()
}
