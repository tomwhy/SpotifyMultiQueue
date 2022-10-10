package spotify


type authErr struct {
	Error string `json:"error"`
	Desc string `json:"error_description"`
}

type tokenResponse struct {
	Token string `json:"access_token"`
	Type string `json:"token_type"`
	Scope string `json:"scope"`
	Expires int `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type apiErr struct {
	Error struct {
		Status int `json:"status"`
		Message string `json:"message"`
	} `json:"error"`
}

type currentSongRes struct {
	Item struct {
		Name string `json:"name"`
	} `json:"item"`
}

type searchRes struct {
	Tracks struct {
		Items []struct {
			Id string `json:"id"`
			Name string `json:"name"`
			Duration int `json:"duration_ms"`
			Artists []struct {
				Name string `json:"name"`
			} `json:"artists"`
		} `json:"items"`
	} `json:"tracks"`
}