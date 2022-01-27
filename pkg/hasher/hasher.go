package hasher

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

// HashPassword receives password string, return hashed string by sha256 or error
func HashPassword(password string) (string, error) {
	return hashAndSalt([]byte(password))
}

//CheckPasswordHash receives user password and hashed string, after compare return true or false depend on result
func CheckPasswordHash(password, hashed string) bool {

	byteHash := []byte(hashed)
	bytePass := []byte(password)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePass)
	if err != nil {
		log.Println(err)
		return false
	}
	return true

}

func hashAndSalt(pwd []byte) (string, error) {

	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
