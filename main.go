package main

import (
	"net/http"

	"github.com/tomwhy/SpotifyMultiQueue/client"
	"github.com/tomwhy/SpotifyMultiQueue/client/spotify"
	"github.com/tomwhy/SpotifyMultiQueue/utils/secrets"
	"github.com/tomwhy/SpotifyMultiQueue/web"
)


func main() {
	server := web.NewWebServer(client.Must(spotify.NewSpotifyClient(secrets.NewLiteralSecret(""), secrets.NewLiteralSecret(""), "http://localhost:8080/auth")))

	err := server.Serve("", 8080)
	if(err != http.ErrServerClosed) {
		panic(err)
	}
}