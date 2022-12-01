package core

import (
	"auto-swap/config"
	"context"
	"encoding/json"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
	"time"
)

func init() {
	config.TestInit()
}

func TestName(t *testing.T) {
	a,b := new(big.Float).SetString("1e10")
	t.Log(a.String())
	t.Log(b)
	c,d := a.Int(nil)
	t.Log(c.String())
	t.Log(d)
	f := big.NewInt(18e17)
	t.Log(f.String())
	gas := new(big.Float).Mul(big.NewFloat(100e4),big.NewFloat(10e9))
	gas = new(big.Float).Quo(gas,big.NewFloat(1e18))
	t.Log(gas.String())
}

func TestInitCore(t *testing.T) {
	core,err := InitCore("1")
	require.Empty(t, err)
	defer core.LvDb.Close()
	t.Logf("core master addr: %v",crypto.PubkeyToAddress(core.Master.privKey.PublicKey))
}

func TestInitClient(t *testing.T) {
	core := &Core{
		exit: make(chan struct{}),
	}
	err := core.Client.Init(core.exit)
	require.Empty(t, err)

	err = core.Client.StartSubscribe(&core.Galaxy,&core.WsServer)
	require.Empty(t, err)
	time.Sleep(10*time.Second)
	core.Client.StopSubscribe()
}

func TestDeliverKeys(t *testing.T) {
	core,err := InitCore("sujiao1234")
	require.Empty(t, err)
	defer core.LvDb.Close()
	err = core.DeliverKeys(&core.Master)
	require.Empty(t, err)
	core.Galaxy.Range(func(node *Node) (next bool) {
		b,err := json.Marshal(node)
		require.Empty(t, err)
		t.Log(string(b))
		return true
	})
}

func TestLoadKeys(t *testing.T) {
	core,err := InitCore("sujiao")
	require.Empty(t, err)
	defer core.LvDb.Close()
	err = core.DeliverKeys(&core.Master)
	require.Empty(t, err)
	maxPath2,err := core.GetMaxPath2()
	require.Empty(t, err)
	chNode := core.IterNode()
	err = core.Galaxy.Load(maxPath2,&core.Master,chNode)
	require.Empty(t, err)
	core.Galaxy.Range(func(node *Node) (next bool) {
		b,err := json.Marshal(node)
		require.Empty(t, err)
		t.Log(string(b))
		return true
	})
}



func TestMaxGasLimit(t *testing.T) {
	core,err := InitCore("sujiao")
	require.Empty(t, err)
	defer core.LvDb.Close()
	err = core.DeliverKeys(&core.Master)
	require.Empty(t, err)
	maxPath2,err := core.GetMaxPath2()
	require.Empty(t, err)
	chNode := core.IterNode()
	err = core.Galaxy.Load(maxPath2,&core.Master,chNode)
	require.Empty(t, err)
	gasUsed := uint64(0)
	core.Range(func(node *Node) (next bool) {
		step := node.GetStep()
		if step != nil {
			receipt,_ := core.Client.TransactionReceipt(context.Background(),step.TxHash)
			if receipt != nil {
				t.Logf("Reward tx: %v, gas used: %v",receipt.TxHash.String(),receipt.GasUsed)
				if receipt.GasUsed > gasUsed {
					gasUsed = receipt.GasUsed
				}
			}
		}
		return true
	})
	t.Logf("最大gas步骤：%v",gasUsed)
}


func TestEstimatePrice(t *testing.T) {
	core,err := InitCore("1")
	require.Empty(t, err)
	defer core.LvDb.Close()
	core.Client.Last.Reserved.HasLiquidity = true
	core.Client.Last.Reserved.USDT = new(big.Int).Mul(big.NewInt(858200),big.NewInt(1e18))
	core.Client.Last.Reserved.SAT = new(big.Int).Mul(big.NewInt(85820000),big.NewInt(1e18))
	out := core.Client.GetUsdtAmountOut(new(big.Int).Mul(big.NewInt(10000),big.NewInt(1e18)))
	o := new(big.Int).Div(out,big.NewInt(1e18)).Int64()
	usdt,_ := new(big.Float).Quo(new(big.Float).SetInt(core.Client.Last.Reserved.USDT),big.NewFloat(1e18)).Float64()
	sat,_ := new(big.Float).Quo(new(big.Float).SetInt(core.Client.Last.Reserved.SAT),big.NewFloat(1e18)).Float64()
	price := usdt/sat
	counter := 1
	t.Logf("%v times AmountOut: %v, price: %v",counter,o,price)
	var p1,p2,p3,p4,p5 int
	for {
		counter ++
		core.Client.Last.Reserved.USDT = new(big.Int).Sub(core.Client.Last.Reserved.USDT,out)
		core.Client.Last.Reserved.SAT = new(big.Int).Add(core.Client.Last.Reserved.SAT,big.NewInt(9170))
		out = core.Client.GetUsdtAmountOut(new(big.Int).Mul(big.NewInt(10000),big.NewInt(1e18)))
		o = new(big.Int).Div(out,big.NewInt(1e18)).Int64()
		usdt,_ = new(big.Float).Quo(new(big.Float).SetInt(core.Client.Last.Reserved.USDT),big.NewFloat(1e18)).Float64()
		sat,_ = new(big.Float).Quo(new(big.Float).SetInt(core.Client.Last.Reserved.SAT),big.NewFloat(1e18)).Float64()
		price = usdt/sat
		t.Logf("%v times AmountOut: %v, price: %v",counter,o,price)
		if price < 0.009 && p1 == 0 {
			p1 = counter
		}
		if price < 0.008 && p2 == 0 {
			p2 = counter
		}
		if price < 0.007 && p3 == 0 {
			p3 = counter
		}
		if price < 0.006 && p4 == 0 {
			p4 = counter
		}
		if price < 0.005 && p5 == 0 {
			p5 = counter
		}
		if price < 0.005 {
			break
		}
	}
	t.Logf("价格降至 0.009 需要 %v 笔交易",p1)
	t.Logf("价格降至 0.008 需要 %v 笔交易",p2)
	t.Logf("价格降至 0.007 需要 %v 笔交易",p3)
	t.Logf("价格降至 0.006 需要 %v 笔交易",p4)
	t.Logf("价格降至 0.005 需要 %v 笔交易",p5)
}
