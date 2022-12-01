package core

import (
	"auto-swap/config"
	"auto-swap/contract"
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	wsReadBuffer       = 1024
	wsWriteBuffer      = 1024
	MaxWait            = 5 * time.Second
)
var wsBufferPool = new(sync.Pool)

type Client struct {
	*ethclient.Client
	rpc        *rpc.Client
	usdt       *contract.USDT
	hft        *contract.SAT
	invitation *contract.Invitation
	router     *contract.Router
	pair       *contract.Pair
	rateLimit   chan struct{}
	rateTicker *time.Ticker
	exit        chan struct{}
	gasPrice   *big.Int
	sub struct{
		sub    ethereum.Subscription
		subWg  sync.WaitGroup
		cancel context.CancelFunc
		//closed bool
	}
	Last struct{
		Header *types.Header
		Reserved struct{
			HasLiquidity bool
			USDT *big.Int // reserve0
			SAT *big.Int  // reserve1
		}
	}
}

func (client *Client)Init(exit chan struct{}) error {

	client.exit = exit
	client.gasPrice = big.NewInt(config.Cfg.Node.GasPrice)

	u, err := url.Parse(config.Cfg.Node.Host)
	if err != nil {
		return err
	}
	ctx,_ := context.WithTimeout(context.Background(),config.MaxWait)
	if u.Scheme == "ws" || u.Scheme == "wss" {
		dialer := websocket.Dialer{
			ReadBufferSize:  wsReadBuffer,
			WriteBufferSize: wsWriteBuffer,
			WriteBufferPool: wsBufferPool,
			Proxy: http.ProxyFromEnvironment,
		}
		client.rpc, err = rpc.DialWebsocketWithDialer(ctx,config.Cfg.Node.Host,"",dialer)
		if err != nil {
			return err
		}
	}else{
		client.rpc, err = rpc.DialContext(ctx, config.Cfg.Node.Host)
		if err != nil {
			return err
		}
	}

	client.Client = ethclient.NewClient(client.rpc)

	client.usdt,err = contract.NewUSDT(common.HexToAddress(config.Cfg.Contract.USDTAddr),client.Client)
	if err != nil {
		return err
	}

	client.hft,err = contract.NewSAT(common.HexToAddress(config.Cfg.Contract.SATAddr),client.Client)
	if err != nil {
		return err
	}

	client.invitation,err = contract.NewInvitation(common.HexToAddress(config.Cfg.Contract.InvitationAddr),client.Client)
	if err != nil {
		return err
	}

	client.router,err = contract.NewRouter(common.HexToAddress(config.Cfg.Contract.RouterAddr),client.Client)
	if err != nil {
		return err
	}

	client.pair,err = contract.NewPair(common.HexToAddress(config.Cfg.Contract.PairAddr),client.Client)
	if err != nil {
		return err
	}

	client.initRateLimit()

	return nil
}

func (client *Client) initRateLimit() {
	client.rateLimit = make(chan struct{}, config.Cfg.Node.RateLimit)
	var fullLimit = func() {
		for i := 0; i < config.Cfg.Node.RateLimit; i++ {
			select {
			case client.rateLimit <- struct{}{}:
			default:
			}
		}
	}
	client.rateTicker = time.NewTicker(time.Second)
	go func() {
		defer func() {
			close(client.rateLimit)
			for range client.rateLimit {
			}
		}()
		fullLimit()
		for {
			select {
			case <-client.rateTicker.C:
				fullLimit()
			case <-client.exit:
				return
			}
		}
	}()
}

func (client *Client) StopSubscribe() {
	logrus.Info("关闭订阅服务...")
	client.sub.sub.Unsubscribe()
	logrus.Info("......")
	client.sub.subWg.Wait()
	logrus.Info("订阅服务关闭成功...")
}

