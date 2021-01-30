package hash

import (
	"HypertubeAuth/configurator"
	"HypertubeAuth/errors"
	"HypertubeAuth/logger"
	"HypertubeAuth/model"
	"testing"
)

func TestHash(t *testing.T) {
	if Err := Init("../../conf.json"); Err != nil {
		t.Errorf("%sError during initialize package - %s%s" , logger.RED_BG, Err.Error(), logger.NO_COLOR)
		t.FailNow()
	}

	t.Run("check for correct signature", func(t_ *testing.T){
		var user = &model.UserBasic{}
		user.UserId = 42
		user.Email = "school21@gmail.com"
		accessToken, Err := CreateToken(user)
		if Err != nil {
			t_.Errorf("%sError during creating token - %s%s" , logger.RED_BG, Err.Error(), logger.NO_COLOR)
			t_.FailNow()
		}

		if Err = CheckTokenBase64Signature(accessToken); Err != nil {
			t_.Errorf("%sError cannot unmarshal token - %s%s" , logger.RED_BG, Err.Error(), logger.NO_COLOR)
			t_.FailNow()
		}
		t_.Logf("%sSuccess: token is valid%s", logger.GREEN_BG, logger.NO_COLOR)
	})

	t.Run("check token header data validity", func(t_ *testing.T){
		var user = &model.UserBasic{}
		user.UserId = 42
		user.Email = "school21@gmail.com"
		imageBody := "image_body"
		user.ImageBody = &imageBody
		user.Displayname = "skinnyman"
		fname := "Den"
		user.Fname = &fname
		lname := "QWERTY"
		user.Lname = &lname
		accessToken, Err := CreateToken(user)
		if Err != nil {
			t_.Errorf("%sError during creating token - %s%s" , logger.RED_BG, Err.Error(), logger.NO_COLOR)
			t_.FailNow()
		}

		header, Err := GetHeaderFromToken(accessToken)
		if Err != nil {
			t_.Errorf("%sError cannot unmarshal token - %s%s" , logger.RED_BG, Err.Error(), logger.NO_COLOR)
			t_.FailNow()
		}

		if header.UserId != user.UserId {
			t_.Errorf("%sError: UserId are incorrect after decoding. Expected %d Got %d%s" , logger.RED_BG,
				user.UserId, header.UserId, logger.NO_COLOR)
		}

		if !t_.Failed() {
			t_.Logf("%sSuccess: token is valid%s", logger.GREEN_BG, logger.NO_COLOR)
		}
	})
	
}

func Init(configFileName string) *errors.Error {
	print("Считываю конфигурационный файл\t\t- ")
	if err := configurator.SetConfigFile(configFileName); err != nil {
		println(logger.RED + "ошибка" + logger.NO_COLOR)
		return errors.ConfigurationFail.SetArgs("Не могу считать файл "+configFileName,
			"Cant read file "+configFileName).SetOrigin(err)
	}
	println(logger.GREEN + "успешно" + logger.NO_COLOR)
	/*
	**	logger
	 */
	print("Настраиваю пакет logger\t\t\t- ")
	cfgLogger := logger.GetConfig()
	if err := configurator.ParsePackageConfig(cfgLogger, "logger"); err != nil {
		println(logger.RED + "ошибка" + logger.NO_COLOR)
		return errors.ConfigurationFail.SetArgs("logger", "logger").SetOrigin(err)
	}
	println(logger.GREEN + "успешно" + logger.NO_COLOR)
	/*
	**	hash
	 */
	print("Настраиваю пакет hash\t\t\t- ")
	cfgHash := GetConfig()
	if err := configurator.ParsePackageConfig(cfgHash, "hash"); err != nil {
		println(logger.RED + "ошибка" + logger.NO_COLOR)
		return errors.ConfigurationFail.SetArgs("controller/hash", "controller/hash").SetOrigin(err)
	}
	println(logger.GREEN + "успешно" + logger.NO_COLOR)
	return nil
} 