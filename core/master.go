package core

import (
	"auto-swap/config"
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/hex"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/sirupsen/logrus"
	"math/big"
)

type Master struct {
	Path1 int // 一级衍生，为了方便扩展，目前 master 服务里固定为 1000
	hdKey        *hdkeychain.ExtendedKey  // 钱包的私钥
	privKey      *ecdsa.PrivateKey
	auth         *bind.TransactOpts
	ReservedGas  *big.Int
	ReservedUSDT *big.Int
	ReservedHFT  *big.Int
}

func (master *Master) Init(pwd string) error {

	g, _  := new(big.Float).SetString(config.Cfg.Node.ReservedGas)
	master.ReservedGas, _ = g.Int(nil)

	seed,err := base64.StdEncoding.DecodeString(config.Cfg.Node.Key)
	if err != nil {
		return err
	}
	fak := FormatAESKey(pwd)
	seed,err = AESCBCDecrypt(seed,fak)
	if err != nil {
		return  err
	}
	netType := chaincfg.MainNetParams //这个参数下面只取了版本号及string格式化的参数，和链无关
	master.hdKey,err = hdkeychain.NewMaster(seed, &netType)

	master.privKey,err = crypto.ToECDSA(seed)
	if err != nil {
		return err
	}

	if crypto.PubkeyToAddress(master.privKey.PublicKey) != common.HexToAddress(config.Cfg.Node.KeyAddress) {
		return  err
	}

	master.auth, err = bind.NewKeyedTransactorWithChainID(master.privKey,big.NewInt(config.Cfg.Node.ChainId))
	if err != nil {
		return  err
	}
	master.auth.GasPrice = big.NewInt(config.Cfg.Node.GasPrice)
	return nil
}

func (master *Master) UpdateReserveConfig(client *Client, sendReq *SendReq){
	balance,ok := client.GetBalance(master.auth.From)
	if !ok ||balance == nil {
		return
	}
	min := new(big.Int).Add(sendReq.Value,master.ReservedGas)
	if balance.Cmp(min) == -1 {
		logrus.Errorf("master SendGas, balance not sufficient")
		return
	}
	sendReq.Tx = client.SendGas(master.auth,sendReq.Addr,sendReq.Value)
	close(sendReq.Down)
}

func (master *Master) DeriveAuth(path Path) (*bind.TransactOpts,string,error) {
	paths,err := path.ParsePath()
	if err != nil {
		return nil,"",err
	}
	ck := master.hdKey
	for _,p := range paths {
		hp := uint32(p) + hdkeychain.HardenedKeyStart
		ck, err = ck.Child(hp)
		if err != nil {
			return nil,"",err
		}
	}
	bek,err := ck.ECPrivKey()
	if err != nil {
		return nil,"",err
	}
	ek := (ecdsa.PrivateKey)(*bek)
	auth, err := bind.NewKeyedTransactorWithChainID(&ek,big.NewInt(config.Cfg.Node.ChainId))
	auth.GasPrice = master.auth.GasPrice
	if err != nil {
		return  nil,"",err
	}
	return auth,hex.EncodeToString(crypto.FromECDSA(&ek)),nil
}

type Coin int
const (
	Gas Coin = iota
	USDT
	HFT
)

func (coin Coin) String() string {
	switch coin {
	case Gas:
		return "BNB"
	case USDT:
		return "USDT"
	case HFT:
		return "BFT"
	default:
		return "NOT KNOW"
	}
}

type SendReq struct {
	Addr  common.Address
	Coin  Coin
	Value *big.Int
	Tx    *types.Transaction
	Down  chan struct{}
}

func (master *Master) ListenSendReq(client *Client) chan *SendReq {
	chSendReq := make(chan *SendReq,100)
	go func() {
		for sendReq := range chSendReq {
			switch sendReq.Coin {
			case  Gas: master.SendGas (client, sendReq)
			case USDT: master.SendUSDT(client, sendReq)
			case  HFT: master.SendHFT (client, sendReq)
			default: continue
			}
			close(sendReq.Down)
		}
	}()
	return chSendReq
}

func (master *Master) SendGas(client *Client, sendReq *SendReq){
	balance,ok := client.GetBalance(master.auth.From)
	if !ok ||balance == nil {
		return
	}
	min := new(big.Int).Add(sendReq.Value,master.ReservedGas)
	if balance.Cmp(min) == -1 {
		logrus.Errorf("master SendGas, balance not sufficient")
		return
	}
	sendReq.Tx = client.SendGas(master.auth,sendReq.Addr,sendReq.Value)
}

func (master *Master) SendUSDT(client *Client, sendReq *SendReq){
	balance,ok := client.GetUSDTBalance(master.auth.From)
	if !ok ||balance == nil {
		return
	}
	if balance.Cmp(sendReq.Value) == -1 {
		logrus.Errorf("master SendUSDT, balance not sufficient")
		return
	}
	sendReq.Tx = client.SendUSDT(master.auth,sendReq.Addr,sendReq.Value)
}

func (master *Master) SendHFT(client *Client, sendReq *SendReq){
	balance,ok := client.GetHFTBalance(master.auth.From)
	if !ok ||balance == nil {
		return
	}
	if balance.Cmp(sendReq.Value) == -1 {
		logrus.Errorf("master SendHFT, balance not sufficient")
		return
	}
	sendReq.Tx = client.SendHFT(master.auth,sendReq.Addr,sendReq.Value)
}
