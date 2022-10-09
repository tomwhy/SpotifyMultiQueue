package spotify


type AuthErr struct {
	Error string `json:"error"`
	Desc string `json:"error_description"`
}

type TokenResponse struct {
	Token string `json:"access_token"`
	Type string `json:"token_type"`
	Scope string `json:"scope"`
	Expires int `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type ApiErr struct {
	Status int `json:"status"`
	Message string `json:"message"`
}