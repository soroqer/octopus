package cmd

import (
	"auto-swap/config"
	"auto-swap/core"
	"auto-swap/lib"
	"errors"
	"fmt"
	"github.com/gogf/gf/os/gproc"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"time"
)

const PIDFILE = config.ServerName + ".lock"

var (
	Child  bool
	Daemon bool
)

//指令的用法，提示以及执行的函数
//这里才是服务真正的入口
//startServer相当于main函数入口
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start " + config.ServerName + " service",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !Child {
			strb, _ := ioutil.ReadFile(PIDFILE)
			if strb != nil {
				pid, err := strconv.Atoi(string(strb))
				if err == nil {
					if checkProcess(pid) {
						return errors.New(config.ServerName + " service already exist")
					}
				}

			}
		}
		return startServer(cmd, args)
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop" + config.ServerName + " service",
	Run: func(cmd *cobra.Command, args []string) {
		b, err := ioutil.ReadFile(path.Join(config.GetBaseDir(), PIDFILE))
		if err != nil {
			fmt.Printf("stop %s failed, error: %v\n", os.Args[0], err)
			return
		}
		strb := string(b)
		pid,err := strconv.ParseInt(strb,10,64)
		if err != nil {
			fmt.Printf("stop %s failed, error: %v\n", os.Args[0], err)
			return
		}
		command := exec.Command("kill", strb)
		if err := command.Start(); err != nil {
			fmt.Printf("stop %s failed, error: %v\n", os.Args[0], err)
		} else {
			println(config.ServerName + " stop...")
			exit := make(chan struct{})
			ticker := time.NewTicker(100*time.Millisecond)
			go func() {
				for range ticker.C {
					if checkProcess(int(pid)) {
						continue
					}else{
						println(config.ServerName + " stop success")
						close(exit)
					}
				}
			}()
			waitForExit(exit)
			_ = os.Remove(path.Join(config.GetBaseDir(), PIDFILE))
		}
	},
}


func startServer(_ *cobra.Command, args []string) error {
	if Daemon {
		name := os.Args[0]
		command := exec.Command(name, "start", "--child")
		if err := command.Start(); err != nil {
			fmt.Printf("start %s failed, error: %v\n", os.Args[0], err)
			return err
		}

		password, err := core.GetPassword()
		if err != nil {
			return err
		}
		pwd,err := core.AESCBCEncrypt([]byte(password),core.FormatAESKey(""))
		if err != nil {
			return err
		}
		err = gproc.Send(command.Process.Pid, pwd)
		if err != nil {
			logrus.Info("send err:", err)
			return err
		}

		fmt.Printf("%s start, [PID] %d running...\n", config.ServerName, command.Process.Pid)
		err = ioutil.WriteFile(path.Join(config.GetBaseDir(), PIDFILE), []byte(strconv.FormatInt(int64(command.Process.Pid),10)), 0666)
		if err != nil {
			return err
		}
		fmt.Println("service started in the background.")
	} else {
		var password string
		if Child {
			msg := gproc.Receive()
			pwd := msg.Data
			b,err := core.AESCBCDecrypt(pwd,core.FormatAESKey(""))
			if err != nil {
				return err
			}
			password = string(b)
		}else{
			var err error
			password, err = core.GetPassword()
			if err != nil {
				return err
			}
		}
		c,err := core.InitCore(password)
		if err != nil {
			return err
		}
		err = c.StartServer()
		if err != nil {
			return err
		}
		defer c.StopServer()

		exit := make(chan struct{})
		_,err = lib.ListenCommand(config.AttachPort,&c.Command,exit)
		if err != nil {
			return err
		}
		waitForExit(exit)
	}
	return nil
}


