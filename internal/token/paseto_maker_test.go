package token

import (
	"Rejuv/internal/utils"
	"crypto/ed25519"
	"encoding/pem"
	"log"
	"os"
	"testing"
	"time"

	"github.com/o1egl/paseto"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	// maker := NewPasetoMaker(utils.RandomString(32))

	maker := NewMockPasetoMaker()

	username := utils.RandomString(6)
	duration := time.Minute

	issueAt := time.Now()
	expired := issueAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
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
	maker := NewPasetoMaker(utils.RandomString(32))

	token, err := maker.CreateToken(utils.RandomString(6), -time.Minute)
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

	// 将私钥编码为PEM格式
	privateKeyPem := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKey,
	}
	privateKeyPemBytes := pem.EncodeToMemory(privateKeyPem)

	// 将公钥编码为PEM格式
	publicKeyPem := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKey,
	}
	publicKeyPemBytes := pem.EncodeToMemory(publicKeyPem)

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

func TestGenerate(t *testing.T) {
	generateEd25519Key()
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
	jsonToken := paseto.JSONToken{
		Issuer:     "qianxia",
		Subject:    "Test",
		Jti:        "123456",
		Expiration: time.Now().Add(30 * time.Minute),
	}

	// 签名
	token, err := paseto.NewV2().Sign(privateKey, jsonToken, nil)
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
	var payload paseto.JSONToken
	err = paseto.NewV2().Verify(token, publicKey, &payload, nil)
	require.NoError(t, err)

	t.Logf("Issuer: %s", payload.Issuer)
	t.Logf("Subject: %s", payload.Subject)
	t.Logf("Jti: %s", payload.Jti)
}

type MockPasetoMaker struct {
	paseto     *paseto.V2
	privateKey ed25519.PrivateKey
	publicKey  ed25519.PublicKey
}

func NewMockPasetoMaker() Maker {

	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		log.Fatalf("generate key failed: %v", err)
	}

	maker := &MockPasetoMaker{
		paseto:     paseto.NewV2(),
		privateKey: privateKey,
		publicKey:  publicKey,
	}

	return maker
}

// 创建Token
func (maker *MockPasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload := NewPayload(username, duration)

	token, err := maker.paseto.Sign(maker.privateKey, payload, nil)

	return token, err
}

// 校验Token
func (maker *MockPasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := maker.paseto.Verify(token, maker.publicKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// err = payload.Valid()
	// if err != nil {
	// 	return nil, err
	// }
	return payload, nil
}
