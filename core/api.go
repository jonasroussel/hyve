package core

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jonasroussel/hyve/tools"
	"gorm.io/gorm"
)

type ContextKey string

const SESSION_CTX ContextKey = "session"

func ServeAPI(rootHandler *tools.HTTPHandler, app *HyveApp) {
	apiHandler := http.NewServeMux()

	apiHandler.HandleFunc("POST /accounts/root", createRootAccount(app))
	apiHandler.Handle("GET /accounts/me", mustBeConnected(app, http.HandlerFunc(getConnectedAccount())))

	apiHandler.HandleFunc("POST /auth/login", authLogin(app))
	apiHandler.Handle("POST /auth/logout", mustBeConnected(app, http.HandlerFunc(authLogout(app))))

	apiHandler.Handle("GET /dashboard", mustBeConnected(app, http.HandlerFunc(getDashboard(app))))

	rootHandler.Handle("/api", apiHandler)
}

func getDashboard(app *HyveApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if app.Store == nil {
			tools.DetailedError(w, http.StatusConflict, "NO_STORE")
			return
		}

		tools.JSONResp(w, http.StatusOK, map[string]any{
			"dashboard_id": 0,
		})
	}
}

func mustBeConnected(app *HyveApp, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var rootAccount DBAccount
		err := app.DB.Where("role = ?", DBRole_SuperAdmin).Select("id").First(&rootAccount).Error
		if err == gorm.ErrRecordNotFound {
			tools.DetailedError(w, http.StatusUnauthorized, "NO_ROOT_ACCOUNT")
			return
		}

		cookie, err := r.Cookie("session-token")
		if cookie == nil || err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		sessionID, err := uuid.Parse(cookie.Value)
		if err != nil {
			http.SetCookie(w, &http.Cookie{
				Name:     "session-token",
				Value:    "",
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
				MaxAge:   -1,
			})
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var session DBSession

		err = app.DB.Where("id = ? AND expires_at > ?", sessionID, time.Now().UTC()).First(&session).Error
		if err != nil {
			http.SetCookie(w, &http.Cookie{
				Name:     "session-token",
				Value:    "",
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
				MaxAge:   -1,
			})
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		session.Account.Password = ""

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), SESSION_CTX, &session)))
	})
}
