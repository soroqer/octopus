package test

import (
	"auto-swap/config"
	"auto-swap/core"
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"math/big"
	"net/http"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"testing"
	"time"
)

func init() {
	config.TestInit()
	var err error
	key,err = hex.DecodeString("98bcce261e4e88c9f0fc128a3aca48da9ce085f4651971857424badcb4fd8853")
	if err != nil {
		panic(err)
	}
}

var key []byte

const HOST = "localhost:5002"

func TestClient(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	c, _, err := websocket.Dial(ctx, fmt.Sprintf("ws://%v/connect",HOST), nil)
	require.Empty(t, err)
	defer func() {
		_ = c.Close(websocket.StatusInternalError, "the sky is falling")
	}()

	err = wsjson.Write(ctx, c, "hi")
	require.Empty(t, err)

	err = c.Close(websocket.StatusNormalClosure, "")
	require.Empty(t, err)
}

func TestGetSet(t *testing.T) {

	res,err := http.Get(fmt.Sprintf("http://%v/set",HOST))
	require.Empty(t, err)

	defer func() {
		_ = res.Body.Close()
	}()

	b,err := ioutil.ReadAll(res.Body)
	require.Empty(t, err)
	hb,err := hex.DecodeString(string(b))
	require.Empty(t, err)
	db,err := core.AESCBCDecrypt(hb,key)
	require.Empty(t, err)

	t.Log(string(db))
}

func TestSetSlippage(t *testing.T) {

	req := fmt.Sprintf("%v|%v=%v",1,"slippage",30)
	eb,err := core.AESCBCEncrypt([]byte(req),key)
	require.Empty(t, err)

	res,err := http.Post(fmt.Sprintf("http://%v/set",HOST),
		"text/plain",
		bytes.NewReader(eb))
	require.Empty(t, err)

	defer func() {
		_ = res.Body.Close()
	}()

	t.Log(res.StatusCode)
	b,err := ioutil.ReadAll(res.Body)
	require.Empty(t, err)
	if res.StatusCode/100 == 2 {
		db,err := core.AESCBCDecrypt(b,key)
		require.Empty(t, err)

		t.Log(string(db))
	}
}

func TestSetAddNode(t *testing.T) {

	req := fmt.Sprintf("%v|%v=%v",19,"add_node",10)
	eb,err := core.AESCBCEncrypt([]byte(req),key)
	require.Empty(t, err)
	heb := []byte(hex.EncodeToString(eb))

	res,err := http.Post(fmt.Sprintf("http://%v/set",HOST),
		"text/plain",
		bytes.NewReader(heb))
	require.Empty(t, err)

	defer func() {
		_ = res.Body.Close()
	}()

	t.Log(res.StatusCode)

	hb,err := ioutil.ReadAll(res.Body)
	require.Empty(t, err)

	if res.StatusCode/100 == 2 {
		b,err := hex.DecodeString(string(hb))
		require.Empty(t, err)
		db,err := core.AESCBCDecrypt(b,key)
		require.Empty(t, err)

		t.Log(string(db))
	}
}

func TestNodeInfo(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	ts := fmt.Sprintf("%v",time.Now().UnixNano()/1e6)
	bts,err := core.AESCBCEncrypt([]byte(ts),key)
	require.Empty(t, err)
	hbts := hex.EncodeToString(bts)

	c, res, err := websocket.Dial(ctx, fmt.Sprintf("ws://%v/connect?%v",HOST,hbts), nil)
	require.Empty(t, err)
	t.Log(res.StatusCode)
	t.Log(res.Body)
	defer func() {
		_ = c.Close(websocket.StatusInternalError, "the sky is falling")
	}()
	for {
		msgType,reader,err := c.Reader(context.Background())
		require.Empty(t, err)
		require.Equal(t, websocket.MessageText,msgType)
		b,err := ioutil.ReadAll(reader)
		require.Empty(t, err)
		bb,err := hex.DecodeString(string(b))
		require.Empty(t, err)
		db,err := core.AESCBCDecrypt(bb,key)
		require.Empty(t, err)
		t.Log(string(db))
	}
}

