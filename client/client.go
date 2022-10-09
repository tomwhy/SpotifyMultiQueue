package client

type Song struct {
	Name string
	Image string 
};

type Client interface {
	SearchSongs(search_phrase string) ([]Song, error)
	QueueSong(song Song) error
	GetPlayingSong() (Song, error)
};