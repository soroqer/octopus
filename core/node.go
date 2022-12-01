package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
	"math/big"
	"strconv"
	"strings"
	"time"
)

type TxState int
const (
	DefaultTS TxState = iota
	Pending
	Success
	Fail
	FailWaiting  //失败了但是不重试，等待重试命令
)

func (txState TxState) String() string{
	switch txState {
	case DefaultTS    : return "等待生成交易"
	case Pending      : return "等待打包"
	case Success      : return "交易成功"
	case Fail         : return "交易失败"
	case FailWaiting  : return "交易失败，等待新指令"
	default           : return "未知状态"
	}
}

//StepType 注意，这里的 StepType 在上线执行后不能随改，否则会导致数据和流程错乱
type StepType  int
const (
	DefaultST    StepType = iota // 0
	ApproveUSDT           // 1  授权USDT给 Swap Router
	ApproveHFT            // 2  授权HFT给 Swap Routerf
	Bind                  // 3  绑定关系
	Standby               // 4  待命状态，随时可以参与买卖
	BuyHFT                // 5  买HFT
	BuyHFTWithSpeedUp     // 6  买HFT
	SellHFT               // 7  卖HFT
	SellHFTWithSpeedUp    // 8  卖HFT
	CollectUSDT           // 9  归集USDT
	CollectBFT            // 10  归集BFT
	CollectBNB            // 11  归集BNB
	TransferGas           // 12  转gas费用，每次转移的gas数量通过配置文件取
	TransferUSDT          // 13 转usdt，每次转移的usdt数量通过配置文件取
	TransferHFT           // 14 转hft，每次转移的hft数量通过配置文件取
)

func (stepType StepType) String() string{
	switch stepType {
	case DefaultST    : return "等待初始化"
	case ApproveUSDT  : return "向 Router 授权 USDT"
	case ApproveHFT   : return "向 Router 授权 BFT"
	case Bind         : return "绑定关系"
	case Standby      : return "待命状态"
	case BuyHFT       : return "购买 BFT"
	case BuyHFTWithSpeedUp  : return "购买 BFT，带加速"
	case SellHFT      : return "售卖 BFT"
	case SellHFTWithSpeedUp : return "售卖 BFT，带加速"
	case CollectUSDT  : return "归集 USDT"
	case CollectBFT   : return "归集 BFT"
	case CollectBNB   : return "归集 BNB"
	case TransferGas  : return "等待主账户转移 gas 费用"
	case TransferUSDT : return "等待主账户转移 USDT 费用"
	case TransferHFT  : return "等待主账户转移 BFT 费用"
	default           : return "未知状态"
	}
}


//Command
const (
	CmdInit            = "init"         // 初始化节点命令，收到改命令后会进行HFT和USDT的授权交易
	CmdBuy             = "buy"          // 买HFT
	CmdBuyWithSpeedUp  = "buy_speedup"  // 买HFT，会在每个新的区块出现后会提高gas费进行加速交易
	CmdSell            = "sell"         // 卖HFT
	CmdSellWithSpeedUp = "sell_speedup" // 卖HFT，会在每个新的区块出现后会提高gas费进行加速交易
	CmdCollectUsdt     = "collect_usdt" // 归集 USDT
	CmdCollectBft      = "collect_bft"  // 归集 BFT
	CmdCollectBnb      = "collect_bnb"  // 归集 BNB
	CmdRetry           = "retry"        // 重新执行当前指令（在失败的状态下）
	CmdBalance         = "balance"      // 刷新余额
	CmdReset           = "reset"        // 重置节点
)

type Step struct {
	StepType
	TxState
	Coin
	ReqId uint64
	Number *big.Int
	OprAddr common.Address
	Expected *big.Int //预期的余额，达到了表示成功
	TxHash common.Hash
	Retry int //重试次数
	HisHashes []common.Hash //如果有重试的话，记录的历史哈希
}

type Path string

func NewPath(ps ...int) Path {
	path := PathPre
	for _,p := range ps {
		path += PathSeparator + strconv.Itoa(p)
	}
	return Path(path)
}

const PathPre = "m"
const PathSeparator = "-"

func (path Path) ParsePath() ([]int,error) {
	strs := strings.Split(string(path),PathSeparator)
	l := len(strs)
	if l < 2 {
		return nil,errors.New(fmt.Sprintf("path format error: %v", path))
	}
	paths := make([]int,l-1,l-1)
	var err error
	for i := range paths {
		paths[i],err = strconv.Atoi(strs[i+1])
		if err != nil {
			return nil,errors.New(fmt.Sprintf("child level %v format error: %v", i+1, path))
		}
	}
	return paths,nil
}

type Node struct {
	Active bool //是否处于活跃状态
	Path
	Addr common.Address
	StepType        // 当前的step状态
	ReqId  uint64   // 当前的请求ID
	Number *big.Int // 当前的请求操作数量
	OprAddr common.Address // 当前的请求操作地址
	PreStepType StepType  // 上一个step状态，触发gas请求完成后，会跳转回之前的状态
	Steps []*Step  // 所有的step，记录所有，方便查询
	Balance struct { //这个只是给显示用，不作为判断依据
		needUpdate bool
		HFT *big.Int
		USDT *big.Int
		BNB *big.Int
	}
	needCache bool
	sunAddr common.Address
	auth *bind.TransactOpts // 小写的不会被json序列化
	key string //todo 这个是临时用的，要注意一下
	cmd string
}

func (node *Node) GetKey() string{
	return node.key
}

