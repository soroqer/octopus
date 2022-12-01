package main

import (
	"auto-swap/cmd"
	"auto-swap/config"
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path"
	"time"
)

//程序的入口，但是并不是逻辑的入口，使用了Cobra把程序改为指令启动的形式，可以执行多个逻辑入口
//通常这个main基本都是一样的，不需要修改，除非有特殊的全局传参
//Cobra is a library providing a simple interface to create powerful modern CLI interfaces similar to git & go tools.
func main() {
	// OnInitialize sets the passed functions to be run when each command's
	// 这里用来系统启动时初始化各种状态，这里只是加载，会在execute之后执行。
	// 传入的方法执行的顺序在指令取参数之后，所以可以根据指令传入的参数对系统做不同的初始化。
	cobra.OnInitialize(initConfig)
	//加载指令，这里是总入口，不同的指令会在cmd.go中加载进来
	command := cmd.GetRootCmd()
	//全局的参数，执行指令时可以通过标记传入参数
	//下面代码可以通过 -c ./fack/config.toml 指令把路径赋给configFile字段
	command.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file (default is ./config/config.toml)")
	//下面代码可以通过 -logdir ./fack 指令把路径赋给logDir字段
	command.PersistentFlags().StringVarP(&logDir, "logdir", "", "", "log file path (default is ./logrus)")
	//下面代码可以通过 -basedir ./fack 指令把路径赋给baseDir字段
	command.PersistentFlags().StringVarP(&baseDir, "basedir", "b", "", "base project path (default is .)")
	//下面代码可以通过 -datadir ./fack 指令把路径赋给dataDir字段
	command.PersistentFlags().StringVarP(&dataDir, "datadir", "", "", "data file path (default is ./data)")
	//下面代码可以通过 -w false 指令设定日志不落盘
	command.PersistentFlags().BoolVarP(&writeLog, "writelog", "w", true, "is needed write logrus to the disk (default is true)")
	//下面代码可以通过 -d true 指令设定后台运行
	command.PersistentFlags().BoolVarP(&cmd.Daemon, "daemon", "d", false, "is running in daemon (default is true)")
	//下面代码可以通过 --child 指令表示是子进程
	command.PersistentFlags().BoolVarP(&cmd.Child, "child", "", false, "is running as child (default is false)")
	//执行指令
	if err := command.Execute(); err != nil {
		logrus.Error("Start server failed, err: ", err)
		os.Exit(1)
	}
}

//该处定义的变量会在执行指令时通过标记传参赋值（可以没有，则为空）
var baseDir string
var configFile string
var logDir string
var dataDir string
var writeLog bool

// 初始化配置文件
func initConfig() {
	//设置路径
	config.SetBaseDir(baseDir)
	if dataDir != "" {
		config.SetDataDir(dataDir)
	}
	if logDir != "" {
		config.SetLogDir(logDir)
	}
	if configFile != "" {
		config.SetLogDir(logDir)
	}

	//指定viper需要加载的配置文件的路径
	viper.SetConfigFile(config.GetConfFile())
	//打开后viper会自动捕捉环境变量的kv
	viper.AutomaticEnv()
	//读取文件并加载到内存中，这里并不进行反序列化解析，在需要的时候才会解析
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("load config file err:", err)
		os.Exit(-1)
	}

	if err := config.Load(); err != nil {
		fmt.Println("load config file err:", err)
		os.Exit(-1)
	}
	//默认开启日志落盘，不需要落盘的话使用-w false关闭
	if writeLog == true {
		err := ConfigLocalFilesystemLogger()
		if err != nil {
			fmt.Println("config log persistence err:", err)
			os.Exit(-1)
		}
	}
}

func ConfigLocalFilesystemLogger() error {

	conf := config.Cfg.LogRotate
	baseLogPaht := path.Join(config.GetLogDir(), config.ServerName)
	writer, err := rotatelogs.New(
		baseLogPaht+".%Y%m%d",
		rotatelogs.WithLinkName(baseLogPaht),                       // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(conf.MaxAge*time.Minute),             // 文件最大保存时间
		rotatelogs.WithRotationTime(conf.RotationTime*time.Minute), // 日志切割时间间隔
	)
	if err != nil {
		return err
	}
	errLogPath := path.Join(config.GetLogDir(), config.ErrLogName)
	errWriter, err := rotatelogs.New(
		errLogPath+".%Y%m%d",
		rotatelogs.WithLinkName(errLogPath),                        // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(conf.MaxAge*time.Minute),             // 文件最大保存时间
		rotatelogs.WithRotationTime(conf.RotationTime*time.Minute), // 日志切割时间间隔
	)
	if err != nil {
		return err
	}
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer, // 为不同级别设置不同的输出目的
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: errWriter,
		logrus.FatalLevel: errWriter,
		logrus.PanicLevel: errWriter,
	}, &logrus.TextFormatter{DisableQuote: true})
	logrus.AddHook(lfHook)
	return nil
}
