package main

import (
	"bytes"
	"fmt"
	"os"
	"runtime"

	"github.com/urfave/cli/v2"

	bfh "github.com/ajenpan/battlefield/handler"
	"github.com/ajenpan/surf/logger"
	"github.com/ajenpan/surf/tcp"
	"github.com/ajenpan/surf/utils/rsagen"
	utilSignal "github.com/ajenpan/surf/utils/signal"
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
	if err := Run(); err != nil {
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
	app.Action = RealMain
	err := app.Run(os.Args)
	return err
}

func RealMain(c *cli.Context) error {
	pk, err := rsagen.LoadRsaPublicKeyFromFile("public.pem")
	if err != nil {
		return err
	}

	h := bfh.New()

	svr, err := tcp.NewServer(tcp.ServerOptions{
		Address:   ":12001",
		OnMessage: h.OnMessage,
		OnConn:    h.OnConnect,
		AuthFunc:  tcp.RsaTokenAuth(pk),
	})
	if err != nil {
		return err
	}

	svr.Start()
	defer svr.Stop()

	s := utilSignal.WaitShutdown()
	logger.Infof("recv signal: %v", s.String())
	return nil
}
