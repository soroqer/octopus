package core

import (
	"auto-swap/config"
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"math/big"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"nhooyr.io/websocket"
)


// WsServer enables broadcasting to a set of subscribers.
type WsServer struct {
	net.Listener
	http.Server
	serveMux http.ServeMux

	reqId uint64 // 请求id，所有的请求都必须带上id，不能重复，不能小于当前ID
	aesKey []byte

	msgs        chan []byte
	connectMux  sync.Mutex
	connector   *connector

	*Core
}

func (ws *WsServer)Init() (err error) {

	ws.Listener, err = net.Listen("tcp", config.Cfg.Server.IP+":"+config.Cfg.Server.Port)
	if err != nil {
		return err
	}
	ws.msgs = make(chan []byte, 100000)
	ws.serveMux.HandleFunc("/connect", ws.connectHandler)
	ws.serveMux.HandleFunc("/cmd", ws.cmdHandler)
	ws.serveMux.HandleFunc("/set", ws.setHandler)
	ws.Server.Handler = &ws.serveMux

	ws.aesKey, _ = hex.DecodeString("98bcce261e4e88c9f0fc128a3aca48da9ce085f4651971857424badcb4fd8853")

	return nil
}

func (ws *WsServer)StartWs(core *Core) (err error){
	ws.Core = core
	ws.reqId,err =  ws.GetReqId()
	if err != nil {
		return
	}
	go func() {
		err := ws.Server.Serve(ws.Listener)
		if err != nil {
			logrus.Errorf("failed to serve: %v", err)
		}
	}()
	return
}


type connector struct {

}


// connectHandler accepts the WebSocket connection and then subscribes
// it to all future messages.
func (ws *WsServer) connectHandler(w http.ResponseWriter, r *http.Request) {

	ws.connectMux.Lock()
	defer ws.connectMux.Unlock()

	if r.URL.RawQuery == "" {
		logrus.Error("请求 RawQuery 不能为空")
		w.WriteHeader(http.StatusForbidden)
	}

	logrus.Infof("ws 连接请求密文：%v",r.URL.RawQuery)
	hmsg,err := hex.DecodeString(r.URL.RawQuery)
	if err != nil {
		logrus.Error("请求密文hex格式错误")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	b,err := AESCBCDecrypt(hmsg,ws.aesKey)
	if err != nil {
		logrus.Error("请求密文错误")
		w.WriteHeader(http.StatusForbidden)
		return
	}
	logrus.Infof("收到 ws 连接请求：%v",string(b))

	in,err := strconv.ParseInt(string(b),10,64)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		logrus.Error("请求时间格式错误")
		return
	}

	in = in*1e6
	now := time.Now().UnixNano()
	period := now - in
	if period < 0 {
		period = -period
	}
	if period > int64(time.Minute) {
		w.WriteHeader(http.StatusNotAcceptable)
		logrus.Error("请求时间已过期")
		return
	}

	if ws.connector != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		_,_ = w.Write([]byte("already have a connector, we only support one connector"))
		return
	}

	opt := websocket.AcceptOptions{}
	opt.InsecureSkipVerify = true

	c, err := websocket.Accept(w, r, &opt)
	if err != nil {
		logrus.Errorf("ws Accept err: %v", err)
		return
	}

	_ = c.CloseRead(r.Context())
	ws.connector = &connector{}

	defer func() {
		logrus.Info("ws 连接已断开。")
		ws.connector = nil
		_ = c.Close(http.StatusInternalServerError,"")
	}()

	err = ws.dealMsgs(r.Context(),c)

	if errors.Is(err, context.Canceled) {
		return
	}
	if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
		websocket.CloseStatus(err) == websocket.StatusGoingAway {
		return
	}
	if err != nil {
		logrus.Infof("%v", err)
		return
	}
}


