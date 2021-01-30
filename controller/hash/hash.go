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

func PasswdHash(pass string) (*string, *errors.Error) {
	conf, Err := getConfig()
	if Err != nil {
		return nil, Err
	}
	pass += conf.PasswdSalt
	crcH := crc32.ChecksumIEEE([]byte(pass))
	passHash := strconv.FormatUint(uint64(crcH), 20)
	return &passHash, nil
}

func CreateTokenSignature(header string) (string, *errors.Error) {
	conf, Err := getConfig()
	if Err != nil {
		return "", Err
	}
	header += conf.PasswdSalt
	crcH := crc32.ChecksumIEEE([]byte(header))
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

func CreateToken(user *model.UserBasic) (string, *errors.Error) {
	var header model.TokenHeader
	header.UserId = user.UserId

	headerJson, err := json.Marshal(header)
	if err != nil {
		return "", errors.MarshalError.SetOrigin(err)
	}
	headerBase64 := base64.StdEncoding.EncodeToString(headerJson)
	signature, Err := CreateTokenSignature(headerBase64)
	if Err != nil {
		return "", Err
	}
	return base64.StdEncoding.EncodeToString([]byte(headerBase64 + "." + signature)), nil
}

func CheckTokenBase64Signature(accessTokenBase64 string) *errors.Error {
	decodedAccessToken, err := base64.StdEncoding.DecodeString(accessTokenBase64)
	if err != nil {
		return errors.InvalidToken.SetHidden("Провал декодирования base64").SetOrigin(err)
	}
	tokenParts := strings.Split(string(decodedAccessToken), ".")
	if len(tokenParts) != 2 {
		return errors.InvalidToken.SetHidden("Токен должен состоять из 2 частей - но содержит "+strconv.Itoa(len(tokenParts)))
	}
	signature, Err := CreateTokenSignature(tokenParts[0])
	if Err != nil {
		return Err
	}
	if signature != tokenParts[1] {
		return errors.InvalidToken.SetHidden("подпись содержит ошибку")
	}
	return nil
}

func CheckTokenPartsSignature(header, origSignature string) *errors.Error {
	signature, Err := CreateTokenSignature(header)
	if Err != nil {
		return Err
	}
	if signature != origSignature {
		return errors.InvalidToken.SetHidden("подпись содержит ошибку")
	}
	return nil
}

func GetHeaderFromToken(accessTokenBase64 string) (model.TokenHeader, *errors.Error) {
	var header model.TokenHeader

	decodedAccessToken, err := base64.StdEncoding.DecodeString(accessTokenBase64)
	if err != nil {
		return header, errors.InvalidToken.SetHidden("Провал декодирования base64").SetOrigin(err)
	}
	tokenParts := strings.Split(string(decodedAccessToken), ".")
	if len(tokenParts) != 2 {
		return header, errors.InvalidToken.SetHidden("Токен должен состоять из 2 частей - но содержит "+strconv.Itoa(len(tokenParts)))
	}

	if Err := CheckTokenPartsSignature(tokenParts[0], tokenParts[1]); Err != nil {
		return header, Err
	}

	decodedHeader, err := base64.StdEncoding.DecodeString(tokenParts[0])
	if err != nil {
		return header, errors.InvalidToken.SetHidden("Провал декодирования base64").SetOrigin(err)
	}
	if err = json.Unmarshal(decodedHeader, &header); err != nil {
		return header, errors.InvalidToken.SetHidden("Провал декодирования json").SetOrigin(err)
	}
	return header, nil
}
