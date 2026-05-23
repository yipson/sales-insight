package security

import (
	"testing"
)

func TestEncrypter_Roundtrip(t *testing.T) {
	key := "01234567890123456789012345678901" // exactly 32 bytes
	enc, err := NewEncrypter(key)
	if err != nil {
		t.Fatalf("failed to create encrypter: %v", err)
	}

	plaintext := "my-secret-token-value"
	ciphertext, err := enc.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}
	if ciphertext == "" {
		t.Error("expected non-empty ciphertext")
	}
	if ciphertext == plaintext {
		t.Error("ciphertext should not equal plaintext")
	}

	decrypted, err := enc.Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("decrypt failed: %v", err)
	}
	if decrypted != plaintext {
		t.Errorf("decrypted = %q, want %q", decrypted, plaintext)
	}
}

func TestEncrypter_DifferentPlaintexts(t *testing.T) {
	key := "01234567890123456789012345678901"
	enc, err := NewEncrypter(key)
	if err != nil {
		t.Fatalf("failed to create encrypter: %v", err)
	}

	cases := []string{
		"",
		"short",
		"a-very-long-token-with-many-characters-and-symbols-!@#$%",
		`{"access_token":"abc123","refresh_token":"xyz789"}`,
	}

	for _, want := range cases {
		cipher, err := enc.Encrypt(want)
		if err != nil {
			t.Errorf("encrypt(%q) failed: %v", want, err)
			continue
		}
		got, err := enc.Decrypt(cipher)
		if err != nil {
			t.Errorf("decrypt(%q) failed: %v", want, err)
			continue
		}
		if got != want {
			t.Errorf("roundtrip(%q) = %q, want %q", want, got, want)
		}
	}
}

func TestEncrypter_InvalidKeyLength(t *testing.T) {
	cases := []string{
		"short",
		"exactly31byteslongexactly31byte",
		"exactly33byteslongexactly33bytesl",
		"",
	}

	for _, key := range cases {
		_, err := NewEncrypter(key)
		if err == nil {
			t.Errorf("NewEncrypter(%q) expected error, got nil", key)
		}
	}
}

func TestEncrypter_DecryptWithWrongKey(t *testing.T) {
	key1 := "01234567890123456789012345678901"
	key2 := "98765432109876543210987654321098"

	enc1, _ := NewEncrypter(key1)
	enc2, _ := NewEncrypter(key2)

	ciphertext, _ := enc1.Encrypt("secret")
	_, err := enc2.Decrypt(ciphertext)
	if err == nil {
		t.Error("expected decryption with wrong key to fail")
	}
}

func TestEncrypter_DecryptInvalidCiphertext(t *testing.T) {
	key := "01234567890123456789012345678901"
	enc, _ := NewEncrypter(key)

	_, err := enc.Decrypt("not-valid-base64!!!")
	if err == nil {
		t.Error("expected decrypt of invalid base64 to fail")
	}

	_, err = enc.Decrypt("dG9vLXNob3J0")
	if err == nil {
		t.Error("expected decrypt of too-short ciphertext to fail")
	}
}
