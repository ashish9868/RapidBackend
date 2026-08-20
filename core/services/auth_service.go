package services

import (
	"net/http"

	"github.com/ashish9868/rapidbackend/core"
	respository "github.com/ashish9868/rapidbackend/core/repository"
	"github.com/ashish9868/rapidbackend/dto"
	"github.com/ashish9868/rapidbackend/models"
	"github.com/ashish9868/rapidbackend/utils"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type authService struct {
	App *core.App
}

func NewAuthService(app *core.App) *authService {
	return &authService{App: app}
}

func (a *authService) ValidateLogin(w http.ResponseWriter, r *http.Request, checkSuperadmin bool) (*dto.LoginForm, *models.AccessKeyToken, map[string]any) {
	form := &dto.LoginForm{}
	if err := a.App.BindSafely(w, r, form); err != nil {
		return form, nil, a.App.FormatErrors(err)
	}
	err := validation.ValidateStruct(form,
		validation.Field(&form.Email,
			validation.Required.Error("Email is required"),
			is.Email.Error("Please Provide a valid Email"),
		),
		validation.Field(&form.Password,
			validation.Required.Error("Password is required"),
		),
	)
	collection := utils.IFElse(checkSuperadmin,
		respository.COLLECTION_SUPERADMINS,
		respository.COLLECTION_USERS,
	)
	if err == nil {

		user := &models.User{
			Collection: collection,
		}
		a.App.BaseRepository.GetByColumn(
			collection,
			"email",
			form.Email,
			user,
		)
		if utils.IsTruthy(user.ID) && utils.CheckPassword(user.Password, form.Password) && user.EmailVerifiedAt != nil && user.IsActive {
			token := a.App.AccessTokenRepository.CreateNewAccessToken(user)
			if token != nil {
				return form, token, nil
			}
		}
		return form, nil, map[string]any{
			"global": "Email and/or password is invalid.",
		}
	}
	return form, nil, a.App.FormatErrors(err)
}
