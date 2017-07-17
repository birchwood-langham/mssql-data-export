package dataexport

import "testing"

func TestEncryptNoSecret(t *testing.T) {
	secret := ""
	encrypt_text := "this is a test"

	encrypted_text := Encrypt(encrypt_text, secret)

	expected_result := "Lpl1hUiXKo6IIq1H+hAX/3Lwbz/2oBaFH0XDmHMrxQw="

	if expected_result != encrypted_text {
		t.Logf("Encrypt text failed, expecting %s, got %s", expected_result, encrypted_text)
		t.Fail()
	}

	t.Logf("Encrypting %s with no salt = %s", encrypt_text, encrypted_text)
}

func TestEncryptWithSecret(t *testing.T) {
	secret := "test_secret"
	encrypt_text := "this is a test"

	encrypted_text := Encrypt(encrypt_text, secret)

	expected_result := "4e9coPTEinoC5mSzaGQnoUxt1zDcgIGdntsyQupJ/Sg="

	if expected_result != encrypted_text {
		t.Logf("Encrypt text failed, expecting %s, got %s", expected_result, encrypted_text)
		t.Fail()
	}

	t.Logf("Encrypting %s with no salt = %s", encrypt_text, encrypted_text)
}
