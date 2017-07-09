package jsonutil_test

import (
	"encoding/json"
	"io"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/shurcooL/githubql"
	"github.com/shurcooL/githubql/internal/jsonutil"
)

func TestUnmarshalGraphQL(t *testing.T) {
	/*
		query {
			viewer {
				login
				createdAt
			}
		}
	*/
	type query struct {
		Viewer struct {
			Login     githubql.String
			CreatedAt githubql.DateTime
		}
	}
	var got query
	err := jsonutil.UnmarshalGraphQL([]byte(`{
		"viewer": {
			"login": "shurcooL-test",
			"createdAt": "2017-06-29T04:12:01Z"
		}
	}`), &got)
	if err != nil {
		t.Fatal(err)
	}
	var want query
	want.Viewer.Login = "shurcooL-test"
	want.Viewer.CreatedAt = githubql.DateTime{Time: time.Unix(1498709521, 0).UTC()}
	if !reflect.DeepEqual(got, want) {
		t.Error("not equal")
	}
}

func BenchmarkUnmarshalGraphQL(b *testing.B) {
	type query struct {
		Viewer struct {
			Login     githubql.String
			CreatedAt githubql.DateTime
		}
	}
	for i := 0; i < b.N; i++ {
		now := time.Now().UTC()
		var got query
		err := jsonutil.UnmarshalGraphQL([]byte(`{
			"viewer": {
				"login": "shurcooL-test",
				"createdAt": "`+now.Format(time.RFC3339Nano)+`"
			}
		}`), &got)
		if err != nil {
			b.Fatal(err)
		}
		var want query
		want.Viewer.Login = "shurcooL-test"
		want.Viewer.CreatedAt = githubql.DateTime{Time: now}
		if !reflect.DeepEqual(got, want) {
			b.Error("not equal")
		}
	}
}

func BenchmarkJSONUnmarshal(b *testing.B) {
	type query struct {
		Viewer struct {
			Login     githubql.String
			CreatedAt githubql.DateTime
		}
	}
	for i := 0; i < b.N; i++ {
		now := time.Now().UTC()
		var got query
		err := json.Unmarshal([]byte(`{
			"viewer": {
				"login": "shurcooL-test",
				"createdAt": "`+now.Format(time.RFC3339Nano)+`"
			}
		}`), &got)
		if err != nil {
			b.Fatal(err)
		}
		var want query
		want.Viewer.Login = "shurcooL-test"
		want.Viewer.CreatedAt = githubql.DateTime{Time: now}
		if !reflect.DeepEqual(got, want) {
			b.Error("not equal")
		}
	}
}

func BenchmarkJSONTokenize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		now := time.Now().UTC()
		dec := json.NewDecoder(strings.NewReader(`{
			"viewer": {
				"login": "shurcooL-test",
				"createdAt": "` + now.Format(time.RFC3339Nano) + `"
			}
		}`))
		var tokens int
		for {
			_, err := dec.Token()
			if err == io.EOF {
				break
			} else if err != nil {
				b.Error(err)
			}
			tokens++
		}
		if tokens != 9 {
			b.Error("not 9 tokens")
		}
	}
}

func TestUnmarshalGraphQL_graphqlTag(t *testing.T) {
	type query struct {
		Foo githubql.String `graphql:"baz"`
	}
	var got query
	err := jsonutil.UnmarshalGraphQL([]byte(`{
		"baz": "bar"
	}`), &got)
	if err != nil {
		t.Fatal(err)
	}
	want := query{
		Foo: "bar",
	}
	if !reflect.DeepEqual(got, want) {
		t.Error("not equal")
	}
}

func TestUnmarshalGraphQL_jsonTag(t *testing.T) {
	type query struct {
		Foo githubql.String `json:"baz"`
	}
	var got query
	err := jsonutil.UnmarshalGraphQL([]byte(`{
		"foo": "bar"
	}`), &got)
	if err != nil {
		t.Fatal(err)
	}
	want := query{
		Foo: "bar",
	}
	if !reflect.DeepEqual(got, want) {
		t.Error("not equal")
	}
}

