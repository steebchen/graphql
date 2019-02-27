package mutation

import (
	"context"
	"github.com/pkg/errors"
	"github.com/steebchen/graphql/gqlgen"
	"github.com/steebchen/graphql/lib/auth"
	"github.com/steebchen/graphql/lib/session_cookie"
	"github.com/steebchen/graphql/prisma"
	"golang.org/x/crypto/bcrypt"
)

var UserNotFoundError = errors.New("user not found")
var IncorrectPasswordError = errors.New("password is incorrect")

func (m *Mutation) Login(ctx context.Context, email string, password string) (gqlgen.LoginResult, error) {
	user, err := m.Prisma.User(prisma.UserWhereUniqueInput{
		Email: &email,
	}).Exec(ctx)

	if err == prisma.ErrNoResult {
		return gqlgen.LoginResult{}, UserNotFoundError
	}

	if err != nil {
		panic(err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return gqlgen.LoginResult{}, IncorrectPasswordError
	}

	session, err := m.Prisma.CreateSession(prisma.SessionCreateInput{
		User: prisma.UserCreateOneWithoutSessionsInput{
			Connect: &prisma.UserWhereUniqueInput{
				ID: &user.ID,
			},
		},
		Token: auth.GenerateToken(),
	}).Exec(ctx)

	session_cookie.Set(ctx, session)

	return gqlgen.LoginResult{
		Session: *session,
		User:    *user,
	}, nil
}
