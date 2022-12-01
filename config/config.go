package config

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"path"
	"time"
)

var Cfg *Config

type Config struct {
	Server      *Server      `mapstructure:"server"`
	Node   		*Node   	 `mapstructure:"node"`
	Contract    *Contract    `mapstructure:"contract"`
	LogRotate   *LogRotate   `mapstructure:"log_rotate"`
}

type Server struct {
	IP   string `mapstructure:"ip"`
	Port string `mapstructure:"port"`
}

type Node struct {
	ChainId   		int64   `mapstructure:"chain_id"`
	Host  	  		string  `mapstructure:"host"`
	Key             string  `mapstructure:"key"`
	KeyAddress      string  `mapstructure:"key_address"`
	GasPrice        int64   `mapstructure:"gas_price"`
	StepPrice       int64   `mapstructure:"step_price"`
	MaxPrice        int64   `mapstructure:"max_price"`
	ReservedGas     string  `mapstructure:"reserved_gas"`
	ReservedUSDT    string  `mapstructure:"reserved_usdt"`
	ApprovedUSDT    string  `mapstructure:"approved_usdt"`
	ReservedBFT     string  `mapstructure:"reserved_bft"`
	ApprovedBFT     string  `mapstructure:"approved_bft"`
	RateLimit       int     `mapstructure:"rate_limit"`
	IsReward        bool    `mapstructure:"is_reward"`
	Slippage        int64   `mapstructure:"slippage"`
	MinUsdt         string  `mapstructure:"min_usdt"`
}

type Contract struct {
	USDTAddr       string  `mapstructure:"usdt_addr"`
	SATAddr        string  `mapstructure:"sat_addr"`
	InvitationAddr string  `mapstructure:"invitation_addr"`
	RouterAddr     string  `mapstructure:"router_addr"`
	PairAddr       string  `mapstructure:"pair_addr"`
}

type Sun struct {
	Address   string  `mapstructure:"address"`
	Path      int     `mapstructure:"path"`
	Number    int     `mapstructure:"number"`
}

type LogRotate struct {
	MaxAge       time.Duration `mapstructure:"max_age"`
	RotationTime time.Duration `mapstructure:"rotation_time"`
}

func Load() error {
	Cfg = &Config{}
	err := viper.Unmarshal(Cfg, func(mcf *mapstructure.DecoderConfig) {
		mcf.ZeroFields = true
		mcf.ErrorUnused = true
	})
	return err
}

func TestInit() {
	SetBaseDir(path.Join(MainPkgPath(), "build"))
	viper.SetConfigFile(GetConfFile())
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("load config file err:", err)
		os.Exit(-1)
	}
	Load()
	logrus.SetFormatter(&logrus.TextFormatter{DisableQuote: true})
}

