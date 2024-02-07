package crypto

import (
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
		"",
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
