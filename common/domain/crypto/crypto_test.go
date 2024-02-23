package crypto

import (
	"fmt"
	"testing"
)

func TestEncDec(t *testing.T) {
	// init a 32byte key
	enc := NewEncryption([]byte("12345678123456781234567812345678"))

	test := []string{
		"hello",
		"as;dlkfjas;dlhfa;sdhfas;df",
		"psdsasfsadfasdfas@asdfsadfl.com",
		"13339849223",
		"1",
	}

	for _, v := range test {
		data, err := enc.Encrypt(v)
		if err != nil {
			t.Fatal(err)
		}

		text, err := enc.Decrypt(data)
		t.Logf("enc is %s, plain is %s\n", string(data), text)
		if err != nil {
			t.Fatal(err)
		}

		if text != v {
			t.Fatal("encrypt and decrypt not equal")
		}
	}
}

func TestLenPlain(t *testing.T) {
	enc := NewEncryption([]byte("12345678123456781234567812345678"))
	// 准备测试数据
	encrypted := "efb"

	// 执行解密
	decrypted, err := enc.Decrypt(encrypted)
	t.Logf("%s\n", decrypted)
	if err != nil {
		fmt.Println("index is negative ")
	} else {
		t.Fatal("Decrypted data")
	}

}
