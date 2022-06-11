package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"git.eletrotupi.com/momento/database"
)

type Account struct {
	Email string
	Password string
	Bio string
	Url string
}

func New() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/profile", handleProfile)
	mux.HandleFunc("/api/register", handleRegister)

	return WithAuth(mux)
}

func handleProfile(w http.ResponseWriter, r *http.Request) {
	user := Auth(r.Context())
	encoder := json.NewEncoder(w)

	err := encoder.Encode(user)
	if err != nil {
		panic(err)
	}
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)

		return
	}

	var acc Account
	ctx := r.Context()

	// TODO: Add some security measures and deal gracefully with problems like
	// malformed body, etc
	err := json.NewDecoder(r.Body).Decode(&acc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if acc.Email == "" || acc.Password == "" {
		http.Error(w, "Missing required fields", http.StatusUnprocessableEntity)
		return
	}

	if !strings.ContainsRune(acc.Email, '@') {
		http.Error(w, "Invalid Email Address", http.StatusUnprocessableEntity)
		return
	}

	pwhash, err := bcrypt.GenerateFromPassword([]byte(acc.Password), bcrypt.DefaultCost)

	if err != nil {
		http.Error(w, "Problem when hashing", http.StatusInternalServerError)
		return
	}

	var userID int
	if err := database.WithTx(ctx, nil, func(tx *sql.Tx) error {
		row := tx.QueryRowContext(ctx, `
		INSERT INTO users (
			created_at, email, password, bio, url
		) VALUES (
			NOW() at time zone 'utc',
			$1, $2, $3, $4
		)
		RETURNING id;
		`, acc.Email, string(pwhash), acc.Bio, acc.Url)
		// TODO: Detect duplicate users
		return row.Scan(&userID)
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
