package core

import (
	"auto-swap/config"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	//"github.com/wumansgy/goEncrypt"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"strings"
	"syscall"
	"time"
)

const (
	AesKeyLen24 = 24
	AesKeyLen32 = 32
)

func padToBlockSize(payload []byte, blockSize int) (aligned, finalBlock []byte) {
	overrun := len(payload) % blockSize
	paddingLen := blockSize - overrun
	aligned = payload[:len(payload)-overrun]
	finalBlock = make([]byte, blockSize)
	copy(finalBlock, payload[len(payload)-overrun:])
	for i := overrun; i < blockSize; i++ {
		finalBlock[i] = byte(paddingLen)
	}
	return
}

// AESCBCEncrypt Encrypt
func AESCBCEncrypt(data []byte, cbcKey []byte) ([]byte, error) {
	var buffer bytes.Buffer
	buffer.Write(data)
	key := make([]byte, len(cbcKey))
	copy(key, cbcKey)
	c, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("NewCipher error:", err)
		return nil, err
	}

	// set random IV len : aes.BlockSize(16)
	iv := make([]byte, aes.BlockSize)
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return nil, err
	}
	encrypter := cipher.NewCBCEncrypter(c, iv)

	//规整数据长度
	aligned, finalBlock := padToBlockSize(buffer.Bytes(), encrypter.BlockSize())
	src := append(aligned, finalBlock...)
	length := len(src)

	enData := make([]byte, length)
	encrypter.CryptBlocks(enData, src)
	// add iv
	retBytes := append(iv, enData...)
	return retBytes, nil
}

// AESCBCDecrypt Decrypt
func AESCBCDecrypt(data []byte, cbcKey []byte) ([]byte, error) {
	var buffer bytes.Buffer
	buffer.Write(data)
	dataCp := buffer.Bytes()
	if len(data) < aes.BlockSize {
		return nil, fmt.Errorf("data is too short")
	}
	iv := dataCp[:aes.BlockSize]
	deData := dataCp[aes.BlockSize:]
	length := len(deData)

	c, err := aes.NewCipher(cbcKey)
	if err != nil {
		err = fmt.Errorf("NewCipher[%s] err: %s", cbcKey, err)
		return nil, err
	}

	deCrypter := cipher.NewCBCDecrypter(c, iv)
	// check block size
	if len(deData)%deCrypter.BlockSize() != 0 {
		return nil, fmt.Errorf("crypto/cipher: input not full blocks")
	}
	deCrypter.CryptBlocks(deData, deData)
	endIndex := length - int(deData[length-1])
	if len(deData) < endIndex || endIndex < 0 {
		return nil, fmt.Errorf("data error[datalen=%v][index=%v]", len(deData), endIndex)
	}
	outData := deData[:endIndex]

	return outData, nil
}


func FormatAESKey(key string) []byte {

	var paddings = []byte("大漠孤烟直长河落日圆")
	bKey := []byte(key)
	// 不足24位补足24位
	if len(bKey) < AesKeyLen24 {
		step := AesKeyLen24 -len(bKey)
		return append(bKey,paddings[:step]...)
	}
	// 24~32位补足32位
	if len(bKey) > AesKeyLen24 && len(bKey)< AesKeyLen32 {
		step := AesKeyLen32 -len(bKey)
		return append(bKey,paddings[:step]...)
	}
	// 大于32位截取前32位
	if len(bKey) > AesKeyLen32 {
		return bKey[:AesKeyLen32]
	}

	return bKey
}

func IsPwdCorrect(pwd string) bool {
	bKey,err := base64.StdEncoding.DecodeString(config.Cfg.Node.Key)
	if err != nil {
		return false
	}
	ak := FormatAESKey(pwd)
	bKey,err = AESCBCDecrypt(bKey,ak)
	if err != nil {
		return false
	}
	privKey,err := crypto.ToECDSA(bKey)
	if err != nil {
		return false
	}

	if crypto.PubkeyToAddress(privKey.PublicKey) != common.HexToAddress(config.Cfg.Node.KeyAddress) {
		return false
	}
	return true
}


func GetPassword() (password string, err error) {
	time.Sleep(time.Millisecond * 100)
	var bytePassword []byte
	for i := 0; i < config.RetryTimes; i++ {
		fmt.Print("➜ Enter Password: ")
		bytePassword, err = terminal.ReadPassword(syscall.Stdin)
		// next line
		fmt.Println()
		if err != nil {
			err = fmt.Errorf("➜ get password error:%v", err)
			return "", err
		}
		password = strings.TrimSpace(string(bytePassword))
		passLen := len(password)
		if passLen < config.PwdMinLen {
			fmt.Printf("➜ pass word too short,at least[%v]\n", config.PwdMinLen)
			err = fmt.Errorf("➜ get pass word failed")
			password = ""
			continue
		}
		if IsPwdCorrect(password) {
			return
		}
		fmt.Println("  Password incorrect.")
	}
	return "", errors.New("password retry exhaust")
}

