package authentication

import (
	"checkpoint/modules/account"
	utils "checkpoint/utils"

	"github.com/graph-gophers/graphql-go"
)

type AuthenticationResolver struct{}

var authenticationService = AuthenticationService{}

func (AuthenticationResolver) Signup(args SignUpData) (*account.Account, error) {
	accountInfo, code, err := authenticationService.SignUp(SignUpData{
		Username: args.Username,
		Password: args.Password,
	})

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    code,
			Message: utils.ForbiddenOperation.Error(),
		}
	}

	return &account.Account{
		Id:        graphql.ID(accountInfo.ID.String()),
		Username:  accountInfo.Username,
		CreatedAt: graphql.Time{Time: accountInfo.CreatedAt},
	}, nil
}

func (AuthenticationResolver) Signin(args SignInData) (*Authentication, error) {
	auth, code, err := authenticationService.SignIn(SignInData{
		Username: args.Username,
		Password: args.Password,
	})

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    code,
			Message: utils.ForbiddenOperation.Error(),
		}
	}

	return auth, nil
}

func (AuthenticationResolver) RefreshToken(args RefreshTokenData) (*Authentication, error) {

	auth, code , err := authenticationService.GetAuthenticationTokenByRefreshToken(RefreshTokenData{
		RefreshToken: args.RefreshToken,
	})

	if err != nil {
		return nil, utils.GraphqlError{
			Code:    code,
			Message: utils.ForbiddenOperation.Error(),
		}
	}

	return auth, nil
}
