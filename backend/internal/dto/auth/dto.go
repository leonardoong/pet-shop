package auth

type CustomerRegisterRequest struct {
	FullName string `json:"full_name" binding:"required,min=2,max=100"`
	Email    string `json:"email"     binding:"required,email"`
	Password string `json:"password"  binding:"required,min=8"`
	Phone    string `json:"phone"     binding:"required,indonesian_phone"`
}

type CustomerLoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AdminLoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type CustomerProfile struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

type AdminProfile struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

type CustomerAuthResponse struct {
	Customer CustomerProfile `json:"customer"`
	Tokens   TokenPair       `json:"tokens"`
}

type AdminAuthResponse struct {
	Admin       AdminProfile `json:"admin"`
	Tokens      TokenPair    `json:"tokens"`
	Permissions []string     `json:"permissions"`
}
