package spotify

import (
	"github.com/pkg/errors"
	"github.com/tomwhy/SpotifyMultiQueue/client"
)


type SpotifyAccessToken struct {
	Token string
	Expires int
	RefreshTokenString string
}

func (t *SpotifyAccessToken) GetToken() string {
	return t.Token
}

func (t *SpotifyAccessToken) RefreshToken(client client.Client) (client.AccessToken, error) {
	spotifyClient, success := client.(*SpotifyClient)

	if(!success) {
		return errors.New("Cannot refresh a spotify token with a non spotify client")
	}

	return spotifyClient.GetAccessToken(t.RefreshTokenString)
}