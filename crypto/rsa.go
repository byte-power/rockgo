package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
)

type RSAEncoder interface {
	Encrypt([]byte) ([]byte, error)
	VerifySign(msg, sign []byte) bool
}

type RSADecoder interface {
	Decrypt([]byte) ([]byte, error)
	Sign(msg []byte) ([]byte, error)
}

func NewRSAKeys(bits int) (publicKey []byte, privateKey []byte, err error) {
	newKey, err := rsa.GenerateKey(rand.Reader, bits)

	if err != nil {
		return nil, nil, err
	}
	privateKey = pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(newKey),
		},
	)

	pubASN1, err := x509.MarshalPKIXPublicKey(newKey.Public())
	if err != nil {
		return nil, nil, err
	}

	publicKey = pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubASN1,
		},
	)
	return
}

func NewRSAEncoder(pubKey []byte) RSAEncoder {
	return rsaEncoder{
		publicKey: pubKey,
	}
}

type rsaEncoder struct {
	publicKey []byte
}

func (encoder rsaEncoder) Encrypt(src []byte) ([]byte, error) {
	block, _ := pem.Decode([]byte(encoder.publicKey))

	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	publicKey := publicKeyInterface.(*rsa.PublicKey)

	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, src)
	if err != nil {
		return nil, err
	}

	return cipherText, nil
}

func (encoder rsaEncoder) VerifySign(msg, sign []byte) bool {
	block, _ := pem.Decode([]byte(encoder.publicKey))

	publicInterface, _ := x509.ParsePKIXPublicKey(block.Bytes)
	publicKey := publicInterface.(*rsa.PublicKey)
	myHash := sha256.New()
	myHash.Write([]byte(msg))
	hashed := myHash.Sum(nil)
	result := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed, sign)
	return result == nil
}

func NewRSADecoder(privKey []byte) RSADecoder {
	return rsaDecoder{
		privateKey: privKey,
	}
}

type rsaDecoder struct {
	privateKey []byte
}

func (decoder rsaDecoder) Decrypt(src []byte) ([]byte, error) {
	block, _ := pem.Decode([]byte(decoder.privateKey))

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	plainText, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, src)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}

func (decoder rsaDecoder) Sign(msg []byte) ([]byte, error) {
	block, _ := pem.Decode([]byte(decoder.privateKey))

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	myHash := sha256.New()
	myHash.Write([]byte(msg))
	hashed := myHash.Sum(nil)
	sign, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed)
	if err != nil {
		return nil, err
	}

	return sign, nil
}
