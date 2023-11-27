package main

import (
	"bytes"
	"crypto/rsa"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/urfave/cli/v2"

	"route/auth"
	"route/server"

	bfh "github.com/ajenpan/battle/handler"
	"github.com/ajenpan/battle/msg"
	"github.com/ajenpan/surf/logger"
	"github.com/ajenpan/surf/utils/rsagen"
	utilSignal "github.com/ajenpan/surf/utils/signal"

	_ "github.com/ajenpan/battle/logic/niuniu"
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

	// h.OnReqStartBattle(nil, &msg.ReqStartBattle{
	// 	LogicName:    "niuniu",
	// 	LogicVersion: "1.0.0",
	// 	// CallbackUrl:  "node://1234",
	// 	BattleConf: &msg.BattleConfig{
	// 		BattleId: 1001,
	// 	},
	// 	LogicConf: []byte("{}"),
	// 	PlayerInfos: []*msg.PlayerInfo{
	// 		{
	// 			Uid:       1001,
	// 			SeatId:    1,
	// 			MainScore: 1000,
	// 		},
	// 		{
	// 			Uid:       1002,
	// 			SeatId:    2,
	// 			MainScore: 1000,
	// 		},
	// 	},
	// })

	// tb := h.GetBattleById(1001)
	// tb.OnPlayerJoin(1001)
	// tb.OnPlayerJoin(1002)

	opts := &server.TcpClientOptions{
		RemoteAddress:        "localhost:8080",
		AuthToken:            jwt,
		ReconnectDelaySecond: 10,
		OnMessage: func(tc *server.TcpClient, m *server.Message) {
			h.OnTcpMessage(tc, m)
		},
		OnStatus: func(tc *server.TcpClient, b bool) {
			if b {
				tc.Send(&server.Message{})
			}
		},
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
	defer client.Close()

	client2 := start2(pk)
	client2.Start()
	defer client2.Close()

	s := utilSignal.WaitShutdown()
	logger.Infof("recv signal: %v", s.String())
	return nil
}

func start2(pk *rsa.PrivateKey) *server.TcpClient {

	jwt, _ := auth.GenerateToken(pk, &auth.UserInfo{
		UId:   1010002,
		UName: "user1010002",
		URole: "User",
	}, 24*time.Hour)

	opts := &server.TcpClientOptions{
		RemoteAddress:        "localhost:8080",
		AuthToken:            jwt,
		ReconnectDelaySecond: 10,

		OnMessage: func(s *server.TcpClient, m *server.Message) {
		},
		OnStatus: func(s *server.TcpClient, enable bool) {
			if enable {
				req := &msg.ReqStartBattle{
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
				}
				// s.Send(&server.Message{})
				s.SendReqMsg(110002, req, server.NewTcpRespCallbackFunc(func(err error, c *server.TcpClient, resp *msg.RespStartBattle) {
					if err != nil {
						logger.Errorf("send req failed: %v", err)
						return
					}
					logger.Infof("recv resp: %v", resp)
				}))
			}
		},
	}

	client := server.NewTcpClient(opts)

	return client
}
