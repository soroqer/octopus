package core

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
	priv := "65d4a52824d91a7fc6e547f6414dea477a4e2813cbe683758e944caaec85515d"
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

func TestIsPwdCorrect(t *testing.T) {
	t.Log(IsPwdCorrect("sujiao1234"))
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

	plainText := []byte("1234567890abcdefg")

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

	key,err := hex.DecodeString("6530498a63963d05902dd7a6fd31dbfe7010ab59435829537f915ad6617436f5")
	require.Empty(t, err)

	enc,err := hex.DecodeString("3362313334316437633365336532393086fec29863b3964d46e7ad66abe8c8a6c2f5b91ddd2c7b7e548b1635005a2d64")
	require.Empty(t, err)

	b,err := AESCBCDecrypt(enc,key)
	require.Empty(t, err)

	t.Log(string(b))

	strs := strings.SplitN(string(b),"|",2)
	t.Log(len(strs))
	t.Log(strs)
}