func (client *Client) StartSubscribe(galaxy *Galaxy,ws *WsServer) error {
	var err error
	logrus.Info("开始订阅最新 Header...")
	ch := make(chan *types.Header, 100)
	ctx, cancel := context.WithTimeout(context.Background(), config.MaxWait)
	client.sub.sub, err = client.SubscribeNewHead(ctx, ch)
	cancel()
	if err != nil {
		return err
	}
	reconnect := false
	go func() {
		client.sub.subWg.Add(1)
		defer func() {
			logrus.Info("结束订阅随 Header...")
			client.sub.subWg.Done()
			if reconnect {
				reconnect = false
				logrus.Info("订阅异常，一秒后开始重新订阅")
			RETRY:
				time.Sleep(time.Second)
				err := client.StartSubscribe(galaxy,ws)
				if err != nil {
					logrus.Errorf("Recursive SubscribeHeaders err: %v", err)
					goto RETRY
				}
			}
		}()
		for header := range ch {
			var price float64
			for i := 3; i > 0; i -- {
				var ok bool
				client.Last.Reserved.USDT, client.Last.Reserved.SAT, ok = client.GetReserve()
				if ok {
					if client.Last.Reserved.SAT != nil &&
						client.Last.Reserved.SAT.Cmp(big.NewInt(0)) ==1 &&
						client.Last.Reserved.USDT != nil &&
						client.Last.Reserved.USDT.Cmp(big.NewInt(0)) ==1{
						client.Last.Reserved.HasLiquidity = true
						usdt,_ := new(big.Float).Quo(new(big.Float).SetInt(client.Last.Reserved.USDT),big.NewFloat(1e18)).Float64()
						sat,_ := new(big.Float).Quo(new(big.Float).SetInt(client.Last.Reserved.SAT),big.NewFloat(1e18)).Float64()
						price = usdt/sat
					}else{
						client.Last.Reserved.HasLiquidity = false
					}
					break
				}
			}
			client.Last.Header = header
			logrus.Infof("New header number: %v, Reserve SAT: %v, USDT: %v, Price: %v",
				client.Last.Header.Number.String(),
				new(big.Int).Div(client.Last.Reserved.SAT,big.NewInt(1e18)),
				new(big.Int).Div(client.Last.Reserved.USDT,big.NewInt(1e18)),
				price)
		}
	}()
	go func() {
		logrus.Info("开始监听错误信息...")
		client.sub.subWg.Add(1)
		defer func() {
			logrus.Info("结束监听错误信息...")
			close(ch)
			client.sub.subWg.Done()
		}()
		for err := range client.sub.sub.Err() {
			logrus.Errorf("ListenErr sub err: %v", err)
			reconnect = true
			client.sub.sub.Unsubscribe()
		}
	}()
	return nil
}


func (client *Client) isOk() bool {
	_, ok := <- client.rateLimit
	return ok
}

func (client *Client) GetTx(txHash common.Hash) (tx *types.Transaction, isPending, ok bool) {
	if !client.isOk() {
		return nil, false,false
	}
	ctx1,cancel1 := context.WithTimeout(context.Background(), MaxWait)
	defer cancel1()
	tx,isPending,err := client.TransactionByHash(ctx1,txHash)
	if err != nil {
		logrus.Errorf("client GetTxResult TransactionReceipt txHash: %v, err: %v",txHash.String(),err)
		return nil, false,false
	}
	return tx, isPending, true
}

func (client *Client) GetTxResult(txHash common.Hash) (success, ok bool) {
	if !client.isOk() {
		return false,false
	}
	ctx1,cancel1 := context.WithTimeout(context.Background(), MaxWait)
	defer cancel1()
	receipt,err := client.TransactionReceipt(ctx1,txHash)
	if err != nil {
		if !strings.Contains(err.Error(),"not found") {
			logrus.Errorf("client GetTxResult TransactionReceipt txHash: %v, err: %v",txHash.String(),err)
		}
		return false,false
	}
	if receipt.Status == types.ReceiptStatusSuccessful {
		return true,true
	}else{
		return false,true
	}
}

func (client *Client) GetBalance(addr common.Address) (*big.Int,bool) {
	if !client.isOk() {
		return nil, false
	}
	ctx1,cancel1 := context.WithTimeout(context.Background(), MaxWait)
	defer cancel1()
	b,err := client.PendingBalanceAt(ctx1,addr)
	if err != nil {
		logrus.Errorf("client GetBalance BalanceAt addr: %v, err: %v",addr.String(),err)
		return nil, false
	}
	return b, true
}


