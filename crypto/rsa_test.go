package crypto

import "testing"

func TestRsa(t *testing.T) {
	plain := []byte("Hello")
	pbKey, pvKey, err := NewRSAKeys(512)
	if err != nil {
		t.Fatal("fail", err)
	}

	encoder := NewRSAEncoder(pbKey)
	decoder := NewRSADecoder(pvKey)

	encryptData, err := encoder.Encrypt(plain)
	if err != nil {
		t.Fatal("fail", err)
	}

	decryptData, err := decoder.Decrypt(encryptData)
	if err != nil {
		t.Fatal("fail", err)
	}

	if string(decryptData) != string(plain) {
		t.Fatal("the text decrypted dose not match the text before encrypted")
	}

	signature, err := decoder.Sign(plain)
	if err != nil {
		t.Fatal("fail", err)
	}

	if !encoder.VerifySign(plain, signature) {
		t.Fatal("rsa verify signature fail")
	}

	t.Log("success")
}
