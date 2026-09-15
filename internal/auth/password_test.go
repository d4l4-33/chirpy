package auth

import (
	"testing"
)

func TestPasswordMatch(t *testing.T) {
	password := "passW0rd123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("Error hashing: %s", err)
	}
	match, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Errorf("Error checking: %s", err)
	}

	if !match {
		t.Fail()
	}
}

func TestPasswordFail(t *testing.T) {
	password := "passW0rd123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("Error hashing: %s", err)
	}
	password = "wrongPassword"
	match, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Errorf("Error checking: %s", err)
	}

	if match {
		t.Fail()
	}
}

func TestPasswordEmpty(t *testing.T) {
	password := ""
	_, err := HashPassword(password)
	if err.Error() != "Password is too short" {
		t.Fail()
	}

}
