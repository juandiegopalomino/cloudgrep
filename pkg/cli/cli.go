package cli

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/run-x/cloudgrep/pkg/api"
	"github.com/run-x/cloudgrep/pkg/config"
	"github.com/run-x/cloudgrep/pkg/datastore"
	"github.com/run-x/cloudgrep/pkg/engine"
	"github.com/run-x/cloudgrep/pkg/util"
)

type cli struct {
	cfg    config.Config
	logger *zap.Logger
	ds     datastore.Datastore
}

func Run(ctx context.Context, cfg config.Config, logger *zap.Logger) error {
	if logger.Core().Enabled(zap.DebugLevel) {
		util.StartProfiler()
	}

	//send amplitude event
	util.SendEvent(ctx, logger, util.BaseEvent, nil)

	//init the storage to contain cloud data
	var err error
	cli := cli{cfg: cfg, logger: logger}
	cli.ds, err = datastore.NewDatastore(ctx, cfg, logger)
	if err != nil {
		return fmt.Errorf("failed to setup datastore: %w", err)
	}

	//start the providers to collect cloud data
	if !cfg.Datastore.SkipRefresh {
		if err := cli.runEngine(ctx); err != nil {
			return err
		}
	}
	api.StartWebServer(ctx, cfg, logger, cli.ds, cli.runEngine)

	url := fmt.Sprintf("http://%v:%v/%v", cfg.Web.Host, cfg.Web.Port, cfg.Web.Prefix)
	fmt.Println("To view Cloudgrep UI, open ", url, "in browser")

	if !cfg.Web.SkipOpen {
		openPage(url)
	}
	handleSignals()
	return nil
}

//runEngine runs the providers to collect the cloud resources
//it returns when it's done fetching
func (cli *cli) runEngine(ctx context.Context) error {
	//TODO send engine start event
	eng, err := engine.NewEngine(ctx, cli.cfg, cli.logger, cli.ds)
	if err == nil {
		if err = eng.Run(ctx); err != nil {
			stats, _ := cli.ds.Stats(ctx)
			if stats.ResourcesCount > 0 {
				//log the error but the api can still server with the datastore
				cli.logger.Sugar().Errorw("some error(s) when running the provider engine", "error", err)
			} else {
				// nothing to view - exit
				return fmt.Errorf("can't run the provider engine: %w", err)
			}
		}
	}
	//TODO send engine end event using err as param
	return nil
}

func handleSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}

func openPage(url string) {
	_, err := exec.Command("which", "open").Output()
	if err != nil {
		return
	}

	err = exec.Command("open", url).Run()
	if err != nil {
		return
	}

}
