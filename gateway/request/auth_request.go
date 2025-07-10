package request

type SignUpRequest struct {
	Username string `json:"username" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type VerifySignUpRequest struct {
	RegistrationToken string `json:"registration_token" binding:"required,uuid4"`
	Otp               string `json:"otp" binding:"required,len=6,numeric"`
}

type SignInRequest struct {
	Username string `json:"username" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=6"`
}
