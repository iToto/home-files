package userhdl

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/wingocard/braavos/internal/domain"
	"github.com/wingocard/braavos/internal/domain/service/usersvc"
	"github.com/wingocard/braavos/internal/rest/handler/form"
	"github.com/wingocard/braavos/internal/rest/handler/form/errs"
	"github.com/wingocard/braavos/internal/wlog"
	"github.com/wingocard/braavos/pkg/render"
	"github.com/wingocard/braavos/pkg/sanitize"
)

const (
	minTeenAgeYears = 13
)

type createTeenForm struct {
	FirstName    string    `json:"first_name" sanitize:"trim"`
	LastName     string    `json:"last_name" sanitize:"trim"`
	Email        string    `json:"email" sanitize:"trim,tolower"`
	Phone        string    `json:"phone" sanitize:"trim"`
	Password     string    `json:"password"`
	DateOfBirth  form.Date `json:"date_of_birth" sanitize:"trim"`
	ReferralCode string    `json:"referral_code" sanitize:"trim"`
	SignupToken  string    `json:"sign_up_token" sanitize:"trim"`
}

func validateTeenAge(v interface{}) error {
	dob, ok := v.(form.Date)
	if !ok {
		return errs.ErrInvalidDate
	}
	dobTime := time.Time(dob)
	minAgeTime := time.Now().AddDate(-minTeenAgeYears, 0, 0)

	if dobTime.After(minAgeTime) {
		return errs.ErrDateOfBirthAge
	}
	return nil
}

func (c createTeenForm) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(
			&c.FirstName,
			validation.Required.ErrorObject(errs.ErrRequired),
		),
		validation.Field(
			&c.LastName,
			validation.Required.ErrorObject(errs.ErrRequired),
		),
		validation.Field(
			&c.Email,
			validation.Required.ErrorObject(errs.ErrRequired),
			is.EmailFormat.ErrorObject(errs.ErrEmail),
		),
		validation.Field(
			&c.Phone,
			validation.Required.ErrorObject(errs.ErrRequired),
			form.IsNAPhoneNumber.ErrorObject(errs.ErrNAPhoneNumber),
		),
		validation.Field(
			&c.Password,
			validation.Required.ErrorObject(errs.ErrRequired),
			validation.Length(minPassLength, maxPassLength).ErrorObject(
				errs.NewPasswordLengthErr(minPassLength, maxPassLength),
			),
			form.IsPassword.ErrorObject(errs.ErrPassword),
		),
		validation.Field(
			&c.DateOfBirth,
			validation.Required.ErrorObject(errs.ErrRequired),
			validation.By(validateTeenAge),
		),
		validation.Field(
			&c.SignupToken,
			validation.Required.ErrorObject(errs.ErrRequired),
		),
		validation.Field(
			&c.ReferralCode,
			validation.By(validateReferralCode),
		),
	)
}

type createTeenResponse struct {
	ID             string             `json:"id"`
	Status         usersvc.UserStatus `json:"status"`
	Type           usersvc.UserType   `json:"type"`
	Email          string             `json:"email"`
	FirstName      string             `json:"first_name"`
	LastName       string             `json:"last_name"`
	Phone          string             `json:"phone"`
	DateOfBirth    form.DateTime      `json:"date_of_birth"`
	EmailValidated bool               `json:"email_validated"`
}

// CreateTeen handles a HTTP request for creating
// a teen user.
func CreateTeen(wl wlog.Logger, us usersvc.SVC) http.HandlerFunc { //nolint:dupl
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		wl := wlog.WithServiceRequest(ctx, wl, "user")

		req := &createTeenForm{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			render.BadRequest(ctx, wl, w, render.ErrJSONDecode)
			return
		}
		if err := sanitize.Struct(req); err != nil {
			render.InternalError(ctx, wl, w, err)
			return
		}
		if err := req.Validate(); err != nil {
			render.BadRequest(ctx, wl, w, err)
			return
		}

		teenDetails := &usersvc.TeenDetails{
			FirstName:    req.FirstName,
			LastName:     req.LastName,
			Password:     req.Password,
			Phone:        req.Phone,
			Email:        req.Email,
			DateOfBirth:  time.Time(req.DateOfBirth),
			ReferralCode: req.ReferralCode,
			SignupToken:  req.SignupToken,
		}

		u, err := us.CreateTeen(ctx, wl, teenDetails)
		if err != nil {
			if errors.Is(err, domain.ErrInvalidSignupToken) {
				wl.Info(err.Error())
				render.BadRequest(ctx, wl, w, render.NewErrorStr("invalid sign up token"))
				return
			}

			render.InternalError(ctx, wl, w, err)
			return
		}

		res := &createTeenResponse{
			ID:             u.ID,
			Status:         u.Status,
			Type:           u.Type,
			Email:          u.Email,
			FirstName:      u.FirstName,
			LastName:       u.LastName,
			Phone:          u.Phone,
			DateOfBirth:    form.DateTime(u.DateOfBirth),
			EmailValidated: u.EmailValidated,
		}
		render.JSON(ctx, wl, w, res, http.StatusCreated)
	}
}