func (client *Client) SendGas(auth *bind.TransactOpts,to common.Address, value *big.Int) (tx *types.Transaction) {
	if !client.isOk() {
		return nil
	}
	ctx1,cancel1 := context.WithTimeout(context.Background(), MaxWait)
	defer cancel1()
	// create the transaction
	nonce, err := client.PendingNonceAt(ctx1,auth.From)
	if err != nil {
		logrus.Errorf("client SendGas PendingNonceAt from: %v, err: %v",auth.From.String(),err)
		return nil
	}
	baseTx := &types.LegacyTx{
		To:       &to,
		Nonce:    nonce,
		GasPrice: client.gasPrice,
		Gas:      21000,
		Value:    value,
	}

	tx = types.NewTx(baseTx)
	tx,err = auth.Signer(auth.From,tx)
	if err != nil {
		logrus.Errorf("client SendGas Signer from: %v, err: %v",auth.From.String(),err)
		return nil
	}
	ctx2,cancel2 := context.WithTimeout(context.Background(), MaxWait)
	defer cancel2()
	err = client.SendTransaction(ctx2,tx)
	if err != nil {
		logrus.Errorf("client SendGas SendTransaction from: %v, err: %v",auth.From.String(),err)
		return nil
	}
	return
}


func (client *Client) GetUSDTBalance(addr common.Address) (*big.Int,bool) {
	if !client.isOk() {
		return nil,false
	}
	b,err := client.usdt.BalanceOf(nil,addr)
	if err != nil {
		logrus.Errorf("client GetUSDTBalance BalanceOf addr: %v, err: %v",addr.String(),err)
		return nil,false
	}
	return b,true
}

func (client *Client) IsUSDTApproved(owner,spender common.Address) (approved bool, ok bool) {
	if !client.isOk() {
		return false,false
	}
	b,err := client.usdt.Allowance(nil,owner,spender)
	if err != nil {
		logrus.Errorf("client IsUSDTApproved Allowance addr: %v, err: %v",owner.String(),err)
		return false,false
	}
	// a < b
	if b == nil || b.Cmp(nodeParams.NeedUSDT) == -1 {
		return false,true
	}else{
		return true,true
	}

}

func (client *Client) ApproveUSDT(auth *bind.TransactOpts,spender common.Address, amount *big.Int) (needGas bool,tx *types.Transaction,ok bool) {
	ok = client.isOk()
	if !ok {
		return false,nil,false
	}
	var err error
	tx,err = client.usdt.Approve(auth, spender, amount)
	if err != nil {
		if strings.Contains(err.Error(),"insufficient") {
			return true,nil,false
		}
		logrus.Errorf("client ApproveUSDT owner: %v, spender: %v, err: %v",auth.From.String(),spender.String(),err)
		return false,nil,false
	}
	return
}

func (client *Client) SendUSDT(auth *bind.TransactOpts, to common.Address, amount *big.Int) (tx *types.Transaction){
	ok := client.isOk()
	if !ok {
		return nil
	}
	var err error
	tx,err = client.usdt.Transfer(auth, to, amount)
	if err != nil {
		logrus.Errorf("client SendUSDT from: %v, to: %v, amount: %v, err: %v",
			auth.From.String(),to.String(),amount.String(),err)
		return nil
	}
	return
}

func (client *Client) GetHFTBalance(addr common.Address) (*big.Int,bool) {
	if !client.isOk() {
		return nil,false
	}
	b,err := client.hft.BalanceOf(nil,addr)
	if err != nil {
		logrus.Errorf("client GetHFTBalance BalanceOf addr: %v, err: %v",addr.String(),err)
		return nil,false
	}
	return b,true
}

func (client *Client) IsHFTApproved(owner,spender common.Address) (approved bool, ok bool) {
	if !client.isOk() {
		return false,false
	}
	b,err := client.hft.Allowance(nil,owner,spender)
	if err != nil {
		logrus.Errorf("client IsSATApproved Allowance addr: %v, err: %v",owner.String(),err)
		return false,false
	}
	// a < b
	if b == nil || b.Cmp(nodeParams.NeedUSDT) == -1 {
		return false,true
	}else{
		return true,true
	}

}

