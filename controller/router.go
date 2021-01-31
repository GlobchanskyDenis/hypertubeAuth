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
	mux.Handle("/api/auth/basic", corsGet(http.HandlerFunc(authBasic)))
	mux.Handle("/api/auth/oauth42", corsGet(http.HandlerFunc(authOauth42)))
	mux.Handle("/api/info", corsGet(http.HandlerFunc(info)))
	mux.Handle("/api/profile/get", corsGet(authMW(http.HandlerFunc(profileGet))))

	// POST
	mux.Handle("/api/auth/check", corsPost(http.HandlerFunc(authCheck)))
	mux.Handle("/api/email/resend", corsPost(http.HandlerFunc(emailResend)))
	mux.Handle("/api/email/confirm", corsGetPost(http.HandlerFunc(emailConfirm)))

	// PUT
	mux.Handle("/api/profile/create", corsPut(http.HandlerFunc(profileCreate)))

	// PATCH
	mux.Handle("/api/profile/patch", corsPatch(authMW(http.HandlerFunc(profilePatch))))
	mux.Handle("/api/passwd/patch", corsPatch(authMW(http.HandlerFunc(passwordPatch))))

	// DELETE
	mux.Handle("/api/profile/delete", corsDelete(authMW(http.HandlerFunc(profileDelete))))

	// /email/patch
	// /passwd/repair

	serveMux := panicRecover(mux)

	return serveMux
}
