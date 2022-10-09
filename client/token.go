package client

type AccessToken interface {
	GetToken() string
	RefreshToken(client Client) (AccessToken, error)
}