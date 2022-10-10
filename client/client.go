package client

import "net/url"

type Song struct {
	Name string
	ImageUrl string 
};

type Client interface {
	GetAuthorizationURL() (*url.URL, error)
	CompleteAuthorization(urlParams url.Values) error

	SearchSongs(search_phrase string) ([]Song, error)
	QueueSong(song Song) error
	GetPlayingSong() (Song, error)
};


func Must(c Client, err error) Client {
	if(err != nil) {
		panic(err)
	}

	return c
}