package lib

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

// This can be overwriten within go build (32 bit)

var CipherKey string = "impEdfCJyeIpDtUOSwMysMIaQYdeONTZ"

// Used to encrypte the (to base64 converted) AFAS token.
func EncryptAES(message string) string {
	//Create byte array from the input string
	plainText := []byte(message)

	//Create a new AES cipher using the key
	block, err := aes.NewCipher([]byte(CipherKey))

	//If NewCipher failed, exit:
	if err != nil {
		fmt.Println("error creating AES cipher key")
		os.Exit(1)
	}

	//Make the cipher text a byte array of size BlockSize + the length of the message
	cipherText := make([]byte, aes.BlockSize+len(plainText))

	//If is the ciphertext up to the blocksize (16)
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		fmt.Println("error encoding AFAS token")
		os.Exit(1)
	}

	//Encrypt the data:
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	//Return string encoded in base64
	return base64.RawStdEncoding.EncodeToString(cipherText)
}

// Decrypte the AFAS token so it can be used in the API call
func DecryptAES(secure string) string {
	//Remove base64 encoding:
	cipherText, err := base64.RawStdEncoding.DecodeString(secure)

	//If DecodeString failed, exit:
	if err != nil {
		fmt.Println("error decoding secure string to base64")
		os.Exit(1)
	}

	//Create a new AES cipher with the key and encrypted message
	block, err := aes.NewCipher([]byte(CipherKey))

	//If NewCipher failed, exit:
	if err != nil {
		fmt.Println("new cipher in decoding AFAS token failed")
		os.Exit(1)
	}

	//If the length of the cipherText is less than 16 Bytes:
	if len(cipherText) < aes.BlockSize {
		fmt.Println("ciphertext block size is too short")
		os.Exit(1)
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	//Decrypt the message
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText)
}
