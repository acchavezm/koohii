// This example demonstrates how to authenticate with Spotify using the authorization code flow.
// In order to run this example yourself, you'll need to:
//
//  1. Register an application at: https://developer.spotify.com/my-applications/
//       - Use "http://localhost:8080/callback" as the redirect URI
//  2. Set the SPOTIFY_ID environment variable to the client ID you got in step 1.
//  3. Set the SPOTIFY_SECRET environment variable to the client secret from step 1.
package main

import (
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

// album represents data about a record album.
type album struct {
    ID     string  `json:"id"`
    Title  string  `json:"title"`
    Artist string  `json:"artist"`
    Price  float64 `json:"price"`
}

const redirectURI = "http://localhost:9001/callback"

var (
	auth   = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate), spotifyauth.WithClientID(viperEnvVariable("SPOTIFY_ID")), spotifyauth.WithClientSecret(viperEnvVariable("SPOTIFY_SECRET")))
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
		c.HTML(http.StatusOK, "user.html", gin.H{
			"user": user,
		})
	})
	r.GET("/climate", func(c *gin.Context) {

	}
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
