package middlewares

import (
	"context"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/ashish9868/rapidbackend/constants"
	"github.com/ashish9868/rapidbackend/core/repository"
	"github.com/ashish9868/rapidbackend/utils"
)

type Middleware func(http.Handler) http.Handler

func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range slices.Backward(middlewares) {
		handler = middleware(handler)
	}
	return handler
}

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				utils.LogF(
					"panic: %v %s %s",
					err,
					r.Method,
					r.URL.Path,
				)

				http.Error(
					w,
					"Internal Server Error",
					http.StatusInternalServerError,
				)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		utils.LogF(
			"%s %s %s",
			r.Method,
			r.URL.Path,
			time.Since(start),
		)
	})
}

func AuthMiddleware(repo *repository.AccessTokenRepository, redirect bool, throw bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accessToken := ""
			cookie, _ := r.Cookie("__auth")

			auth := r.Header.Get("Authorization")

			if !(len(auth) < 7 || !strings.EqualFold(auth[:7], "Bearer ")) {
				accessToken = strings.TrimSpace(auth[7:])
			}

			if len(accessToken) < 1 && cookie != nil {
				accessToken = cookie.Value
			}

			isHtmx := utils.IsHtmx(r)
			user := repo.GetUserFromToken(accessToken)
			if user != nil {
				newCtx := context.WithValue(r.Context(), constants.USER_CONTEXT_KEY, user)
				if strings.HasPrefix(r.URL.Path, "/login") {
					w.Header().Add("Hx-Redirect", "/dashboard")
					http.Redirect(w, r, "/dashboard", utils.IfElse(isHtmx, http.StatusOK, http.StatusTemporaryRedirect))
					return
				}
				next.ServeHTTP(w, r.WithContext(newCtx))
				return
			}
			if redirect {
				w.Header().Add("Hx-Redirect", "/dashboard")
				http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
				return
			} else if throw {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
