package controller

import (
	"HypertubeAuth/controller/hash"
	"HypertubeAuth/errors"
	"HypertubeAuth/logger"
	"HypertubeAuth/model"
	"HypertubeAuth/postgres"
	"io/ioutil"
	"encoding/json"
	"strconv"
	"net/http"
	"net/url"
	"net"
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
	UserId      uint   `json:"id"`
	Email       string `json:"email"`
	Fname       string `json:"first_name"`
	Lname       string `json:"last_name"`
	Displayname string `json:"displayname"`
	ImageBody   string `json:"image_url"`
}

/*
**	Endpoint function
*/
func userAuthOauth42(w http.ResponseWriter, r *http.Request) {
	params, Err := parseRequestParams(r)
	if Err != nil {
		logger.Warning(r, Err.Error())
		errorResponse(w, Err)
		return
	}
	token, Err := getTokenFrom42(params)
	if Err != nil {
		logger.Error(r, Err)
		errorResponse(w, Err)
		return
	}
	user, Err := getUser42(token)
	if Err != nil {
		logger.Error(r, Err)
		errorResponse(w, Err)
		return
	}
	// fmt.Printf("%#v\n", user)
	_, Err = postgres.UserGet42ById(user.UserId)
	if Err != nil && !errors.UserNotExist.IsOverlapWithError(Err) {
		logger.Error(r, Err.SetArgs("1", "1"))
		errorResponse(w, Err)
		return
	}
	if Err != nil && errors.UserNotExist.IsOverlapWithError(Err) {
		if Err = postgres.UserSet42(user); Err != nil {
			logger.Error(r, Err.SetArgs("2", "2"))
			errorResponse(w, Err)
			return
		}
	} else {
		if Err = postgres.UserUpdate42(user); Err != nil {
			logger.Error(r, Err.SetArgs("3", "3"))
			errorResponse(w, Err)
			return
		}
	}

	// tokenRaw, Err := hash.CreateToken(user.TransformToUser(), "user_42_strategy")
	// if Err != nil {
	// 	logger.Warning(r, "cannot get password hash - " + Err.Error())
	// 	errorResponse(w, Err)
	// 	return
	// }

	accessTokenRaw, Err := hash.CreateAccessTokenBase64(user.TransformToUser(), "user_42_strategy")
	if Err != nil {
		logger.Warning(r, "cannot create access token - " + Err.Error())
		errorResponse(w, Err)
		return
	}

	// successResponse(w, tokenRaw)
	logger.Success(r, "user #" + strconv.Itoa(int(user.UserId)) + " was authenticated")
	cookie := &http.Cookie{ Name: "tokenRaw", Value: accessTokenRaw}
	// 	"/",
	// 	"www.domain.com",
	// 	expire,
	// 	expire.Format(time.UnixDate),
	// 	86400,
	// 	true,
	// 	true,
	// 	"test=tcookie",
	// 	[]string{"test=tcookie"}
	// }
	// req.AddCookie(&cookie)
	http.SetCookie(w, cookie)
	w.Header().Add("Authenticate", accessTokenRaw)
	http.Redirect(w, r,
		"http://localhost:8008/",
		// "http://file:///home/skinny/Documents/go/src/HypertubeAuth/client/client.html",
		http.StatusTemporaryRedirect)
}

/*
**	Parsing GET params from request
*/
func parseRequestParams(r *http.Request) (requestParams, *errors.Error) {
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
		"client_id": {"96975efecfd0e5efee67c9ac4cc350ac9372ae559b2fb8a08feba6841a33fb53",},
		"client_secret": {"bdcbe28874ab05962b50430b1466a8ebcbda45ba8c3c1beee600699478ad2a4d",},
		"code": {params.Code,},
		"state": {params.State,},
		// "redirect_uri": {"file:///home/skinny/Documents/go/src/HypertubeAuth/client/client.html",},
		"redirect_uri": {"http://localhost:"+portString+"/user/auth/oauth42",},
		"grant_type": {"authorization_code",},
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
		"client_id": {"96975efecfd0e5efee67c9ac4cc350ac9372ae559b2fb8a08feba6841a33fb53",},
		"client_secret": {"bdcbe28874ab05962b50430b1466a8ebcbda45ba8c3c1beee600699478ad2a4d",},
		"refresh_token": {refreshToken,},
		"redirect_uri": {"http://localhost:"+portString+"/user/auth/oauth42",},
		"grant_type": {"refresh_token",},
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
			Timeout: 30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
			MaxIdleConns: 100,
			IdleConnTimeout: 90 * time.Second,
			TLSHandshakeTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
	}
	client := &http.Client{
		Timeout: time.Second * 10,
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
		Email: profile.Email,
		Fname: profile.Fname,
		Lname: profile.Lname,
		Displayname: profile.Displayname,
		ImageBody: profile.ImageBody,
		User42Model: model.User42Model{
			UserId: profile.UserId,
			AccessToken: &token.AccessToken,
			RefreshToken: &token.RefreshToken,
			ExpiresAt: &token.ExpiresAt,
		},
	}, nil
}