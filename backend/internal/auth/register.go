package auth

import "context"

type RegisterRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type RegisterHandler struct {
	auth AuthService
}

func NewRegisterHandler(auth AuthService) *RegisterHandler {
	return &RegisterHandler{auth: auth}
}

func (h *RegisterHandler) Method() string  { return "POST" }
func (h *RegisterHandler) Pattern() string { return "/auth/register" }

func (h *RegisterHandler) Validate(req RegisterRequest) error {
	if req.FirstName == "" {
		return errValidation("first_name is required")
	}
	if req.LastName == "" {
		return errValidation("last_name is required")
	}
	if req.Email == "" {
		return errValidation("email is required")
	}
	if req.Password == "" {
		return errValidation("password is required")
	}
	return nil
}

func (h *RegisterHandler) Serve(ctx context.Context, req RegisterRequest) (LoginResponse, error) {
	access, refresh, err := h.auth.Register(ctx, req.FirstName, req.LastName, req.Email, req.Password)
	if err != nil {
		return LoginResponse{}, err
	}
	return LoginResponse{Token: access, RefreshToken: refresh}, nil
}
