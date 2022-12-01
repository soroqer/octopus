package core

import (
	"auto-swap/config"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
	"math/big"
	"sync"
)

type NodeParams struct {
	MinGas   *big.Int
	NeedUSDT *big.Int
	NeedHFT *big.Int
	USDTApproveAddr common.Address
	SATApproveAddr common.Address
	ApproveUSDT  *big.Int
	ApproveSAT  *big.Int
	MinSwapUsdt     *big.Int
}

var nodeParams NodeParams

type Core struct {
	LvDb
	Client
	Galaxy
	Master
	WsServer
	Command
	sync.WaitGroup
	chSendReq chan *SendReq
	exit chan struct{}
}

func InitCore(pwd string) (core *Core,err error) {

	g, _  := new(big.Float).SetString(config.Cfg.Node.ReservedGas)
	nodeParams.MinGas, _ = g.Int(nil)
	ru, _  := new(big.Float).SetString(config.Cfg.Node.ReservedUSDT)
	nodeParams.NeedUSDT, _  = ru.Int(nil)
	rb, _  := new(big.Float).SetString(config.Cfg.Node.ReservedBFT)
	nodeParams.NeedHFT, _  = rb.Int(nil)
	au, _  := new(big.Float).SetString(config.Cfg.Node.ApprovedUSDT)
	nodeParams.ApproveUSDT, _  = au.Int(nil)
	nodeParams.USDTApproveAddr = common.HexToAddress(config.Cfg.Contract.RouterAddr)
	nodeParams.SATApproveAddr = common.HexToAddress(config.Cfg.Contract.RouterAddr)
	as, _  := new(big.Float).SetString(config.Cfg.Node.ApprovedBFT)
	nodeParams.ApproveSAT, _ = as.Int(nil)
	mu, _  := new(big.Float).SetString(config.Cfg.Node.MinUsdt)
	nodeParams.MinSwapUsdt, _ = mu.Int(nil)

	core = &Core{
		exit: make(chan struct{}),
	}
	err = core.LvDb.Init("node")
	if err != nil {
		logrus.Info("lvdb init error.")
		return
	}
	err = core.Client.Init(core.exit)
	if err != nil {
		logrus.Info("Client init error.")
		return
	}
	err = core.Galaxy.Init()
	if err != nil {
		logrus.Info("Galaxy init error.")
		return
	}
	err = core.Master.Init(pwd)
	if err != nil {
		logrus.Info("Master init error.")
		return
	}
	err = core.WsServer.Init()
	if err != nil {
		logrus.Info("WsServer init error.")
		return
	}
	core.Command.Init(core)
	return
}

func (core *Core)StartServer() error {
	err := core.StartWs(core)
	if err != nil {
		logrus.Info("StartWs error.")
		return err
	}
	maxPath2,err := core.GetMaxPath2()
	if err != nil {
		logrus.Info("GetMaxPath2 error.")
		return err
	}
	chNode := core.IterNode()
	err = core.Galaxy.Load(maxPath2,&core.Master,chNode)
	if err != nil {
		logrus.Info("Load Galaxy error.")
		return err
	}
	err = core.Client.StartSubscribe(&core.Galaxy, &core.WsServer)
	if err != nil {
		logrus.Info("StartSubscribe error.")
		return err
	}
	core.chSendReq = core.Master.ListenSendReq(&core.Client)
	core.Run()
	return nil
}

func (core *Core)StopServer() {
	core.Client.StopSubscribe()
	err := core.LvDb.Close()
	if err != nil {
		logrus.Error("core lvdb close err: ",err)
	}
}

func (core *Core)Run() {
	if core.chSendReq == nil {
		panic("run chSendReq is nil")
	}
	core.Range(func(node *Node) (next bool) {
		core.Add(1)
		go func() {
			defer core.Done()
			node.StartTask(core.exit,&core.LvDb,&core.WsServer,&core.Client,core.chSendReq)
		}()
		return true
	})
}

func (core *Core)RunNodeTask(node *Node) {
	core.Add(1)
	go func() {
		defer core.Done()
		node.StartTask(core.exit,&core.LvDb,&core.WsServer,&core.Client,core.chSendReq)
	}()
}



