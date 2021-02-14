package controller

import (
	"HypertubeAuth/controller/hash"
	"HypertubeAuth/errors"
	"HypertubeAuth/logger"
	"HypertubeAuth/model"
	"HypertubeAuth/postgres"
	// "encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
	// "fmt"
)

type requestParams struct {
	Code             string
	State            string
	Error            string
	ErrorDescription string
}

type token42 struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
	ExpiresIn    uint   `json:"expires_in"`
	ExpiresAt    time.Time
}

type profile42 struct {
	User42Id      uint   `json:"id"`
	Email       string `json:"email"`
	Fname       string `json:"first_name"`
	Lname       string `json:"last_name"`
	Displayname string `json:"displayname"`
	ImageBody   string `json:"image_url"`
}

/*
**	/api/auth/oauth42
 */
func authOauth42(w http.ResponseWriter, r *http.Request) {
	conf, Err := getConfig()
	if Err != nil {
		logger.Error(r, Err)
		// errorResponse(w, Err)
		http.Redirect(w, r,
			conf.SocketRedirect + conf.ErrorRedirect + "?error=" + url.QueryEscape(string(Err.ToJson())),//base64.StdEncoding.EncodeToString(Err.ToJson())
			http.StatusTemporaryRedirect)
		return
	}

	params, Err := parseRequestParams42(r)
	if Err != nil {
		logger.Warning(r, Err.Error())
		// errorResponse(w, Err)
		http.Redirect(w, r,
			conf.SocketRedirect + conf.ErrorRedirect + "?error=" + url.QueryEscape(string(Err.ToJson())),
			http.StatusTemporaryRedirect)
		return
	}
	/*
	**	Getging intra42 api token
	 */
	token, Err := getTokenFrom42(params)
	if Err != nil {
		logger.Error(r, Err)
		// errorResponse(w, Err)
		http.Redirect(w, r,
			conf.SocketRedirect + conf.ErrorRedirect + "?error=" + url.QueryEscape(string(Err.ToJson())),
			http.StatusTemporaryRedirect)
		return
	}
	/*
	**	Getting user profile from intra api and fills it into *model.User42
	 */
	user, Err := getUser42(token)
	if Err != nil {
		logger.Error(r, Err)
		// errorResponse(w, Err)
		http.Redirect(w, r,
			conf.SocketRedirect + conf.ErrorRedirect + "?error=" + url.QueryEscape(string(Err.ToJson())),
			http.StatusTemporaryRedirect)
		return
	}

	var userBasic *model.UserBasic
	/*
	**	getting user from db if it exists
	 */
	userFromDb, Err := postgres.UserGet42ById(user.User42Id)
	if Err != nil {
		if errors.UserNotExist.IsOverlapWithError(Err) {
			// user not exists
			logger.Log(r, "User42 with user42Id "+strconv.Itoa(int(user.User42Id))+" not found in database. Creating new one")
			if userBasic, Err = postgres.UserSet42(user); Err != nil {
				logger.Error(r, Err.SetArgs("1", "1"))
				// errorResponse(w, Err)
				http.Redirect(w, r,
					conf.SocketRedirect + conf.ErrorRedirect + "?error=" + url.QueryEscape(string(Err.ToJson())),
					http.StatusTemporaryRedirect)
				return
			}
		} else {
			// database error
			logger.Error(r, Err.SetArgs("2", "2"))
			// errorResponse(w, Err)
			http.Redirect(w, r,
				conf.SocketRedirect + conf.ErrorRedirect + "?error=" + url.QueryEscape(string(Err.ToJson())),
				http.StatusTemporaryRedirect)
			return
		}
	} else {
		user.UserId = userFromDb.UserId
		if Err = postgres.UserUpdate42(user); Err != nil {
			logger.Error(r, Err.SetArgs("3", "3"))
			// errorResponse(w, Err)
			http.Redirect(w, r,
				conf.SocketRedirect + conf.ErrorRedirect + "?error=" + url.QueryEscape(string(Err.ToJson())),
				http.StatusTemporaryRedirect)
			return
		}
		userBasic, Err = postgres.UserGetBasicById(user.UserId)
		if Err = postgres.UserUpdate42(user); Err != nil {
			logger.Error(r, Err.SetArgs("4", "4"))
			// errorResponse(w, Err)
			http.Redirect(w, r,
				conf.SocketRedirect + conf.ErrorRedirect + "?error=" + url.QueryEscape(string(Err.ToJson())),
				http.StatusTemporaryRedirect)
			return
		}
	}

	accessToken, Err := hash.CreateToken(userBasic)
	if Err != nil {
		logger.Warning(r, "cannot get password hash - "+Err.Error())
		// errorResponse(w, Err)
		http.Redirect(w, r,
			conf.SocketRedirect + conf.ErrorRedirect + "?error=" + url.QueryEscape(string(Err.ToJson())),
			http.StatusTemporaryRedirect)
		return
	}

	logger.Success(r, "user #"+strconv.Itoa(int(user.UserId))+" was authenticated")
	cookie := &http.Cookie{Name: "accessToken", Value: accessToken, HttpOnly: false}

	http.SetCookie(w, cookie)
	// w.Header().Add("accessToken", accessToken)
	http.Redirect(w, r,
		conf.SocketRedirect + conf.OauthRedirect + "?accessToken=" + accessToken,
		http.StatusTemporaryRedirect)
}

/*
**	Parsing GET params from request
 */