func (client *Client) ApproveHFT(auth *bind.TransactOpts,spender common.Address, amount *big.Int) (needGas bool, tx *types.Transaction,ok bool) {
	ok = client.isOk()
	if !ok {
		return false,nil,false
	}
	var err error
	tx,err = client.hft.Approve(auth, spender, amount)
	if err != nil {
		if strings.Contains(err.Error(),"insufficient") {
			return true,nil,false
		}
		logrus.Errorf("client ApproveSAT owner: %v, spender: %v, err: %v",auth.From.String(),spender.String(),err)
		return false,nil,false
	}
	return
}

func (client *Client) SendHFT(auth *bind.TransactOpts, to common.Address, amount *big.Int) (tx *types.Transaction){
	ok := client.isOk()
	if !ok {
		return nil
	}
	var err error
	tx,err = client.hft.Transfer(auth, to, amount)
	if err != nil {
		logrus.Errorf("client SendHFT from: %v, to: %v, amount: %v, err: %v",
			auth.From.String(),to.String(),amount.String(),err)
		return nil
	}
	return
}

func (client *Client) GetInviter(user common.Address) (inviter common.Address, ok bool) {
	ok = client.isOk()
	if !ok {
		return
	}
	info,err := client.invitation.GetInvitation(nil,user)
	if err != nil {
		logrus.Errorf("client GetInviter user: %v, err: %v",user.String(),err)
		return common.Address{},false
	}
	return info.Inviter,true
}

func (client *Client) Bind(auth *bind.TransactOpts,inviter common.Address) (needTrans StepType,tx *types.Transaction,ok bool) {
	ok = client.isOk()
	if !ok {
		return DefaultST,nil,false
	}
	var err error
	tx,err = client.invitation.Bind(auth,inviter)
	if err != nil {
		if strings.Contains(err.Error(),"insufficient") {
			return TransferGas,nil,false
		}
		logrus.Errorf("client Bind user: %v, inviter: %v, err: %v",auth.From.String(),inviter.String(),err)
		return DefaultST,nil,false
	}
	return
}

func (client *Client) GetReserve() (reserve0,reserve1 *big.Int, ok bool) {
	ok = client.isOk()
	if !ok {
		return
	}
	info,err := client.pair.GetReserves(nil)
	if err != nil {
		logrus.Errorf("client GetReserve err: %v",err)
		return nil,nil,false
	}
	return info.Reserve0, info.Reserve1, true
}

func (client *Client) GetUsdtAmountOut(satIn *big.Int) (usdtOut *big.Int) {
	usdtOut = big.NewInt(0)
	if satIn.Cmp(usdtOut) != 1 {
		return
	}
	if !client.Last.Reserved.HasLiquidity {
		return
	}
	amountInWithFee := new(big.Int).Mul(satIn,big.NewInt(1000-3-80)) // 0.3%是swap的手续费，8%是sat的手续费
	numerator := new(big.Int).Mul(amountInWithFee,client.Last.Reserved.USDT)
	denominator :=  new(big.Int).Mul(client.Last.Reserved.SAT,big.NewInt(1000))
	denominator = denominator.Add(denominator,amountInWithFee)
	usdtOut = usdtOut.Div(numerator, denominator)
	return
}


func (client *Client) GetHFTAmountOut(usdtIn *big.Int) (hftOut *big.Int) {
	hftOut = big.NewInt(0)
	if usdtIn.Cmp(hftOut) != 1 {
		return
	}
	if !client.Last.Reserved.HasLiquidity {
		return
	}
	amountInWithFee := new(big.Int).Mul(usdtIn,big.NewInt(1000-3)) // 0.3%是swap的手续费
	numerator := new(big.Int).Mul(amountInWithFee,client.Last.Reserved.SAT)
	denominator :=  new(big.Int).Mul(client.Last.Reserved.USDT,big.NewInt(1000))
	denominator = denominator.Add(denominator,amountInWithFee)
	hftOut = hftOut.Div(numerator, denominator)
	hftOut = hftOut.Mul(hftOut,big.NewInt(92)) // 合约会收8%的手续费
	hftOut = hftOut.Div(hftOut,big.NewInt(100))
	return
}

