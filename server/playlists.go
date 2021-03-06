package server

import (
    "fmt"
    "bytes"
    "encoding/json"
    "io/ioutil"
    "net/http"
    "strings"

    "spotify-playlist-generator/errors"

    "github.com/labstack/echo"
)

type GenerationInfo struct {
    User         string   `json:"user"`
    Name         string   `json:"name"`
    Description  string   `json:"description"`
    Tracks       string   `json:"tracks"`
    Artists      string   `json:"artists"`
    Limit        *int     `json:"limit,omitempty"`
    Danceability *float64 `json:"danceability,omitempty"`
    Energy       *float64 `json:"energy,omitempty"`
    Liveness     *float64 `json:"liveness,omitempty"`
    Loudness     *float64 `json:"loudness,omitempty"`
    Mode         *int     `json:"mode,omitempty"`
    Popularity   *int     `json:"popularity,omitempty"`
    Valence      *float64 `json:"valence,omitempty"`
}

type CreatePlaylistRequest struct {
    Name        string `json:"name"`
    Description string `json:"description"`
}

type CreatePlaylistResponse struct {
    // The response from Spotify has a lot more information
    // but all we need is the id
    PlaylistId  string `json:"id"`
}

type RecommendationsResponse struct {
    Tracks []struct {
        Name string `json:"name"`
        Uri  string `json:"uri"`
    } `json:"tracks"`
}


type AddTracksRequest struct {
    Uris []string `json:"uris"`
}

func (s *Server) generatePlaylist(c echo.Context) (err error) {
    // We need to be passed the user's access token through the Authorization header
    authorization := c.Request().Header.Get("Authorization")
    if authorization == "" {
        s.e.Logger.Error("No authorization provided, could not generate playlist")
        return c.JSON(http.StatusBadRequest, Error{errors.ErrNoAuthorizationHeader.Error()})
    }

    // Get details from POST request
    // NOTE: The user must have logged in and given our app the
    // permissions it needs
    // We are currently passing access token through body...
    generationInfo := new(GenerationInfo)
    if err = c.Bind(generationInfo); err != nil {
        return
    }

    // Create http client
    client := &http.Client{}

    // Get track recommendations
    recommendationsResponse, err := s.getRecommendations(authorization, generationInfo, client)
    if err != nil {
        if err == errors.ErrSpotifyUnauthorized {
            return c.JSON(http.StatusUnauthorized, Error{err.Error()})
        }
        return c.JSON(http.StatusBadRequest, Error{"Could not get recommendations"})
    }

    // Generate playlist creation request
    createPlaylistRequest := new(CreatePlaylistRequest)
    createPlaylistRequest.Name = generationInfo.Name;
    createPlaylistRequest.Description = generationInfo.Description;

    createPlaylistResponse, err := s.createPlaylist(authorization, createPlaylistRequest, generationInfo, client)
    if err != nil {
        if err == errors.ErrSpotifyUnauthorized {
            return c.JSON(http.StatusUnauthorized, Error{err.Error()})
        }
        return c.JSON(http.StatusBadRequest, Error{"Could not create playlist"})
    }

    // Generate add tracks to playlist request
    addTracksRequest := new(AddTracksRequest)
    givenTracks := strings.Split(generationInfo.Tracks, ",")
    for _, track := range givenTracks {
        addTracksRequest.Uris = append(addTracksRequest.Uris, "spotify:track:" + track)
    }
    for _, track := range recommendationsResponse.Tracks {
        addTracksRequest.Uris = append(addTracksRequest.Uris, track.Uri)
    }

    if err = s.addTracksToPlaylist(authorization, createPlaylistResponse.PlaylistId, addTracksRequest, generationInfo, client); err != nil {
        if err == errors.ErrSpotifyUnauthorized {
            return c.JSON(http.StatusUnauthorized, Error{err.Error()})
        }
        return c.JSON(http.StatusBadRequest, Error{"Could not add tracks to playlist"})
    }

    return c.JSON(http.StatusOK, Message{"Successfully created playlist"})
}

func (s *Server) createPlaylist(authorization string, createPlaylistRequest *CreatePlaylistRequest,
    generationInfo *GenerationInfo, client *http.Client) (createPlaylistResponse *CreatePlaylistResponse, err error) {
    // POST to Spotify's playlists endpoint to create playlist
    s.e.Logger.Debug("Creating Spotify playlist")

    // Generate request url
    createPlaylistUrl := "https://api.spotify.com/v1/users/" + generationInfo.User + "/playlists"
    b, err := json.Marshal(createPlaylistRequest)
    if err != nil {
        s.e.Logger.Errorf("Could not marshal create playlist request: ", err)
        return
    }
    req, _ := http.NewRequest("POST", createPlaylistUrl, bytes.NewBuffer(b))
    req.Header.Set("Authorization", authorization)
    req.Header.Set("Content-Type", "application/json")
    s.e.Logger.Debug(req.URL.String())
    res, err := client.Do(req)
    if err != nil {
        s.e.Logger.Errorf("Could not POST to playlists endpoint: ", err)
        return
    }

    // Check response from playlist creation
    if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
        s.e.Logger.Error("Could not create playlist, status code: " + http.StatusText(res.StatusCode))
        if res.StatusCode == http.StatusUnauthorized {
            return nil, errors.ErrSpotifyUnauthorized
        }
        return nil, errors.ErrSpotifyBadStatus
	}
    // Unmarshal response and return
    body, err := ioutil.ReadAll(res.Body)
	if err != nil {
        s.e.Logger.Errorf("Could not read create playlist response body: ", err)
		return
	}
    createPlaylistResponse = new(CreatePlaylistResponse)
    if err = json.Unmarshal(body, &createPlaylistResponse); err != nil {
        s.e.Logger.Errorf("Could not unmarshal create playlist response: ", err)
        return
    }
    return
}

