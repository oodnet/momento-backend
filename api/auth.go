package api

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"git.eletrotupi.com/momento/database"
)

type User struct {
	ID			int       `json:"id"`
	CreatedAt	time.Time `json:"created_at"`
	Email		string    `json:"email"`
}

var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

func Context(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userCtxKey, user)
}

func Auth(ctx context.Context) *User {
	user, ok := ctx.Value(userCtxKey).(*User)

	if !ok {
		panic(errors.New("Invalid authentication context"))
	}

	return user
}

func WithAuth(h http.Handler) http.Handler {
    middleware := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/register" {
			h.ServeHTTP(w, r)
			return
		}

		// TODO: Replace this with the final auth mechanism
		email, password, ok := r.BasicAuth()
		if !ok {
			http.Error(w, "Authorization required", http.StatusUnauthorized)
			return
		}

		var (
			user         User
			pwhash       string
		)

		if err := database.WithTx(r.Context(), &sql.TxOptions{
			Isolation: 0,
			ReadOnly:  true,
		}, func(tx *sql.Tx) error {
			row := tx.QueryRowContext(r.Context(), `
				SELECT
					id, created_at, email, password
				FROM users
				WHERE email = $1;
			`, email)
			return row.Scan(&user.ID, &user.CreatedAt, &user.Email, &pwhash)
		}); err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Invalid username or password", http.StatusUnauthorized)
				return
			}

			panic(err)
		}

		err := bcrypt.CompareHashAndPassword([]byte(pwhash), []byte(password))

		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		ctx := Context(r.Context(), &user)
		r = r.WithContext(ctx)
        h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(middleware)
}

