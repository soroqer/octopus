package lib

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestCrypto(t *testing.T) {
	pwd := "ldifgy-9"
	t.Log(len(pwd),pwd)
	data := "abcd"
	key := FormatAESKey(pwd)
	t.Log(len(key),string(key))
	enc,err := AESCBCEncrypt([]byte(data),key)
	require.Empty(t, err)
	t.Log(hex.EncodeToString(enc))
	dec,err := AESCBCDecrypt(enc,key)
	require.Empty(t, err)
	require.Equal(t, string(dec), data)
}

func TestExportKey(t *testing.T) {
	pwd := "1"
	t.Log(len(pwd),pwd)
	priv := "fbdd18beb4c945882436fc67f79e3136cdcd2fa7160d793be4c89dfcdfbe5982"
	data,err := hex.DecodeString(priv)
	require.Empty(t, err)
	key := FormatAESKey(pwd)
	t.Log(len(key),string(key))
	enc,err := AESCBCEncrypt(data,key)
	require.Empty(t, err)
	t.Log(base64.StdEncoding.EncodeToString(enc))
	dec,err := AESCBCDecrypt(enc,key)
	require.Empty(t, err)
	require.Equal(t, hex.EncodeToString(dec), priv)
}


func TestECIES(t *testing.T) {

	plainText := []byte("1234567890abcdefg")


	priv1,err := crypto.GenerateKey()
	require.Empty(t, err)

	privStr := priv1.D.Text(16)
	t.Log("秘钥：",privStr)

	publicKey := ecies.ImportECDSAPublic(&priv1.PublicKey)

	cryptText, err := ecies.Encrypt(rand.Reader,publicKey,plainText,nil,nil)
	require.Empty(t, err)
	fmt.Println("ECC传入公钥加密的密文为：", hex.EncodeToString(cryptText))

	priv2,err := crypto.HexToECDSA(privStr)
	require.Empty(t, err)

	privateKey := ecies.ImportECDSA(priv2)

	msg, err := privateKey.Decrypt(cryptText, nil,nil)
	require.Empty(t, err)
	fmt.Println("ECC传入私钥解密后的明文为：", string(msg))


}

func TestAES(t *testing.T) {

	plainText := []byte("1|m-1001-1,buy,22")

	key := make([]byte,32,32)
	n,err := rand.Reader.Read(key)
	require.Empty(t, err)
	t.Log(n)
	t.Log("秘钥 hex: ",hex.EncodeToString(key))
	enc,err := AESCBCEncrypt(plainText,key)
	require.Empty(t, err)
	t.Log("原文 utf-8：",string(plainText))
	t.Log("原文 hex：",hex.EncodeToString(plainText))
	t.Log("密文 hex：",hex.EncodeToString(enc))
	dec,err := AESCBCDecrypt(enc,key)
	require.Empty(t, err)
	require.Equal(t, dec, plainText)


}

func TestAesDecrypt(t *testing.T) {

	key,err := hex.DecodeString("98bcce261e4e88c9f0fc128a3aca48da9ce085f4651971857424badcb4fd8853")
	require.Empty(t, err)

	enc,err := hex.DecodeString("2a3a8aab098eee539a3e4f4c0f1e7d5442670ab5a6827c89a6895910ec24bd5632fcc492ae468bd48e91c72cc4937c2d")
	require.Empty(t, err)

	b,err := AESCBCDecrypt(enc,key)
	require.Empty(t, err)

	t.Log(string(b))

	strs := strings.SplitN(string(b),"|",2)
	t.Log(len(strs))
	t.Log(strs)
}

