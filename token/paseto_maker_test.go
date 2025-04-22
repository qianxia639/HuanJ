package token

import (
	"Rejuv/utils"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/o1egl/paseto"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker := NewPasetoMaker(utils.RandomString(32))

	username := utils.RandomString(6)
	duration := time.Minute

	issueAt := time.Now()
	expired := issueAt.Add(duration)

	token, err := maker.CreateToken(Token{Username: username, Duration: duration})
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issueAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expired, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {

	key := utils.RandomString(32)

	maker := NewPasetoMaker(key)

	token, err := maker.CreateToken(Token{Username: utils.RandomString(6), Duration: -time.Minute})
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

const (
	privateKeyFile = "private_key.pem"
	publicKeyFile  = "public_key.pem"

	privateKeyBlockType = "PRIVATE KEY"
	publicKeyBlockType  = "PUBLIC KEY"
)

func generateEd25519Key() {
	// 生成Ed25519密钥对
	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		log.Fatalf("Failed to generate key: %v", err)
	}

	log.Printf("len(publicKey): %d, len(privateKey): %d", len(publicKey), len(privateKey))

	log.Printf("public key: %x\n", publicKey)
	log.Printf("private key: %x\n", privateKey)

	log.Print("-----------------------------")

	fmt.Println("Private Key (hex):", hex.EncodeToString(privateKey))
	fmt.Println("Public Key (hex):", hex.EncodeToString(publicKey))

	// 将私钥编码为PEM格式
	// privateKeyPem := &pem.Block{
	// 	Type:  "PRIVATE KEY",
	// 	Bytes: privateKey,
	// }
	// privateKeyPemBytes := pem.EncodeToMemory(privateKeyPem)
	privateKeyPemBytes := pemEncode(privateKey, privateKeyBlockType)

	// 将公钥编码为PEM格式
	// publicKeyPem := &pem.Block{
	// 	Type:  "PUBLIC KEY",
	// 	Bytes: publicKey,
	// }
	// publicKeyPemBytes := pem.EncodeToMemory(publicKeyPem)
	publicKeyPemBytes := pemEncode(privateKey, publicKeyBlockType)

	// 将私钥和公钥保存到文件
	err = os.WriteFile(privateKeyFile, privateKeyPemBytes, 0644)
	if err != nil {
		log.Fatalf("Failed to save private key: %v", err)
	}

	err = os.WriteFile(publicKeyFile, publicKeyPemBytes, 0644)
	if err != nil {
		log.Fatalf("Failed to save public key: %v", err)
	}

}

// 将密钥编码为PEM格式
func pemEncode(key []byte, typeLabel string) []byte {
	block := &pem.Block{
		Type:  typeLabel,
		Bytes: key,
	}

	return pem.EncodeToMemory(block)
}

func TestGenerate(t *testing.T) {
	generateEd25519Key()
}

// 从 PEM 文件加载密钥
func TestLoadKey(t *testing.T) {

	// 读取私钥文件
	privateKeyPEM, err := os.ReadFile(privateKeyFile)
	if err != nil {
		fmt.Println("Error reading private key file:", err)
		return
	}

	// 读取公钥文件
	publicKeyPEM, err := os.ReadFile(publicKeyFile)
	if err != nil {
		fmt.Println("Error reading public key file:", err)
		return
	}

	// 解码私钥
	privateKeyBlock, _ := pem.Decode(privateKeyPEM)
	if privateKeyBlock == nil || privateKeyBlock.Type != privateKeyBlockType {
		fmt.Println("Invalid private key format.")
		return
	}

	// 解码公钥
	publicKeyBlock, _ := pem.Decode(publicKeyPEM)
	if publicKeyBlock == nil || publicKeyBlock.Type != publicKeyBlockType {
		fmt.Println("Invalid public key format.")
		return
	}

	// 转换为字节数组
	privateKey := privateKeyBlock.Bytes
	publicKey := publicKeyBlock.Bytes

	fmt.Println("Private Key (hex):", hex.EncodeToString(privateKey))
	fmt.Println("Public Key (hex):", hex.EncodeToString(publicKey))
}

func TestVerify(t *testing.T) {

	// 读取私钥
	privateKeyPem, err := os.ReadFile(privateKeyFile)
	require.NoError(t, err)

	// 解析Pem格式的私钥
	block, _ := pem.Decode(privateKeyPem)
	if block == nil || block.Type != "PRIVATE KEY" {
		t.Fatal("无效的PEM文件")
	}

	privateKey := ed25519.PrivateKey(block.Bytes)

	// 创建 JSON Token
	// jsonToken := paseto.JSONToken{
	// 	Issuer:     "qianxia",
	// 	Subject:    "Test",
	// 	Jti:        "123456",
	// 	Expiration: time.Now().Add(30 * time.Minute),
	// }

	payload := Payload{
		ID:        "123456",
		Username:  "Test",
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(5 * time.Minute),
	}

	// 签名
	token, err := paseto.NewV2().Sign(privateKey, payload, nil)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	t.Logf("Token: %s", token)

	// 读取公钥
	publicKeyPem, err := os.ReadFile(publicKeyFile)
	require.NoError(t, err)

	block, _ = pem.Decode(publicKeyPem)
	if block == nil || block.Type != "PUBLIC KEY" {
		t.Fatal("无效的PEM文件")
	}

	publicKey := ed25519.PublicKey(block.Bytes)

	// 校验
	var newPayload Payload
	err = paseto.NewV2().Verify(token, publicKey, &newPayload, nil)
	require.NoError(t, err)
	require.NotEmpty(t, newPayload)

	t.Logf("new payload: %+v", newPayload)
}

func TestPasetoMakerV2(t *testing.T) {

	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)

	maker := NewPasetoMakerV2(privateKey, publicKey)

	username := utils.RandomString(6)
	duration := time.Minute
	issueAt := time.Now()
	expired := issueAt.Add(duration)

	token, err := maker.CreateToken(Token{Username: username, Duration: time.Minute})
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issueAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expired, payload.ExpiredAt, time.Second)
}
