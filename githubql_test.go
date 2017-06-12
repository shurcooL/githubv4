package githubql_test

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/shurcooL/githubql"
)

func TestNewClient_nil(t *testing.T) {
	// Shouldn't panic with nil parameter.
	client := githubql.NewClient(nil)
	_ = client
}

func TestClient_Query(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/graphql", func(w http.ResponseWriter, req *http.Request) {
		if got, want := req.Method, http.MethodPost; got != want {
			t.Errorf("got request method: %v, want: %v", got, want)
		}
		body := mustRead(req.Body)
		if got, want := body, `{"query":"{viewer{login}}"}`+"\n"; got != want {
			t.Errorf("got body: %v, want %v", got, want)
		}
		mustWrite(w, `{"data": {"viewer": {"login": "gopher"}}}`)
	})
	client := githubql.NewClient(&http.Client{Transport: localRoundTripper{mux: mux}})

	type query struct {
		Viewer struct {
			Login githubql.String
		}
	}

	var q query
	err := client.Query(context.Background(), &q, nil)
	if err != nil {
		t.Fatal(err)
	}
	got := q

	var want query
	want.Viewer.Login = "gopher"
	if !reflect.DeepEqual(got, want) {
		t.Errorf("client.Query got: %v, want: %v", got, want)
	}
}

func TestClient_Query_errorResponse(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/graphql", func(w http.ResponseWriter, req *http.Request) {
		mustWrite(w, `{
			"data": null,
			"errors": [
				{
					"message": "Field 'bad' doesn't exist on type 'Query'",
					"locations": [
						{
							"line": 7,
							"column": 3
						}
					]
				}
			]
		}`)
	})
	client := githubql.NewClient(&http.Client{Transport: localRoundTripper{mux: mux}})

	var q struct {
		Bad githubql.String
	}
	err := client.Query(context.Background(), &q, nil)
	if err == nil {
		t.Fatal("got error: nil, want: non-nil")
	}
	if got, want := err.Error(), "Field 'bad' doesn't exist on type 'Query'"; got != want {
		t.Errorf("got error: %v, want: %v", got, want)
	}
}

func TestClient_Query_errorStatusCode(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/graphql", func(w http.ResponseWriter, req *http.Request) {
		http.Error(w, "404 Not Found", http.StatusNotFound)
	})
	client := githubql.NewClient(&http.Client{Transport: localRoundTripper{mux: mux}})

	var q struct {
		Viewer struct {
			Login githubql.String
		}
	}
	err := client.Query(context.Background(), &q, nil)
	if err == nil {
		t.Fatal("got error: nil, want: non-nil")
	}
	if got, want := err.Error(), "unexpected status: Not Found"; got != want {
		t.Errorf("got error: %v, want: %v", got, want)
	}
}

func TestClient_Mutate(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/graphql", func(w http.ResponseWriter, req *http.Request) {
		if got, want := req.Method, http.MethodPost; got != want {
			t.Errorf("got request method: %v, want: %v", got, want)
		}
		body := mustRead(req.Body)
		if got, want := body, `{"query":"mutation($input:AddReactionInput!){addReaction(input:$input){reaction{content},subject{id,reactionGroups{users{totalCount}}}}}","variables":{"input":{"subjectId":"MDU6SXNzdWUyMTc5NTQ0OTc=","content":"HOORAY"}}}`+"\n"; got != want {
			t.Errorf("got body: %v, want %v", got, want)
		}
		mustWrite(w, `{"data": {
			"addReaction": {
				"reaction": {
					"content": "HOORAY"
				},
				"subject": {
					"id": "MDU6SXNzdWUyMTc5NTQ0OTc=",
					"reactionGroups": [
						{
							"users": {"totalCount": 3}
						}
					]
				}
			}
		}}`)
	})
	client := githubql.NewClient(&http.Client{Transport: localRoundTripper{mux: mux}})

	type reactionGroup struct {
		Users struct {
			TotalCount githubql.Int
		}
	}
	type mutation struct {
		AddReaction struct {
			Reaction struct {
				Content githubql.ReactionContent
			}
			Subject struct {
				ID             githubql.ID
				ReactionGroups []reactionGroup
			}
		} `graphql:"addReaction(input:$input)"`
	}

	var m mutation
	input := githubql.AddReactionInput{
		SubjectID: "MDU6SXNzdWUyMTc5NTQ0OTc=",
		Content:   githubql.ReactionContentHooray,
	}
	err := client.Mutate(context.Background(), &m, input, nil)
	if err != nil {
		t.Fatal(err)
	}
	got := m

	var want mutation
	want.AddReaction.Reaction.Content = githubql.ReactionContentHooray
	want.AddReaction.Subject.ID = "MDU6SXNzdWUyMTc5NTQ0OTc="
	var rg reactionGroup
	rg.Users.TotalCount = 3
	want.AddReaction.Subject.ReactionGroups = []reactionGroup{rg}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("client.Query got: %v, want: %v", got, want)
	}
}

// localRoundTripper is an http.RoundTripper that executes HTTP transactions
// by using mux directly, instead of going over an HTTP connection.
type localRoundTripper struct {
	mux *http.ServeMux
}

func (l localRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	w := httptest.NewRecorder()
	l.mux.ServeHTTP(w, req)
	return w.Result(), nil
}

func mustRead(r io.Reader) string {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func mustWrite(w io.Writer, s string) {
	_, err := io.WriteString(w, s)
	if err != nil {
		panic(err)
	}
}
