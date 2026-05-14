package envfile

import (
	"strings"
	"testing"
)

func testKey() []byte {
	// 32 bytes for AES-256
	return []byte("01234567890123456789012345678901")
}

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	key := testKey()
	plaintext := "super-secret-value"

	enc, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}

	if enc == plaintext {
		t.Error("expected ciphertext to differ from plaintext")
	}

	dec, err := Decrypt(enc, key)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}

	if dec != plaintext {
		t.Errorf("expected %q, got %q", plaintext, dec)
	}
}

func TestEncrypt_ProducesUniqueOutputs(t *testing.T) {
	key := testKey()
	a, _ := Encrypt("value", key)
	b, _ := Encrypt("value", key)
	if a == b {
		t.Error("expected different ciphertexts due to random nonce")
	}
}

func TestDecrypt_InvalidBase64(t *testing.T) {
	_, err := Decrypt("not-valid-base64!!!", testKey())
	if err != ErrInvalidCiphertext {
		t.Errorf("expected ErrInvalidCiphertext, got %v", err)
	}
}

func TestDecrypt_TamperedCiphertext(t *testing.T) {
	enc, _ := Encrypt("hello", testKey())
	tampered := strings.ToUpper(enc[:4]) + enc[4:]
	_, err := Decrypt(tampered, testKey())
	if err == nil {
		t.Error("expected error decrypting tampered ciphertext")
	}
}

func TestEncryptMap_RoundTrip(t *testing.T) {
	key := testKey()
	env := map[string]string{
		"DB_PASSWORD": "hunter2",
		"API_KEY":     "abc123",
	}

	enc, err := EncryptMap(env, key)
	if err != nil {
		t.Fatalf("EncryptMap: %v", err)
	}

	for k, v := range env {
		if enc[k] == v {
			t.Errorf("key %s: expected encrypted value to differ", k)
		}
	}

	dec, err := DecryptMap(enc, key)
	if err != nil {
		t.Fatalf("DecryptMap: %v", err)
	}

	for k, want := range env {
		if got := dec[k]; got != want {
			t.Errorf("key %s: expected %q, got %q", k, want, got)
		}
	}
}

func TestDecryptMap_BadValue(t *testing.T) {
	env := map[string]string{"KEY": "!!!not-encrypted!!!"}
	_, err := DecryptMap(env, testKey())
	if err == nil {
		t.Error("expected error for invalid ciphertext in map")
	}
}
