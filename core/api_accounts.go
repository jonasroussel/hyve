package core

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jonasroussel/hyve/tools"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// @POST /api/accounts/root
type createRootAccountBody struct {
	Username string  `json:"username"`
	Email    *string `json:"email,omitempty"`
	Password string  `json:"password"`
}

func createRootAccount(app *HyveApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body createRootAccountBody
		if err := tools.ParseBody(r.Body, &body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var account DBAccount
		err := app.DB.Where("role = ?", DBRole_SuperAdmin).Select("id").First(&account).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err == nil {
			tools.DetailedError(w, http.StatusConflict, "ROOT_ACCOUNT_EXISTS", "Root account already exists")
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = app.DB.Create(&DBAccount{
			Username: body.Username,
			Email:    body.Email,
			Password: string(hashedPassword),
			Role:     DBRole_SuperAdmin,
		}).Error
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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

		tools.JSONResp(w, http.StatusCreated, nil)
	}
}

// @GET /api/accounts/me
func getConnectedAccount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := r.Context().Value(SESSION_CTX).(*DBSession)

		tools.JSONResp(w, http.StatusOK, session.Account)
	}
}