func (client *Client) IsBuyExecuted(account common.Address,expect *big.Int) (hftBalance *big.Int,executed bool, ok bool) {
	if !client.isOk() {
		return nil,false,false
	}
	if expect == nil {
		return nil,false,true
	}
	b,err := client.hft.BalanceOf(nil,account)
	if err != nil {
		logrus.Errorf("client IsBuyExecuted BalanceOf addr: %v, err: %v",account.String(),err)
		return nil,false,false
	}
	// a < b
	if b == nil || b.Cmp(expect) == -1 {
		return b,false,true
	}else{
		return b,true,true
	}

}

func (client *Client) BuyHFT(auth *bind.TransactOpts, usdtAmount, hftBalance *big.Int) (needTrans StepType,tx *types.Transaction,expect *big.Int,ok bool) {
	ok = client.isOk()
	if !ok {
		return DefaultST,nil, nil,false
	}
	ub,err := client.usdt.BalanceOf(nil,auth.From)
	if err != nil {
		logrus.Errorf("client BuyHFT BalanceOf usdt addr: %v, err: %v",auth.From.String(),err)
		return DefaultST,nil, nil,false
	}
	if ub.Cmp(usdtAmount) == -1 {
		return TransferUSDT,nil,nil,false
	}

	hftAmount := client.GetHFTAmountOut(usdtAmount)
	hftAmount = hftAmount.Mul(hftAmount,big.NewInt(100-config.Cfg.Node.Slippage))
	hftAmount = hftAmount.Div(hftAmount,big.NewInt(100))

	if hftBalance == nil {
		expect = hftAmount
	}else{
		expect = new(big.Int).Add(hftBalance,hftAmount)
	}

	tx,err = client.router.SwapExactTokensForTokensSupportingFeeOnTransferTokens(
		auth,usdtAmount,hftAmount,[]common.Address{
			common.HexToAddress(config.Cfg.Contract.USDTAddr),
			common.HexToAddress(config.Cfg.Contract.SATAddr)},
		auth.From,big.NewInt(int64(client.Last.Header.Time)+int64(time.Hour)))
	if err != nil {
		if strings.Contains(err.Error(),"insufficient") {
			return TransferGas,nil,nil,false
		}
		logrus.Errorf("client BuyHFT user: %v, err: %v",auth.From.String(),err)
		return DefaultST,nil,nil,false
	}
	return
}

func (client *Client) RefreshBuyHFT(auth *bind.TransactOpts, tx *types.Transaction, usdtAmount , hftBalance *big.Int) (newTx *types.Transaction,expect *big.Int,ok bool) {
	ok = client.isOk()
	if !ok {
		return nil,nil,false
	}
	// a >= b 已到最高gas费，不需要再刷新
	if tx.GasPrice().Cmp(big.NewInt(config.Cfg.Node.MaxPrice)) != -1 {
		return nil,nil,false
	}

	hftAmount := client.GetHFTAmountOut(usdtAmount)
	hftAmount = hftAmount.Mul(hftAmount,big.NewInt(100-config.Cfg.Node.Slippage))
	hftAmount = hftAmount.Div(hftAmount,big.NewInt(100))

	if hftBalance == nil {
		expect = hftAmount
	}else{
		expect = new(big.Int).Add(hftBalance,hftAmount)
	}

	auth.Nonce = big.NewInt(int64(tx.Nonce()))
	auth.GasPrice = new(big.Int).Add(tx.GasPrice(),big.NewInt(config.Cfg.Node.StepPrice))
	defer func() {
		auth.Nonce = nil
		auth.GasPrice = client.gasPrice
	}()
	var err error
	newTx,err = client.router.SwapExactTokensForTokensSupportingFeeOnTransferTokens(
		auth,usdtAmount,hftAmount,[]common.Address{
			common.HexToAddress(config.Cfg.Contract.USDTAddr),
			common.HexToAddress(config.Cfg.Contract.SATAddr)},
		auth.From,big.NewInt(int64(client.Last.Header.Time)+int64(time.Hour)))
	if err != nil {
		logrus.Errorf("client RefreshProfit user: %v, err: %v",auth.From.String(),err)
		return nil,nil,false
	}
	return
}


