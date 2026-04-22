package auth

import "context"

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginHandler struct {
	auth AuthService
}

func NewLoginHandler(auth AuthService) *LoginHandler {
	return &LoginHandler{auth: auth}
}

func (h *LoginHandler) Method() string  { return "POST" }
func (h *LoginHandler) Pattern() string { return "/auth/login" }

func (h *LoginHandler) Validate(req LoginRequest) error {
	if req.Email == "" {
		return errValidation("email is required")
	}
	if req.Password == "" {
		return errValidation("password is required")
	}
	return nil
}

func (h *LoginHandler) Serve(ctx context.Context, req LoginRequest) (LoginResponse, error) {
	access, refresh, err := h.auth.Login(ctx, req.Email, req.Password)
	if err != nil {
		return LoginResponse{}, err
	}
	return LoginResponse{Token: access, RefreshToken: refresh}, nil
}
