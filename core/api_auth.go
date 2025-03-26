package core

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jonasroussel/hyve/tools"
	"golang.org/x/crypto/bcrypt"
)

// @POST /api/auth/login
type authLoginBody struct {
	UsernameEmail string `json:"username_email"`
	Password      string `json:"password"`
}

func authLogin(app *HyveApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body authLoginBody
		if err := tools.ParseBody(r.Body, &body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var account DBAccount
		err := app.DB.Where("username = ? OR email = ?", body.UsernameEmail, body.UsernameEmail).First(&account).Error
		if err != nil {
			tools.DetailedError(w, http.StatusUnauthorized, "WRONG_CREDENTIALS")
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(body.Password))
		if err != nil {
			tools.DetailedError(w, http.StatusUnauthorized, "WRONG_CREDENTIALS")
			return
		}

		sessionID, err := uuid.NewRandom()
		if err != nil {
			tools.DetailedError(w, http.StatusInternalServerError, "UNABLE_TO_CREATE_SESSION")
			return
		}

		expiresIn := time.Hour * 24

		session := DBSession{
			ID:        sessionID.String(),
			AccountID: account.ID,
			CreatedAt: time.Now().UTC(),
			ExpiresAt: time.Now().UTC().Add(expiresIn),
		}
		err = app.DB.Create(&session).Error
		if err != nil {
			tools.DetailedError(w, http.StatusInternalServerError, "UNABLE_TO_CREATE_SESSION")
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session-token",
			Value:    session.ID,
			HttpOnly: true,
			Path:     "/",
			SameSite: http.SameSiteStrictMode,
			MaxAge:   int(expiresIn.Seconds()),
		})

		tools.JSONResp(w, http.StatusOK, account)
	}
}

// @POST /api/auth/logout
func authLogout(app *HyveApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := r.Context().Value(SESSION_CTX).(*DBSession)
		if session != nil {
			app.DB.Delete(&session)
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session-token",
			Value:    "",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   -1,
		})

		tools.JSONResp(w, http.StatusOK, nil)
	}
}
