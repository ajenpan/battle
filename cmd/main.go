package main

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/urfave/cli/v2"

	"route/auth"
	"route/server"

	bfh "github.com/ajenpan/battlefield/handler"
	"github.com/ajenpan/battlefield/msg"
	"github.com/ajenpan/surf/logger"
	"github.com/ajenpan/surf/utils/rsagen"
	utilSignal "github.com/ajenpan/surf/utils/signal"

	_ "github.com/ajenpan/battlefield/logic/niuniu"
)

var (
	Name       string = "battlefield"
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

	pk, err := rsagen.LoadRsaPrivateKeyFromFile("private.pem")
	if err != nil {
		return err
	}

	jwt, _ := auth.GenerateToken(pk, &auth.UserInfo{
		UId:   110002,
		UName: "battle",
		URole: "battle",
	}, 24*time.Hour)

	h := bfh.New()

	h.OnReqStartBattle(nil, &msg.ReqStartBattle{
		LogicName:    "niuniu",
		LogicVersion: "1.0.0",
		// CallbackUrl:  "node://1234",
		BattleConf: &msg.BattleConfig{
			BattleId: 1001,
		},
		LogicConf: []byte("{}"),
		PlayerInfos: []*msg.PlayerInfo{
			{
				Uid:       1001,
				SeatId:    1,
				MainScore: 1000,
			},
			{
				Uid:       1002,
				SeatId:    2,
				MainScore: 1000,
			},
		},
	})

	tb := h.GetBattleById(1001)

	tb.OnPlayerJoin(1001)
	tb.OnPlayerJoin(1002)

	opts := &server.TcpClientOptions{
		RemoteAddress:        "101.133.149.209:8080",
		AuthToken:            jwt,
		ReconnectDelaySecond: 10,
		OnSessionMessage:     h.OnSessionMessage,
		OnSessionStatus:      h.OnSessionStatus,
	}

	client := server.NewTcpClient(opts)
	// opts := &server.TcpServerOptions{
	// 	ListenAddr:       ":12002",
	// 	AuthPublicKey:    pk,
	// 	OnSessionMessage: h.OnSessionMessage,
	// 	OnSessionStatus:  h.OnSessionStatus,
	// }
	// svr, err := server.NewTcpServer(opts)

	err = client.Start()
	if err != nil {
		return err
	}

	defer client.Stop()

	s := utilSignal.WaitShutdown()
	logger.Infof("recv signal: %v", s.String())
	return nil
}
