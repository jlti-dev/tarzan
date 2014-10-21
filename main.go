package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/davecgh/go-spew/spew"
	apiClient "github.com/fsouza/go-dockerclient"
	"github.com/gorilla/mux"
	libgit "github.com/libgit2/git2go"
	"github.com/unrolled/render"
)

type GithubPushEventPayload struct {
	Hook struct {
		Config struct {
			Secret string `json:"secret"`
		} `json:"config"`
	} `json:"hook"`
	Repository struct {
		FullName string `json:"full_name"`
		HtmlUrl  string `json:"html_url"`
	} `json:"repository"`
}

func main() {
	dockerHost := os.Getenv("DOCKER_HOST")
	if dockerHost == "" {
		dockerHost = "unix:///var/run/docker.sock"
	}

	_, err := apiClient.NewClient(dockerHost)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating client!!", err)
	}

	r := render.New(render.Options{})

	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		var (
			payload GithubPushEventPayload
		)
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&payload)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error decoding Github push payload:", err)
		}
		spew.Dump(payload)
		repoPath := fmt.Sprintf("./repos/%s", payload.Repository.FullName)
		repoUrl := payload.Repository.HtmlUrl
		_, err = libgit.Clone(repoUrl, repoPath, &libgit.CloneOptions{
			Bare: false,
		})
		r.JSON(w, http.StatusOK, "")
	}).Methods("POST")

	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(":3000")
}
