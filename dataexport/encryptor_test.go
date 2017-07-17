package dataexport

import "testing"

func TestEncryptNoSecret(t *testing.T) {
	secret := ""
	encrypt_text := "this is a test"

	encrypted_text := Encrypt(encrypt_text, secret)

	expected_result := "2e99758548972a8e8822ad47fa1017ff72f06f3ff6a016851f45c398732bc50c"

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

	expected_result := "e1ef5ca0f4c48a7a02e664b3686427a14c6dd730dc80819d9edb3242ea49fd28"

	if expected_result != encrypted_text {
		t.Logf("Encrypt text failed, expecting %s, got %s", expected_result, encrypted_text)
		t.Fail()
	}

	t.Logf("Encrypting %s with no salt = %s", encrypt_text, encrypted_text)
}