func TestSetInitNode(t *testing.T) {

	req := fmt.Sprintf("%v|%v,%v,%v",23,"m-1001-10","init",0)
	eb,err := core.AESCBCEncrypt([]byte(req),key)
	require.Empty(t, err)
	heb := []byte(hex.EncodeToString(eb))

	res,err := http.Post(fmt.Sprintf("http://%v/cmd",HOST),
		"text/plain",
		bytes.NewReader(heb))
	require.Empty(t, err)

	defer func() {
		_ = res.Body.Close()
	}()

	t.Log(res.StatusCode)

	hb,err := ioutil.ReadAll(res.Body)
	require.Empty(t, err)

	if res.StatusCode/100 == 2 {
		b,err := hex.DecodeString(string(hb))
		require.Empty(t, err)
		db,err := core.AESCBCDecrypt(b,key)
		require.Empty(t, err)

		t.Log(string(db))
	}
}

func TestBuyHFT(t *testing.T) {

	req := fmt.Sprintf("%v|%v,%v,%v",63,"m-1001-4","buy",30)
	eb,err := core.AESCBCEncrypt([]byte(req),key)
	require.Empty(t, err)
	heb := []byte(hex.EncodeToString(eb))

	res,err := http.Post(fmt.Sprintf("http://%v/cmd",HOST),
		"text/plain",
		bytes.NewReader(heb))
	require.Empty(t, err)

	defer func() {
		_ = res.Body.Close()
	}()

	hb,err := ioutil.ReadAll(res.Body)
	require.Empty(t, err)

	if res.StatusCode/100 == 2 {

		b,err := hex.DecodeString(string(hb))
		require.Empty(t, err)
		db,err := core.AESCBCDecrypt(b,key)
		require.Empty(t, err)

		t.Log(string(db))
	}else{
		t.Log(string(hb))
	}
}

func TestUpdateBalance(t *testing.T) {

	req := fmt.Sprintf("%v|%v,%v,%v",61,"m-1001-4","balance",100)
	eb,err := core.AESCBCEncrypt([]byte(req),key)
	require.Empty(t, err)
	heb := []byte(hex.EncodeToString(eb))

	res,err := http.Post(fmt.Sprintf("http://%v/cmd",HOST),
		"text/plain",
		bytes.NewReader(heb))
	require.Empty(t, err)

	defer func() {
		_ = res.Body.Close()
	}()

	t.Log(res.StatusCode)

	hb,err := ioutil.ReadAll(res.Body)
	require.Empty(t, err)

	if res.StatusCode/100 == 2 {

		b,err := hex.DecodeString(string(hb))
		require.Empty(t, err)
		db,err := core.AESCBCDecrypt(b,key)
		require.Empty(t, err)

		t.Log(string(db))
	}else{
		t.Log(string(hb))
	}
}

func TestBigInt(t *testing.T) {
	var a *big.Int
	t.Log(fmt.Sprintf("a=%v",a))
}

func getAmountOut(usdtIn *big.Int) (bftOut *big.Int) {
	bftOut = big.NewInt(0)
	bftReserve  := big.NewInt(1e18)
	usdtReserve := big.NewInt(1e16)
	amountInWithFee := new(big.Int).Mul(usdtIn,big.NewInt(1000-3)) // 0.3%是swap的手续费
	numerator := new(big.Int).Mul(amountInWithFee,bftReserve)
	denominator :=  new(big.Int).Mul(usdtReserve,big.NewInt(1000))
	denominator = denominator.Add(denominator,amountInWithFee)
	bftOut = bftOut.Div(numerator, denominator)
	bftOut = bftOut.Div(bftOut,big.NewInt(100))
	return
}

func TestEstimatePrice(t *testing.T) {
		usdtIn := big.NewInt(1e18)
		bftOut := getAmountOut(usdtIn)
		fIn := new(big.Float).SetInt(usdtIn)
		fOut := new(big.Float).SetInt(bftOut)
		t.Log(new(big.Float).Quo(fIn,fOut).String())
}