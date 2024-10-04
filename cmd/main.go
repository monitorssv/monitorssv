package main

import (
	"context"
	"errors"
	"fmt"
	logging "github.com/ipfs/go-log/v2"
	"github.com/monitorssv/monitorssv/alert"
	"github.com/monitorssv/monitorssv/config"
	"github.com/monitorssv/monitorssv/eth1/client"
	"github.com/monitorssv/monitorssv/eth1/ssv"
	"github.com/monitorssv/monitorssv/eth2"
	client2 "github.com/monitorssv/monitorssv/eth2/client"
	"github.com/monitorssv/monitorssv/service"
	"github.com/monitorssv/monitorssv/store"
	"github.com/urfave/cli/v2"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var log = logging.Logger("monitor-ssv")

func main() {
	_ = logging.SetLogLevel("*", "INFO")
	app := &cli.App{
		Name:                 "monitorssv",
		Usage:                "monitor ssv",
		Version:              "v0.1",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			importCmd,
			runCmd,
			fixDBOperatorValidatorCountCmd,
			fixDBClusterFeeAddressCmd,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Warnf("%+v", err)
		os.Exit(1)
		return
	}
}

var runCmd = &cli.Command{
	Name:  "run",
	Usage: "run ssv monitor",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "conf-path",
			Usage: "config.yaml path",
			Value: "",
		},
		&cli.Uint64Flag{
			Name:  "monitor-api",
			Usage: "monitor api port",
			Value: 8890,
		},
	},
	Action: func(ctx *cli.Context) error {
		cfg, err := config.InitConfig(ctx.String("conf-path"))
		if err != nil {
			log.Errorw("InitConfig", "err", err)
			return err
		}

		db, err := store.NewStore(cfg)
		if err != nil {
			log.Errorw("NewStore", "err", err)
			return err
		}

		eth1Client, err := client.NewEth1Client(cfg)
		if err != nil {
			log.Errorw("NewEth1Client", "err", err)
			return err
		}

		password := os.Getenv("ENCRYPTION_KEY")
		if password == "" {
			return errors.New("no encrypted password")
		}

		alarmDaemon, err := alert.NewAlarmDaemon(db, eth1Client, password)
		if err != nil {
			log.Errorw("NewAlarmDaemon", "err", err)
			return err
		}

		ssv, err := ssv.NewSSV(cfg, eth1Client, db, alarmDaemon)
		if err != nil {
			log.Errorw("NewEth1Client", "err", err)
			return err
		}

		eth2Client := client2.NewClient(cfg.Eth2Rpc)
		beaconMonitor, err := eth2.NewBeaconMonitor(cfg, eth2Client, db, alarmDaemon)
		if err != nil {
			log.Errorw("NewBeaconMonitor", "err", err)
			return err
		}

		monitorService, err := service.NewMonitorSSV(db, ssv, beaconMonitor, alarmDaemon, password)
		if err != nil {
			log.Errorw("NewMonitorSSV", "err", err)
			return err
		}

		router := monitorService.NewRouter()
		port := ctx.Uint64("monitor-api")
		s := &http.Server{
			Addr:         fmt.Sprintf("0.0.0.0:%d", port),
			Handler:      router,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		}
		log.Infow("start monitor server", "network", cfg.Network, "port", port)
		go func() {
			if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("s.ListenAndServe err: %v", err)
			}
		}()

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		monitorService.Stop()
		ctxT, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.Shutdown(ctxT); err != nil {
			log.Fatal("server forced to shutdown:", err)
		}

		log.Info("monitor server exit")

		return nil
	},
}
