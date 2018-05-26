package server

import (
    "bytes"
    "encoding/json"
    "io/ioutil"
    "net/http"

    "github.com/labstack/echo"
)

type GenerationInfo struct {
    User        string   `json:"user"`
    Token       string   `json:"token"`
    Name        string   `json:"name"`
    Description string   `json:"description"`
    Songs       []string `json:"songs"`
    Artists     []string `json:"artists"`
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

/*
type RecommendationsRequest struct {

}

type AddSongsRequest struct {

}*/

func (s *Server) generatePlaylist(c echo.Context) (err error) {
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

    // Generate playlist creation request
    createPlaylistRequest := new(CreatePlaylistRequest)
    createPlaylistRequest.Name = generationInfo.Name;
    createPlaylistRequest.Description = generationInfo.Description;

    createPlaylistResponse, err := s.createPlaylist(createPlaylistRequest, generationInfo, client)
    if err != nil {
        c.JSON(http.StatusBadRequest, Error{"Could not create playlist"})
    }
    return c.JSON(http.StatusOK, createPlaylistResponse)
}

func (s *Server) createPlaylist(createPlaylistRequest *CreatePlaylistRequest,
    generationInfo *GenerationInfo, client *http.Client) (createPlaylistResponse *CreatePlaylistResponse, err error) {
    // POST to Spotify's playlists endpoint to create playlist
    createPlaylistUrl := "https://api.spotify.com/v1/users/" + generationInfo.User + "/playlists"
    b, err := json.Marshal(createPlaylistRequest)
    if err != nil {
        s.e.Logger.Errorf("Could not marshal create playlist request: ", err)
        return
    }
    req, _ := http.NewRequest("POST", createPlaylistUrl, bytes.NewBuffer(b))
    req.Header.Set("Authorization", "Bearer " + generationInfo.Token)
    req.Header.Set("Content-Type", "application/json")
    res, err := client.Do(req)
    if err != nil {
        s.e.Logger.Errorf("Could not POST to playlists endpoint: ", err)
        return
    }

    // Check response from playlist creation
    if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
        s.e.Logger.Error("Bad status code from Spotify API")
		return
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
