package crypto

import "testing"

func TestEncryptDecrypt(t *testing.T) {
	key := []byte("12345678901234567890123456789012") // 32 bytes
	c, err := NewAESCrypto(key)
	if err != nil {
		t.Fatalf("erro ao criar crypto: %v", err)
	}

	original := "https://google.com"

	encrypted, err := c.Encrypt(original)
	if err != nil {
		t.Fatalf("erro ao criptografar: %v", err)
	}

	if string(encrypted) == original {
		t.Fatal("texto nao foi criptografado")
	}

	decrypted, err := c.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("erro ao descriptografar: %v", err)
	}

	if decrypted != original {
		t.Fatalf("esperava %q, obteve %q", original, decrypted)
	}
}

func TestKeyValidation(t *testing.T) {
	_, err := NewAESCrypto([]byte("chave-curta"))
	if err == nil {
		t.Fatal("deveria falhar com chave de tamanho errado")
	}
}

func TestEncryptDifferentEachTime(t *testing.T) {
	key := []byte("12345678901234567890123456789012")
	c, _ := NewAESCrypto(key)

	a, _ := c.Encrypt("mesmo texto")
	b, _ := c.Encrypt("mesmo texto")

	if string(a) == string(b) {
		t.Fatal("nonce nao esta funcionando: criptografias identicas")
	}
}