func (node *Node) ExecuteCmd(reqId uint64,cmd string,number int64, oprAddr string) (info string) {
	node.Number = new(big.Int).Mul(big.NewInt(number),big.NewInt(1e18))
	node.OprAddr = common.HexToAddress(oprAddr)
	node.ReqId = reqId
	switch cmd {
	case CmdInit:
		if node.StepType == DefaultST {
			node.StepType = ApproveUSDT //开始初始化，执行授权指令
			node.PreStepType = DefaultST
		}else{
			return "当前状态不可执行初始化指令。"
		}
	case CmdBuy :
		if node.StepType == Standby {
			node.StepType = BuyHFT
			node.PreStepType = Standby
		}else{
			return "当前状态不可执行买卖指令。"
		}
	case CmdBuyWithSpeedUp  :
		if node.StepType == Standby {
			node.StepType = BuyHFTWithSpeedUp
			node.PreStepType = Standby
		}else{
			return "当前状态不可执行买卖指令。"
		}
	case CmdSell            :
		if node.StepType == Standby {
			node.StepType = SellHFT
			node.PreStepType = Standby
		}else{
			return "当前状态不可执行买卖指令。"
		}
	case CmdSellWithSpeedUp :
		if node.StepType == Standby {
			node.StepType = SellHFTWithSpeedUp
			node.PreStepType = Standby
		}else{
			return "当前状态不可执行买卖指令。"
		}
	case CmdCollectUsdt :
		if node.StepType == Standby {
			node.StepType = CollectUSDT
			node.PreStepType = Standby
		}else{
			return "当前状态不可执行归集指令。"
		}
	case CmdCollectBft :
		if node.StepType == Standby {
			node.StepType = CollectBFT
			node.PreStepType = Standby
		}else{
			return "当前状态不可执行归集指令。"
		}
	case CmdCollectBnb :
		if node.StepType == Standby {
			node.StepType = CollectBNB
			node.PreStepType = Standby
		}else{
			return "当前状态不可执行归集指令。"
		}
	case CmdRetry :
		step := node.GetStep()
		if step != nil && step.TxState == FailWaiting {
			step.TxState = DefaultTS
		}else{
			return "当前状态不可执行重试指令。"
		}
	case CmdBalance :
		node.Balance.needUpdate = true
	case CmdReset :
		node.ReSet()
	default :
		return "无效指令。"
	}
	return "命令发布成功，等待节点执行。"
}

func (node *Node) ReSet(){
	// 该重置方法有一定的概率导致初始化过的节点退回未初始化状态
	// 只需再初始化一次即可，并不影响流程
	if node.StepType == TransferGas {
		node.StepType = DefaultST
	}else if node.StepType > Standby  {
		node.StepType = Standby
	}else if node.StepType > DefaultST {
		node.StepType = DefaultST
	}
}

func (node *Node) ReSetToStandBy(){
	// TODO 最后关闭服务时临时使用，正常用上面那个
	node.StepType = Standby
}

func (node *Node) UpdateBalance(client *Client){
	ub,ok := client.GetUSDTBalance(node.Addr)
	if ok {
		node.Balance.USDT = ub
		node.needCache = true
	}
	hb,ok := client.GetHFTBalance(node.Addr)
	if ok {
		node.Balance.HFT = hb
		node.needCache = true
	}
	node.Balance.needUpdate = false
}

func (node *Node) StartTask(exit chan struct{},lvdb *LvDb,ws *WsServer, client *Client, chSendReq chan *SendReq) {
	if node.auth == nil {
		b,_ := json.Marshal(node)
		logrus.Error(string(b))
		panic("node auth 为空")
	}
	//node.UpdateBalance(client)
	ticker := time.NewTicker(100*time.Millisecond) //每100毫秒检查一次，是否有最新header
	counter := 0
	var tempHeader *types.Header
	defer ticker.Stop()
	node.Active = true
	defer func() {
		node.Active = false
	}()
	for {
		select {
		case <- exit:
			return
		case <- ticker.C:
			counter ++
			// 最多等待4秒就开始循环
			if counter > 40 {
				goto EXE
			}
			// 或者有新的区块,就开始循环
			if tempHeader != nil && client.Last.Header != nil &&
				client.Last.Header.Number.Cmp(tempHeader.Number) == 1 {
				goto EXE
			}
			continue
		EXE: tempHeader = client.Last.Header
			counter = 0 //重置计数器
			if node.Balance.needUpdate {
				node.UpdateBalance(client)
			}
			switch node.StepType {
			case DefaultST    : node.StandBy()      // 准备状态（默认状态）
			case ApproveUSDT  : node.USDTApprove(client) // 授权USDT给Router，buy的时候需要
			case ApproveHFT   : node.HFTApprove(client) // 授权HFT给Router，sell的时候需要
			case Bind         : node.Bind(client) // 绑定关系
			case Standby      : node.StandBy()      // 待命状态
			case BuyHFT       : node.BuyHFT(client)      // 卖出SAT
			case BuyHFTWithSpeedUp  : node.BuyHFTWithSpeedUp(client)      // 卖出SAT
			case SellHFT      : node.SellHFT(client)      // 卖出SAT
			case SellHFTWithSpeedUp : node.SellHFTWithSpeedUp(client)      // 卖出SAT
			case CollectUSDT  : node.CollectUsdt(client)
			case CollectBFT   : node.CollectBft(client)
			case CollectBNB   : node.CollectBnb(client)
			case TransferGas  : node.GasCheck(client, chSendReq)  // 转gas费用，每个gas费用通过配置文件取
			case TransferUSDT : node.USDTCheck(client, chSendReq)  // 转gas费用，每个gas费用通过配置文件取
			case TransferHFT  : node.HFTCheck(client, chSendReq)  // 转gas费用，每个gas费用通过配置文件取
			default           : return
			}
		}
		if node.needCache {
			lvdb.CacheNode(node)
			ws.SendMsg(client.Last.Header.Number.String(),node)
			node.needCache = false
		}
	}
}

