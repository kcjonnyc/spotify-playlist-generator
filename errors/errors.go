package errors

import (
    "errors"
)

var (
    ErrNoAuthorizationHeader = errors.New("No authorization header provided")
    ErrSpotifyBadStatus = errors.New("Bad status code from Spotify API")
    ErrSpotifyUnauthorized = errors.New("Unauthorized, bad access token")
)
