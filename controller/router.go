package controller

import (
	"HypertubeAuth/errors"
	"net/http"
)

func errorResponse(w http.ResponseWriter, err *errors.Error) {
	w.WriteHeader(err.GetHttpStatus())
	w.Write(err.ToJson())
}

func successResponse(w http.ResponseWriter, response []byte) {
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func Router() http.Handler {

	mux := http.NewServeMux()

	// GET
	mux.Handle("/user/auth/basic", corsGet(http.HandlerFunc(userAuthBasic)))
	mux.Handle("/user/auth/oauth42", corsGet(http.HandlerFunc(userAuthOauth42)))
	mux.Handle("/info", corsGet(http.HandlerFunc(info)))

	// POST
	mux.Handle("/token/decode", corsPost(http.HandlerFunc(tokenCheck)))
	mux.Handle("/user/profile", corsPost(http.HandlerFunc(userProfile)))

	// PUT
	mux.Handle("/user/create/basic", corsPut(http.HandlerFunc(userCreateBasic)))
	// mux.HandleFunc("/oauth42/", oauth42)

	serveMux := panicRecover(mux)

	return serveMux
}
