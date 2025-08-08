package model

import "time"

type PushEvent struct {
	Ref        string `json:"ref"`
	Repository struct {
		CloneURL string `json:"clone_url"`
	} `json:"repository"`
	HeadCommit struct {
		Message string `json:"message"`
		URL       string    `json:"url"`
		Timestamp time.Time `json:"timestamp"`
		Author  struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"author"`
	} `json:"head_commit"`
}




// package model
//
// import (
// 	"encoding/json"
// )
//
// // GitHubPayload represents the common structure of GitHub webhook payloads
// type GitHubPayload struct {
// 	Action     string          `json:"action,omitempty"`
// 	Repository Repository      `json:"repository"`
// 	Sender     User            `json:"sender"`
// 	PullRequest *PullRequest   `json:"pull_request,omitempty"`
// 	Push       *PushPayload    `json:"push,omitempty"`
// 	RawPayload json.RawMessage `json:"-"`
// }
//
// type Repository struct {
// 	ID       int    `json:"id"`
// 	Name     string `json:"name"`
// 	FullName string `json:"full_name"`
// 	CloneURL string `json:"clone_url"`
// 	SSHURL   string `json:"ssh_url"`
// }
//
// type User struct {
// 	ID    int    `json:"id"`
// 	Login string `json:"login"`
// }
//
// type PullRequest struct {
// 	ID     int    `json:"id"`
// 	Number int    `json:"number"`
// 	Title  string `json:"title"`
// 	State  string `json:"state"`
// 	Head   Branch `json:"head"`
// 	Base   Branch `json:"base"`
// }
//
// type PushPayload struct {
// 	Ref     string   `json:"ref"`
// 	Before  string   `json:"before"`
// 	After   string   `json:"after"`
// 	Commits []Commit `json:"commits"`
// }
//
// type Branch struct {
// 	Ref string `json:"ref"`
// 	SHA string `json:"sha"`
// }
//
// type Commit struct {
// 	ID      string `json:"id"`
// 	Message string `json:"message"`
// 	Author  Author `json:"author"`
// }
//
// type Author struct {
// 	Name  string `json:"name"`
// 	Email string `json:"email"`
// }
//
