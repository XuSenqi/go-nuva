/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"go-nuva/internal/config"
	"go-nuva/internal/log"

	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:     "server",
	Short:   "",
	Long:    ``,
	PreRunE: preRunE,
	RunE:    server,
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

func server(_ *cobra.Command, _ []string) error {
	defer log.Stop()

	lvl, err := log.ParseLevel(config.GetLogLevel()) // zap log rotate目前写死一个小时一个文件，后面改成使用lumberjack包
	if err != nil {
		return err
	}
	log.Init(lvl, config.GetLogPath())

	//todo: run http server; and handlers
	log.Debug("root server running")
	return nil
}