func TestUnmarshalGraphQL_array(t *testing.T) {
	type query struct {
		Foo []githubql.String
	}
	var got query
	err := jsonutil.UnmarshalGraphQL([]byte(`{
		"foo": [
			"bar",
			"baz"
		]
	}`), &got)
	if err != nil {
		t.Fatal(err)
	}
	want := query{
		Foo: []githubql.String{"bar", "baz"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Error("not equal")
	}
}

// When unmarshaling into an array, its initial value should be overwritten
// (rather than appended to).
func TestUnmarshalGraphQL_arrayReset(t *testing.T) {
	var got = []string{"initial"}
	err := jsonutil.UnmarshalGraphQL([]byte(`["bar", "baz"]`), &got)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"bar", "baz"}
	if !reflect.DeepEqual(got, want) {
		t.Error("not equal")
	}
}

func TestUnmarshalGraphQL_objectArray(t *testing.T) {
	type query struct {
		Foo []struct {
			Name githubql.String
		}
	}
	var got query
	err := jsonutil.UnmarshalGraphQL([]byte(`{
		"foo": [
			{"name": "bar"},
			{"name": "baz"}
		]
	}`), &got)
	if err != nil {
		t.Fatal(err)
	}
	want := query{
		Foo: []struct{ Name githubql.String }{
			{"bar"},
			{"baz"},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Error("not equal")
	}
}

func TestUnmarshalGraphQL_pointer(t *testing.T) {
	type query struct {
		Foo *githubql.String
		Bar *githubql.String
	}
	var got query
	got.Bar = new(githubql.String) // Test that got.Bar gets set to nil.
	err := jsonutil.UnmarshalGraphQL([]byte(`{
		"foo": "foo",
		"bar": null
	}`), &got)
	if err != nil {
		t.Fatal(err)
	}
	want := query{
		Foo: githubql.NewString("foo"),
		Bar: nil,
	}
	if !reflect.DeepEqual(got, want) {
		t.Error("not equal")
	}
}

func TestUnmarshalGraphQL_objectPointerArray(t *testing.T) {
	type query struct {
		Foo []*struct {
			Name githubql.String
		}
	}
	var got query
	err := jsonutil.UnmarshalGraphQL([]byte(`{
		"foo": [
			{"name": "bar"},
			null,
			{"name": "baz"}
		]
	}`), &got)
	if err != nil {
		t.Fatal(err)
	}
	want := query{
		Foo: []*struct{ Name githubql.String }{
			{"bar"},
			nil,
			{"baz"},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Error("not equal")
	}
}

func TestUnmarshalGraphQL_multipleValues(t *testing.T) {
	type query struct {
		Foo githubql.String
	}
	err := jsonutil.UnmarshalGraphQL([]byte(`{"foo": "bar"}{"foo": "baz"}`), new(query))
	if err == nil {
		t.Fatal("got error: nil, want: non-nil")
	}
	if got, want := err.Error(), "invalid token '{' after top-level value"; got != want {
		t.Errorf("got error: %v, want: %v", got, want)
	}
}

func TestUnmarshalGraphQL_union(t *testing.T) {
	/*
		{
			__typename
			... on ClosedEvent {
				createdAt
				actor {login}
			}
			... on ReopenedEvent {
				createdAt
				actor {login}
			}
		}
	*/
	type githubqlActor struct{ Login githubql.String }
	type closedEvent struct {
		Actor     githubqlActor
		CreatedAt githubql.DateTime
	}
	type reopenedEvent struct {
		Actor     githubqlActor
		CreatedAt githubql.DateTime
	}
	type issueTimelineItem struct {
		Typename      string        `graphql:"__typename"`
		ClosedEvent   closedEvent   `graphql:"... on ClosedEvent"`
		ReopenedEvent reopenedEvent `graphql:"... on ReopenedEvent"`
	}
	var got issueTimelineItem
	err := jsonutil.UnmarshalGraphQL([]byte(`{
		"__typename": "ClosedEvent",
		"createdAt": "2017-06-29T04:12:01Z",
		"actor": {
			"login": "shurcooL-test"
		}
	}`), &got)
	if err != nil {
		t.Fatal(err)
	}
	want := issueTimelineItem{
		Typename: "ClosedEvent",
		ClosedEvent: closedEvent{
			Actor: githubqlActor{
				Login: "shurcooL-test",
			},
			CreatedAt: githubql.DateTime{Time: time.Unix(1498709521, 0).UTC()},
		},
		ReopenedEvent: reopenedEvent{
			Actor: githubqlActor{
				Login: "shurcooL-test",
			},
			CreatedAt: githubql.DateTime{Time: time.Unix(1498709521, 0).UTC()},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Error("not equal")
	}
}
