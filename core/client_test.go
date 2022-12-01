package core

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func InitClientForTest(t *testing.T,core *Core)  {
	var err error
	r0,r1,ok := core.Client.GetReserve()
	require.True(t, ok)
	core.Client.Last.Reserved.HasLiquidity = true
	core.Client.Last.Reserved.USDT = r0
	core.Client.Last.Reserved.SAT = r1
	core.Client.Last.Header,err = core.Client.HeaderByNumber(context.Background(),nil)
	require.Empty(t, err)
	maxPath2,err := core.GetMaxPath2()
	require.Empty(t, err)
	chNode := core.IterNode()
	err = core.Galaxy.Load(maxPath2,&core.Master,chNode)
	require.Empty(t, err)
}

func TestSendGas(t *testing.T) {
	core,err := InitCore("1")
	require.Empty(t, err)
	defer core.LvDb.Close()
	sendReq := &SendReq{
		Addr:  common.HexToAddress("0xa6f79b60359f141df90a0c745125b131caaffd12"),
		Coin:  Gas,
		Value: nodeParams.MinGas,
		Tx:    nil,
		Down:  make(chan struct{}),
	}
	core.Master.SendGas(&core.Client,sendReq)
	<- sendReq.Down
	t.Log(sendReq.Tx.Hash().String())
}

func TestSendUSDT(t *testing.T) {
	core,err := InitCore("1")
	require.Empty(t, err)
	defer core.LvDb.Close()
	sendReq := &SendReq{
		Addr:  common.HexToAddress("0xa6f79b60359f141df90a0c745125b131caaffd12"),
		Coin:  USDT,
		Value: nodeParams.NeedUSDT,
		Tx:    nil,
		Down:  make(chan struct{}),
	}
	core.Master.SendUSDT(&core.Client,sendReq)
	<- sendReq.Down
	require.NotEmpty(t, sendReq.Tx)
	t.Log(sendReq.Tx.Hash().String())
}

func TestListenSendReqGas(t *testing.T) {
	core,err := InitCore("1")
	require.Empty(t, err)
	defer core.LvDb.Close()
	chSendReq := core.Master.ListenSendReq(&core.Client)
	sendReq := &SendReq{
		Addr:  common.HexToAddress("0xa6f79b60359f141df90a0c745125b131caaffd12"),
		Coin:  Gas,
		Value: nodeParams.MinGas,
		Tx:    nil,
		Down:  make(chan struct{}),
	}
	chSendReq <- sendReq
	<- sendReq.Down
	require.NotEmpty(t, sendReq.Tx)
	t.Log(sendReq.Tx.Hash().String())
}

func TestListenSendReqUsdt(t *testing.T) {
	core,err := InitCore("1")
	require.Empty(t, err)
	defer core.LvDb.Close()
	chSendReq := core.Master.ListenSendReq(&core.Client)
	sendReq := &SendReq{
		Addr:  common.HexToAddress("0xa6f79b60359f141df90a0c745125b131caaffd12"),
		Coin:  USDT,
		Value: nodeParams.NeedUSDT,
		Tx:    nil,
		Down:  make(chan struct{}),
	}
	chSendReq <- sendReq
	<- sendReq.Down
	require.NotEmpty(t, sendReq.Tx)
	t.Log(sendReq.Tx.Hash().String())
}

func TestIsHFTApproved(t *testing.T) {
	core,err := InitCore("1")
	require.Empty(t, err)
	defer core.LvDb.Close()
	approved,ok := core.Client.IsHFTApproved(common.HexToAddress("0xa6f79b60359f141df90a0c745125b131caaffd12"),nodeParams.SATApproveAddr)
	t.Log("approved: ",approved)
	t.Log("ok: ",ok)
}

func TestApproveHFT(t *testing.T) {
	core,err := InitCore("1")
	require.Empty(t, err)
	defer core.LvDb.Close()
	needGas,tx,ok := core.Client.ApproveHFT(
		core.Master.auth,nodeParams.SATApproveAddr,nodeParams.ApproveSAT)
	t.Log(ok)
	t.Log(needGas)
	if tx != nil {
		b,err := tx.MarshalJSON()
		t.Log(string(b),err)
	}
}

