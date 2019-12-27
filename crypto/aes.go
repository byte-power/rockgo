package crypto

import (
	"bytes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"strconv"
)
import "crypto/aes"

type Coder interface {
	Encrypt([]byte) ([]byte, error)
	Decrypt([]byte) ([]byte, error)
}

func pkcs7pad(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func pkcs7strip(origData []byte) ([]byte, error) {
	length := len(origData)
	if length == 0 {
		return nil, errors.New("pkcs7: Data is empty")
	}
	padSize := int(origData[length-1])
	if length < padSize {
		return nil, errors.New("pkcs7: Invalid padding")
	}
	return origData[:length-padSize], nil
}

func NewAESCoderWithECB(key []byte) (Coder, error) {
	c, e := aes.NewCipher(key)
	if e != nil {
		return nil, e
	}
	return aesECBCoder{
		cipher: c,
	}, nil
}

type aesECBCoder struct {
	cipher cipher.Block
}

func (coder aesECBCoder) Encrypt(src []byte) ([]byte, error) {
	block := coder.cipher

	data := pkcs7pad(src, block.BlockSize())
	encrypted := make([]byte, len(data))
	size := block.BlockSize()

	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		block.Encrypt(encrypted[bs:be], data[bs:be])
	}
	return encrypted, nil
}

func (coder aesECBCoder) Decrypt(src []byte) (dst []byte, e error) {
	block := coder.cipher

	length := len(src)
	decrypted := make([]byte, length)
	size := block.BlockSize()

	//分组分块加密
	for bs, be := 0, size; bs < length; bs, be = bs+size, be+size {
		block.Decrypt(decrypted[bs:be], src[bs:be])
	}

	cipherText, err := pkcs7strip(decrypted)
	if err != nil {
		return nil, fmt.Errorf("Unpadding failed. %w", err)
	}
	return cipherText, nil
}

func NewAESCoderWithCBC(key, iv []byte) (Coder, error) {
	c, e := aes.NewCipher(key)
	if e != nil {
		return nil, e
	}
	return aesCBCCoder{
		cipher: c,
		iv:     iv,
	}, nil
}

type aesCBCCoder struct {
	cipher cipher.Block
	iv     []byte
}

func (coder aesCBCCoder) Encrypt(src []byte) ([]byte, error) {
	block := coder.cipher

	//填充原文
	blockSize := block.BlockSize()
	rawData := pkcs7pad(src, blockSize)
	cipherText := make([]byte, blockSize+len(rawData))

	var iv []byte

	if len(coder.iv) != blockSize {
		return nil, errors.New("The length of iv should be " + strconv.Itoa(blockSize))
	} else {
		iv = coder.iv
	}
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	//block大小和初始向量大小一定要一致
	encrypt := cipher.NewCBCEncrypter(block, iv)
	encrypt.CryptBlocks(cipherText[blockSize:], rawData)

	return cipherText, nil
}

func (coder aesCBCCoder) Decrypt(src []byte) ([]byte, error) {
	block := coder.cipher

	blockSize := block.BlockSize()
	encryptData := src
	if len(encryptData) < blockSize {
		return nil, errors.New("Cipher text too short ")
	}

	var iv []byte
	if len(coder.iv) != blockSize {
		return nil, errors.New("The length of iv should be " + strconv.Itoa(blockSize))
	} else {
		iv = coder.iv
	}

	encryptData = encryptData[blockSize:]

	// CBC mode always works in whole blocks.
	if len(encryptData)%blockSize != 0 {
		return nil, errors.New("Cipher text is not a multiple of the block size ")
	}

	decrypt := cipher.NewCBCDecrypter(block, iv)

	// CryptBlocks can work in-place if the two arguments are the same.
	decrypt.CryptBlocks(encryptData, encryptData)
	//解填充
	plainText, err := pkcs7strip(encryptData)

	if err != nil {
		return nil, err
	}
	return plainText, nil
}
