/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"go-nuva/internal/log"
	"go-nuva/internal/log/types"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-nuva",
	Short: "A framework for launching http server quickly",
	Long:  ``,
}

func Execute() {
	defer log.Stop()
	log.Init(types.DebugLevel, "./logs/go-nuva.log") //todo: 手动指定level / 日志路径，改成配置文件读取
	                                                 //      zap log rotate目前写死一个小时一个文件，后面改成使用lumberjack包

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
