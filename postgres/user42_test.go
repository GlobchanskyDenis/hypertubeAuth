package postgres

import (
	"HypertubeAuth/errors"
	"HypertubeAuth/logger"
	"HypertubeAuth/model"
	"testing"
	"time"
)

func TestSetUser42(t *testing.T) {
	var (
		user1 = &model.User42{}
		user2 = &model.User42{}
	)
	user1.UserId = 42
	user2.UserId = 21
	accessToken := "access_token"
	user1.AccessToken = &accessToken
	user2.AccessToken = nil
	refreshToken := "refresh_token"
	user1.RefreshToken = &refreshToken
	user2.RefreshToken = nil
	t1 := time.Now()
	user1.ExpiresAt = &t1
	user2.ExpiresAt = nil

	initTest(t)

	defer func(t *testing.T, user1, user2 *model.User42) {
		t.Run("Delete test user #1", func(t_ *testing.T) {
			if Err := UserDelete42(user1); Err != nil {
				t_.Errorf("%sError: cannot delete test user - %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
			} else {
				t_.Logf("%sSuccess%s", logger.GREEN_BG, logger.NO_COLOR)
			}
		})
		t.Run("Delete test user #2", func(t_ *testing.T) {
			if Err := UserDelete42(user2); Err != nil {
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
		if Err := UserSet42(user1); Err != nil {
			t_.Errorf("%sError: %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
		} else {
			t_.Logf("%sSuccess: user was created successfully%s", logger.GREEN_BG, logger.NO_COLOR)
		}
	})

	// /!!!!!
	t.Run("valid update user #1", func(t_ *testing.T) {
		newAccessToken := "new access_token"
		user1.AccessToken = &newAccessToken
		if Err := UserUpdate42(user1); Err != nil {
			t_.Errorf("%sError: %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
		} else {
			t_.Logf("%sSuccess: user was created successfully%s", logger.GREEN_BG, logger.NO_COLOR)
		}
	})

	t.Run("invalid get user #1 by id", func(t_ *testing.T) {
		_, Err := UserGet42ById(0)
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
		newUser, Err := UserGet42ById(user1.UserId)
		if Err != nil {
			t_.Errorf("%sError: %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
		} else if user1.UserId != newUser.UserId || newUser.AccessToken == nil || 
		newUser.RefreshToken == nil || newUser.ExpiresAt == nil || *user1.AccessToken != *newUser.AccessToken ||
		*user1.RefreshToken != *newUser.RefreshToken || user1.ExpiresAt.Format(time.StampMilli) != newUser.ExpiresAt.Format(time.StampMilli) {
			t_.Errorf("%sError: received user differs from original%s\nexpected %#v\ngot %#v", logger.RED_BG, logger.NO_COLOR,
				user1.User42Model, newUser.User42Model)
			t_.Errorf("%s\n%s", user1.ExpiresAt.Format(time.StampMilli), newUser.ExpiresAt.Format(time.StampMilli))
		} else {
			t_.Logf("%sSuccess: user was received successfully%s", logger.GREEN_BG, logger.NO_COLOR)
		}
	})

	t.Run("valid create user #2", func(t_ *testing.T) {
		if Err := UserSet42(user2); Err != nil {
			t_.Errorf("%sError: %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
		} else {
			t_.Logf("%sSuccess: user was created successfully%s", logger.GREEN_BG, logger.NO_COLOR)
		}
	})

	t.Run("invalid create user #1", func(t_ *testing.T) {
		if Err := UserSet42(user1); Err != nil {
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
		if Err := UserDelete42(user2); Err != nil {
			t_.Errorf("%sError: %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
		} else {
			t_.Logf("%sSuccess: user was deleted successfully%s", logger.GREEN_BG, logger.NO_COLOR)
		}
	})

	// /!!!!!
	t.Run("invalid update user #2", func(t_ *testing.T) {
		newAccessToken := "new access_token"
		user2.AccessToken = &newAccessToken
		if Err := UserUpdate42(user2); Err != nil {
			if errors.UserNotExist.IsOverlapWithError(Err) {
				t_.Logf("%sSuccess: user not exists as it expected%s", logger.GREEN_BG, logger.NO_COLOR)
			} else {
				t_.Errorf("%sError: %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
			}
		} else {
			t_.Errorf("%sError: expected but not found error%s", logger.RED_BG, logger.NO_COLOR)
		}
	})

	t.Run("invalid user delete #2", func(t_ *testing.T) {
		if Err := UserDelete42(user2); Err != nil {
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
		if Err := UserSet42(user2); Err != nil {
			t_.Errorf("%sError: %s%s", logger.RED_BG, Err.Error(), logger.NO_COLOR)
		} else {
			t_.Logf("%sSuccess: user was created successfully%s", logger.GREEN_BG, logger.NO_COLOR)
		}
	})

	t.Run("invalid recreate user #2", func(t_ *testing.T) {
		if Err := UserSet42(user2); Err != nil {
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

	if !t.Failed() {
		t.Logf("%sSuccess%s", logger.GREEN_BG, logger.NO_COLOR)
	}
}
