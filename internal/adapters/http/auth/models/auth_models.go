package models

import (
	pkgValidator "villainrsty-ecommerce-server/pkg/validator"
)

type (
	LoginRequest struct {
		Email      string `json:"email" validate:"required,email"`
		Password   string `json:"password" validate:"required,min=8"`
		RememberMe bool   `json:"remember_me"`
	}

	Login2FARequest struct {
		Email      string `json:"email" validate:"required,email"`
		Password   string `json:"password" validate:"required,min=8"`
		RememberMe bool   `json:"remember_me"`
	}

	VerifyLogin2FARequest struct {
		ChallengeID string `json:"challenge_id" validate:"required"`
		OTPCode     string `json:"otp_code" validate:"required,len=6"`
		RememberMe  bool   `json:"remember_me"`
	}

	RegisterRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
		Name     string `json:"name" validate:"required,min=1,max=100"`
	}

	RefreshTokenRequest struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	LogoutRequest struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	ForgotPasswordRequest struct {
		Email string `json:"email" validate:"required,email"`
	}

	ResetPasswordRequest struct {
		Token       string `json:"token" validate:"required"`
		NewPassword string `json:"new_password" validate:"required,min=8"`
	}

	UserDTO struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	RegisterResponse struct {
		User UserDTO `json:"user"`
	}

	LoginResponse struct {
		User         UserDTO `json:"user"`
		Token        string  `json:"token"`
		RefreshToken string  `json:"refresh_token"`
	}

	Login2FAResponse struct {
		ChallengeID string `json:"challenge_id"`
	}

	RefreshTokenResponse struct {
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}
)

func (r *LoginRequest) Validate() error {
	v := pkgValidator.NewValidate()
	if err := v.ValidateStruct(r); err != nil {
		return err
	}

	return v.ValidatePassword(r.Password)
}

func (r *Login2FARequest) Validate() error {
	v := pkgValidator.NewValidate()
	if err := v.ValidateStruct(r); err != nil {
		return err
	}

	return v.ValidatePassword(r.Password)
}

func (r *VerifyLogin2FARequest) Validate() error {
	v := pkgValidator.NewValidate()
	return v.ValidateStruct(r)
}

func (r *RegisterRequest) Validate() error {
	v := pkgValidator.NewValidate()
	if err := v.ValidateStruct(r); err != nil {
		return err
	}

	return v.ValidatePassword(r.Password)
}

func (r *RefreshTokenRequest) Validate() error {
	v := pkgValidator.NewValidate()
	return v.ValidateStruct(r)
}

func (r *LogoutRequest) Validate() error {
	v := pkgValidator.NewValidate()
	return v.ValidateStruct(r)
}

func (r *ForgotPasswordRequest) Validate() error {
	v := pkgValidator.NewValidate()
	return v.ValidateStruct(r)
}

func (r *ResetPasswordRequest) Validate() error {
	v := pkgValidator.NewValidate()
	if err := v.ValidateStruct(r); err != nil {
		return err
	}

	return v.ValidatePassword(r.NewPassword)
}