func TestIsUSDTApproved(t *testing.T) {
	core,err := InitCore("1")
	require.Empty(t, err)
	defer core.LvDb.Close()
	approved,ok := core.Client.IsUSDTApproved(common.HexToAddress("0xa6f79b60359f141df90a0c745125b131caaffd12"),nodeParams.USDTApproveAddr)
	t.Log("approved: ",approved)
	t.Log("ok: ",ok)
}

func TestApproveUSDT(t *testing.T) {
	core,err := InitCore("1")
	require.Empty(t, err)
	defer core.LvDb.Close()
	needGas,tx,ok := core.Client.ApproveUSDT(
		core.Master.auth,nodeParams.USDTApproveAddr,nodeParams.ApproveUSDT)
	t.Log(ok)
	t.Log(needGas)
	if tx != nil {
		b,err := tx.MarshalJSON()
		t.Log(string(b),err)
	}
}

func TestGetUSDTBalance(t *testing.T) {
	core,err := InitCore("1")
	require.Empty(t, err)
	defer core.LvDb.Close()
	balance,ok := core.Client.GetUSDTBalance(common.HexToAddress("0xa6f79b60359f141df90a0c745125b131caaffd12"))
	t.Log("balance: ",balance)
	t.Log("ok: ",ok)
}


func TestIsBuyExecuted(t *testing.T) {
	core,err := InitCore("1")
	require.Empty(t, err)
	defer core.LvDb.Close()
	hftBalance,executed,ok := core.Client.IsBuyExecuted(common.HexToAddress("0xa6f79b60359f141df90a0c745125b131caaffd12"),nil)
	t.Log("hftBalance: ",hftBalance)
	t.Log("executed: ",executed)
	t.Log("ok: ",ok)
}

func TestGetAmountOut(t *testing.T) {
	core,err := InitCore("1")
	require.Empty(t, err)
	defer core.LvDb.Close()
	r0,r1,ok := core.Client.GetReserve()
	require.True(t, ok)
	core.Client.Last.Reserved.HasLiquidity = true
	core.Client.Last.Reserved.USDT = r0
	core.Client.Last.Reserved.SAT = r1
	out := core.Client.GetUsdtAmountOut(new(big.Int).Mul(big.NewInt(1000),big.NewInt(1e18)))
	o := new(big.Int).Div(out,big.NewInt(1e18)).Int64()
	t.Log("AmountOut: ",o)
	for i:=0; i< 40; i ++ {
		core.Client.Last.Reserved.USDT = new(big.Int).Sub(core.Client.Last.Reserved.USDT,out)
		core.Client.Last.Reserved.SAT = new(big.Int).Add(core.Client.Last.Reserved.SAT,big.NewInt(917))
		out1 := core.Client.GetUsdtAmountOut(new(big.Int).Mul(big.NewInt(1000),big.NewInt(1e18)))
		o1 := new(big.Int).Div(out1,big.NewInt(1e18)).Int64()
		t.Logf("%v times AmountOut: %v, slippage: %v",i+1,o1,float64(o-o1)/float64(o)*100)
	}

}

func TestGetReserve(t *testing.T) {
	core,err := InitCore("1")
	require.Empty(t, err)
	defer core.LvDb.Close()
	r0,r1,ok := core.Client.GetReserve()
	t.Log("r0: ",r0.String())
	t.Log("r1: ",r1.String())
	t.Log("ok: ",ok)
}

func TestBuyHft(t *testing.T) {
	core,err := InitCore("1")
	require.Empty(t, err)
	defer core.LvDb.Close()
	InitClientForTest(t,core)
	//node,err := core.GetNode("m-1001-3")
	//require.Empty(t, err)
	//needTrans,tx,expected,ok := core.Client.BuyHFT(node.auth,new(big.Int).Mul(big.NewInt(100),big.NewInt(1e18)))
	//t.Log(needTrans)
	//t.Log(ok)
	//t.Log(expected)
	//if tx != nil {
	//	b,err := tx.MarshalJSON()
	//	t.Log(string(b),err)
	//}
}

func TestSellHFT(t *testing.T) {
	core,err := InitCore("1")
	require.Empty(t, err)
	defer core.LvDb.Close()

	//tx,expected,ok := core.Client.SellHFT(core.Master.auth,new(big.Int).Mul(big.NewInt(1000),big.NewInt(1e18)))
	//t.Log(ok)
	//t.Log(expected)
	//if tx != nil {
	//	b,err := tx.MarshalJSON()
	//	t.Log(string(b),err)
	//}
}
