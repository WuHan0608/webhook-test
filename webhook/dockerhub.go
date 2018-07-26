package webhook

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"gopkg.in/birkirb/loggers.v1/log"
)

type pushData struct {
	Images   []string `json:"images"`
	PushedAt int64    `json:"pushed_at"`
	Pusher   string   `json:"pusher"`
	Tag      string   `json:"tag"`
}

type repository struct {
	CommentCount    int64  `json:"comment_count"`
	DataCreated     int64  `json:"date_created"`
	Description     string `json:"description"`
	Dockefile       string `json:"dockerfile"`
	FullDescription string `json:"full_description"`
	IsOfficial      bool   `json:"is_official"`
	IsPrivate       bool   `json:"is_private"`
	IsTrusted       bool   `json:"is_trusted"`
	Name            string `json:"name"`
	NameSpace       string `json:"namespace"`
	Owner           string `json:"owner"`
	RepoName        string `json:"repo_name"`
	RepoURL         string `json:"repo_url"`
	StarCount       int64  `json:"star_count"`
	Status          string `json:"status"`
}

type payload struct {
	CallbackURL string      `json:"callback_url"`
	PushData    *pushData   `json:"push_data"`
	Repository  *repository `json:"repository"`
}

// DockerHubHandler handlers webhook push data from docker hub
func DockerHubHandler() http.Handler {
	dockerHubMux := http.NewServeMux()
	dockerHubMux.HandleFunc("/dockerhub", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			fmt.Fprintln(w, "OK")
		} else if r.Method == http.MethodPost {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Errorf("read request error: %v", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			defer r.Body.Close()
			data := payload{}
			if err := json.Unmarshal(body, &data); err != nil {
				log.Errorf("decode payload error: %v", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			log.Printf("%#v\n", data)
		} else {
			http.Error(w, "Allowed methods: GET, POST", http.StatusMethodNotAllowed)
			return
		}
	})
	return dockerHubMux
}
