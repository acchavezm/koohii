// This example demonstrates how to authenticate with Spotify using the authorization code flow.
// In order to run this example yourself, you'll need to:
//
//  1. Register an application at: https://developer.spotify.com/my-applications/
//       - Use "http://localhost:8080/callback" as the redirect URI
//  2. Set the SPOTIFY_ID environment variable to the client ID you got in step 1.
//  3. Set the SPOTIFY_SECRET environment variable to the client secret from step 1.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"context"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

// class that represents data about a the weather on a city.
/* Sample JSON
{
"coord": {
	"lon": -79.9,
	"lat": -2.1667
},
"weather": [
	{
	"id": 804,
	"main": "Clouds",
	"description": "overcast clouds",
	"icon": "04n"
	}
],
"base": "stations",
"main": {
	"temp": 296.31,
	"feels_like": 296.9,
	"temp_min": 295.9,
	"temp_max": 296.31,
	"pressure": 1012,
	"humidity": 85
},
"visibility": 10000,
"wind": {
	"speed": 2.68,
	"deg": 285,
	"gust": 8.05
},
"clouds": {
	"all": 90
},
"dt": 1629516875,
"sys": {
	"type": 2,
	"id": 2008064,
	"country": "EC",
	"sunrise": 1629458495,
	"sunset": 1629501876
},
"timezone": -18000,
"id": 3657509,
"name": "Guayaquil",
"cod": 200
}
*/
type CityWeather struct {
	Coord      Coordinate `json:"coord"`
	Weather    []Weather  `json:"weather"`
	Base       string     `json:"base"`
	Main       Main       `json:"main"`
	Visibility int64      `json:"visibility"`
	Wind       Wind       `json:"wind"`
	Clouds     Clouds     `json:"clouds"`
	Dt         int64      `json:"dt"`
	Sys        Sys        `json:"sys"`
	Timezone   int64      `json:"timezone"`
	Id         int64      `json:"id"`
	Name       string     `json:"name"`
	Cod        int64      `json:"cod"`
}

type Coordinate struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

