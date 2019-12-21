package crypto

import "crypto/cipher"

type Coder interface {
	Encrypt([]byte) ([]byte, error)
	Decrypt([]byte) ([]byte, error)
}

func NewAESCoderWithECB(key []byte) (Coder, error) {
	return nil, nil
}

type aesECBCoder struct {
	cipher cipher.Block
}

func NewAESCoderWithCBC(key, iv []byte) (Coder, error) {
	return nil, nil
}

type aesCBCCoder struct {
	cipher cipher.Block
	iv     []byte
}
