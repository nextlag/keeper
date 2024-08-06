package usecase

import (
	"context"
	"fmt"
	"net/mail"

	"github.com/nextlag/keeper/internal/entity"
	"github.com/nextlag/keeper/internal/utils"
	"github.com/nextlag/keeper/internal/utils/errs"
	"github.com/nextlag/keeper/pkg/logger/l"
)

const minutesPerHour = 60

// SignUpUser registers a new user with the provided email and hashed password.
// It validates the email format and hashes the password before storing it.
func (uc *UseCase) SignUpUser(ctx context.Context, email, password string) (user entity.User, err error) {
	if _, err = mail.ParseAddress(email); err != nil {
		err = errs.ErrWrongEmail
		return user, err
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return user, l.WrapErr(err)
	}

	return uc.repo.AddUser(ctx, email, hashedPassword)
}

// SignInUser authenticates a user with the provided email and password.
// It first validates the email format using mail.ParseAddress from the net/mail package.
// If the email format is incorrect, it returns ErrWrongEmail.
// It then retrieves the user from the repository using GetUserByEmail method.
// If an error occurs during user retrieval, it propagates the error.
// Upon successful authentication, it generates JWT tokens (access token and refresh token) for the user.
// The access token and refresh token are created using utils.CreateToken function with configured expiration times and private keys.
// It sets the max age, domain, and other properties for the JWT tokens based on the configuration.
// If any error occurs during token generation, it returns the error.
// Finally, it returns the generated JWT tokens and a nil error if the operation is successful.
func (uc *UseCase) SignInUser(ctx context.Context, email, password string) (token entity.JWT, err error) {
	if _, err = mail.ParseAddress(email); err != nil {
		err = errs.ErrWrongEmail
		return token, err
	}

	user, err := uc.repo.GetUserByEmail(ctx, email, password)
	if err != nil {
		return token, l.WrapErr(err)
	}

	token.AccessToken, err = utils.CreateToken(
		uc.cfg.Security.AccessTokenExpiresIn,
		user.ID,
		uc.cfg.Security.AccessTokenPrivateKey)
	if err != nil {
		return token, l.WrapErr(err)
	}

	token.RefreshToken, err = utils.CreateToken(
		uc.cfg.Security.RefreshTokenExpiresIn,
		user.ID,
		uc.cfg.Security.RefreshTokenPrivateKey)

	if err != nil {
		return token, l.WrapErr(err)
	}

	token.AccessTokenMaxAge = uc.cfg.Security.AccessTokenMaxAge * minutesPerHour
	token.RefreshTokenMaxAge = uc.cfg.Security.RefreshTokenMaxAge * minutesPerHour
	token.Domain = uc.cfg.Security.Domain

	return
}

// CheckAccessToken verifies the validity of the provided access token.
// If the token is valid, it retrieves and returns the associated user details.
// It first checks a local cache for the user corresponding to the access token.
// If not found in cache, it validates the token using a public key and retrieves
// the user details from the repository using the extracted userID from the token's subject.
// Upon successful validation, it caches the user details for future requests with the same token.
func (uc *UseCase) CheckAccessToken(ctx context.Context, accessToken string) (user entity.User, err error) {
	if userFromCache, found := uc.cache.Get(accessToken); found {
		checkedUser, ok := userFromCache.(entity.User)
		if ok {
			return checkedUser, nil
		}
	}

	sub, err := utils.ValidToken(accessToken, uc.cfg.Security.AccessTokenPublicKey)
	if err != nil {
		err = errs.ErrTokenValidation
		return user, err
	}

	userID := fmt.Sprint(sub)
	user, err = uc.repo.GetUserByID(ctx, userID)
	if err != nil {
		err = errs.ErrTokenValidation
		return user, err
	}

	uc.cache.Set(accessToken, user)
	return
}

// RefreshAccessToken validates the provided refresh token, retrieves the corresponding user,
// and generates a new access token.
func (uc *UseCase) RefreshAccessToken(ctx context.Context, refreshToken string) (token entity.JWT, err error) {

	userID, err := utils.ValidToken(refreshToken, uc.cfg.Security.RefreshTokenPublicKey)
	if err != nil {
		err = errs.ErrTokenValidation
		return token, err
	}

	user, err := uc.repo.GetUserByID(ctx, fmt.Sprint(userID))
	if err != nil {
		err = errs.ErrTokenValidation
		return token, err
	}

	token.RefreshToken = refreshToken
	token.AccessToken, err = utils.CreateToken(
		uc.cfg.Security.AccessTokenExpiresIn,
		user.ID,
		uc.cfg.Security.AccessTokenPrivateKey)

	if err != nil {
		return token, l.WrapErr(err)
	}

	token.AccessTokenMaxAge = uc.cfg.Security.AccessTokenMaxAge * minutesPerHour
	token.RefreshTokenMaxAge = uc.cfg.Security.RefreshTokenMaxAge * minutesPerHour
	token.Domain = uc.cfg.Security.Domain
	return
}
