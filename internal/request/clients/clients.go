// Package clients contains code related to the clients used to perform requests.
package clients

import "github.com/google/go-github/v79/github"

// Clients stores the clients used to perform requests.
var Clients Type

// Type is the type of the clients data.
type Type struct {
	GitHub *github.Client
}
