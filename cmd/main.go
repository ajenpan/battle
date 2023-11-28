package main

import (
	"bytes"
	"crypto/rsa"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/ajenpan/surf/auth"
	"github.com/ajenpan/surf/server"

	"github.com/ajenpan/battle/msg"
	"github.com/ajenpan/surf/log"
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

	opts := &server.TcpClientOptions{
		RemoteAddress:        "localhost:8080",
		AuthToken:            jwt,
		ReconnectDelaySecond: 10,
		OnMessage: func(tc *server.TcpClient, m *server.MsgWraper) {
		},
		OnStatus: func(tc *server.TcpClient, b bool) {

		},
	}

	client := server.NewTcpClient(opts)

	err = client.Connect()
	if err != nil {
		return err
	}
	defer client.Close()

	client2 := start2(pk)
	client2.Connect()
	defer client2.Close()

	s := utilSignal.WaitShutdown()
	log.Infof("recv signal: %v", s.String())
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

		OnMessage: func(s *server.TcpClient, m *server.MsgWraper) {
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
				s.SendReqMsg(110002, req, server.NewTcpRespCallbackFunc(func(c server.Session, resp *msg.RespStartBattle, err error) {
					if err != nil {
						log.Errorf("send req failed: %v", err)
						return
					}
					log.Infof("recv resp: %v", resp)

					if resp.BattleId != 0 {
						req := &msg.ReqJoinBattle{
							BattleId: resp.BattleId,
						}
						s.SendReqMsg(110002, req, server.NewTcpRespCallbackFunc(func(c server.Session, resp *msg.RespJoinBattle, err error) {
							if err != nil {
								log.Errorf("send req failed: %v", err)
								return
							}
							log.Infof("recv resp: %v", resp)
						}))

					}

				}))
			}
		},
	}

	client := server.NewTcpClient(opts)

	return client
}
