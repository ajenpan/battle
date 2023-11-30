package main

import (
	"bytes"
	"crypto/rsa"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/ajenpan/surf"
	"github.com/ajenpan/surf/auth"
	"github.com/ajenpan/surf/log"
	"github.com/ajenpan/surf/route"
	"github.com/ajenpan/surf/server"
	"github.com/ajenpan/surf/utils/rsagen"
	utilSignal "github.com/ajenpan/surf/utils/signal"

	"github.com/ajenpan/battle/handler"
	"github.com/ajenpan/battle/msg"

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

var battleuid = uint32(1001)

func RealMain(c *cli.Context) error {
	pk, err := rsagen.LoadRsaPrivateKeyFromFile("private.pem")
	if err != nil {
		return err
	}

	usr1, _ := auth.GenerateToken(pk, &auth.UserInfo{
		UId:   10001,
		UName: "user1",
		URole: "user",
	}, 24*time.Hour)
	fmt.Println(usr1)

	jwt, _ := auth.GenerateToken(pk, &auth.UserInfo{
		UId:   battleuid,
		UName: "battle",
		URole: "battle",
	}, 24*time.Hour)

	svr := server1(pk.Public().(*rsa.PublicKey), ":8080")
	svr.Start()
	defer svr.Stop()

	bh := handler.New()
	opts := &surf.Options{
		JWToken:    jwt,
		RouteAddrs: []string{"localhost:8080"},
		UnhandleFunc: func(s *server.TcpClient, m *server.MsgWraper) {
			log.Infof("recv msg: %v", m)
		},
	}
	core := surf.New(opts)

	surf.RegisterRequestHandle(core, "ReqStartBattle", bh.OnReqStartBattle)
	surf.RegisterRequestHandle(core, "ReqJoinBattle", bh.OnReqJoinBattle)
	surf.RegisterAysncHandle(core, "BattleMessageWrap", bh.OnBattleMessageWrap)

	core.Start()
	defer core.Stop()

	// cli := conn2(pk)
	// cli.Connect()
	// defer cli.Close()

	s := utilSignal.WaitShutdown()
	log.Infof("recv signal: %v", s.String())
	return nil
}

func server1(pk *rsa.PublicKey, listenAt string) *server.TcpServer {
	var err error
	h, err := route.NewRouter()
	if err != nil {
		panic(err)
	}
	svropt := &server.TcpServerOptions{
		AuthPublicKey:    pk,
		ListenAddr:       listenAt,
		OnSessionMessage: h.OnSessionMessage,
		OnSessionStatus:  h.OnSessionStatus,
	}
	svr, err := server.NewTcpServer(svropt)
	if err != nil {
		panic(err)
	}
	return svr
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

func conn2(pk *rsa.PrivateKey) *server.TcpClient {
	myuid := uint32(1002)

	jwt, _ := auth.GenerateToken(pk, &auth.UserInfo{
		UId:   myuid,
		UName: "user2",
		URole: "User",
	}, 24*time.Hour)

	opts := &server.TcpClientOptions{
		RemoteAddress:        "localhost:8080",
		AuthToken:            jwt,
		ReconnectDelaySecond: 10,

		OnMessage: func(s *server.TcpClient, m *server.MsgWraper) {
		},
		OnStatus: func(s *server.TcpClient, enable bool) {
			log.Debug("OnStatus: ", enable)
			if !enable {
				return
			}
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
						Uid:       myuid,
						SeatId:    1,
						MainScore: 1000,
					},
					{
						Uid:       2,
						SeatId:    2,
						Role:      1,
						MainScore: 1000,
					},
				},
			}
			s.SendReqMsg(battleuid, req, server.NewTcpRespCallbackFunc(func(c server.Session, resp *msg.RespStartBattle, err error) {
				if err != nil {
					log.Errorf("send req failed: %v", err)
					return
				}
				log.Infof("recv resp: %v", resp)

				if resp.BattleId != 0 {
					req := &msg.ReqJoinBattle{
						BattleId: resp.BattleId,
					}
					s.SendReqMsg(battleuid, req, server.NewTcpRespCallbackFunc(func(c server.Session, resp *msg.RespJoinBattle, err error) {
						if err != nil {
							log.Errorf("send req failed: %v", err)
							return
						}
						log.Infof("recv resp: %v", resp)
					}))
				}
			}))
		},
	}

	client := server.NewTcpClient(opts)

	return client
}
