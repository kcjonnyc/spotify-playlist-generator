# Spotify Playlist Generator API

This API serves as a wrapper around Spotify's web API. It provides a simplified search endpoint as well as a playlists endpoint that will create a playlist for a user and add recommended tracks to it.

## Getting Started

Clone the repo in your GOPATH. We use the echo framework so you may need to get dependencies. Run `go get -u github.com/labstack/echo/...`.

To start up up the service, run `go run main.go`. The service should start up on port 8080 (accessible from localhost:8080)
```
   ____    __
  / __/___/ /  ___
 / _// __/ _ \/ _ \
/___/\__/_//_/\___/ v3.3.dev
High performance, minimalist Go web framework
https://echo.labstack.com
____________________________________O/_______
                                    O\
⇨ http server started on [::]:8080
```

## Endpoints

All endpoints require an `Authorization` header with `Bearer <access token>`, where the access token is an access token that belongs to the user. This means that you must use the Authorization Code Grant or Implicit Grant to get a token on behalf of the user. We also require the `playlist-modify-public` scope to create and modify playlists.

Notes on access tokens: The API will return a 500 error if the API fails at any point due to authorization. The logs will contain more information however this should be changed in the future to provide better responses.

### GET /tracks

The request requires query parameters q, type and limit. These parameters match up with the query parameters required for Spotify's search endpoint. q is the query, type a list of types you want to search for (album, artist, playlist, track), limit is the maximum number of results to return.

Sample: /tracks?q=the%201975&type=track&limit=3
```
{
    "tracks": [
        {
            "name": "Chocolate",
            "imageUrl": "https://i.scdn.co/image/281926da293f453a9f486c81c059c20278c87195",
            "trackId": "44Ljlpy44mHvLJxcYUvTK0",
            "artistId": "3mIj9lX2MWuHmhNCA7LSCW"
        },
        {
            "name": "The Sound",
            "imageUrl": "https://i.scdn.co/image/a2aba574af86865eee08624042b1bf75d15cc362",
            "trackId": "316r1KLN0bcmpr7TZcMCXT",
            "artistId": "3mIj9lX2MWuHmhNCA7LSCW"
        },
        {
            "name": "Somebody Else",
            "imageUrl": "https://i.scdn.co/image/a2aba574af86865eee08624042b1bf75d15cc362",
            "trackId": "5hc71nKsUgtwQ3z52KEKQk",
            "artistId": "3mIj9lX2MWuHmhNCA7LSCW"
        }
    ]
}
```

### POST /playlists

The request requires a body
```
{
	"user": "kcjonnyc",
	"name": "generated playlist",
	"description": "auto-generated playlist",
	"tracks": "5hc71nKsUgtwQ3z52KEKQk",
	"artists": "3mIj9lX2MWuHmhNCA7LSCW",
	"limit": 50
}
```
Where, tracks and artists are comma separated lists of Spotify ID's. For more on Spotify ID's visit: https://beta.developer.spotify.com/documentation/web-api/#spotify-uris-and-ids

If the operation is successful, the response will the a 200 with message "Successfully created playlist". 
