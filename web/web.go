package web

import (
	"fmt"
	"net/http"

	echo "github.com/labstack/echo/v4"
	"github.com/tomwhy/SpotifyMultiQueue/client"
)

type WebServer struct {
	server *echo.Echo
	client client.Client
}


func NewWebServer(client client.Client) *WebServer {
	server := &WebServer {
		server: echo.New(),
		client: client,
	}

	server.registerRoutes()

	return server
}

func (s *WebServer) registerRoutes() {
	s.server.GET("/", s.homePage)
	s.server.GET("/auth", s.clientAuthCallback)
	s.server.GET("/share", s.adminSharePage)
}

func (s *WebServer) homePage(c echo.Context) error {
	auth_url, err := s.client.GetAuthorizationURL()
	if(err != nil) {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.Redirect(http.StatusTemporaryRedirect, auth_url.String())
}

func (s *WebServer) clientAuthCallback(c echo.Context) error {
	err := s.client.CompleteAuthorization(c.QueryParams())
	if(err != nil) {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.Redirect(http.StatusTemporaryRedirect, "/share")
}

func (s *WebServer) adminSharePage(c echo.Context) error {
	song, err := s.client.GetPlayingSong()
	if(err != nil) {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, song.Name)
}

func (s *WebServer) Serve(host string, port int16) error {
	return s.server.Start(fmt.Sprintf("%s:%d", host, port))
}

func (s *WebServer) ServeTLS(host string, port int16, cert string, key string) error {
	return s.server.StartTLS(fmt.Sprintf("%s:%d", host, port), cert, key)	
}