func (node *Node)AddStep(step *Step)  {
	l := len(node.Steps)
	if l >= 10 {
		node.Steps = node.Steps[:l-1]
	}
	node.Steps = append(node.Steps,step)
}

func (node *Node)RemoveLastStep()  {
	l := len(node.Steps)
	if l > 0 {
		node.Steps = node.Steps[:l-1]
	}
}

func (node *Node)GetStep() *Step {
	// step 里会有重复的步骤
	// 只需取最后一个，和当前状态对比
	l := len(node.Steps)
	if l == 0 {
		return nil
	}
	step := node.Steps[l-1]
	if step.StepType != node.StepType ||
		step.ReqId != node.ReqId{
		return nil
	}
	return step
}

func (node *Node) SwitchToReq(stepType StepType)  {

	if stepType != TransferGas &&
		stepType != TransferUSDT &&
		stepType != TransferHFT {
		panic("SwitchToReq StepType err")
	}
	node.PreStepType = node.StepType
	node.RemoveLastStep()
	node.StepType = stepType
	node.needCache = true
}

func (node *Node) StandBy() {
	return //不需要执行任何内容
}

func (node *Node) USDTApprove(client *Client) {
	//logrus.Info("执行 USDTApprove")
	//defer logrus.Info("USDTApprove 执行完毕")
	//先检查步骤状态
	step := node.GetStep()
	if step == nil {
		step = &Step{
			StepType:  ApproveUSDT,
			TxState:   DefaultTS,
			Coin:      USDT,
			ReqId:     node.ReqId,
			Number:    node.Number,
			OprAddr:   node.OprAddr,
			TxHash:    common.Hash{},
			Retry:     0,
			HisHashes: make([]common.Hash, 0),
		}
		node.Steps = append(node.Steps,step)
		node.needCache = true
	}

	if step.TxState == FailWaiting {
		// 直接退出等待新的指令
		return
	}

	if step.TxState == DefaultTS {
		needTrans,success,next := step.CheckApproveAndSendTx(node,client)
		if needTrans > 0 {
			node.SwitchToReq(needTrans)
			return
		}
		if !success {
			return
		}
		if next {
			node.PreStepType = node.StepType
			node.StepType = ApproveHFT
			node.needCache = true
			return
		}
	}

	if step.TxState == Pending {
		success,ok := client.GetTxResult(step.TxHash)
		if !ok  {
			return
		}
		if success {
			step.TxState = Success
		}else{
			step.TxState = Fail
		}
		node.needCache = true
	}

	if step.TxState == Success {
		needTrans,success,next := step.CheckApproveAndSendTx(node,client)
		if needTrans > 0 {
			node.SwitchToReq(needTrans)
			return
		}
		if !success {
			return
		}
		if next {
			node.PreStepType = node.StepType
			node.StepType = ApproveHFT
			node.needCache = true
			return
		}
	}

	if step.TxState == Fail {
		//失败了，就存下该次操作，重置操作步骤
		step.FailWaitingReset(node)
	}
}

func (node *Node) HFTApprove(client *Client) {
	//logrus.Info("执行 HFTApprove")
	//defer logrus.Info("HFTApprove 执行完毕")
	//先检查步骤状态
	step := node.GetStep()
	if step == nil {
		step = &Step{
			StepType:  ApproveHFT,
			TxState:   DefaultTS,
			Coin :     HFT,
			ReqId:     node.ReqId,
			Number:    node.Number,
			OprAddr:   node.OprAddr,
			TxHash:    common.Hash{},
			Retry:     0,
			HisHashes: make([]common.Hash,0),
		}
		node.Steps = append(node.Steps,step)
		node.needCache = true
	}

	if step.TxState == FailWaiting {
		// 直接退出等待新的指令
		return
	}

	if step.TxState == DefaultTS {
		needTrans,success,next := step.CheckApproveAndSendTx(node,client)
		if needTrans > 0 {
			node.SwitchToReq(needTrans)
			return
		}
		if !success {
			return
		}
		if next {
			node.PreStepType = node.StepType
			node.StepType = Bind
			node.needCache = true
			return
		}
	}

	if step.TxState == Pending {
		success,ok := client.GetTxResult(step.TxHash)
		if !ok  {
			return
		}
		if success {
			step.TxState = Success
		}else{
			step.TxState = Fail
		}
		node.needCache = true
	}

	if step.TxState == Success {
		needTrans,success,next := step.CheckApproveAndSendTx(node,client)
		if needTrans > 0 {
			node.SwitchToReq(needTrans)
			return
		}
		if !success {
			return
		}
		if next {
			node.PreStepType = node.StepType
			node.StepType = Bind
			node.needCache = true
			return
		}
	}

	if step.TxState == Fail {
		step.FailWaitingReset(node)
	}
}

func (node *Node) Bind(client *Client) {
	//logrus.Info("执行 Bind")
	//defer logrus.Info("执行完毕")
	//先检查步骤状态
	step := node.GetStep()
	if step == nil {
		step = &Step{
			StepType:  Bind,
			TxState:   DefaultTS,
			Coin :     Gas,
			ReqId:     node.ReqId,
			TxHash:    common.Hash{},
			Retry:     0,
			HisHashes: make([]common.Hash,0),
		}
		node.Steps = append(node.Steps,step)
		node.needCache = true
	}

	// 多个 if 可以从任意地方开始顺序执行
	if step.TxState == DefaultTS {
		needTrans,success,next := step.CheckBindAndSendTx(node,client)
		if needTrans > 0 {
			node.SwitchToReq(needTrans)
			return
		}
		if !success {
			return
		}
		if next {
			node.PreStepType = node.StepType
			node.StepType = Standby
			node.needCache = true
			return
		}
	}

	if step.TxState == Pending {
		success,ok := client.GetTxResult(step.TxHash)
		if !ok {
			return
		}
		if success {
			step.TxState = Success
		}else{
			step.TxState = Fail
		}
		node.needCache = true
	}

	if step.TxState == Success {
		needTrans,success,next := step.CheckBindAndSendTx(node,client)
		if needTrans > 0 {
			node.SwitchToReq(needTrans)
			return
		}
		if !success {
			return
		}
		if next {
			node.PreStepType = node.StepType
			node.StepType = Standby
			node.needCache = true
			return
		}
	}

	if step.TxState == Fail {
		step.FailWaitingReset(node)
	}
}

