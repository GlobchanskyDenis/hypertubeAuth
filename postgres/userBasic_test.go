package postgres

import (
	"HypertubeAuth/errors"
	"HypertubeAuth/logger"
	"HypertubeAuth/model"
	"testing"
)

func TestSetUserBasic(t *testing.T) {
	var (
		user1 = &model.UserBasic{}
		user2 = &model.UserBasic{}
	)
	user1.Email = emailValid1
	user1.EncryptedPass = encryptedPass
	user1.Fname = "Denis"
	user1.Lname = "Globchansky"
	user1.Displayname = displayName
	user1.EmailConfirmHash = emailConfirmHash
	user2.Email = emailValid2
	user2.EncryptedPass = encryptedPass
	user2.Displayname = displayName
	user2.EmailConfirmHash = emailConfirmHash

	initTest(t)

	defer func(t *testing.T, user1, user2 *model.UserBasic) {
		t.Run("Delete test user #1", func(t_ *testing.T) {
			if Err := UserDeleteBasic(user1); Err != nil {
				t_.Errorf("%sError: cannot delete test user - %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
			} else {
				t_.Logf("%sSuccess%s", logger.GREEN_BG, logger.NO_COLOR)
			}
		})
		t.Run("Delete test user #2", func(t_ *testing.T) {
			if Err := UserDeleteBasic(user2); Err != nil {
				t_.Errorf("%sError: cannot delete test user - %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
			} else {
				t_.Logf("%sSuccess%s", logger.GREEN_BG, logger.NO_COLOR)
			}
		})
		t.Run("Close connection", func(t_ *testing.T) {
			if Err := Close(); Err != nil {
				t_.Errorf("%sError: cannot close connection - %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
			} else {
				t_.Logf("%sSuccess%s", logger.GREEN_BG, logger.NO_COLOR)
			}
		})
	}(t, user1, user2)

	t.Run("valid create user #1", func(t_ *testing.T) {
		if Err := UserSetBasic(user1); Err != nil {
			t_.Errorf("%sError: %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
		} else {
			t_.Logf("%sSuccess: user was created successfully%s", logger.GREEN_BG, logger.NO_COLOR)
		}
	})

	t.Run("valid get user #1 by email", func(t_ *testing.T) {
		newUser, Err := UserGetBasicByEmail(user1.Email)
		if Err != nil {
			t_.Errorf("%sError: %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
		} else if user1.UserId != newUser.UserId || user1.Email != newUser.Email ||
			user1.EncryptedPass != newUser.EncryptedPass || user1.Fname != user1.Fname ||
			user1.Lname != newUser.Lname || user1.ImageBody != newUser.ImageBody ||
			user1.Displayname != newUser.Displayname || user1.EmailConfirmHash != newUser.EmailConfirmHash {
			t_.Errorf("%sError: received user differs from original%s\nexpected %#v got %#v", logger.RED_BG, logger.NO_COLOR,
				user1, newUser)
		} else {
			t_.Logf("%sSuccess: user was received successfully%s", logger.GREEN_BG, logger.NO_COLOR)
		}
	})

	t.Run("invalid get user #1 by email", func(t_ *testing.T) {
		_, Err := UserGetBasicByEmail("not existing email")
		if Err != nil {
			if errors.UserNotExist.IsOverlapWithError(Err) {
				t_.Logf("%sSuccess: user not exists as it expected%s", logger.GREEN_BG, logger.NO_COLOR)
			} else {
				t_.Errorf("%sError: %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
			}
		} else {
			t_.Errorf("%sError: expected but not found error%s", logger.RED_BG, logger.NO_COLOR)
		}
	})

	t.Run("invalid get user #1 by id", func(t_ *testing.T) {
		_, Err := UserGetBasicById(0)
		if Err != nil {
			if errors.UserNotExist.IsOverlapWithError(Err) {
				t_.Logf("%sSuccess: user not exists as it expected%s", logger.GREEN_BG, logger.NO_COLOR)
			} else {
				t_.Errorf("%sError: %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
			}
		} else {
			t_.Errorf("%sError: expected but not found error%s", logger.RED_BG, logger.NO_COLOR)
		}
	})

	t.Run("valid get user #1 by user id", func(t_ *testing.T) {
		newUser, Err := UserGetBasicById(user1.UserId)
		if Err != nil {
			t_.Errorf("%sError: %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
		} else if user1.UserId != newUser.UserId || user1.Email != newUser.Email ||
			user1.EncryptedPass != newUser.EncryptedPass || user1.Fname != user1.Fname ||
			user1.Lname != newUser.Lname || user1.ImageBody != newUser.ImageBody ||
			user1.Displayname != newUser.Displayname || user1.EmailConfirmHash != newUser.EmailConfirmHash {
			t_.Errorf("%sError: received user differs from original%s\nexpected %#v got %#v", logger.RED_BG, logger.NO_COLOR,
				user1, newUser)
		} else {
			t_.Logf("%sSuccess: user was received successfully%s", logger.GREEN_BG, logger.NO_COLOR)
		}
	})

	t.Run("valid create user #2", func(t_ *testing.T) {
		if Err := UserSetBasic(user2); Err != nil {
			t_.Errorf("%sError: %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
		} else {
			t_.Logf("%sSuccess: user was created successfully%s", logger.GREEN_BG, logger.NO_COLOR)
		}
	})

	t.Run("invalid create user #1", func(t_ *testing.T) {
		if Err := UserSetBasic(user1); Err != nil {
			if errors.ImpossibleToExecute.IsOverlapWithError(Err) {
				t_.Logf("%sSuccess: found error that was expected%s", logger.GREEN_BG, logger.NO_COLOR)
			} else {
				t_.Errorf("%sError: expected %s found %s error%s", logger.RED_BG,
					errors.ImpossibleToExecute, Err.Error(), logger.NO_COLOR)
			}
		} else {
			t_.Errorf("%sError: expected but not found error%s", logger.RED_BG, logger.NO_COLOR)

		}
	})

	t.Run("valid user delete #2", func(t_ *testing.T) {
		if Err := UserDeleteBasic(user2); Err != nil {
			t_.Errorf("%sError: %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
		} else {
			t_.Logf("%sSuccess: user was deleted successfully%s", logger.GREEN_BG, logger.NO_COLOR)
		}
	})

	t.Run("invalid user delete #2", func(t_ *testing.T) {
		if Err := UserDeleteBasic(user2); Err != nil {
			if errors.ImpossibleToExecute.IsOverlapWithError(Err) {
				t_.Logf("%sSuccess: found error that was expected%s", logger.GREEN_BG, logger.NO_COLOR)
			} else {
				t_.Errorf("%sError: expected %s found %s error%s", logger.RED_BG,
					errors.ImpossibleToExecute, Err.Error(), logger.NO_COLOR)
			}
		} else {
			t_.Errorf("%sError: expected but not found error%s", logger.RED_BG, logger.NO_COLOR)

		}
	})

	t.Run("valid recreate user #2", func(t_ *testing.T) {
		if Err := UserSetBasic(user2); Err != nil {
			t_.Errorf("%sError: %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
		} else {
			t_.Logf("%sSuccess: user was created successfully%s", logger.GREEN_BG, logger.NO_COLOR)
		}
	})

	t.Run("invalid recreate user #2", func(t_ *testing.T) {
		if Err := UserSetBasic(user2); Err != nil {
			if errors.ImpossibleToExecute.IsOverlapWithError(Err) {
				t_.Logf("%sSuccess: found error that was expected%s", logger.GREEN_BG, logger.NO_COLOR)
			} else {
				t_.Errorf("%sError: expected %s found %s error%s", logger.RED_BG,
					errors.ImpossibleToExecute, Err.Error(), logger.NO_COLOR)
			}
		} else {
			t_.Errorf("%sError: expected but not found error%s", logger.RED_BG, logger.NO_COLOR)
		}
	})

	t.Run("valid update user #2", func(t_ *testing.T) {
		if Err := UserUpdateBasic(user2); Err != nil {
			t_.Errorf("%sError: %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
		} else {
			t_.Logf("%sSuccess: user was updated successfully%s", logger.GREEN_BG, logger.NO_COLOR)
		}
	})

	t.Run("valid get user #2 by id, check update status", func(t_ *testing.T) {
		newUser, Err := UserGetBasicById(user2.UserId)
		if Err != nil {
			t_.Errorf("%sError: %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
		} else if newUser.IsEmailConfirmed != true {
			t_.Errorf("%sError: user email confirm status didnt change%s", logger.RED_BG, logger.NO_COLOR)
		} else {
			t_.Logf("%sSuccess: user was received and checked email confirm status%s", logger.GREEN_BG, logger.NO_COLOR)
		}
	})

	if !t.Failed() {
		t.Logf("%sSuccess%s", logger.GREEN_BG, logger.NO_COLOR)
	}
}
