package main

import (
	"encoding/json"
	"net/http"
	"bytes"

	"github.com/aquiffoo/goblet"
)

type ChatReq struct {
	Model		string		`json:"model"`
	Messages	[]Message	`json:"messages"`
	Stream		bool		`json:"stream"`
}

type Message struct {
	Role	string	`json:"role"`
	Content	string	`json:"content"`
}

type ChatRes struct {
	ID      string    `json:"id"`
    Object  string    `json:"object"`
    Created int       `json:"created"`
    Model   string    `json:"model"`
    Choices []Choice  `json:"choices"`
    Usage   UsageInfo `json:"usage"`
}

type Choice struct {
    Index        int      `json:"index"`
    Message      Message  `json:"message"`
    FinishReason string   `json:"finish_reason"`
}

type UsageInfo struct {
	PromptTokens     int `json:"prompt_tokens"`
    CompletionTokens int `json:"completion_tokens"`
    TotalTokens      int `json:"total_tokens"`
}

func main() {
	app := goblet.New(true)

	app.Handle("/", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"title": "ChatSeek",
		}
		app.Render(w, "index.html", data)
	})

    app.Handle("/chat/completions", func(w http.ResponseWriter, r *http.Request) {
        var req ChatReq
        err := json.NewDecoder(r.Body).Decode(&req)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        cablyURL := "https://cablyai.com/v1/chat/completions"
        jsonReq, err := json.Marshal(req)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        resp, err := http.Post(cablyURL, "application/json", bytes.NewBuffer(jsonReq))
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
            http.Error(w, "CablyAI error: "+resp.Status, resp.StatusCode)
            return
        }

        var res ChatRes
        err = json.NewDecoder(resp.Body).Decode(&res)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        response := map[string]interface{}{
            "id":      res.ID,
            "object":  res.Object,
            "created": res.Created,
            "model":   res.Model,
            "choices": res.Choices,
            "usage":   res.Usage,
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    })

	err := app.Serve("6969")
	if err != nil {
		panic(err)
	}
}
