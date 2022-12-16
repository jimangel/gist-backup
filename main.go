package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/google/go-github/v29/github"
	"golang.org/x/oauth2"
)

func main() {
	// Get the GitHub access token from the environment variable
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("Error: GITHUB_TOKEN environment variable is not set")
	}

	// Create an OAuth2 client using the token
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	// Create a new GitHub client
	client := github.NewClient(httpClient)

	// Set the GistListOptions to get 100 gists per page
	// and start at page 1
	opts := &github.GistListOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
			Page:    1,
		},
	}

	// Loop through the pages of gists
	for i := 0; i < 5; i++ {

		// Get the authenticated user's gists
		gists, _, err := client.Gists.List(context.Background(), "", opts)
		if err != nil {
			log.Fatal(err)
		}

		// Clone each gist to the local backup destination
		for _, gist := range gists {
			fmt.Printf("Backing up gist %s...\n", *gist.ID)
			cloneGist(client, gist)
		}

		// Increment the page number
		opts.Page++
	}

}

func cloneGist(client *github.Client, gist *github.Gist) {
	// Get the URL of the git repository for the gist
	gitURL := *gist.GitPullURL

	cmd := exec.Command("git", "clone", gitURL)
	err := cmd.Run()
	if err != nil {
		log.Printf("Failed to clone gist %s: %v", gist, err)
	}
}
