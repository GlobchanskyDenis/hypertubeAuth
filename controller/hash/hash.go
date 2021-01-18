package hash

import (
	"HypertubeAuth/errors"
	"HypertubeAuth/model"
	"hash/crc32"
	"strings"
	"strconv"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"encoding/base64"
	"io"
	// "fmt"
)

type Config struct {
	PasswdSalt string `conf:"passwordSalt"`
	MasterKey  string `conf:"masterKey"`
}

var cfg *Config

func GetConfig() *Config {
	if cfg == nil {
		cfg = &Config{}
	}
	return cfg
}

func getConfig() (*Config, *errors.Error) {
	if cfg == nil {
		return nil, errors.NotConfiguredPackage.SetArgs("controller/hash", "controller/hash")
	}
	return cfg, nil
}

func PasswdHash(pass string) (string, *errors.Error) {
	conf, Err := getConfig()
	if Err != nil {
		return "", Err
	}
	pass += conf.PasswdSalt
	crcH := crc32.ChecksumIEEE([]byte(pass))
	return strconv.FormatUint(uint64(crcH), 20), nil
}

func EmailHash(mail string) (string, *errors.Error) {
	conf, Err := getConfig()
	if Err != nil {
		return "", Err
	}
	c, err := aes.NewCipher([]byte(conf.MasterKey))
	if err != nil {
		return "", errors.HashError.SetOrigin(err)
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", errors.HashError.SetOrigin(err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", errors.HashError.SetOrigin(err)
	}
	token := gcm.Seal(nonce, nonce, []byte(mail), nil)
	return base64.URLEncoding.EncodeToString(token), nil
}

func CreateToken(user model.User, accountType string) ([]byte, *errors.Error) {
	var header model.TokenHeader
	header.UserId = user.UserId
	header.AccountType = accountType
	headerJson, err := json.Marshal(header)
	if err != nil {
		return nil, errors.MarshalError.SetOrigin(err)
	}
	payloadJson, err := json.Marshal(user)
	if err != nil {
		return nil, errors.MarshalError.SetOrigin(err)
	}

	headerPayload := base64.StdEncoding.EncodeToString(headerJson) + "." + base64.StdEncoding.EncodeToString(payloadJson)
	signature, Err := PasswdHash(headerPayload)
	if Err != nil {
		return nil, Err
	}

	var token = model.Token{
		AccessToken: base64.StdEncoding.EncodeToString([]byte(headerPayload + "." + signature)),
		Profile: &user,
	}
	tokenJson, err := json.Marshal(token)
	if err != nil {
		return nil, errors.MarshalError.SetOrigin(err)
	}
	
	return tokenJson, nil
}

func CreateAccessTokenBase64(user model.User, accountType string) (string, *errors.Error) {
	var header model.TokenHeader
	header.UserId = user.UserId
	header.AccountType = accountType
	headerJson, err := json.Marshal(header)
	if err != nil {
		return "", errors.MarshalError.SetOrigin(err)
	}
	payloadJson, err := json.Marshal(user)
	if err != nil {
		return "", errors.MarshalError.SetOrigin(err)
	}

	headerPayload := base64.StdEncoding.EncodeToString(headerJson) + "." + base64.StdEncoding.EncodeToString(payloadJson)
	signature, Err := PasswdHash(headerPayload)
	if Err != nil {
		return "", Err
	}
	
	return base64.StdEncoding.EncodeToString([]byte(headerPayload + "." + signature)), nil
}

func CheckTokenBase64Signature(accessTokenBase64 string) *errors.Error {
	decodedAccessToken, err := base64.StdEncoding.DecodeString(accessTokenBase64)
	if err != nil {
		return errors.InvalidToken.SetHidden("Провал декодирования base64").SetOrigin(err)
	}
	tokenParts := strings.Split(string(decodedAccessToken), ".")
	if len(tokenParts) != 3 {
		return errors.InvalidToken.SetHidden("Токен должен состоять из 3 частей - но содержит "+strconv.Itoa(len(tokenParts)))
	}
	signature, Err := PasswdHash(tokenParts[0] + "." + tokenParts[1])
	if Err != nil {
		return Err
	}
	if signature != tokenParts[2] {
		return errors.InvalidToken.SetHidden("подпись содержит ошибку")
	}
	return nil
}

func CheckTokenPartsSignature(headerPayload, origSignature string) *errors.Error {
	signature, Err := PasswdHash(headerPayload)
	if Err != nil {
		return Err
	}
	if signature != origSignature {
		return errors.InvalidToken.SetHidden("подпись содержит ошибку")
	}
	return nil
}

/*
func parseHeaderBase64(headerBase64 string) (model.TokenHeader, *errors.Error) {
	var header model.TokenHeader
	decodedHeader, err := base64.StdEncoding.DecodeString(headerBase64)
	if err != nil {
		return header, errors.InvalidToken.SetHidden("Провал декодирования base64").SetOrigin(err)
	}
	if err = json.Unmarshal(decodedHeader, &header); err != nil {
		return header, errors.InvalidToken.SetHidden("Провал декодирования json").SetOrigin(err)
	}
	return header, nil
}
*/

func GetUserFromToken(accessTokenBase64 string) (model.User, *errors.Error) {
	var user model.User

	decodedAccessToken, err := base64.StdEncoding.DecodeString(accessTokenBase64)
	if err != nil {
		return user, errors.InvalidToken.SetHidden("Провал декодирования base64").SetOrigin(err)
	}
	tokenParts := strings.Split(string(decodedAccessToken), ".")
	if len(tokenParts) != 3 {
		return user, errors.InvalidToken.SetHidden("Токен должен состоять из 3 частей - но содержит "+strconv.Itoa(len(tokenParts)))
	}
	if Err := CheckTokenPartsSignature(tokenParts[0] + "." + tokenParts[1], tokenParts[2]); Err != nil {
		return user, Err
	}
	jsonUser, err := base64.StdEncoding.DecodeString(tokenParts[1])
	if err != nil {
		return user, errors.InvalidToken.SetHidden("Провал декодирования base64").SetOrigin(err)
	}
	
	if err = json.Unmarshal(jsonUser, &user); err != nil {
		return user, errors.InvalidToken.SetHidden("Провал декодирования json").SetOrigin(err)
	}
	return user, nil
}