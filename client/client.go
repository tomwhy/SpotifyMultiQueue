package client

import "net/url"

type Song struct {
	Name string
	ImageUrl string 
};

type Client interface {
	GetAuthorizationURL() (*url.URL, error)
	CompleteAuthorization(urlParams map[string]string) error

	SearchSongs(search_phrase string) ([]Song, error)
	QueueSong(song Song) error
	GetPlayingSong() (Song, error)
};