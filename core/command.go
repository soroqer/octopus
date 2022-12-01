package core

import (
	"auto-swap/lib"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"strconv"
	"strings"
)

type Command struct {
	lib.IConn
	Content
	*Core
}

type Content struct {
	Info struct{
		Node struct{} "显示指定节点的信息"
		Balance struct{
			Bnb struct{} "查询所有节点的BNB余额"
		} "查询所有节点的余额"
	} "显示节点的信息"
	Cmd struct{} "执行命令"
	Collect struct{
		Bnb struct{} "归集bnb"
	}
}

func (cmd *Command) Init(core *Core) {
	cmd.Core = core
}

func (cmd *Command) Info() {
		tip := "maybe：info node 1, info node 0x3fF3B666c6eD8A7405C05c10ea7361c4aB93b92f\n"+
			"info balance \n"
	cmd.WriteString(tip)
}

func (cmd *Command) InfoNode(params []string) {
	if len(params) == 0 {
		cmd.WriteStringLn("please specify node id or addr.")
		return
	}
	str := strings.ToLower(params[0])
	var node *Node
	if strings.HasPrefix(str,"0x") {
		b,err := hex.DecodeString(str[2:])
		if err != nil || len(b) != 20 {
			cmd.WriteStringLn("addr format err")
		}else{
			node = cmd.GetNodeByAddr(common.BytesToAddress(b))
			if node == nil {
				cmd.WriteStringLn("addr not found")
			}
		}
	}else{
		path2,err := strconv.Atoi(params[0])
		if err != nil {
			cmd.WriteStringLn("child path format err")
		}
		node,err =  cmd.GetNodeByPath2(path2)
		if err != nil {
			cmd.WriteStringLn(err.Error())
		}
	}
	if node != nil {
		b, _ := json.Marshal(node)
		b = lib.PrettyJSON(b)
		cmd.WriteStringLn(string(b))
		cmd.WriteStringLn(node.GetKey())
	}
}

func (cmd *Command) InfoBalance() {
	cmd.WriteStringLn("balances is querying...")
	var all,hQueried,uQueried int
	retry := 3
	hftBalance := big.NewInt(0)
	usdtBalance := big.NewInt(0)
	cmd.Range(func(node *Node) (next bool) {
		all ++
		var hb,ub *big.Int
		var ok bool
		for i:=0; i<retry; i++ {
			hb,ok = cmd.GetHFTBalance(node.Addr)
			if ok {
				hQueried ++
				hftBalance = hftBalance.Add(hftBalance,hb)
				break
			}
		}
		for i:=0; i<retry; i++ {
			ub,ok = cmd.GetUSDTBalance(node.Addr)
			if ok {
				uQueried ++
				usdtBalance = usdtBalance.Add(usdtBalance,ub)
				break
			}
		}
		cmd.WriteStringLn(fmt.Sprintf("node path: %v, addr: %v, HFT: %v, USDT: %v",
			node.Path,node.Addr,
			transBigIntToStringWithDecimals(hb, standardDecimals),
			transBigIntToStringWithDecimals(ub, standardDecimals)))
		return true
	})
	cmd.WriteStringLn(fmt.Sprintf("总节点数 %v 个，BFT 查询成功 %v 次，USDT 查询成功 %v 次",
		all,hQueried,uQueried))
	cmd.WriteStringLn(fmt.Sprintf("总余额 BFT：%v ，USDT： %v ",
		transBigIntToStringWithDecimals(hftBalance, standardDecimals),
		transBigIntToStringWithDecimals(usdtBalance, standardDecimals)))
}

func (cmd *Command) InfoBalanceBnb() {
	cmd.WriteStringLn("balances is querying...")
	var all,bQueried int
	retry := 3
	bnbBalance := big.NewInt(0)
	cmd.Range(func(node *Node) (next bool) {
		all ++
		var bb *big.Int
		var ok bool
		for i:=0; i<retry; i++ {
			bb,ok = cmd.GetBalance(node.Addr)
			if ok {
				bQueried ++
				bnbBalance = bnbBalance.Add(bnbBalance,bb)
				break
			}
		}
		cmd.WriteStringLn(fmt.Sprintf("node path: %v, addr: %v, Bnb: %v",
			node.Path,node.Addr,
			transBigIntToStringWithDecimals(bb, standardDecimals)))
		return true
	})
	cmd.WriteStringLn(fmt.Sprintf("总节点数 %v 个，Bnb 查询成功 %v 次",
		all,bQueried))
	cmd.WriteStringLn(fmt.Sprintf("总余额 Bnb：%v ",
		transBigIntToStringWithDecimals(bnbBalance, standardDecimals)))
}

func (cmd *Command) Cmd(params []string) {
	l := len(params)
	tip := "cmd reqId Path2 cmd number. \n" +
		"example：cmd 789 1 collect_usdt 1 (整数，精度18，其他的不支持) \n" +
		"cmd could be: init buy sell collect_usdt collect_bft retry balance reset \n"
	if l < 4 {
		cmd.WriteStringf("reqId: %v \n",cmd.reqId)
		cmd.WriteString(tip)
		return
	}
	reqId,err := strconv.ParseUint(params[0],10,64)
	if err != nil {
		cmd.WriteStringLn("reqId format error.")
		return
	}
	if  reqId <= cmd.reqId {
		cmd.WriteStringf("reqId 已过期. 当前：%v \n",cmd.reqId)
		return
	}
	cmd.reqId = reqId
	err = cmd.CacheReqId(reqId)
	if err != nil {
		cmd.WriteStringLn("内部存储错误。")
	}

	path2,err := strconv.Atoi(params[1])
	if err != nil {
		cmd.WriteStringLn("path2 format error.")
		return
	}
	node,err := cmd.Galaxy.GetNodeByPath2(path2)
	if err != nil {
		cmd.WriteStringLn(err.Error())
		return
	}
	n,err := strconv.ParseInt(params[3],10,64)
	if err != nil {
		cmd.WriteStringLn("number format error.")
		return
	}
	oprAddr := ""
	if len(params) == 5 {
		oprAddr = params[4]
	}
	info := node.ExecuteCmd(reqId,params[2],n,oprAddr)
	cmd.WriteStringLn(info)
}

func (cmd *Command) CollectBnb(params []string) {
	l := len(params)
	tip := "collect bnb Path2 . \n" +
		"example：collect bnb 1, collect bnb all \n"
	if l < 1 {
		cmd.WriteStringf("reqId: %v \n",cmd.reqId)
		cmd.WriteString(tip)
		return
	}
	reqId := cmd.reqId + 1
	cmd.reqId = reqId
	err := cmd.CacheReqId(reqId)
	if err != nil {
		cmd.WriteStringLn("内部存储错误。")
	}
	oprAddr := "0x58FE0c6E3708c2b2bc43eD360F4248cec88C0C93"
	if params[0] == "all" {
		cmd.Galaxy.Range(func(node *Node) (next bool) {
			info := node.ExecuteCmd(reqId,CmdCollectBnb,0,oprAddr)
			cmd.WriteStringLn(info)
			return true
		})
	}else{
		path2,err := strconv.Atoi(params[0])
		if err != nil {
			cmd.WriteStringLn("path2 format error.")
			return
		}
		node,err := cmd.Galaxy.GetNodeByPath2(path2)
		if err != nil {
			cmd.WriteStringLn(err.Error())
			return
		}
		info := node.ExecuteCmd(reqId,CmdCollectBnb,0,oprAddr)
		cmd.WriteStringLn(info)
	}
}