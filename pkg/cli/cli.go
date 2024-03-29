package cli

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/hashicorp/go-multierror"
	"github.com/juandiegopalomino/cloudgrep/pkg/model"
	"github.com/juandiegopalomino/cloudgrep/pkg/util/amplitude"

	"go.uber.org/zap"

	"github.com/juandiegopalomino/cloudgrep/pkg/api"
	"github.com/juandiegopalomino/cloudgrep/pkg/config"
	"github.com/juandiegopalomino/cloudgrep/pkg/datastore"
	"github.com/juandiegopalomino/cloudgrep/pkg/engine"
	"github.com/juandiegopalomino/cloudgrep/pkg/util"
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
	amplitude.SendEvent(logger, amplitude.EventLoad, nil)

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
	url = strings.Trim(url, "/")
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
	if err := cli.ds.WriteEvent(ctx, model.NewEngineEventStart()); err != nil {
		return err
	}
	eng, errors := engine.NewEngine(ctx, cli.cfg, cli.logger, cli.ds)
	if errors == nil {
		err := eng.Run(ctx)
		if err != nil {
			stats, _ := cli.ds.Stats(ctx)
			if stats.ResourcesCount > 0 {
				//log the error but the api can still server with the datastore
				cli.logger.Sugar().Errorw("some error(s) when running the provider engine", "error", err)
			} else {
				// nothing to view - exit
				errors = multierror.Append(errors, fmt.Errorf("can't run the provider engine: %w", err))
			}
		}
	}
	if err := cli.ds.WriteEvent(ctx, model.NewEngineEventEnd(errors)); err != nil {
		errors = multierror.Append(errors, err)
	}
	return errors
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
