package server

import (
    "encoding/json"
    "errors"
    "io/ioutil"
    "net/http"

    "github.com/labstack/echo"
)

type TracksResponse struct {
    Tracks   []Track `json:"tracks"`
}

type Track struct {
    // NOTE: For both image url and artist id, we will return
    // the element at index 0 for simplicity
    Name     string `json:"name"`
    ImageUrl string `json:"imageUrl"`
    TrackId  string `json:"trackId"`
    ArtistId string `json:"artistId"`
}

type SearchResponse struct {
    Tracks struct {
        Items []struct {
            Album struct {
                Artists []struct {
                    Id string `json:"id"`
                } `json:"artists"`
                Images []struct {
                    Url string `json:"url"`
                }
            } `json:"album"`
            Id   string `json:"id"`
            Name string `json:"name"`
        } `json:"items"`
    } `json:"tracks"`
}

func (s *Server) searchTracks(c echo.Context) (err error) {
    // We need to be passed the user's access token through the Authorization header
    authorization := c.Request().Header.Get("Authorization")
    if authorization == "" {
        s.e.Logger.Error("No authorization provided, could not search tracks")
        err = errors.New("No authorization provided")
        return
    }

    // Create http client
    client := &http.Client{}

    // Generate request url
    searchUrl := "https://api.spotify.com/v1/search"
    req, _ := http.NewRequest("GET", searchUrl, nil)
    req.Header.Set("Authorization", authorization)
    req.URL.RawQuery = c.QueryString()
    s.e.Logger.Debug(req.URL.String())
    res, err := client.Do(req)
    if err != nil {
        s.e.Logger.Errorf("Could not GET from search endpoint: ", err)
        return
    }

    // Check search response
    if res.StatusCode != http.StatusOK {
        s.e.Logger.Error("Could not get tracks, status code: " + http.StatusText(res.StatusCode))
        err = errors.New("Bad status code from Spotify API, could not get tracks")
		return
	}
    // Unmarshal response and return
    body, err := ioutil.ReadAll(res.Body)
	if err != nil {
        s.e.Logger.Errorf("Could not read search response body: ", err)
		return
	}
    searchResponse := new(SearchResponse)
    if err = json.Unmarshal(body, &searchResponse); err != nil {
        s.e.Logger.Errorf("Could not unmarshal search response: ", err)
        return
    }

    // Generate our nicely formatted response
    tracksResponse := new(TracksResponse)
    for _, item := range searchResponse.Tracks.Items {
        track := Track{}
        track.Name = item.Name
        track.ImageUrl = item.Album.Images[0].Url
        track.TrackId = item.Id
        track.ArtistId = item.Album.Artists[0].Id
        tracksResponse.Tracks = append(tracksResponse.Tracks, track)
    }

    return c.JSON(http.StatusOK, tracksResponse)
}
