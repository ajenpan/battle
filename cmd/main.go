package main

import (
	"bytes"
	"fmt"
	"os"
	"runtime"

	"github.com/urfave/cli/v2"

	battleHandler "github.com/ajenpan/battle/handler"
	"github.com/ajenpan/battle/logger"
	"github.com/ajenpan/battle/proto"
	"github.com/ajenpan/battle/utils/calltable"

	utilSignal "github.com/ajenpan/battle/utils/signal"
)

var (
	Name       string = "battle"
	Version    string = "unknow"
	GitCommit  string = "unknow"
	BuildAt    string = "unknow"
	BuildBy    string = runtime.Version()
	RunnningOS string = runtime.GOOS + "/" + runtime.GOARCH
)

func longVersion() string {
	buf := bytes.NewBuffer(nil)
	fmt.Fprintln(buf, "project:", Name)
	fmt.Fprintln(buf, "version:", Version)
	fmt.Fprintln(buf, "git commit:", GitCommit)
	fmt.Fprintln(buf, "build at:", BuildAt)
	fmt.Fprintln(buf, "build by:", BuildBy)
	fmt.Fprintln(buf, "running OS/Arch:", RunnningOS)
	return buf.String()
}

func main() {
	err := Run()
	if err != nil {
		fmt.Println(err)
	}
}

func Run() error {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Println(longVersion())
	}
	app := cli.NewApp()
	app.Version = Version
	app.Name = Name
	err := app.Run(os.Args)
	return err
}

func RealMain(c *cli.Context) error {
	h := battleHandler.New()

	h.CT = calltable.ExtractAsyncMethod(proto.File_proto_battle_client_proto.Messages(), h)

	s := utilSignal.WaitShutdown()
	logger.Infof("recv signal: %v", s.String())
	return nil
}
