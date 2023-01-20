package hash

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pass []byte) (hashedPass []byte, err error) {
	hashedPass, err = bcrypt.GenerateFromPassword(pass, bcrypt.MinCost)
	return
}
