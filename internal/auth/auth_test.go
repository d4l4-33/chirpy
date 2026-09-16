package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
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

func TestCreateAndValidate(t *testing.T) {
	userId, err := uuid.NewRandom()
	if err != nil {
		t.Error(err)
	}
	expiresIn, err := time.ParseDuration("5m")
	if err != nil {
		t.Error(err)
	}
	tokenSecret := "verySecretC0de"

	tokenString, err := MakeJWT(userId, tokenSecret, expiresIn)
	if err != nil {
		t.Error(err)
	}

	returedUuid, err := ValidateJWT(tokenString, tokenSecret)
	if err != nil {
		t.Error(err)
	}
	if returedUuid != userId {
		t.Fail()
	}
}

func TestWrongValidateCode(t *testing.T) {
	userId, err := uuid.NewRandom()
	if err != nil {
		t.Error(err)
	}
	expiresIn, err := time.ParseDuration("5m")
	if err != nil {
		t.Error(err)
	}
	tokenSecret := "ExtraSEcretCode3"

	tokenString, err := MakeJWT(userId, tokenSecret, expiresIn)
	if err != nil {
		t.Error(err)
	}

	_, err = ValidateJWT(tokenString, "tokenSecret")
	if err == nil {
		t.Fail()
	}
}

func TestTimeoutError(t *testing.T) {
	userId, err := uuid.NewRandom()
	if err != nil {
		t.Error(err)
	}
	expiresIn, err := time.ParseDuration("1s")
	if err != nil {
		t.Error(err)
	}
	tokenSecret := "shouldTimeOut"

	tokenString, err := MakeJWT(userId, tokenSecret, expiresIn)
	if err != nil {
		t.Error(err)
	}

	duration, err := time.ParseDuration("2s")
	if err != nil {
		t.Error(err)
	}
	time.Sleep(duration)

	_, err = ValidateJWT(tokenString, tokenSecret)
	if err == nil {
		t.Fail()
	}
}

func TestBearerToken(t *testing.T) {
	token := "fortyfiveandahalflitersofvodka"
	header := http.Header{}
	header.Add("Authorization", "Bearer "+token)

	response, err := GetBearerToken(header)
	if err != nil {
		t.Error(err)
	}

	if response != token {
		t.Fail()
	}
}

func TestBearerNotFound(t *testing.T) {
	header := http.Header{}
	_, err := GetBearerToken(header)
	if err.Error() != "Bearer token not found" {
		t.Fail()
	}

}