func (client *Client) IsSellExecuted(account common.Address,expect *big.Int) (usdtBalance *big.Int,executed bool, ok bool) {
	if !client.isOk() {
		return nil,false,false
	}
	if expect == nil {
		return nil,false,true
	}
	b,err := client.usdt.BalanceOf(nil,account)
	if err != nil {
		logrus.Errorf("client IsSellExecuted BalanceOf addr: %v, err: %v",account.String(),err)
		return nil,false,false
	}
	// a < b
	if b == nil || b.Cmp(expect) == -1 {
		return b,false,true
	}else{
		return b,true,true
	}

}

func (client *Client) SellHFT(auth *bind.TransactOpts, hftAmount, usdtBalance *big.Int) (needTrans StepType,tx *types.Transaction,expect *big.Int,ok bool) {
	ok = client.isOk()
	if !ok {
		return DefaultST,nil,nil,false
	}

	hb,err := client.hft.BalanceOf(nil,auth.From)
	if err != nil {
		logrus.Errorf("client SellHFT BalanceOf HFT addr: %v, err: %v",auth.From.String(),err)
		return DefaultST,nil, nil,false
	}
	if hb.Cmp(hftAmount) == -1 {
		return TransferHFT,nil,nil,false
	}

	usdtAmount := client.GetUsdtAmountOut(hftAmount)
	usdtAmount = usdtAmount.Mul(usdtAmount,big.NewInt(100-config.Cfg.Node.Slippage))
	usdtAmount = usdtAmount.Div(usdtAmount,big.NewInt(100))
	if usdtAmount.Cmp(nodeParams.MinSwapUsdt) != 1 {
		return DefaultST,nil, nil,false
	}

	if usdtBalance == nil {
		expect = usdtAmount
	}else{
		expect = new(big.Int).Add(usdtBalance,usdtAmount)
	}

	tx,err = client.router.SwapExactTokensForTokensSupportingFeeOnTransferTokens(
		auth,hftAmount,usdtAmount,[]common.Address{
			common.HexToAddress(config.Cfg.Contract.SATAddr),
			common.HexToAddress(config.Cfg.Contract.USDTAddr)},
		auth.From,big.NewInt(int64(client.Last.Header.Time)+int64(time.Hour)))
	if err != nil {
		if strings.Contains(err.Error(),"insufficient") {
			return TransferGas,nil,nil,false
		}
		logrus.Errorf("client Profit user: %v, err: %v",auth.From.String(),err)
		return DefaultST,nil, nil,false
	}
	return
}

func (client *Client) RefreshSellHFT(auth *bind.TransactOpts, tx *types.Transaction, hftAmount, usdtBalance *big.Int) (newTx *types.Transaction,expect *big.Int,ok bool) {
	ok = client.isOk()
	if !ok {
		return nil,nil,false
	}
	// a >= b 已到最高gas费，不需要再刷新
	if tx.GasPrice().Cmp(big.NewInt(config.Cfg.Node.MaxPrice)) != -1 {
		return nil,nil,false
	}

	usdtAmount := client.GetUsdtAmountOut(hftAmount)
	usdtAmount = usdtAmount.Mul(usdtAmount,big.NewInt(100-config.Cfg.Node.Slippage))
	usdtAmount = usdtAmount.Div(usdtAmount,big.NewInt(100))
	if usdtAmount.Cmp(nodeParams.MinSwapUsdt) != 1 {
		return nil,nil,false
	}

	if usdtBalance == nil {
		expect = usdtAmount
	}else{
		expect = new(big.Int).Add(usdtBalance,usdtAmount)
	}

	auth.Nonce = big.NewInt(int64(tx.Nonce()))
	auth.GasPrice = new(big.Int).Add(tx.GasPrice(),big.NewInt(config.Cfg.Node.StepPrice))
	defer func() {
		auth.Nonce = nil
		auth.GasPrice = client.gasPrice
	}()
	var err error
	newTx,err = client.router.SwapExactTokensForTokensSupportingFeeOnTransferTokens(
		auth,hftAmount,usdtAmount,[]common.Address{
			common.HexToAddress(config.Cfg.Contract.SATAddr),
			common.HexToAddress(config.Cfg.Contract.USDTAddr)},
		auth.From,big.NewInt(int64(client.Last.Header.Time)+int64(time.Hour)))
	if err != nil {
		logrus.Errorf("client RefreshProfit user: %v, err: %v",auth.From.String(),err)
		return nil,nil,false
	}
	return
}

