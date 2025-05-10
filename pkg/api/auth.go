package api

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// смотрим наличие пароля
		if len(PASS) > 0 {
			var tokenString string // JWT-токен из куки
			// получаем куку
			cookie, err := r.Cookie("token")
			if err == nil {
				tokenString = cookie.Value
			}
			var valid bool
			secretKey := []byte("13_go_basic")
			jwtToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return secretKey, nil
			})
			if err != nil {
				//output := "ошибка парсинга токена"
				//writeJson(w, Out{Error: output})
				//http.Redirect(w, r, "/login.html", http.StatusUnauthorized)
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}

			// приводим поле Claims к типу jwt.MapClaims
			res, ok := jwtToken.Claims.(jwt.MapClaims)
			if !ok {
				//output := "ошибка парсинга токена"
				//writeJson(w, Out{Error: output})
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				//http.Redirect(w, r, "/login.html", http.StatusUnauthorized)
				return
			}

			hashRaw := res["hash"]
			// loginRaw — интерфейс, так как тип значения в jwt.Claims — интерфейс. Чтобы получить строку, нужно
			// снова сделать приведение типа к строке.
			hash, ok := hashRaw.(string)
			if !ok {
				//output := "ошибка парсинга токена"
				//writeJson(w, Out{Error: output})
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				//http.Redirect(w, r, "/login.html", http.StatusUnauthorized)
				return
			}

			h := sha1.New()
			h.Write([]byte(PASS))
			sha1_hash := hex.EncodeToString(h.Sum(nil))

			if hash != sha1_hash {
				//output := "ошибка авторизации, введён неправильный пароль"
				//writeJson(w, Out{Error: output})
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				//http.Redirect(w, r, "/login.html", http.StatusUnauthorized)
				return
			} else {
				valid = true
			}

			if !valid {
				// возвращаем ошибку авторизации 401
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				//http.Redirect(w, r, "/login.html", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