func (node *Node) BuyHFT(client *Client) {

	//先检查步骤状态
	step := node.GetStep()
	if step == nil {
		step = &Step{
			StepType:  BuyHFT,
			TxState:   DefaultTS,
			Coin :     HFT,
			ReqId:     node.ReqId,
			Number:    node.Number,
			OprAddr:   node.OprAddr,
			TxHash:    common.Hash{},
			Retry:     0,
			HisHashes: make([]common.Hash,0),
		}
		node.Steps = append(node.Steps,step)
		node.needCache = true
	}

	if step.TxState == FailWaiting {
		// 直接退出等待新的指令
		return
	}

	// 多个 if 可以从任意地方开始顺序执行
	if step.TxState == DefaultTS {
		needTrans,success,next := step.CheckBuyAndSendTx(node,client)
		if needTrans > 0 {
			node.SwitchToReq(needTrans)
			return
		}
		if !success {
			return
		}
		if next {
			node.PreStepType = node.StepType
			node.StepType = Standby
			node.Balance.needUpdate = true
			node.needCache = true
			return
		}
	}

	if step.TxState == Pending {
		success,ok := client.GetTxResult(step.TxHash)
		if !ok {
			return
		}
		if success {
			step.TxState = Success
		}else{
			step.TxState = Fail
		}
		node.needCache = true
	}

	if step.TxState == Success {
		needTrans,success,next := step.CheckBuyAndSendTx(node,client)
		if needTrans > 0 {
			node.SwitchToReq(needTrans)
			return
		}
		if !success {
			return
		}
		if next {
			node.PreStepType = node.StepType
			node.StepType = Standby
			node.Balance.needUpdate = true
			node.needCache = true
			return
		}
	}

	if step.TxState == Fail {
		step.FailWaitingReset(node)
	}

}

func (node *Node) BuyHFTWithSpeedUp(client *Client) {

	//先检查步骤状态
	step := node.GetStep()
	if step == nil {
		step = &Step{
			StepType:  BuyHFTWithSpeedUp,
			TxState:   DefaultTS,
			Coin :     HFT,
			ReqId:     node.ReqId,
			Number:    node.Number,
			OprAddr:   node.OprAddr,
			TxHash:    common.Hash{},
			Retry:     0,
			HisHashes: make([]common.Hash,0),
		}
		node.Steps = append(node.Steps,step)
		node.needCache = true
	}

	if step.TxState == FailWaiting {
		// 直接退出等待新的指令
		return
	}

	// 多个 if 可以从任意地方开始顺序执行
	if step.TxState == DefaultTS {
		needTrans,success,next := step.CheckBuyAndSendTx(node,client)
		if needTrans > 0 {
			node.SwitchToReq(needTrans)
			return
		}
		if !success {
			return
		}
		if next {
			node.PreStepType = node.StepType
			node.StepType = Standby
			node.Balance.needUpdate = true
			node.needCache = true
			return
		}
	}

	if step.TxState == Pending {
		success,ok := step.CheckBuyPackAndRefreshTx(node,client)
		if !ok {
			return
		}
		if success {
			step.TxState = Success
		}else{
			step.TxState = Fail
		}
		node.needCache = true
	}

	if step.TxState == Success {
		needTrans,success,next := step.CheckBuyAndSendTx(node,client)
		if needTrans > 0 {
			node.SwitchToReq(needTrans)
			return
		}
		if !success {
			return
		}
		if next {
			node.PreStepType = node.StepType
			node.StepType = Standby
			node.Balance.needUpdate = true
			node.needCache = true
			return
		}
	}

	if step.TxState == Fail {
		step.FailWaitingReset(node)
	}

}

func (node *Node) SellHFT(client *Client) {

	//先检查步骤状态
	step := node.GetStep()
	if step == nil {
		step = &Step{
			StepType:  SellHFT,
			TxState:   DefaultTS,
			Coin :     HFT,
			ReqId:     node.ReqId,
			Number:    node.Number,
			OprAddr:   node.OprAddr,
			TxHash:    common.Hash{},
			Retry:     0,
			HisHashes: make([]common.Hash,0),
		}
		node.Steps = append(node.Steps,step)
		node.needCache = true
	}

	if step.TxState == FailWaiting {
		// 直接退出等待新的指令
		return
	}

	// 多个 if 可以从任意地方开始顺序执行
	if step.TxState == DefaultTS {
		needTrans,success,next := step.CheckSellAndSendTx(node,client)
		if needTrans > 0 {
			node.SwitchToReq(needTrans)
			return
		}
		if !success {
			return
		}
		if next {
			node.PreStepType = node.StepType
			node.StepType = Standby
			node.Balance.needUpdate = true
			node.needCache = true
			return
		}
	}

	if step.TxState == Pending {
		success,ok := client.GetTxResult(step.TxHash)
		if !ok {
			return
		}
		if success {
			step.TxState = Success
		}else{
			step.TxState = Fail
		}
		node.needCache = true
	}

	if step.TxState == Success {
		needTrans,success,next := step.CheckSellAndSendTx(node,client)
		if needTrans > 0 {
			node.SwitchToReq(needTrans)
			return
		}
		if !success {
			return
		}
		if next {
			node.PreStepType = node.StepType
			node.StepType = Standby
			node.Balance.needUpdate = true
			node.needCache = true
			return
		}
	}

	if step.TxState == Fail {
		step.FailWaitingReset(node)
	}

}

