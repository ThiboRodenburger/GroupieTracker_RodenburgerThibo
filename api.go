package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

var videoID string

func main() {
	css := http.FileServer(http.Dir("css"))
	http.Handle("/css/", http.StripPrefix("/css/", css))
	http.HandleFunc("/", homePage)
	http.HandleFunc("/video", playVideo)
	http.HandleFunc("/search", searchVideos)
	http.HandleFunc("/select", selectVideo)
	http.HandleFunc("/favicon.ico", faviconHandler)
	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", nil)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("C:/Users/roden/OneDrive/Documents/coursb1/Groupie_Tracker/projet api/templates/Index.html"))
	tmpl.Execute(w, nil)
}

func playVideo(w http.ResponseWriter, r *http.Request) {
	videoURL := "https://www.youtube.com/embed/" + videoID
	http.Redirect(w, r, videoURL, http.StatusSeeOther)
}

func searchVideos(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("q")
	if query == "" {
		http.Error(w, "Please specify a search query", http.StatusBadRequest)
		return
	}

	youtubeService, err := youtube.New(&http.Client{
		Transport: &transport.APIKey{Key: "AIzaSyDCRxQxNNwJRt4JhHjd4EyxNaKHoKZIpsY"},
	})
	if err != nil {
		log.Fatalf("Failed to create YouTube service: %v", err)
	}

	searchCall := youtubeService.Search.List([]string{"id", "snippet"}).Q(query).MaxResults(10)

	searchResponse, err := searchCall.Do()
	if err != nil {
		log.Fatalf("Failed to search for videos: %v", err)
	}

	videoList := make(map[string]string)

	for _, item := range searchResponse.Items {
		switch item.Id.Kind {
		case "youtube#video":
			videoList[item.Id.VideoId] = item.Snippet.Title
		}
	}

	tmpl := template.Must(template.ParseFiles("search.html"))
	tmpl.Execute(w, videoList)
}

func selectVideo(w http.ResponseWriter, r *http.Request) {
	videoID = r.FormValue("videoID")
	http.Redirect(w, r, "/video", http.StatusSeeOther)
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "favicon.ico")
}
