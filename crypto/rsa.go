package crypto

type RSAEncoder interface {
	Encrypt([]byte) ([]byte, error)
	VerifySign(msg, sign []byte) bool
}

type RSADecoder interface {
	Decrypt([]byte) ([]byte, error)
	Sign(msg []byte) ([]byte, error)
}

func NewRSAKeys(bits int) (publicKey []byte, privateKey []byte, err error) {
	return nil, nil, nil
}

func NewRSAEncoder(pubKey []byte) RSAEncoder {
	return nil
}

type rsaEncoder struct {
	publicKey []byte
}

func NewRSADecoder(privKey []byte) RSADecoder {
	return nil
}

type rsaDecoder struct {
	privateKey []byte
}
