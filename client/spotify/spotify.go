package spotify

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/tomwhy/SpotifyMultiQueue/client"
	"github.com/tomwhy/SpotifyMultiQueue/utils"
	"github.com/tomwhy/SpotifyMultiQueue/utils/secrets"
)

const (
	API_BASE_URL = "https://api.spotify.com/v1"
	AUTH_BASE_URL = "https://accounts.spotify.com"
	APP_SCOPES = ""
)

type SpotifyClient struct {
	clientId secrets.Secret 
	clientSecret secrets.Secret
	redirectUri string
	state string

	accessToken client.AccessToken

};


func NewSpotifyClient(clientId secrets.Secret, clientSecret secrets.Secret, redirectUri string) (*SpotifyClient, error) {
	client := new(SpotifyClient);
	
	client.accessToken = nil
	client.clientId = clientId
	client.clientSecret = clientSecret
	client.redirectUri = redirectUri

	state_bytes, err := utils.RandBytes(16)
	if(err != nil) {
		return nil, errors.Wrap(err, "Failed generating state")
	}
	client.state = base64.StdEncoding.EncodeToString(state_bytes)


	return client, nil
}

func (c *SpotifyClient) GetAuthorizationURL() (*url.URL, error) {
	spotifyClientId, err := c.clientId.Read()
	if(err != nil) {
		return nil, errors.Warp(err, "Failed reading clientId")
	}


	return utils.BuildRequestURL(
		AUTH_BASE_URL, 
		"authorize", 
		map[string]string{
			"client_id": string(spotifyClientId),
			"redirect_uri": c.redirectUri,
			"state": c.state,
		},
	)
}


func (c *SpotifyClient) getClientAuthString() (string, error) {
	clientId, err := c.clientId.Read()
	if(err != nil) {
		return "", errors.Warp(err, "Failed reading clientId")
	}

	clientSecret, err := c.clientSecret.Read()
	if(err != nil) {
		return "", errors.Warp(err, "Failed reading clientSecret")
	}

	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", clientId, clientSecret))), nil
}

func (c *SpotifyClient) GetAccessToken(code string) (client.AccessToken, error) {
	auth, err := c.getClientAuthString()
	if(err != nil) {
		return errors.Wrap("failed getting client auth string")
	} 

	authRes, err := utils.SendPostApiRequest(
		AUTH_BASE_URL, 
		"api/token", 
		map[string]string{
			"grant_type": "authorization_code",
			"code": code,
			"redirect_uri": c.redirectUri,
		},
		map[string]string{
			"Authorization": fmt.Sprintf("Basic %s", auth),
			"Content-Type": "application/x-www-form-urlencoded",
		},
	)

	if(err != nil) {
		return errors.Wrap(err, "failed getting access token")
	}
	resDecoder := json.NewDecoder(authRes.Body)

	if(authRes.StatusCode != http.StatusOK) {
		var errRes AuthErr;
		resDecoder.Decode(&errRes)

		return errors.New(fmt.Sprintf("%s: %s", errRes.Error, errRes.Desc))
	}

	var tokenRes TokenResponse;
	resDecoder.Decode(&tokenRes)

	return &SpotifyAccessToken{
		Token: tokenRes.Token,
		Expires: tokenRes.Expires,
		RefreshTokenString: tokenRes.RefreshToken,
	}, nil
}

func (c *SpotifyClient) CompleteAuthorization(urlParams map[string]string) error {
	state, exists := urlParams["state"]
	if(!exists || state != c.state) {
		return errors.New("Unexpected state was passed to callback");
	}

	errorMsg, exists := urlParams["error"]
	if(exists) {
		return errors.New(errorMsg)
	}

	code, exists := urlParams["code"]
	if(!exists) {
		return errors.New("Missing authorization code")
	}

	token, err := c.GetAccessToken(code)
	if (err != nil) {
		return errors.Wrap(err, "failed getting access token")
	}

	c.accessToken = token
	return nil
}

func (c *SpotifyClient) getApiEndpoint(endpoint string, params map[string]string, expectedCode int, resp interface{}, errResp interface{}) (bool, error) {
	apiRes, err := utils.SendGetApiRequest(
		API_BASE_URL, 
		endpoint, 
		params,
		map[string]string {
			"Authorization": fmt.Sprintf("Bearer %s", c.accessToken.GetToken()),
			"Content-Type": "application/json",
		},
	)

	if(err != nil) {
		return false, errors.Wrap("failed getting api result")
	}

	if(apiRes.StatusCode != expectedCode) {
		json.NewDecoder(apiRes.Body).Decode(errResp)
	} else {
		json.NewDecoder(apiRes.Body).Decode(resp)
	}

	return apiRes.StatusCode == expectedCode, nil
}


func (c *SpotifyClient) SearchSongs(search_phrase string) ([]client.Song, error) {
	return nil, errors.New("Unimplemnted")
}

func (c *SpotifyClient) QueueSong(song client.Song) error {
	return errors.New("Unimplemnted")

}


func (c *SpotifyClient) GetPlayingSong() (client.Song, error) {

	var currentSong map[string]interface{}
	var errMsg ApiErr

	success, err := c.getApiEndpoint("/me/player/currently-playing", nil, http.StatusOK, &currentSong, &errMsg)
	if(err != nil) {
		return client.Song{}, errors.Wrap("failed getting current song")
	}

	if(!success) {
		return client.Song{}, errors.New(fmt.Sprintf("%d: %s", errMsg.Status, errMsg.Message))
	}

	// TEMP
	return client.Song{
		Name: currentSong["item"].(map[string]interface{})["name"].(string),
	}, nil
}