func (node *Node) SellHFTWithSpeedUp(client *Client) {

	//先检查步骤状态
	step := node.GetStep()
	if step == nil {
		step = &Step{
			StepType:  SellHFTWithSpeedUp,
			TxState:   DefaultTS,
			Coin :     HFT,
			ReqId:     node.ReqId,
			Number:    node.Number,
			OprAddr:   node.OprAddr,
			TxHash:    common.Hash{},
			Retry:     0,
			HisHashes: make([]common.Hash,0),
		}
		node.Steps = append(node.Steps,step)
		node.needCache = true
	}

	if step.TxState == FailWaiting {
		// 直接退出等待新的指令
		return
	}

	// 多个 if 可以从任意地方开始顺序执行
	if step.TxState == DefaultTS {
		needTrans,success,next := step.CheckSellAndSendTx(node,client)
		if needTrans > 0 {
			node.SwitchToReq(needTrans)
			return
		}
		if !success {
			return
		}
		if next {
			node.PreStepType = node.StepType
			node.StepType = Standby
			node.Balance.needUpdate = true
			node.needCache = true
			return
		}
	}

	if step.TxState == Pending {
		success,ok := step.CheckSellPackAndRefreshTx(node,client)
		if !ok {
			return
		}
		if success {
			step.TxState = Success
		}else{
			step.TxState = Fail
		}
		node.needCache = true
	}

	if step.TxState == Success {
		needTrans,success,next := step.CheckSellAndSendTx(node,client)
		if needTrans > 0 {
			node.SwitchToReq(needTrans)
			return
		}
		if !success {
			return
		}
		if next {
			node.PreStepType = node.StepType
			node.StepType = Standby
			node.Balance.needUpdate = true
			node.needCache = true
			return
		}
	}

	if step.TxState == Fail {
		step.FailWaitingReset(node)
	}

}

func (node *Node) CollectUsdt(client *Client) {

	//先检查步骤状态
	step := node.GetStep()
	if step == nil {
		step = &Step{
			StepType:  CollectUSDT,
			TxState:   DefaultTS,
			Coin :     USDT,
			ReqId:     node.ReqId,
			Number:    node.Number,
			OprAddr:   node.OprAddr,
			TxHash:    common.Hash{},
			Retry:     0,
			HisHashes: make([]common.Hash,0),
		}
		node.Steps = append(node.Steps,step)
		node.needCache = true
	}

	if step.TxState == FailWaiting {
		// 直接退出等待新的指令
		return
	}

	// 多个 if 可以从任意地方开始顺序执行
	if step.TxState == DefaultTS {
		needTrans,success,next := step.CheckCollectUsdtAndSendTx(node,client)
		if needTrans > 0 {
			node.SwitchToReq(needTrans)
			return
		}
		if !success {
			return
		}
		if next {
			node.PreStepType = node.StepType
			node.StepType = Standby
			node.Balance.needUpdate = true
			node.needCache = true
			return
		}
	}

	if step.TxState == Pending {
		success,ok := client.GetTxResult(step.TxHash)
		if !ok {
			return
		}
		if success {
			step.TxState = Success
		}else{
			step.TxState = Fail
		}
		node.needCache = true
	}

	if step.TxState == Success {
		needTrans,success,next := step.CheckCollectUsdtAndSendTx(node,client)
		if needTrans > 0 {
			node.SwitchToReq(needTrans)
			return
		}
		if !success {
			return
		}
		if next {
			node.PreStepType = node.StepType
			node.StepType = Standby
			node.Balance.needUpdate = true
			node.needCache = true
			return
		}
	}

	if step.TxState == Fail {
		step.FailWaitingReset(node)
	}

}

func (node *Node) CollectBft(client *Client) {

	//先检查步骤状态
	step := node.GetStep()
	if step == nil {
		step = &Step{
			StepType:  CollectBFT,
			TxState:   DefaultTS,
			Coin :     HFT,
			ReqId:     node.ReqId,
			Number:    node.Number,
			OprAddr:   node.OprAddr,
			TxHash:    common.Hash{},
			Retry:     0,
			HisHashes: make([]common.Hash,0),
		}
		node.Steps = append(node.Steps,step)
		node.needCache = true
	}

	if step.TxState == FailWaiting {
		// 直接退出等待新的指令
		return
	}

	// 多个 if 可以从任意地方开始顺序执行
	if step.TxState == DefaultTS {
		needTrans,success,next := step.CheckCollectBftAndSendTx(node,client)
		if needTrans > 0 {
			node.SwitchToReq(needTrans)
			return
		}
		if !success {
			return
		}
		if next {
			node.PreStepType = node.StepType
			node.StepType = Standby
			node.Balance.needUpdate = true
			node.needCache = true
			return
		}
	}

	if step.TxState == Pending {
		success,ok := client.GetTxResult(step.TxHash)
		if !ok {
			return
		}
		if success {
			step.TxState = Success
		}else{
			step.TxState = Fail
		}
		node.needCache = true
	}

	if step.TxState == Success {
		needTrans,success,next := step.CheckCollectBftAndSendTx(node,client)
		if needTrans > 0 {
			node.SwitchToReq(needTrans)
			return
		}
		if !success {
			return
		}
		if next {
			node.PreStepType = node.StepType
			node.StepType = Standby
			node.Balance.needUpdate = true
			node.needCache = true
			return
		}
	}

	if step.TxState == Fail {
		step.FailWaitingReset(node)
	}

}

