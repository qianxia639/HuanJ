package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	// 定义命令行参数
	var action string
	flag.StringVar(&action, "action", "", "Action to perform: 'generate', 'load',  'sign', or 'verify'")
	flag.Parse()

	switch action {
	case "generate":
		flag.CommandLine.Parse(os.Args[2:])
		generateKeyPair()
	case "load":
		loadKeysFromFile()
	case "sign":
		signMessage()
	case "verify":
		verifySignature()
	default:
		fmt.Println("Invalid action. Use 'generate', 'sign', or 'verify'.")
	}

}

const (
	// privateKeyFile = "private_key.pem"
	// publicKeyFile  = "public_key.pem"

	privateKeyBlockType = "PRIVATE KEY"
	publicKeyBlockType  = "PUBLIC KEY"
)

// 生成 ED25519 密钥对
func generateKeyPair() {

	var publicKeyFile, privateKeyFile string
	flagSet := flag.NewFlagSet("generate", flag.ExitOnError)
	flagSet.StringVar(&publicKeyFile, "public-key-file", "public_key.pem", "公钥文件输出路径")
	flagSet.StringVar(&privateKeyFile, "private-key-file", "private_key.pem", "公钥文件输出路径")
	flagSet.Parse(os.Args[:])

	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		log.Fatalf("Failed to generate key: %v", err)
	}

	fmt.Println("Private Key (hex):", hex.EncodeToString(privateKey))
	fmt.Println("Public Key (hex):", hex.EncodeToString(publicKey))

	// 将私钥编码为PEM格式
	privateKeyPem := &pem.Block{
		Type:  privateKeyBlockType,
		Bytes: privateKey,
	}
	privateKeyPemBytes := pem.EncodeToMemory(privateKeyPem)

	// 将公钥编码为PEM格式
	publicKeyPem := &pem.Block{
		Type:  publicKeyBlockType,
		Bytes: publicKey,
	}
	publicKeyPemBytes := pem.EncodeToMemory(publicKeyPem)

	// 将私钥和公钥保存到文件
	err = os.WriteFile(privateKeyFile, privateKeyPemBytes, 0600)
	if err != nil {
		log.Fatalf("Failed to save private key: %v", err)
	}

	err = os.WriteFile(publicKeyFile, publicKeyPemBytes, 0644)
	if err != nil {
		log.Fatalf("Failed to save public key: %v", err)
	}

	log.Print("successfully...")
}

func loadKeysFromFile() {

	var publicKeyFile, privateKeyFile string
	flagSet := flag.NewFlagSet("load", flag.ExitOnError)
	flagSet.StringVar(&publicKeyFile, "public-key-file", "public_key.pem", "公钥文件读取路径")
	flagSet.StringVar(&privateKeyFile, "private-key-file", "private_key.pem", "公钥文件读取路径")
	flagSet.Parse(os.Args[2:])

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

// 签名消息
func signMessage() {
	var privateKeyHex, message string
	flag.StringVar(&privateKeyHex, "private-key", "", "Private key in hex format")
	flag.StringVar(&message, "message", "", "Message to sign")
	flag.Parse()

	if privateKeyHex == "" || message == "" {
		fmt.Println("Please provide both private key and message.")
		return
	}

	// 将私钥从十六进制转换为字节数组
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		fmt.Println("Invalid private key:", err)
		return
	}

	// 生成签名
	signature := ed25519.Sign(privateKeyBytes, []byte(message))

	fmt.Println("Signature (hex):", hex.EncodeToString(signature))
}

// 验证签名
func verifySignature() {
	var publicKeyHex, message, signatureHex string
	flag.StringVar(&publicKeyHex, "public-key", "", "Public key in hex format")
	flag.StringVar(&message, "message", "", "Message to verify")
	flag.StringVar(&signatureHex, "signature", "", "Signature in hex format")
	flag.Parse()

	if publicKeyHex == "" || message == "" || signatureHex == "" {
		fmt.Println("Please provide public key, message, and signature.")
		return
	}

	// 将公钥和签名从十六进制转换为字节数组
	publicKeyBytes, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		fmt.Println("Invalid public key:", err)
		return
	}

	signatureBytes, err := hex.DecodeString(signatureHex)
	if err != nil {
		fmt.Println("Invalid signature:", err)
		return
	}

	// 验证签名
	isValid := ed25519.Verify(publicKeyBytes, []byte(message), signatureBytes)

	if isValid {
		fmt.Println("Signature is valid.")
	} else {
		fmt.Println("Signature is invalid.")
	}
}
