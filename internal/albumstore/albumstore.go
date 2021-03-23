package albumstore

import (
	"fmt"
	"sync"
)

type Album struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
	URL    string `json:"url"`
}

// AlbumStore is an in-memory database of Albums. AlbumStore methods are safe
// to call concurrently.
type AlbumStore struct {
	sync.Mutex

	albums map[int]Album
	nextID int
}

func New() *AlbumStore {
	as := &AlbumStore{}
	as.albums = make(map[int]Album)
	as.nextID = 0
	return as
}

// LoadSamples populates the album store with sample data. It erases the album
// store before performing the load.
func (as *AlbumStore) LoadSamples() {
	as.Lock()
	as.albums = make(map[int]Album)
	as.nextID = 0
	as.Unlock()

	as.CreateAlbum("The Queen Is Dead", "The Smiths", "https://music.apple.com/gb/album/the-queen-is-dead/800092985")
	as.CreateAlbum("Revolver", "The Beatles", "https://music.apple.com/gb/album/revolver/1441164670")
	as.CreateAlbum("Hunky Dory", "David Bowie", "https://music.apple.com/gb/album/hunky-dory-remastered/1039798000")
}

// CreateAlbum creates a new album in the store.
func (as *AlbumStore) CreateAlbum(title string, artist string, url string) int {
	as.Lock()
	defer as.Unlock()

	album := Album{
		ID:     as.nextID,
		Title:  title,
		Artist: artist,
		URL:    url,
	}

	as.albums[as.nextID] = album
	as.nextID++
	return album.ID
}

// GetAlbum retrieves an album from the store, by id. If no such id exists, an
// error is returned.
func (as *AlbumStore) GetAlbum(id int) (Album, error) {
	as.Lock()
	defer as.Unlock()

	a, ok := as.albums[id]
	if ok {
		return a, nil
	} else {
		return Album{}, fmt.Errorf("album with id=%d not found", id)
	}
}
