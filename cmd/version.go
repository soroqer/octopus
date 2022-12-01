package cmd

import (
	"auto-swap/config"
	"fmt"
	"github.com/spf13/cobra"
)

//这个文件的代码基本不用变，只需要根据具体的服务修改对应的服务信息即可


//在文件加载的时候把指令加进根指令中
//因为版本指令所有服务都需要有，所以在这里隐形的加载到了根指令中
func init() {
	rootCmd.AddCommand(versionCmd)
}

//指令的用法，提示以及执行的函数
//outputVersion相当于main函数入口
var versionCmd = &cobra.Command{
	Use: "version",
	Short: "version of "+ config.ServerName +" service",
	Run: outputVersion,
}

//输出版本信息
func outputVersion(cmd *cobra.Command, args []string) {
	fmt.Println(config.ServerName, config.ServerVersion)
}
