package main

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func main() {

	// cryptoAlgorithm := crp
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"usr": "username",
	})

	if signedToken, err := token.SignedString([]byte(jwt.SigningMethodHS256.Hash.String())); err != nil { // less security as it's not a encrypt key hash
		panic(err)
	} else {
		jwtObject, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
			if token.Method.Alg() != jwt.SigningMethodHS256.Name {
				return nil, errors.New("doesn't have the same signature alg")
			} else {
				return []byte(jwt.SigningMethodHS256.Hash.String()), nil
			}
		})

		if err != nil {
			panic(err)
		} else {
			if claims, ok := jwtObject.Claims.(jwt.MapClaims); ok {	
				fmt.Printf("usr claim: %s", claims["usr"])
			} else {
				panic("cannot found claim usr")
			}
		}
	}
}