type Weather struct {
	Id          int64  `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type Main struct {
	Temp       float64 `json:"temp"`
	Feels_like float64 `json:"feels_like"`
	Temp_min   float64 `json:"temp_min"`
	Temp_max   float64 `json:"temp_max"`
	Pressure   int64   `json:"pressure"`
	Humidity   int64   `json:"humidity"`
}

type Wind struct {
	Speed float64 `json:"speed"`
	Deg   int64   `json:"deg"`
	Gust  float64 `json:"gust"`
}

type Clouds struct {
	All int64 `json:"all"`
}

type Sys struct {
	Id       int64  `json:"id"`
	Sys_type int64  `json:"type"`
	Country  string `json:"country"`
	Sunrise  int64  `json:"sunrise"`
	Sunset   int64  `json:"sunset"`
}

const redirectURI = "http://localhost:9001/callback"

var (
	auth = spotifyauth.New(
		spotifyauth.WithRedirectURL(redirectURI),
		spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate),
		spotifyauth.WithScopes(spotifyauth.ScopeUserTopRead),
		spotifyauth.WithClientID(viperEnvVariable("SPOTIFY_ID")),
		spotifyauth.WithClientSecret(viperEnvVariable("SPOTIFY_SECRET")))
	client *spotify.Client
	state  = "abc123"
)

// use viper package to read .env file
// return the value of the key
func viperEnvVariable(key string) string {

	// SetConfigFile explicitly defines the path, name and extension of the config file.
	// Viper will use this and not check any of the config paths.
	// .env - It will search for the .env file in the current directory
	viper.SetConfigFile(".env")

	// Find and read the config file
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}

	// viper.Get() returns an empty interface{}
	// to get the underlying type of the key,
	// we have to do the type assertion, we know the underlying value is string
	// if we type assert to other type it will throw an error
	value, ok := viper.Get(key).(string)

	// If the type is a string then ok will be true
	// ok will make sure the program not break
	if !ok {
		log.Fatalf("Invalid type assertion")
	}

	return value
}

func main() {

	os.Setenv("SPOTIFY_ID", viperEnvVariable("SPOTIFY_ID"))
	os.Setenv("SPOTIFY_SECRET", viperEnvVariable("SPOTIFY_SECRET"))

	//redirectURI := viperEnvVariable("SPOTIFY_REDIRECT_URI")

	//auth := spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate))

	/*

		fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

		// wait for auth to complete
		client := <-ch

		// use the client to make calls that require authorization
		user, err := client.CurrentUser()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("You are logged in as:", user.ID)

		// search for playlists and albums containing "holiday"
		results, err := client.Search("holiday", spotify.SearchTypePlaylist|spotify.SearchTypeAlbum)
		if err != nil {
			log.Fatal(err)
		}

		// handle album results
		if results.Albums != nil {
			fmt.Println("Albums:")
			for _, item := range results.Albums.Albums {
				fmt.Println("   ", item.Name)
			}
		}
		// handle playlist results
		if results.Playlists != nil {
			fmt.Println("Playlists:")
			for _, item := range results.Playlists.Playlists {
				fmt.Println("   ", item.Name)
			}
		}
	*/

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	r.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Main website",
		})
	})
	r.GET("/callback", func(c *gin.Context) {
		completeAuth(c.Writer, c.Request)
		c.Redirect(http.StatusMovedPermanently, "/user")
	})
	r.GET("/login", func(c *gin.Context) {
		fmt.Println("Got login request")
		//construir la URL
		url := auth.AuthURL(state)
		fmt.Println(url)
		c.Redirect(http.StatusMovedPermanently, url)
	})
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/user", func(c *gin.Context) {
		// use the client to make calls that require authorization
		user, err := client.CurrentUser(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("You are logged in as:", user.ID)
		fmt.Println("User ID:", user.ID)
		fmt.Println("Display name:", user.DisplayName)
		fmt.Println("Spotify URI:", string(user.URI))
		fmt.Println("Endpoint:", user.Endpoint)
		fmt.Println("Followers:", user.Followers.Count)

		resp, err := http.Get("https://api.openweathermap.org/data/2.5/weather?q=Guayaquil&appid=1cdbcd14a6e201f2b5d091e4b1c53b8a&units=metric&lang=es")
		if err != nil {
			log.Fatalln(err)
		}
		decoder := json.NewDecoder(resp.Body)
		var data CityWeather
		err = decoder.Decode(&data)
		if err != nil {
			log.Fatalln(err)
		}
		raw_tracks, err := client.GetPlaylistTracks(context.Background(), spotify.ID("37i9dQZEVXbMDoHDwVN2tF"))
		if err != nil {
			log.Fatalln(err)
		}

		var track_list []spotify.FullTrack
		for _, element := range raw_tracks.Tracks {
			// index is the index where we are
			// element is the element from someSlice for where we are
			track := element.Track
			track_id := track.ID
			audio_features, err := client.GetAudioFeatures(context.Background(), track_id)
			if err != nil {
				log.Fatalln(err)
			}
			if audio_features[len(audio_features)-1].Energy <= 0.5 {
				track_list = append(track_list, track)
			}
		}

		c.HTML(http.StatusOK, "user.html", gin.H{
			"user":         user,
			"city_climate": data,
			"track_list":   track_list,
		})
	})
	r.GET("/climate", func(c *gin.Context) {
		resp, err := http.Get("https://api.openweathermap.org/data/2.5/weather?q=Guayaquil&appid=1cdbcd14a6e201f2b5d091e4b1c53b8a&units=metric&lang=es")
		if err != nil {
			log.Fatalln(err)
		}
		decoder := json.NewDecoder(resp.Body)
		var data CityWeather
		err = decoder.Decode(&data)
		if err != nil {
			log.Fatalln(err)
		}
		c.IndentedJSON(http.StatusOK, data)
	})
	r.Run(":9001") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	log.Println("Callback handler for state:", state)
	fmt.Println(r.FormValue("state"))

	tok, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	// use the token to get an authenticated client
	client = spotify.New(auth.Client(r.Context(), tok))
	fmt.Println("Client a!")
}
