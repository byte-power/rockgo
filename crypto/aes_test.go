package crypto

import (
	"testing"
)

func TestNewAESCoderWithECB(t *testing.T) {
	key := []byte("1234567890abcdef")
	plain := []byte("Hello")

	ecb, err := NewAESCoderWithECB(key)
	if err != nil {
		t.Fatal("fail", err)
	}
	encryptBytes, err := ecb.Encrypt(plain)
	if err != nil {
		t.Fatal("fail", err)
	}
	decryptBytes, err := ecb.Decrypt(encryptBytes)
	if err != nil {
		t.Fatal("fail", err)
	}

	if string(decryptBytes) == string(plain) {
		t.Log("success")
	}
}

func TestNewAESCoderWithCBC(t *testing.T) {
	key := []byte("1234567890abcdef")
	iv := []byte("0123456789abcdef")
	plain := []byte("Hello")

	cbc, err := NewAESCoderWithCBC(key, iv)
	if err != nil {
		t.Fatal("fail", err)
	}
	encryptBytes, err := cbc.Encrypt(plain)
	if err != nil {
		t.Fatal("fail", err)
	}
	decryptBytes, err := cbc.Decrypt(encryptBytes)
	if err != nil {
		t.Fatal("fail", err)
	}

	if string(decryptBytes) == string(plain) {
		t.Log("success")
	}
}
