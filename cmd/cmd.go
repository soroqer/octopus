package cmd

import (
	"auto-swap/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

//该文件是服务的入口，通过go的init机制在文件加载的时候把指令加入根指令中，
//通过不同的指令进入不同的服务入口开始运行服务。
//其中version指令没有在这个文件中加载，而是隐藏到了version.go文件中的init函数中
//因为该指令是所有服务都需要有的。
//该文件所在的文件夹下每一个go文件都代表了一个指令，是一个服务的入口，相当于main函数

var rootCmd = &cobra.Command{
	Use:   config.ServerName,
	Short: config.ServerDescription,
}

//加载文件的时候把指令加入根指令中
func init() {
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
}

// GetRootCmd 返回给main函数使用
func GetRootCmd() *cobra.Command {
	return rootCmd
}

// waitForExit 系统退出信号监听
func waitForExit(exit chan struct{}) {

	// 捕捉指定的信号
	quit := make(chan os.Signal)
	// 前台时，按 ^C 时触发
	signal.Notify(quit, syscall.SIGINT)
	// 后台时，kill 时触发。kill -9 时的信号 SIGKILL 不能捕捉，所以不用添加
	signal.Notify(quit, syscall.SIGTERM)
	go func() {
		// 等待退出信号
		sig := <-quit
		logrus.Printf("received signal: %v, start to shutdown server.\n", sig)
		close(exit)
	}()
	//其他地方关闭exit同样可以触发服务结束
	<- exit
}

// Will return true if the process with PID exists.
func checkProcess(pid int) bool {
	process, _ := os.FindProcess(pid)
	err := process.Signal(syscall.Signal(0))
	if err != nil {
		return false
	} else {
		return true
	}
}
