package crypto

import (
	"fmt"
	"testing"
)

func TestDigest(t *testing.T) {
	plain := []byte("Hello")
	expectedMd5Str := "8b1a9953c4611296a827abf8c47804d7"
	expectedSha1Str := "f7ff9e8b7bb2e09b70935a5d785e0cc5d9d0abf0"
	expectedSha256Str := "185f8db32271fe25f561a6fc938b2e264306ec304eda518007d1764826381969"
	expectedSha512Str := "3615f80c9d293ed7402687f94b22d58e529b8cc7916f8fac7fddf7fbd5af4cf777d3d795a7a00a16bf7e7f3fb9561ee9baae480da9fe7a18769e71886b03f315"

	md5Str := fmt.Sprintf("%x", MD5(plain))
	sha1Str := fmt.Sprintf("%x", SHA1(plain))
	sha256Str := fmt.Sprintf("%x", SHA256(plain))
	sha512Str := fmt.Sprintf("%x", SHA512(plain))

	if md5Str != expectedMd5Str {
		t.Fatal("md5 encryption error")
	}

	if sha1Str != expectedSha1Str {
		t.Fatal("sha1 encryption error")
	}

	if sha256Str != expectedSha256Str {
		t.Fatal("sha256 encryption error")
	}

	if sha512Str != expectedSha512Str {
		t.Fatal("sha512 encryption error")
	}

	t.Log("success")

}