func (node *Node) CollectBnb(client *Client) {

	//先检查步骤状态
	step := node.GetStep()
	if step == nil {
		step = &Step{
			StepType:  CollectBNB,
			TxState:   DefaultTS,
			Coin :     Gas,
			ReqId:     node.ReqId,
			Number:    node.Number,
			OprAddr:   node.OprAddr,
			TxHash:    common.Hash{},
			Retry:     0,
			HisHashes: make([]common.Hash,0),
		}
		node.Steps = append(node.Steps,step)
		node.needCache = true
	}

	if step.TxState == FailWaiting {
		// 直接退出等待新的指令
		return
	}

	// 多个 if 可以从任意地方开始顺序执行
	if step.TxState == DefaultTS {
		success,next := step.CheckCollectBnbAndSendTx(node,client)
		if !success {
			return
		}
		if next {
			node.PreStepType = node.StepType
			node.StepType = Standby
			node.Balance.needUpdate = true
			node.needCache = true
			return
		}
	}

	if step.TxState == Pending {
		success,ok := client.GetTxResult(step.TxHash)
		if !ok {
			return
		}
		if success {
			step.TxState = Success
		}else{
			step.TxState = Fail
		}
		node.needCache = true
	}

	if step.TxState == Success {
		success,next := step.CheckCollectBnbAndSendTx(node,client)
		if !success {
			return
		}
		if next {
			node.PreStepType = node.StepType
			node.StepType = Standby
			node.Balance.needUpdate = true
			node.needCache = true
			return
		}
	}

	if step.TxState == Fail {
		step.FailWaitingReset(node)
	}

}


func (node *Node) GasCheck(client *Client, chSendReq chan *SendReq) {
	//logrus.Info("执行 GasCheck")
	//defer logrus.Info("GasCheck 执行完毕")
	//先检查步骤状态
	step := node.GetStep()
	if step == nil {
		step = &Step{
			StepType:  TransferGas,
			TxState:   DefaultTS,
			Coin :     Gas,
			ReqId:     node.ReqId,
			Number:    node.Number,
			OprAddr:   node.OprAddr,
			TxHash:    common.Hash{},
			Retry:     0,
			HisHashes: make([]common.Hash,0),
		}
		node.Steps = append(node.Steps,step)
		node.needCache = true
	}
	// 先检查是不是默认状态，即之前有没有请求过 master 账号发送 gas
	if step.TxState == DefaultTS {
		success,next := step.CheckBalanceAndReqMaster(node,client,chSendReq)
		if !success {
			return
		}
		if next {
			node.StepType = node.PreStepType
			node.needCache = true
			return
		}
	}

	//多个if可以从任意地方开始顺序执行
	if step.TxState == Pending {
		success,ok := client.GetTxResult(step.TxHash)
		if !ok {
			return
		}
		if success {
			step.TxState = Success
		}else{
			step.TxState = Fail
		}
		node.needCache = true
	}

	if step.TxState == Success {
		success,next := step.CheckBalanceAndReqMaster(node,client,chSendReq)
		if !success {
			return
		}
		if next {
			node.StepType = node.PreStepType
			node.needCache = true
			return
		}
	}

	if step.TxState == FailWaiting {
		//失败了，就存下该次操作，重置操作步骤
		step.FailWaitingReset(node)
	}

}

func (node *Node) USDTCheck(client *Client, chSendReq chan *SendReq) {
	//logrus.Info("执行 USDTCheck")
	//defer logrus.Info("执行完毕")
	//先检查步骤状态
	step := node.GetStep()
	if step == nil {
		step = &Step{
			StepType:  TransferUSDT,
			TxState:   DefaultTS,
			Coin :     USDT,
			ReqId:     node.ReqId,
			Number:    node.Number,
			OprAddr:   node.OprAddr,
			TxHash:    common.Hash{},
			Retry:     0,
			HisHashes: make([]common.Hash,0),
		}
		node.Steps = append(node.Steps,step)
		node.needCache = true
	}

	// 多个 if 可以从任意地方开始顺序执行
	if step.TxState == DefaultTS {
		success,next := step.CheckBalanceAndReqMaster(node,client,chSendReq)
		if !success {
			return
		}
		if next {
			node.StepType = node.PreStepType
			node.needCache = true
			return
		}
	}

	if step.TxState == Pending {
		success,ok := client.GetTxResult(step.TxHash)
		if !ok {
			return
		}
		if success {
			step.TxState = Success
		}else{
			step.TxState = Fail
		}
		node.needCache = true
	}

	if step.TxState == Success {
		success,next := step.CheckBalanceAndReqMaster(node,client,chSendReq)
		if !success {
			return
		}
		if next {
			node.StepType = node.PreStepType
			node.needCache = true
			return
		}
	}

	if step.TxState == Fail {
		//失败了，就存下该次操作，重置操作步骤
		step.FailWaitingReset(node)
	}
}

