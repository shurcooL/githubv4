package githubql

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shurcooL/go/ctxhttp"
)

// Client is a GitHub GraphQL API v4 client.
type Client struct {
	httpClient *http.Client
}

// NewClient create a new GitHub GraphQL API v4 client with the provided http.Client.
// If it's nil, then http.DefaultClient is used.
//
// Note that GitHub GraphQL API v4 requires authentication, so
// the provided http.Client is expected to take care of that.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{
		httpClient: httpClient,
	}
}

// Query executes a single GraphQL query request,
// with a query derived from q, populating the response into it.
// q should be a pointer to struct that corresponds to the GitHub GraphQL schema.
func (c *Client) Query(ctx context.Context, q interface{}, variables map[string]interface{}) error {
	return c.do(ctx, queryOperation, q, variables)
}

// Mutate executes a single GraphQL mutation request,
// with a mutation derived from m, populating the response into it.
// m should be a pointer to struct that corresponds to the GitHub GraphQL schema.
// Provided input will be set as a variable named "input".
func (c *Client) Mutate(ctx context.Context, m interface{}, input Input, variables map[string]interface{}) error {
	if variables == nil {
		variables = map[string]interface{}{"input": input}
	} else {
		variables["input"] = input
	}
	return c.do(ctx, mutationOperation, m, variables)
}

// do executes a single GraphQL operation.
func (c *Client) do(ctx context.Context, op operationType, v interface{}, variables map[string]interface{}) error {
	var query string
	switch op {
	case queryOperation:
		query = constructQuery(v, variables)
	case mutationOperation:
		query = constructMutation(v, variables)
	}
	in := struct {
		Query     string                 `json:"query"`
		Variables map[string]interface{} `json:"variables,omitempty"`
	}{
		Query:     query,
		Variables: variables,
	}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(in)
	if err != nil {
		return err
	}
	resp, err := ctxhttp.Post(ctx, c.httpClient, "https://api.github.com/graphql", "application/json", &buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %v", resp.Status)
	}
	out := struct {
		Data   interface{}
		Errors errors
		//Extensions interface{} // Unused.
	}{Data: v}
	err = json.NewDecoder(resp.Body).Decode(&out)
	if err != nil {
		return err
	}
	if len(out.Errors) > 0 {
		return out.Errors
	}
	return nil
}

// errors represents the "errors" array in a response from a GraphQL server.
// If returned via error interface, the slice is expected to contain at least 1 element.
//
// Specification: https://facebook.github.io/graphql/#sec-Errors.
type errors []struct {
	Message   string
	Locations []struct {
		Line   int
		Column int
	}
}

// Error implements error interface.
func (e errors) Error() string {
	return e[0].Message
}

type operationType uint8

const (
	queryOperation operationType = iota
	mutationOperation
	//subscriptionOperation // Unused.
)
