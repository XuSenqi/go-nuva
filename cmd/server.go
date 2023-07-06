/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"go-nuva/internal/config"
	"go-nuva/internal/log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:     "server",
	Short:   "",
	Long:    ``,
	PreRunE: preRunE,
	RunE:    runServer,
}

func init() {
	serverCmd.Flags().StringP("config", "c", "deploy/conf/config.yaml", "config file path")
}

func preRunE(cmd *cobra.Command, _ []string) error {
	cfg, err := cmd.Flags().GetString("config")
	if err != nil {
		return fmt.Errorf("config path error,%s", err.Error())
	}
	if err = config.InitConfig(cfg); err != nil {
		return fmt.Errorf("init config error,%s", err.Error())
	}
	return nil
}

func runServer(_ *cobra.Command, _ []string) error {
	defer log.Stop()

	lvl, err := log.ParseLevel(config.GetLogLevel()) // zap log rotate目前写死一个小时一个文件，后面改成使用lumberjack包
	if err != nil {
		return err
	}
	log.Init(lvl, config.GetLogPath())

	//todo: run http server; and handlers
	startHttpServer()
	log.Debug("root server running")
	return nil
}

// http server demo
func startHttpServer() error {
	// handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
		fmt.Fprintln(w, "hello")
	})

	srv := http.Server{
		Addr:    config.GetListenIp() + ":" + strconv.Itoa(int(config.GetListenPort())),
		Handler: handler,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := srv.Shutdown(context.Background()); err != nil {
			log.Infof("HTTP server Shutdown: %v", err)
		}
		log.Info("server gracefully shutdown")
		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Errorf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
	return nil
}