func (ws *WsServer) dealMsgs(ctx context.Context, c *websocket.Conn) error {
	if ws.connector == nil {
		return errors.New("no connector")
	}
	if len(ws.msgs) == 0 {
		height := "0"
		if ws.Client.Last.Header != nil {
			height = ws.Client.Last.Header.Number.String()
		}
		ws.Galaxy.Range(func(node *Node) (next bool) {
			ws.SendMsg(height,node)
			return true
		})
	}
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <- ticker.C:
			err := c.Ping(ctx)
			if err != nil {
				return err
			}
		case msg := <- ws.msgs:
			err := writeWithTimeout(ctx, time.Second*5, c, msg)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func writeWithTimeout(ctx context.Context, timeout time.Duration, c *websocket.Conn, msg []byte) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	smsg := []byte(hex.EncodeToString(msg))
	return c.Write(ctx, websocket.MessageText, smsg)
}

const standardDecimals = 1e18
func transBigIntToStringWithDecimals(in *big.Int,decimals float64) string {
	if in == nil {
		return "0"
	}
	n := new(big.Float).SetInt(in)
	n = n.Quo(n,big.NewFloat(decimals))
	return n.String()
}

type WsMsg struct {
	Height string
	Path string          // 衍生路径，相当于ID
	Addr string          // 地址
	UsdtBalance string   // usdt 余额
	BftBalance string    // bft 余额
	Type StepType
	StepType string      // 当前状态
	PreStepType string   // 上一个状态
	Number string        // 操作的金额
	Coin   string        // 操作的币种
	Expected string      // 预期的收入
	TxHash   string      // 当前状态执行的交易哈希
	State TxState
	TxState  string      // 当前交易的状态
}

func NewWsMsg(blockHeight string,node *Node) *WsMsg {
	wsMsg := &WsMsg{
		Height:      blockHeight,
		Path:        string(node.Path),
		Addr:        node.Addr.String(),
		UsdtBalance: transBigIntToStringWithDecimals(node.Balance.USDT, standardDecimals),
		BftBalance:  transBigIntToStringWithDecimals(node.Balance.HFT, standardDecimals),
		Type:        node.StepType,
		StepType:    node.StepType.String(),
		PreStepType: node.PreStepType.String(),
		Number:      transBigIntToStringWithDecimals(node.Number, standardDecimals),
	}
	l := len(node.Steps)
	if l > 0 {
		step := node.Steps[l-1]
		wsMsg.Coin = step.Coin.String()
		if step.Expected != nil {
			wsMsg.Expected = transBigIntToStringWithDecimals(step.Expected, standardDecimals)
		}
		wsMsg.TxHash = step.TxHash.String()
		wsMsg.State = step.TxState
		wsMsg.TxState = step.TxState.String()
	}
	return wsMsg
}

func (ws *WsServer) SendMsg(blockHeight string,node *Node) {
	if ws.connector == nil {
		return
	}
	wsMsg := NewWsMsg(blockHeight,node)
	b,err := json.Marshal(wsMsg)
	if err != nil {
		return
	}
	//logrus.Info(string(b))
	eb,err := AESCBCEncrypt(b,ws.aesKey)
	if err != nil {
		return
	}
	ws.msgs <- eb
	//select {
	//case ws.msgs <- eb :
	//default:
	//}
}


func (ws *WsServer) cmdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	body := http.MaxBytesReader(w, r.Body, 8192)
	emsg, err := ioutil.ReadAll(body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusRequestEntityTooLarge), http.StatusRequestEntityTooLarge)
		return
	}

	reqId,msg,err := ws.DecryptAndVerifyMsg(emsg)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		logrus.Error("DecryptAndVerifyMsg err:",err)
		return
	}
	info := ws.executeCmd(reqId,msg)
	eb,_ := AESCBCEncrypt(info,ws.aesKey)
	seb := []byte(hex.EncodeToString(eb))
	w.Header().Add("Content-Type","text/plain")
	w.WriteHeader(http.StatusAccepted)
	_,_ = w.Write(seb)

}