func (client *Client) IsUsdtCollected(account common.Address,expect *big.Int) (usdtBalance *big.Int,collected bool, ok bool) {
	if !client.isOk() {
		return nil,false,false
	}
	b,err := client.usdt.BalanceOf(nil,account)
	if err != nil {
		logrus.Errorf("client IsUsdtCollected BalanceOf addr: %v, err: %v",account.String(),err)
		return nil,false,false
	}
	if b == nil {
		b = big.NewInt(0)
	}
	if b.IsInt64() && b.Int64() == 0 {
		return b,true,true // 如果没有余额，认为归集成功了
	}
	if expect == nil {
		return b,false,true
	}
	// l > r
	if b.Cmp(expect) == 1 {
		return b,false,true
	}else{
		return b,true,true
	}

}

func (client *Client) CollectUsdt(auth *bind.TransactOpts, recipient common.Address, amount *big.Int) (needTrans StepType,tx *types.Transaction,ok bool) {
	ok = client.isOk()
	if !ok {
		return DefaultST,nil,false
	}

	var err error
	tx,err = client.usdt.Transfer(auth,recipient,amount)
	if err != nil {
		if strings.Contains(err.Error(),"insufficient") {
			return TransferGas,nil,false
		}
		logrus.Errorf("client CollectUsdt user: %v, err: %v",auth.From.String(),err)
		return DefaultST,nil,false
	}
	return DefaultST,tx,true
}

func (client *Client) IsBftCollected(account common.Address,expect *big.Int) (bftBalance *big.Int,collected bool, ok bool) {
	if !client.isOk() {
		return nil,false,false
	}
	b,err := client.hft.BalanceOf(nil,account)
	if err != nil {
		logrus.Errorf("client IsBftCollected BalanceOf addr: %v, err: %v",account.String(),err)
		return nil,false,false
	}
	if b == nil {
		b = big.NewInt(0)
	}
	if b.IsInt64() && b.Int64() == 0 {
		return b,true,true // 如果没有余额，认为归集成功了
	}
	if expect == nil {
		return b,false,true
	}
	// l > r
	if b.Cmp(expect) == 1 {
		return b,false,true
	}else{
		return b,true,true
	}

}

func (client *Client) CollectBft(auth *bind.TransactOpts, recipient common.Address, amount *big.Int) (needTrans StepType,tx *types.Transaction,ok bool) {
	ok = client.isOk()
	if !ok {
		return DefaultST,nil,false
	}

	var err error
	tx,err = client.hft.Transfer(auth,recipient,amount)
	if err != nil {
		if strings.Contains(err.Error(),"insufficient") {
			return TransferGas,nil,false
		}
		logrus.Errorf("client CollectBft user: %v, err: %v",auth.From.String(),err)
		return DefaultST,nil,false
	}
	return DefaultST,tx,true
}

func (client *Client) IsBnbCollected(auth *bind.TransactOpts) (bnbBalance *big.Int,collected bool, ok bool) {
	if !client.isOk() {
		return nil,false,false
	}
	b,err := client.BalanceAt(context.Background(),auth.From,nil)
	if err != nil {
		logrus.Errorf("client IsBftCollected BalanceOf addr: %v, err: %v",auth.From.String(),err)
		return nil,false,false
	}
	if b == nil {
		b = big.NewInt(0)
	}
	// 转账需要的手续费，如果小于该手续费直接认为归集成功了，不需要执行归集动作
	fee := new(big.Int).Mul(client.gasPrice,big.NewInt(21000))
	b = b.Sub(b,fee)
	// l > r
	if b.Cmp(big.NewInt(0)) == 1 {
		return b,false,true
	}else{
		return b,true,true
	}

}

func (client *Client) CollectBnb(auth *bind.TransactOpts, recipient common.Address,amount *big.Int) (tx *types.Transaction,ok bool) {
	ok = client.isOk()
	if !ok {
		return nil,false
	}
	tx = client.SendGas(auth,recipient,amount)
	return tx,true
}