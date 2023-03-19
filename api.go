package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// Define the data structure for a YouTube video
type Video struct {
	Title string
	Url   string
}

func main() {
	apiKey := os.Getenv("YOUTUBE_API_KEY")
	if apiKey == "AIzaSyDvH8v227tyToFWyDGnJwqj--Od5LMY2BM" {
		log.Fatal("YOUTUBE_API_KEY environment variable is required")
	}

	// Create a new gin engine
	r := gin.Default()

	// Define the route for the home page
	r.GET("/", func(c *gin.Context) {
		// Load the HTML template
		tmpl := template.Must(template.ParseFiles("Index.html"))

		// Display the HTML template
		if err := tmpl.Execute(c.Writer, nil); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	})

	// Define the route for the search API
	r.GET("/search", func(c *gin.Context) {
		// Get the search query from the query string
		query := c.Query("q")

		// Create a new YouTube service with the API key
		ctx := context.Background()
		youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
		if err != nil {
			log.Fatalf("Failed to create YouTube service: %v", err)
		}

		// Define the API request to search for videos
		searchCall := youtubeService.Search.List([]string{"id", "snippet"}).Q(query).MaxResults(10)

		// Execute the API request
		searchResponse, err := searchCall.Do()
		if err != nil {
			log.Fatalf("Failed to search for videos: %v", err)
		}

		// Extract the videos from the API response
		videos := make([]Video, 0, len(searchResponse.Items))
		for _, searchResult := range searchResponse.Items {
			if searchResult.Id.Kind == "youtube#video" {
				video := Video{
					Title: searchResult.Snippet.Title,
					Url:   fmt.Sprintf("https://www.youtube.com/watch?v=%v", searchResult.Id.VideoId),
				}
				videos = append(videos, video)
			}
		}

		// Render the video results in the search template
		tmpl := template.Must(template.ParseFiles("search.html"))
		if err := tmpl.Execute(c.Writer, videos); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	})

	// Run the gin engine on port 8080
	r.Run(":8080")
}
