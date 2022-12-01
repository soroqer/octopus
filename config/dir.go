package config

import (
	"os"
	"path"
	"path/filepath"
	"strings"
)

//说明：框架不支持存放及导出普通文件（execl、图片等），建议使用单独的文件服务器

//在文件路径策略中，主要需要制定的文件有三类：一类是配置文件，一类是数据文件，一类是日志文件
//我们并没有采用类似系统级的配置，不同类型的文件在不同的地方
//我们默认所有类型的文件都在同一个目录下，并且已定好默认的文件结构
//只需要设定基础路径即可，默认的基础路径为可执行文件所在目录
//同时为了增加灵活性，我们允许对不同类型的文件根目录进行单独的设定，
//但是同一类型的文件必须在同一根目录下
//可执行文件可以在任意路径下，但是如果可执行文件的路径不是基础路径的话，需要手动指定基础路径

//基础路径
var baseDir string

//除非有特殊需求，否则不建议对目录结构进行修改，因为需要考虑到服务部署的兼容性问题
//对于特殊的文件例如证书类的文件没有做默认设置，但是建议当做配置文件放在confDir下
//配置文件路径
var confDir string

//数据文件路径
var dataDir string

//日志文件路径
var logDir string

//三类文件的默认目录（相对于基础路径）
const (
	ConfDir = "./config"
	DataDir = "./data"
	LogDir  = "./log"
)

//主要配置文件的默认名称ConfigFile
const (
	ConfigFile = "config.toml"
)

//如果baseDir为空则取可执行文件所在的路径
//或者是相对路径（以.开始），则获取当前可执行文件的路径拼接后返回，
//如果是绝对路径则直接返回
func GetBaseDir() string {
	if baseDir == "" {
		baseDir = GetCurrentDirectory()
	} else if strings.HasPrefix(baseDir, ".") {
		baseDir = path.Join(GetCurrentDirectory(), baseDir)
	}
	return baseDir
}

//设置基础路径
func SetBaseDir(path string) {
	baseDir = path
}

//如果dataDir为空则取默认路径
//或者是相对路径（以.开始），则和基础路径拼接后返回，
//如果是绝对路径则直接返回
func GetDataDir() string {
	if dataDir == "" {
		dataDir = path.Join(GetBaseDir(), DataDir)
	} else if strings.HasPrefix(dataDir, ".") {
		dataDir = path.Join(GetBaseDir(), dataDir)
	}
	return dataDir
}

//设置配置路径
func SetDataDir(path string) {
	dataDir = path
}

//如果logDir为空则取默认路径
//或者是相对路径（以.开始），则和基础路径拼接后返回，
//如果是绝对路径则直接返回
func GetLogDir() string {
	if logDir == "" {
		logDir = path.Join(GetBaseDir(), LogDir)
	} else if strings.HasPrefix(logDir, ".") {
		logDir = path.Join(GetBaseDir(), logDir)
	}
	return logDir
}

//设置配置路径
func SetLogDir(path string) {
	logDir = path
}

//如果confDir为空则取默认路径
//或者是相对路径（以.开始），则和基础路径拼接后返回，
//如果是绝对路径则直接返回
func GetConfDir() string {
	if confDir == "" {
		confDir = path.Join(GetBaseDir(), ConfDir)
	} else if strings.HasPrefix(confDir, ".") {
		confDir = path.Join(GetBaseDir(), confDir)
	}
	return confDir
}

//获取配置文件
func GetConfFile() string {
	return path.Join(GetConfDir(), ConfigFile)
}

//设置配置路径
func SetConfDir(path string) {
	confDir = path
}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}
	return strings.Replace(dir, "\\", "/", -1)
}