func (node *Node) HFTCheck(client *Client, chSendReq chan *SendReq) {
	//logrus.Info("执行 HFTCheck")
	//defer logrus.Info("HFTCheck 执行完毕")
	//先检查步骤状态
	step := node.GetStep()
	if step == nil {
		step = &Step{
			StepType:  TransferHFT,
			TxState:   DefaultTS,
			Coin :     HFT,
			ReqId:     node.ReqId,
			Number:    node.Number,
			OprAddr:   node.OprAddr,
			TxHash:    common.Hash{},
			Retry:     0,
			HisHashes: make([]common.Hash,0),
		}
		node.Steps = append(node.Steps,step)
		node.needCache = true
	}

	// 多个 if 可以从任意地方开始顺序执行
	if step.TxState == DefaultTS {
		success,next := step.CheckBalanceAndReqMaster(node,client,chSendReq)
		if !success {
			return
		}
		if next {
			node.StepType = node.PreStepType
			node.needCache = true
			return
		}
	}

	if step.TxState == Pending {
		success,ok := client.GetTxResult(step.TxHash)
		if !ok {
			return
		}
		if success {
			step.TxState = Success
		}else{
			step.TxState = Fail
		}
		node.needCache = true
	}

	if step.TxState == Success {
		success,next := step.CheckBalanceAndReqMaster(node,client,chSendReq)
		if !success {
			return
		}
		if next {
			node.StepType = node.PreStepType
			node.needCache = true
			return
		}
	}

	if step.TxState == Fail {
		//失败了，就存下该次操作，重置操作步骤
		step.FailWaitingReset(node)
	}
}



func (step *Step) FailWaitingReset(node *Node) {
	step.Retry ++
	step.HisHashes = append(step.HisHashes,step.TxHash)
	step.TxHash = common.Hash{}
	step.TxState = FailWaiting
	node.needCache = true
}


func (step *Step) CheckBalanceAndReqMaster(node *Node, client *Client, chSendReq chan *SendReq) (success,next bool) {
	var balance,needBalance *big.Int
	var ok bool
	switch step.Coin {
	case Gas:
		balance,ok = client.GetBalance(node.Addr)
		needBalance = nodeParams.MinGas
	case USDT:
		balance,ok = client.GetUSDTBalance(node.Addr)
		needBalance = nodeParams.NeedUSDT
	case HFT:
		balance,ok = client.GetHFTBalance(node.Addr)
		needBalance = nodeParams.NeedHFT
	default:
		return false,false
	}
	if !ok {
		return false,false
	}
	// a < b
	if balance.Cmp(needBalance) == -1 { //余额不足，请求 master 发送 coin
		sendReq := SendReq{
			Addr:  node.Addr,
			Coin:  step.Coin,
			Value: needBalance,
			Tx:    nil,
			Down:  make(chan struct{}),
		}
		chSendReq <- &sendReq
		<- sendReq.Down //阻塞等待 master 处理
		if sendReq.Tx != nil { //说明 master 请求成功了
			step.TxState = Pending
			step.TxHash = sendReq.Tx.Hash()
			node.needCache = true
			return true,false
		}
		return false,false
	}else{ //余额足够了，表示可以下一步了（下一个大步骤，不是step里的步骤）
		return true,true
	}
}

func (step *Step) CheckApproveAndSendTx(node *Node, client *Client) (needTrans StepType,success,next bool) {
	var approved,ok bool
	switch step.Coin {
	case USDT:
		approved,ok = client.IsUSDTApproved(node.auth.From,nodeParams.USDTApproveAddr)
	case HFT:
		approved,ok = client.IsHFTApproved(node.auth.From,nodeParams.SATApproveAddr)
	default:
		logrus.Infof("CheckApproveAndSendTx coin %v not support",step.Coin)
		return DefaultST,false,false
	}
	if !ok {
		return DefaultST,false,false
	}
	// a < b
	if !approved { //没有授权，发送授权交易
		var tx *types.Transaction
		var needGas bool
		switch step.Coin {
		case USDT:
			needGas,tx,ok = client.ApproveUSDT(node.auth,nodeParams.USDTApproveAddr,nodeParams.ApproveUSDT)
		case HFT:
			needGas,tx,ok = client.ApproveHFT(node.auth,nodeParams.SATApproveAddr,nodeParams.ApproveSAT)
		}
		if needGas {
			return TransferGas,false,false
		}
		if !ok {
			return DefaultST,false,false
		}
		step.TxState = Pending
		step.TxHash  = tx.Hash()
		node.needCache = true
		return DefaultST,true,false
	}else{ //余额足够了，表示可以下一步了（下一个大步骤，不是step里的步骤）
		return DefaultST,true,true
	}
}

func (step *Step) CheckBindAndSendTx(node *Node, client *Client,) (needTrans StepType,success,next bool) {
	inviter,ok := client.GetInviter(node.Addr)
	if !ok {
		return DefaultST,false,false
	}
	// a < b
	if inviter == ZeroAddr { //没有绑定，发送绑定交易
		needTrans,tx,ok := client.Bind(node.auth,node.sunAddr)
		if needTrans > 0 {
			return needTrans,false,false
		}
		if !ok {
			return DefaultST,false,false
		}
		step.TxState = Pending
		step.TxHash  = tx.Hash()
		node.needCache = true
		return DefaultST,true,false
	}else{ //绑定好了，表示可以下一步了（下一个大步骤，不是step里的步骤）
		if inviter != node.sunAddr {
			logrus.Errorf("请注意！地址 %v 应该绑定 %v，但实际绑定 %v ！",
				node.Addr.String(), node.sunAddr.String(), inviter.String())
		}
		return DefaultST,true,true
	}
}


func (step *Step) CheckBuyAndSendTx(node *Node, client *Client) (needTrans StepType,success,next bool) {
	hftBalance,executed,ok := client.IsBuyExecuted(node.auth.From,step.Expected)
	if !ok {
		return DefaultST,false,false
	}
	node.Balance.HFT = hftBalance
	// a < b
	if !executed { //还没有购买成功
		needTrans,tx,expected,ok := client.BuyHFT(node.auth,step.Number,hftBalance)
		if needTrans > 0 {
			return needTrans,false,false
		}
		if !ok {
			return DefaultST,false,false
		}
		step.Expected = expected
		step.TxState = Pending
		step.TxHash  = tx.Hash()
		node.needCache = true
		return DefaultST,true,false
	}else{ //绑定好了，表示可以下一步了（下一个大步骤，不是step里的步骤）
		return DefaultST,true,true
	}
}

