package jwt

import "github.com/golang-jwt/jwt/v5"

type Jwt struct {
	Secret string
}

func NewJwt(secret string) *Jwt {
	return &Jwt{secret}
}

func (j *Jwt) Sign(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
	})
	signedToken, err := token.SignedString([]byte(j.Secret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