func parseRequestParams42(r *http.Request) (requestParams, *errors.Error) {
	var params requestParams

	params.Code = r.FormValue("code")
	params.State = r.FormValue("state")
	params.Error = r.FormValue("error")
	params.ErrorDescription = r.FormValue("error_description")

	if params.Error != "" || params.ErrorDescription != "" {
		return params, errors.AccessDenied.SetHidden("Сервер авторизации 42 ответил: " +
			params.Error + " - " + params.ErrorDescription)
	}
	if params.Code == "" || params.State == "" {
		return params, errors.AccessDenied.SetHidden("Сервер авторизации 42 прислал невалидные данные. code: " +
			params.Code + " state" + params.State)
	}
	return params, nil
}

/*
**	Request to ecole 42 server API for token
 */
func getTokenFrom42(params requestParams) (token42, *errors.Error) {
	var result token42

	conf, Err := getConfig()
	if Err != nil {
		return result, Err
	}
	portString := strconv.FormatUint(uint64(conf.ServerPort), 10)

	formData := url.Values{
		"client_id":     {"96975efecfd0e5efee67c9ac4cc350ac9372ae559b2fb8a08feba6841a33fb53"},
		"client_secret": {"bdcbe28874ab05962b50430b1466a8ebcbda45ba8c3c1beee600699478ad2a4d"},
		"code":          {params.Code},
		"state":         {params.State},
		// "redirect_uri": {"file:///home/skinny/Documents/go/src/HypertubeAuth/client/client.html",},
		"redirect_uri": {"http://localhost:" + portString + "/api/auth/oauth42"},
		"grant_type":   {"authorization_code"},
	}
	resp, err := http.PostForm("https://api.intra.42.fr/oauth/token", formData)
	if err != nil {
		return result, errors.AccessDenied.SetHidden("Запрос токена из 42 провален").SetOrigin(err)
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return result, errors.AccessDenied.SetHidden("Декодирование json дало ошибку").SetOrigin(err)
	}
	if result.ExpiresIn < 4 {
		result, Err := refreshTokenFrom42(result.RefreshToken)
		return result, Err
	}
	duration, err := time.ParseDuration(strconv.FormatUint(uint64(result.ExpiresIn), 10) + "s")
	if err != nil {
		return result, errors.UnknownInternalError.SetArgs("ошибка парсинга времени", "time parse fail").SetOrigin(err)
	}
	result.ExpiresAt = time.Now().Add(duration)
	return result, nil
}

func refreshTokenFrom42(refreshToken string) (token42, *errors.Error) {
	var result token42

	conf, Err := getConfig()
	if Err != nil {
		return result, Err
	}
	portString := strconv.FormatUint(uint64(conf.ServerPort), 10)

	formData := url.Values{
		"client_id":     {"96975efecfd0e5efee67c9ac4cc350ac9372ae559b2fb8a08feba6841a33fb53"},
		"client_secret": {"bdcbe28874ab05962b50430b1466a8ebcbda45ba8c3c1beee600699478ad2a4d"},
		"refresh_token": {refreshToken},
		"redirect_uri":  {"http://localhost:" + portString + "/api/auth/oauth42"},
		"grant_type":    {"refresh_token"},
	}
	resp, err := http.PostForm("https://api.intra.42.fr/oauth/token", formData)
	if err != nil {
		return result, errors.AccessDenied.SetHidden("Запрос токена из 42 провален").SetOrigin(err)
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return result, errors.AccessDenied.SetHidden("Декодирование json дало ошибку").SetOrigin(err)
	}
	duration, err := time.ParseDuration(strconv.FormatUint(uint64(result.ExpiresIn), 10) + "s")
	if err != nil {
		return result, errors.UnknownInternalError.SetArgs("ошибка парсинга времени", "time parse fail").SetOrigin(err)
	}
	result.ExpiresAt = time.Now().Add(duration)
	return result, nil
}

/*
**	Request to ecole 42 server API for user profile
 */
func getUserProfile(accessToken string) (profile42, *errors.Error) {
	var profile profile42
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	client := &http.Client{
		Timeout:   time.Second * 10,
		Transport: transport,
	}
	url := "https://api.intra.42.fr/v2/me"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return profile, errors.AccessDenied.SetHidden("Запрос данных пользователя 42 провален").SetOrigin(err)
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)
	resp, err := client.Do(req)
	if err != nil {
		return profile, errors.AccessDenied.SetHidden("Запрос данных пользователя 42 провален").SetOrigin(err)
	}
	defer resp.Body.Close() // важный пункт!
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return profile, errors.AccessDenied.SetHidden("Чтение данных пользователя 42 провалено").SetOrigin(err)
	}
	if err = json.Unmarshal(respBody, &profile); err != nil {
		return profile, errors.AccessDenied.SetHidden("Декодирование данных пользователя из json дало ошибку").SetOrigin(err)
	}
	return profile, nil
}

/*
**	Forming User42 structure
 */
func getUser42(token token42) (*model.User42, *errors.Error) {
	profile, Err := getUserProfile(token.AccessToken)
	if Err != nil {
		return nil, Err
	}

	return &model.User42{
		Email:       profile.Email,
		Fname:       profile.Fname,
		Lname:       profile.Lname,
		Displayname: profile.Displayname,
		ImageBody:   profile.ImageBody,
		User42Model: model.User42Model{
			User42Id:     profile.User42Id,
			AccessToken:  &token.AccessToken,
			RefreshToken: &token.RefreshToken,
			ExpiresAt:    &token.ExpiresAt,
		},
	}, nil
}
