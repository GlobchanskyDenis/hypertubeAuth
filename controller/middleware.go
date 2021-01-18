package controller

import (
	"HypertubeAuth/errors"
	"HypertubeAuth/logger"
	"net/http"
	"time"
)

func panicRecover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		defer func() {
			if rec := recover(); rec != nil {
				err, ok := rec.(error)
				if ok {
					logger.Error(r, errors.UnknownInternalError.SetArgs("Произошла ПАНИКА", "PANIC happened").SetOrigin(err))
				} else {
					logger.Error(r, errors.UnknownInternalError.SetArgs("Произошла ПАНИКА, отсутствует интерфейс ошибки",
						"PANIC happened, error interface expected"))
				}
				errorResponse(w, errors.UnknownInternalError)
				return
			}
		}()
		next.ServeHTTP(w, r)
		logger.Duration(r, time.Since(t))
	})
}

func corsPut(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "PUT,OPTIONS")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type,Content-Length,Authorization")

		if r.Method == "OPTIONS" {
			logger.Log(r, "client wants to know what methods are allowed")
			return
		} else if r.Method != "PUT" {
			logger.Warning(r, "wrong request method. Should be PUT method")
			w.WriteHeader(http.StatusMethodNotAllowed) // 405
			return
		}
		next.ServeHTTP(w, r)
	})
}

func corsDelete(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "DELETE,OPTIONS")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type,Content-Length,Authorization")

		if r.Method == "OPTIONS" {
			logger.Log(r, "client wants to know what methods are allowed")
			return
		} else if r.Method != "DELETE" {
			logger.Warning(r, "wrong request method. Should be DELETE method")
			w.WriteHeader(http.StatusMethodNotAllowed) // 405
			return
		}
		next.ServeHTTP(w, r)
	})
}

func corsPost(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "POST,OPTIONS")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type,Content-Length,Authorization")

		if r.Method == "OPTIONS" {
			logger.Log(r, "client wants to know what methods are allowed")
			return
		} else if r.Method != "POST" {
			logger.Warning(r, "wrong request method. Should be POST method")
			w.WriteHeader(http.StatusMethodNotAllowed) // 405
			return
		}
		next.ServeHTTP(w, r)
	})
}

func corsPatch(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "PATCH,OPTIONS")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type,Content-Length,Authorization")

		if r.Method == "OPTIONS" {
			logger.Log(r, "client wants to know what methods are allowed")
			return
		} else if r.Method != "PATCH" {
			logger.Warning(r, "wrong request method. Should be PATCH method")
			w.WriteHeader(http.StatusMethodNotAllowed) // 405
			return
		}
		next.ServeHTTP(w, r)
	})
}

func corsGet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "GET,OPTIONS")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type,Content-Length,Authorization")
		if r.Method == "OPTIONS" {
			logger.Log(r, "client wants to know what methods are allowed")
			return
		} else if r.Method != "GET" {
			logger.Warning(r, "wrong request method. Should be GET method")
			w.WriteHeader(http.StatusMethodNotAllowed) // 405
			return
		}
		next.ServeHTTP(w, r)
	})
}
