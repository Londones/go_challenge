// Description: This file contains the functions that are used in the application.
package utils

import "github.com/go-chi/jwtauth/v5"

func MakeToken(name string) string {
	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"name": name})
	return tokenString
}

func init() {
	tokenAuth = jwtauth.New("HS256", []byte(Secret), nil)
}