func (s *Server) getRecommendations(authorization string, generationInfo *GenerationInfo,
    client *http.Client) (recommendationsResponse *RecommendationsResponse, err error) {
    // GET recommendations from Spotify recommendations endpoint
    s.e.Logger.Debug("Getting Spotify track recommendations")

    // Generate request url
    recommendationsUrl := "https://api.spotify.com/v1/recommendations"
    req, _ := http.NewRequest("GET", recommendationsUrl, nil)
    req.Header.Set("Authorization", authorization)
    query := req.URL.Query()
    query.Add("seed_tracks", generationInfo.Tracks)
    numberTracks := len(generationInfo.Tracks)
    numberArtists := len(generationInfo.Artists)
    maxArtists := 5 - numberTracks
    if generationInfo.Artists != "" && numberTracks < 5 {
        if (numberArtists < maxArtists) {
            query.Add("seed_artists", generationInfo.Artists[0:numberArtists])
        } else {
            query.Add("seed_artists", generationInfo.Artists[0:maxArtists])
        }
    }
    // Add optional query parameters
    if generationInfo.Limit != nil {
        query.Add("limit", fmt.Sprint(*generationInfo.Limit))
    }
    if generationInfo.Danceability != nil {
        query.Add("target_danceability", fmt.Sprint(*generationInfo.Danceability))
    }
    if generationInfo.Energy != nil {
        query.Add("target_energy", fmt.Sprint(*generationInfo.Energy))
    }
    if generationInfo.Liveness != nil {
        query.Add("target_liveness", fmt.Sprint(*generationInfo.Liveness))
    }
    if generationInfo.Loudness != nil {
        query.Add("target_loudness", fmt.Sprint(*generationInfo.Loudness))
    }
    if generationInfo.Mode != nil {
        query.Add("target_mode", fmt.Sprint(*generationInfo.Mode))
    }
    if generationInfo.Popularity != nil {
        query.Add("target_popularity", fmt.Sprint(*generationInfo.Popularity))
    }
    if generationInfo.Valence != nil {
        query.Add("target_valence", fmt.Sprint(*generationInfo.Valence))
    }

    req.URL.RawQuery = query.Encode()
    s.e.Logger.Debug(req.URL.String())
    res, err := client.Do(req)
    if err != nil {
        s.e.Logger.Errorf("Could not GET from recommendations endpoint: ", err)
        return
    }

    // Check recommendations response
    if res.StatusCode != http.StatusOK {
        s.e.Logger.Error("Could not get recommendations, status code: " + http.StatusText(res.StatusCode))
        if res.StatusCode == http.StatusUnauthorized {
            return nil, errors.ErrSpotifyUnauthorized
        }
        return nil, errors.ErrSpotifyBadStatus
	}
    // Unmarshal response and return
    body, err := ioutil.ReadAll(res.Body)
	if err != nil {
        s.e.Logger.Errorf("Could not read recommendations response body: ", err)
		return
	}
    recommendationsResponse = new(RecommendationsResponse)
    if err = json.Unmarshal(body, &recommendationsResponse); err != nil {
        s.e.Logger.Errorf("Could not unmarshal recommendations response: ", err)
        return
    }
    return
}

func (s *Server) addTracksToPlaylist(authorization string, playlistId string, addTracksRequest *AddTracksRequest,
    generationInfo *GenerationInfo, client *http.Client) (err error) {
    // POST to Spotify's playlists endpoint (for the specific playlist) to add tracks
    s.e.Logger.Debug("Adding tracks to Spotify playlist")

    // Generate request url
    addTracksUrl := "https://api.spotify.com/v1/users/" + generationInfo.User + "/playlists/" + playlistId + "/tracks"
    b, err := json.Marshal(addTracksRequest)
    if err != nil {
        s.e.Logger.Errorf("Could not marshal add tracks request: ", err)
        return
    }
    req, _ := http.NewRequest("POST", addTracksUrl, bytes.NewBuffer(b))
    req.Header.Set("Authorization", authorization)
    req.Header.Set("Content-Type", "application/json")
    s.e.Logger.Debug(req.URL.String())
    res, err := client.Do(req)
    if err != nil {
        s.e.Logger.Errorf("Could not POST to playlists endpoint to add tracks: ", err)
        return
    }

    // Check response from adding tracks
    if res.StatusCode != http.StatusCreated {
        s.e.Logger.Error("Could not add tracks, status code: " + http.StatusText(res.StatusCode))
        if res.StatusCode == http.StatusUnauthorized {
            return errors.ErrSpotifyUnauthorized
        }
        return errors.ErrSpotifyBadStatus
	}
    // We don't care about the response body in this case
    return
}
