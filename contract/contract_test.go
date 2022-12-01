package contract

import (
	"auto-swap/config"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
	"math/big"
	"net/http"
	"net/url"
	"sync"
	"testing"
)

const (
	wsReadBuffer       = 1024
	wsWriteBuffer      = 1024
)
var wsBufferPool = new(sync.Pool)

func init() {
	config.TestInit()
	u, err := url.Parse(config.Cfg.Node.Host)
	if err != nil {
		panic(err)
	}
	if u.Scheme == "ws" || u.Scheme == "wss" {
		dialer := websocket.Dialer{
			ReadBufferSize:  wsReadBuffer,
			WriteBufferSize: wsWriteBuffer,
			WriteBufferPool: wsBufferPool,
			Proxy: http.ProxyFromEnvironment,
		}
		c, err := rpc.DialWebsocketWithDialer(context.Background(),config.Cfg.Node.Host,"",dialer)
		if err != nil {
			panic(err)
		}
		client = ethclient.NewClient(c)
	}else{
		client, err = ethclient.Dial(config.Cfg.Node.Host)
		if err != nil {
			panic(err)
		}
	}

	usdt,err = NewUSDT(common.HexToAddress(config.Cfg.Contract.USDTAddr),client)
	if err != nil {
		panic(err)
	}

	invitation,err = NewInvitation(common.HexToAddress(config.Cfg.Contract.InvitationAddr),client)
	if err != nil {
		panic(err)
	}

	ido,err = NewIDO(common.HexToAddress(config.Cfg.Contract.IDOAddr),client)
	if err != nil {
		panic(err)
	}


	testPrivKey,err = crypto.HexToECDSA(testPKey)
	if err != nil {
		panic(err)
	}
	if crypto.PubkeyToAddress(testPrivKey.PublicKey) != common.HexToAddress(testAddr) {
		panic("签名和地址不一致")
	}

	test,err = bind.NewKeyedTransactorWithChainID(testPrivKey,big.NewInt(config.Cfg.Node.ChainId))
	if err != nil {
		panic(err)
	}
	test.GasPrice = big.NewInt(1e9)

	ownerPrivKey,err = crypto.HexToECDSA(ownerPKey)
	if err != nil {
		panic(err)
	}
	if crypto.PubkeyToAddress(ownerPrivKey.PublicKey) != common.HexToAddress(ownerAddr) {
		panic("签名和地址不一致")
	}

	owner,err = bind.NewKeyedTransactorWithChainID(ownerPrivKey,big.NewInt(config.Cfg.Node.ChainId))
	if err != nil {
		panic(err)
	}
	owner.GasPrice = big.NewInt(5e9)
}



const testPKey = ""
const testAddr = ""

const ownerPKey = ""
const ownerAddr = ""


var client *ethclient.Client
var test *bind.TransactOpts
var owner *bind.TransactOpts


var usdt *USDT
var invitation *Invitation
var ido *IDO

var testPrivKey *ecdsa.PrivateKey
var ownerPrivKey *ecdsa.PrivateKey

func TestChainInfo(t *testing.T) {
	number := int64(5989940)
	header,err := client.HeaderByNumber(context.Background(),big.NewInt(number))
	require.Empty(t, err)
	t.Logf("number:%v, time:%v, difficult:%v", number, header.Time, header.Difficulty.String())
}

func TestBalanceAt(t *testing.T) {
	balance,err := client.BalanceAt(context.Background(), common.HexToAddress(testAddr), nil)
	require.Empty(t, err)
	t.Log(balance.String())
}

func TestGetInvitation(t *testing.T) {
	zero := common.Address{}
	info,err := invitation.GetInvitation(nil,common.HexToAddress("0x8193f3bcF8E3396de72E8FcdA8253de6d3f93E04"))
	require.Empty(t, err)
	t.Log(info.Inviter != zero)
	b,err := json.Marshal(&info)
	t.Log(string(b))
}

func TestBind(t *testing.T) {
	tx,err := invitation.Bind(owner,common.HexToAddress(testAddr))
	require.Empty(t, err)
	t.Log(tx.Hash().String())
}

func TestApprove(t *testing.T) {
	tx,err := usdt.Approve(owner, common.HexToAddress(config.Cfg.Contract.IDOAddr),new(big.Int).Mul(big.NewInt(1e8),big.NewInt(1e18)))
	require.Empty(t, err)
	t.Log(tx.Hash())
}

func TestInvest1(t *testing.T) {
	tx,err := ido.Invest1(owner)
	require.Empty(t, err)
	t.Log(tx.Hash().String())
}

func TestUserInfo(t *testing.T) {
	info,err := ido.UserInfo(nil,common.HexToAddress("0x8193f3bcF8E3396de72E8FcdA8253de6d3f93E04"))
	require.Empty(t, err)
	b,err := json.Marshal(&info)
	t.Log(string(b))
}

func TestClaim(t *testing.T) {
	tx,err := ido.Claim(owner,common.HexToAddress(ownerAddr))
	require.Empty(t, err)
	t.Log(tx.Hash().String())
}

