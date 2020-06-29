# PlayMo Spotify Cloud Function

Spotify API implementation

## Getting Started

- Install dependencies `go mod tidy`

- Create a **.env** file and fill in values from **.env_example**

- Source your environment variables `source .env`

- To run functions locally `go run cmd/main.go`

## Routes

- **GET**

  - Endpoint: `/?link=spotify_playlist_link`

- **POST**

  - Endpoint: `/`

  - Body:

    ```json
    {
      "name": "string", // The playlist name
      "description": "string", // The playlist description
      // List of tracks to add to playlist. Maximum of 10 tracks.
      "tracks": [
        {
          "title": "string", // Track title
          "artist": "string", // Track artist
          "album": "string", // Track album
          "contributors": ["string"], // Track contributors. Must have atleast one value i.e. the track artist
          "duration": "integer" // The track duration in seconds
        }
      ]
    }
    ```
