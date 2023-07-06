/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"go-nuva/internal/log"
	"go-nuva/internal/log/types"

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
}

func preRunE(_ *cobra.Command, _ []string) error {
	//log.InitLogger("debug", "./logs/go-nuva.log")
	return nil
}

func server(_ *cobra.Command, _ []string) error {
	defer log.Stop()
	log.Init(types.DebugLevel, "./logs/go-nuva.log") //todo: 手动指定level / 日志路径，改成配置文件读取
	                                                 //      zap log rotate目前写死一个小时一个文件，后面改成使用lumberjack包
	log.Debug("root server running")
	return nil
}
