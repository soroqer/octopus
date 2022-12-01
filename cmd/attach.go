package cmd

import (
	"auto-swap/config"
	"auto-swap/lib"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(attachCmd)
}

//指令的用法，提示以及执行的函数
//这里才是服务真正的入口
//startServer相当于main函数入口
var attachCmd = &cobra.Command{
	Use:   "attach",
	Short: "attach to the system .",
	RunE:  attachServer,
}

//服务入口，逻辑代码从这里开始
func attachServer(_ *cobra.Command, _ []string) error {
	lib.Attach(config.AttachPort)
	return nil
}