func (step *Step) CheckBuyPackAndRefreshTx(node *Node, client *Client) (success,ok bool) {
	hftBalance,executed,ok := client.IsBuyExecuted(node.auth.From,step.Expected)
	if !ok {
		return false,false
	}
	node.Balance.HFT = hftBalance
	// a < b
	if !executed {
		tx,isPending,ok := client.GetTx(step.TxHash)
		if !ok {
			return false,false
		}
		if isPending {
			tx,expected,ok := client.RefreshBuyHFT(node.auth,tx,step.Number,hftBalance)
			if ok {
				step.Expected = expected
				step.TxHash  = tx.Hash()
				node.needCache = true
			}
			return false,false
		}else{
			return client.GetTxResult(step.TxHash)
		}
	}else{ //绑定好了，表示可以下一步了（下一个大步骤，不是step里的步骤）
		return true,true
	}
}


func (step *Step) CheckSellAndSendTx(node *Node, client *Client) (needTrans StepType,success,next bool) {
	usdtBalance,executed,ok := client.IsSellExecuted(node.auth.From,step.Expected)
	if !ok {
		return DefaultST,false,false
	}
	node.Balance.USDT = usdtBalance
	// a < b
	if !executed { //还没有购买成功
		needTrans,tx,expected,ok := client.SellHFT(node.auth,step.Number,usdtBalance)
		if needTrans > 0 {
			return needTrans,false,false
		}
		if !ok {
			return DefaultST,false,false
		}
		step.Expected = expected
		step.TxState = Pending
		step.TxHash  = tx.Hash()
		node.needCache = true
		return DefaultST,true,false
	}else{ //绑定好了，表示可以下一步了（下一个大步骤，不是step里的步骤）
		return DefaultST,true,true
	}
}

func (step *Step) CheckSellPackAndRefreshTx(node *Node, client *Client) (success,ok bool) {
	usdtBalance,executed,ok := client.IsSellExecuted(node.auth.From,step.Expected)
	if !ok {
		return false,false
	}
	node.Balance.USDT = usdtBalance
	// a < b
	if !executed {
		tx,isPending,ok := client.GetTx(step.TxHash)
		if !ok {
			return false,false
		}
		if isPending {
			tx,expected,ok := client.RefreshSellHFT(node.auth,tx,step.Number,usdtBalance)
			if ok {
				step.Expected = expected
				step.TxHash  = tx.Hash()
				node.needCache = true
			}
			return false,false
		}else{
			return client.GetTxResult(step.TxHash)
		}
	}else{ //绑定好了，表示可以下一步了（下一个大步骤，不是step里的步骤）
		return true,true
	}
}

func (step *Step) CheckCollectUsdtAndSendTx(node *Node, client *Client) (needTrans StepType,success,next bool) {
	usdtBalance,collected,ok := client.IsUsdtCollected(node.auth.From,step.Expected)
	if !ok {
		return DefaultST,false,false
	}
	node.Balance.USDT = usdtBalance
	// a < b
	if !collected { //还没有归集成功
		if step.Number == nil {
			step.Number = big.NewInt(0)
		}
		if step.Number.Cmp(usdtBalance) == 1 {
			step.Number = usdtBalance
			step.Expected = big.NewInt(0)
		}else{
			step.Expected = new(big.Int).Sub(usdtBalance,step.Number)
		}

		needTrans,tx,ok := client.CollectUsdt(node.auth,step.OprAddr,step.Number)
		if needTrans > 0 {
			return needTrans,false,false
		}
		if !ok {
			return DefaultST,false,false
		}
		step.TxState = Pending
		step.TxHash  = tx.Hash()
		node.needCache = true
		return DefaultST,true,false
	}else{ //表示可以下一步了（下一个大步骤，不是step里的步骤）
		return DefaultST,true,true
	}
}

func (step *Step) CheckCollectBftAndSendTx(node *Node, client *Client) (needTrans StepType,success,next bool) {
	bftBalance,collected,ok := client.IsBftCollected(node.auth.From,step.Expected)
	if !ok {
		return DefaultST,false,false
	}
	node.Balance.HFT = bftBalance
	// a < b
	if !collected { //还没有归集成功
		if step.Number == nil {
			step.Number = big.NewInt(0)
		}
		if step.Number.Cmp(bftBalance) == 1 {
			step.Number = bftBalance
			step.Expected = big.NewInt(0)
		}else{
			step.Expected = new(big.Int).Sub(bftBalance,step.Number)
		}

		needTrans,tx,ok := client.CollectBft(node.auth,step.OprAddr,step.Number)
		if needTrans > 0 {
			return needTrans,false,false
		}
		if !ok {
			return DefaultST,false,false
		}
		step.TxState = Pending
		step.TxHash  = tx.Hash()
		node.needCache = true
		return DefaultST,true,false
	}else{ //表示可以下一步了（下一个大步骤，不是step里的步骤）
		return DefaultST,true,true
	}
}

func (step *Step) CheckCollectBnbAndSendTx(node *Node, client *Client) (success,next bool) {
	bnbBalance,collected,ok := client.IsBnbCollected(node.auth)
	if !ok {
		return false,false
	}
	node.Balance.BNB = bnbBalance
	// a < b
	if !collected { //还没有归集成功
		step.Number = bnbBalance

		tx,ok := client.CollectBnb(node.auth,step.OprAddr,step.Number)
		if !ok {
			return false,false
		}
		step.TxState = Pending
		step.TxHash  = tx.Hash()
		node.needCache = true
		return true,false
	}else{ //表示可以下一步了（下一个大步骤，不是step里的步骤）
		return true,true
	}
}