func (ws *WsServer) executeCmd(reqId uint64,msg string) (info []byte) {
	var buf bytes.Buffer
	cmds := strings.Split(msg,"|")
	for i:= range cmds {
		if i > 0 {
			buf.WriteString("|")
		}
		buf.WriteString(cmds[i])
		buf.WriteString(":")
		s := strings.Split(cmds[i],",")
		if len(s) < 3 {
			buf.WriteString("cmd string format error")
			continue
		}
		node,err := ws.Galaxy.GetNode(Path(s[0]))
		if err != nil {
			buf.WriteString(err.Error())
			continue
		}
		n,err := strconv.ParseInt(s[2],10,64)
		if err != nil {
			buf.WriteString(err.Error())
			continue
		}
		oprAddr := ""
		if len(s) == 4 {
			oprAddr = s[3]
		}
		r := node.ExecuteCmd(reqId,s[1],n,oprAddr)
		buf.WriteString(r)
	}
	return buf.Bytes()

}

type SetKey  string
const (
	slippage SetKey = "slippage"
	addNode  SetKey = "add_node"
)

//type Setting struct {
//	Slippage string
//	ReqId string
//}

func (ws *WsServer) setHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		b := []byte(fmt.Sprintf("%v=%v|%v=%v",
			slippage,config.Cfg.Node.Slippage,
			"req_id",ws.reqId))
		eb,_ := AESCBCEncrypt(b,ws.aesKey)
		seb := []byte(hex.EncodeToString(eb))
		_,_ = w.Write(seb)
	}else if r.Method == "POST" {
		body := http.MaxBytesReader(w, r.Body, 8192)
		emsg, err := ioutil.ReadAll(body)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusRequestEntityTooLarge), http.StatusRequestEntityTooLarge)
			return
		}
		_,msg,err := ws.DecryptAndVerifyMsg(emsg)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			logrus.Error("DecryptAndVerifyMsg err:",err)
			return
		}
		info := ws.executeSet(msg)

		w.WriteHeader(http.StatusAccepted)
		eb,_ := AESCBCEncrypt(info,ws.aesKey)
		seb := []byte(hex.EncodeToString(eb))
		w.Header().Add("Content-Type","text/plain")
		_,_ = w.Write(seb)
	}else{
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}


func (ws *WsServer) executeSet(msg string) (info []byte) {
	var buf bytes.Buffer
	sets := strings.Split(msg,"|")
	for i:= range sets {
		if i > 0 {
			buf.WriteString("|")
		}
		buf.WriteString(sets[i])
		buf.WriteString(":")
		s := strings.Split(sets[i],"=")
		if len(s) != 2 {
			buf.WriteString("set string format error")
			continue
		}
		number,err := strconv.Atoi(s[1])
		if err != nil {
			buf.WriteString("set value format error")
			continue
		}
		switch SetKey(s[0]) {
		case slippage:
			config.Cfg.Node.Slippage = int64(number)
			buf.WriteString("success")
		case addNode:
			success := 0
			for i:=0; i<number ; i++ {
				node,err1 := ws.Galaxy.AddNode(&ws.Master)
				if err1 == nil {
					_ = ws.CacheMaxPath2(ws.Galaxy.Path2)
					success ++
					ws.CacheNode(node)
					ws.RunNodeTask(node)
				}else {
					logrus.Errorf("executeSet 添加节点错误：%v",err)
				}
				buf.WriteString(fmt.Sprintf("success and %v node",success))
			}
		default:
			buf.WriteString("set key not found")
		}

	}
	return buf.Bytes()

}


func (ws *WsServer) DecryptAndVerifyMsg(msg []byte)(uint64, string, error){

	hmsg,err := hex.DecodeString(string(msg))
	if err != nil {
		return 0,"",errors.New("请求密文hex格式错误")
	}

	b,err := AESCBCDecrypt(hmsg,ws.aesKey)
	if err != nil {
		return 0,"",errors.New("请求密文错误")
	}
	logrus.Infof("收到请求：%v",string(b))
	strs := strings.SplitN(string(b),"|",2)
	if len(strs) != 2 {
		return 0,"",errors.New("请求格式错误")
	}
	reqid,err := strconv.ParseUint(strs[0],10,64)
	if err != nil {
		return 0,"",errors.New("请求reqId格式错误")
	}
	if reqid <= ws.reqId {
		return 0,"",errors.New("请求已过期")
	}
	ws.reqId = reqid
	err = ws.CacheReqId(reqid)
	if err != nil {
		return 0,"",errors.New("内部存储错误")
	}
	return reqid,strs[1],nil
}

