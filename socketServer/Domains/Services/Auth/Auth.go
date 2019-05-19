package Auth

import (
	"errors"
	"log"
	"net/http"
	"socket/socketServer/Domains/Repository/Mongodb"
	"socket/socketServer/Domains/Services/Api"
	"strings"

	"gopkg.in/mgo.v2"
)

func AuthMiddleware(next http.Handler, db *mgo.Session) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()

		log.Print(r.Header)

		if !validToken(r.Header.Get("Authorization"), db) {
			Api.ReturnHttpError(errors.New("Unauthorized"), w, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func validToken(token string, db *mgo.Session) bool {
	splitToken := strings.Split(token, "Bearer ")
	token = splitToken[1]

	return Mongodb.ExistsToken(token, db)
}
