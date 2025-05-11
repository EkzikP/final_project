package api

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

var PASS string

func signingHandler(w http.ResponseWriter, r *http.Request) {
	insertedPass := make(map[string]string)
	var buf bytes.Buffer

	//десериализуем полученный в запросе JSON
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, "Authentification required", http.StatusUnauthorized)
		//http.Redirect(w, r, "/login.html", http.StatusUnauthorized)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &insertedPass); err != nil {
		http.Error(w, "Authentification required", http.StatusUnauthorized)
		//http.Redirect(w, r, "/login.html", http.StatusUnauthorized)
		return
	}

	enterPass, ok := insertedPass["password"]
	if !ok || enterPass != PASS {
		http.Error(w, "Authentification required", http.StatusUnauthorized)
		//http.Redirect(w, r, "/login.html", http.StatusUnauthorized)
		return
	}

	secret := []byte("13_go_basic")
	h := sha1.New()
	h.Write([]byte(PASS))
	sha1_hash := hex.EncodeToString(h.Sum(nil))
	// создаём payload
	claims := jwt.MapClaims{"hash": sha1_hash}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString(secret)
	if err != nil {
		http.Error(w, "Authentification required", http.StatusUnauthorized)
		//http.Redirect(w, r, "/login.html", http.StatusUnauthorized)
		return
	}
	writeJson(w, map[string]string{"token": token}, nil)
}
