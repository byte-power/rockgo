package crypto

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
)

func MD5(src []byte) []byte {
	d := md5.New()
	d.Write(src)
	return d.Sum(nil)
}

func SHA1(src []byte) []byte {
	d := sha1.New()
	d.Write(src)
	return d.Sum(nil)
}

func SHA256(src []byte) []byte {
	d := sha256.New()
	d.Write(src)
	return d.Sum(nil)
}

func SHA512(src []byte) []byte {
	d := sha512.New()
	d.Write(src)
	return d.Sum(nil)
}
