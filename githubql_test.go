package githubql

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestNewClient_nil(t *testing.T) {
	// Shouldn't panic with nil parameter.
	client := NewClient(nil)
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
	client := NewClient(&http.Client{Transport: localRoundTripper{mux: mux}})

	type query struct {
		Viewer struct {
			Login String
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

func TestClient_Query_error(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/graphql", func(w http.ResponseWriter, req *http.Request) {
		http.Error(w, "404 Not Found", http.StatusNotFound)
	})
	client := NewClient(&http.Client{Transport: localRoundTripper{mux: mux}})

	type query struct {
		Viewer struct {
			Login String
		}
	}

	var q query
	err := client.Query(context.Background(), &q, nil)
	if err == nil {
		t.Fatal("got error: nil, want: non-nil")
	}
	if got, want := err.Error(), "unexpected status: Not Found"; got != want {
		t.Fatalf("got error: %v, want: %v", got, want)
	}
}

func TestClient_Mutate(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/graphql", func(w http.ResponseWriter, req *http.Request) {
		if got, want := req.Method, http.MethodPost; got != want {
			t.Errorf("got request method: %v, want: %v", got, want)
		}
		body := mustRead(req.Body)
		if got, want := body, `{"query":"mutation($Input:AddReactionInput!){addReaction(input:$Input){reaction{content},subject{id,reactionGroups{users{totalCount}}}}}","variables":{"Input":{"subjectId":"MDU6SXNzdWUyMTc5NTQ0OTc=","content":"HOORAY"}}}`+"\n"; got != want {
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
	client := NewClient(&http.Client{Transport: localRoundTripper{mux: mux}})

	type reactionGroup struct {
		Users struct {
			TotalCount Int
		}
	}
	type mutation struct {
		AddReaction struct {
			Reaction struct {
				Content ReactionContent
			}
			Subject struct {
				ID             ID
				ReactionGroups []reactionGroup
			}
		} `graphql:"addReaction(input:$Input)"`
	}

	var m mutation
	variables := map[string]interface{}{
		"Input": AddReactionInput{
			SubjectID: "MDU6SXNzdWUyMTc5NTQ0OTc=",
			Content:   Hooray,
		},
	}
	err := client.Mutate(context.Background(), &m, variables)
	if err != nil {
		t.Fatal(err)
	}
	got := m

	var want mutation
	want.AddReaction.Reaction.Content = Hooray
	want.AddReaction.Subject.ID = "MDU6SXNzdWUyMTc5NTQ0OTc="
	var rg reactionGroup
	rg.Users.TotalCount = 3
	want.AddReaction.Subject.ReactionGroups = []reactionGroup{rg}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("client.Query got: %v, want: %v", got, want)
	}
}

func TestQueryArguments(t *testing.T) {
	tests := []struct {
		name string
		in   map[string]interface{}
		want string
	}{
		{
			in:   map[string]interface{}{"A": Int(123), "B": NewBoolean(true)},
			want: "$A:Int!$B:Boolean",
		},
	}
	for _, tc := range tests {
		got := queryArguments(tc.in)
		if got != tc.want {
			t.Errorf("%s: got: %q, want: %q", tc.name, got, tc.want)
		}
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

// TODO: These should be converted into real tests. I used examples during
//       rapid prototyping because they're quicker to write (less boilerplate).

func ExampleInternalTest_mixedCapsToLowerCamelCase() {
	fmt.Println(mixedCapsToLowerCamelCase("DatabaseID"))
	fmt.Println(mixedCapsToLowerCamelCase("URL"))
	fmt.Println(mixedCapsToLowerCamelCase("ID"))
	fmt.Println(mixedCapsToLowerCamelCase("CreatedAt"))
	fmt.Println(mixedCapsToLowerCamelCase("Login"))
	fmt.Println(mixedCapsToLowerCamelCase("ResetAt"))

	// Output:
	// databaseId
	// url
	// id
	// createdAt
	// login
	// resetAt
}

func ExampleInternalTest_constructQuery() {
	{
		var v struct {
			Viewer struct {
				Login      String
				CreatedAt  DateTime
				ID         ID
				DatabaseID Int
			}
			RateLimit struct {
				Cost      Int
				Limit     Int
				Remaining Int
				ResetAt   DateTime
			}
		}
		query := constructQuery(v, nil)
		fmt.Println(query)
	}

	{
		var v struct {
			Repository struct {
				DatabaseID Int
				URL        URI

				Issue struct {
					Comments struct {
						Edges []struct {
							Node struct {
								Body   String
								Author struct {
									Login String
								}
								Editor struct {
									Login String
								}
							}
							Cursor String
						}
					} `graphql:"comments(first:1after:\"Y3Vyc29yOjE5NTE4NDI1Ng==\")"`
				} `graphql:"issue(number:1)"`
			} `graphql:"repository(owner:\"shurcooL-test\"name:\"test-repo\")"`
		}
		query := constructQuery(v, nil)
		fmt.Println(query)
	}

	{
		type githubqlActor struct {
			Login     String
			AvatarURL URI
			URL       URI
		}

		var v struct {
			Repository struct {
				DatabaseID Int
				URL        URI

				Issue struct {
					Comments struct {
						Edges []struct {
							Node struct {
								DatabaseID      Int
								Author          githubqlActor
								PublishedAt     DateTime
								LastEditedAt    *DateTime
								Editor          *githubqlActor
								Body            String
								ViewerCanUpdate Boolean
							}
							Cursor String
						}
					} `graphql:"comments(first:1)"`
				} `graphql:"issue(number:1)"`
			} `graphql:"repository(owner:\"shurcooL-test\"name:\"test-repo\")"`
		}
		query := constructQuery(v, nil)
		fmt.Println(query)
	}

	{
		type githubqlActor struct {
			Login     String
			AvatarURL URI `graphql:"avatarUrl(size:72)"`
			URL       URI
		}

		var v struct {
			Repository struct {
				Issue struct {
					Author         githubqlActor
					PublishedAt    DateTime
					LastEditedAt   *DateTime
					Editor         *githubqlActor
					Body           String
					ReactionGroups []struct {
						Content ReactionContent
						Users   struct {
							TotalCount Int
						}
						ViewerHasReacted Boolean
					}
					ViewerCanUpdate Boolean

					Comments struct {
						Nodes []struct {
							DatabaseID     Int
							Author         githubqlActor
							PublishedAt    DateTime
							LastEditedAt   *DateTime
							Editor         *githubqlActor
							Body           String
							ReactionGroups []struct {
								Content ReactionContent
								Users   struct {
									TotalCount Int
								}
								ViewerHasReacted Boolean
							}
							ViewerCanUpdate Boolean
						}
						PageInfo struct {
							EndCursor   String
							HasNextPage Boolean
						}
					} `graphql:"comments(first:1)"`
				} `graphql:"issue(number:1)"`
			} `graphql:"repository(owner:\"shurcooL-test\"name:\"test-repo\")"`
		}
		query := constructQuery(v, nil)
		fmt.Println(query)
	}

	{
		var v struct {
			Repository struct {
				Issue struct {
					Body String
				} `graphql:"issue(number: 1)"`
			} `graphql:"repository(owner:\"shurcooL-test\"name:\"test-repo\")"`
		}
		query := constructQuery(v, nil)
		fmt.Println(query)
	}

	{
		var v struct {
			Repository struct {
				Issue struct {
					Body String
				} `graphql:"issue(number: $IssueNumber)"`
			} `graphql:"repository(owner: $RepositoryOwner, name: $RepositoryName)"`
		}
		variables := map[string]interface{}{
			"RepositoryOwner": String("shurcooL-test"),
			"RepositoryName":  String("test-repo"),
			"IssueNumber":     Int(1),
		}
		query := constructQuery(v, variables)
		fmt.Println(query)
	}

	{
		var v struct {
			Repository struct {
				Issue struct {
					ReactionGroups []struct {
						Users struct {
							Nodes []struct {
								Login String
							}
						} `graphql:"users(first:10)"`
					}
				} `graphql:"issue(number: $IssueNumber)"`
			} `graphql:"repository(owner: $RepositoryOwner, name: $RepositoryName)"`
		}
		variables := map[string]interface{}{
			"RepositoryOwner": String("shurcooL-test"),
			"RepositoryName":  String("test-repo"),
			"IssueNumber":     Int(1),
		}
		query := constructQuery(v, variables)
		fmt.Println(query)
	}

	// Output:
	// {viewer{login,createdAt,id,databaseId},rateLimit{cost,limit,remaining,resetAt}}
	// {repository(owner:"shurcooL-test"name:"test-repo"){databaseId,url,issue(number:1){comments(first:1after:"Y3Vyc29yOjE5NTE4NDI1Ng=="){edges{node{body,author{login},editor{login}},cursor}}}}}
	// {repository(owner:"shurcooL-test"name:"test-repo"){databaseId,url,issue(number:1){comments(first:1){edges{node{databaseId,author{login,avatarUrl,url},publishedAt,lastEditedAt,editor{login,avatarUrl,url},body,viewerCanUpdate},cursor}}}}}
	// {repository(owner:"shurcooL-test"name:"test-repo"){issue(number:1){author{login,avatarUrl(size:72),url},publishedAt,lastEditedAt,editor{login,avatarUrl(size:72),url},body,reactionGroups{content,users{totalCount},viewerHasReacted},viewerCanUpdate,comments(first:1){nodes{databaseId,author{login,avatarUrl(size:72),url},publishedAt,lastEditedAt,editor{login,avatarUrl(size:72),url},body,reactionGroups{content,users{totalCount},viewerHasReacted},viewerCanUpdate},pageInfo{endCursor,hasNextPage}}}}}
	// {repository(owner:"shurcooL-test"name:"test-repo"){issue(number: 1){body}}}
	// query($IssueNumber:Int!$RepositoryName:String!$RepositoryOwner:String!){repository(owner: $RepositoryOwner, name: $RepositoryName){issue(number: $IssueNumber){body}}}
	// query($IssueNumber:Int!$RepositoryName:String!$RepositoryOwner:String!){repository(owner: $RepositoryOwner, name: $RepositoryName){issue(number: $IssueNumber){reactionGroups{users(first:10){nodes{login}}}}}}
}

func ExampleInternalTest_constructMutation() {
	{
		var m struct {
			AddReaction struct {
				Subject struct {
					ReactionGroups []struct {
						Users struct {
							TotalCount Int
						}
					}
				}
			} `graphql:"addReaction(input:$Input)"`
		}
		variables := map[string]interface{}{
			"Input": AddReactionInput{
				SubjectID: "MDU6SXNzdWUyMzE1MjcyNzk=",
				Content:   ThumbsUp,
			},
		}
		mutation := constructMutation(m, variables)
		fmt.Println(mutation)
	}

	// Output:
	// mutation($Input:AddReactionInput!){addReaction(input:$Input){subject{reactionGroups{users{totalCount}}}}}
}
