package githubql

import "testing"

func TestConstructQuery(t *testing.T) {
	tests := []struct {
		inV         interface{}
		inVariables map[string]interface{}
		want        string
	}{
		{
			inV: struct {
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
			}{},
			want: `{viewer{login,createdAt,id,databaseId},rateLimit{cost,limit,remaining,resetAt}}`,
		},
		{
			inV: struct {
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
			}{},
			want: `{repository(owner:"shurcooL-test"name:"test-repo"){databaseId,url,issue(number:1){comments(first:1after:"Y3Vyc29yOjE5NTE4NDI1Ng=="){edges{node{body,author{login},editor{login}},cursor}}}}}`,
		},
		{
			inV: func() interface{} {
				type githubqlActor struct {
					Login     String
					AvatarURL URI
					URL       URI
				}

				return struct {
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
				}{}
			}(),
			want: `{repository(owner:"shurcooL-test"name:"test-repo"){databaseId,url,issue(number:1){comments(first:1){edges{node{databaseId,author{login,avatarUrl,url},publishedAt,lastEditedAt,editor{login,avatarUrl,url},body,viewerCanUpdate},cursor}}}}}`,
		},
		{
			inV: func() interface{} {
				type githubqlActor struct {
					Login     String
					AvatarURL URI `graphql:"avatarUrl(size:72)"`
					URL       URI
				}

				return struct {
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
				}{}
			}(),
			want: `{repository(owner:"shurcooL-test"name:"test-repo"){issue(number:1){author{login,avatarUrl(size:72),url},publishedAt,lastEditedAt,editor{login,avatarUrl(size:72),url},body,reactionGroups{content,users{totalCount},viewerHasReacted},viewerCanUpdate,comments(first:1){nodes{databaseId,author{login,avatarUrl(size:72),url},publishedAt,lastEditedAt,editor{login,avatarUrl(size:72),url},body,reactionGroups{content,users{totalCount},viewerHasReacted},viewerCanUpdate},pageInfo{endCursor,hasNextPage}}}}}`,
		},
		{
			inV: struct {
				Repository struct {
					Issue struct {
						Body String
					} `graphql:"issue(number: 1)"`
				} `graphql:"repository(owner:\"shurcooL-test\"name:\"test-repo\")"`
			}{},
			want: `{repository(owner:"shurcooL-test"name:"test-repo"){issue(number: 1){body}}}`,
		},
		{
			inV: struct {
				Repository struct {
					Issue struct {
						Body String
					} `graphql:"issue(number: $issueNumber)"`
				} `graphql:"repository(owner: $repositoryOwner, name: $repositoryName)"`
			}{},
			inVariables: map[string]interface{}{
				"repositoryOwner": String("shurcooL-test"),
				"repositoryName":  String("test-repo"),
				"issueNumber":     Int(1),
			},
			want: `query($issueNumber:Int!$repositoryName:String!$repositoryOwner:String!){repository(owner: $repositoryOwner, name: $repositoryName){issue(number: $issueNumber){body}}}`,
		},
		{
			inV: struct {
				Repository struct {
					Issue struct {
						ReactionGroups []struct {
							Users struct {
								Nodes []struct {
									Login String
								}
							} `graphql:"users(first:10)"`
						}
					} `graphql:"issue(number: $issueNumber)"`
				} `graphql:"repository(owner: $repositoryOwner, name: $repositoryName)"`
			}{},
			inVariables: map[string]interface{}{
				"repositoryOwner": String("shurcooL-test"),
				"repositoryName":  String("test-repo"),
				"issueNumber":     Int(1),
			},
			want: `query($issueNumber:Int!$repositoryName:String!$repositoryOwner:String!){repository(owner: $repositoryOwner, name: $repositoryName){issue(number: $issueNumber){reactionGroups{users(first:10){nodes{login}}}}}}`,
		},
	}
	for _, tc := range tests {
		got := constructQuery(tc.inV, tc.inVariables)
		if got != tc.want {
			t.Errorf("\ngot:  %q\nwant: %q\n", got, tc.want)
		}
	}
}

func TestConstructMutation(t *testing.T) {
	tests := []struct {
		inV         interface{}
		inVariables map[string]interface{}
		want        string
	}{
		{
			inV: struct {
				AddReaction struct {
					Subject struct {
						ReactionGroups []struct {
							Users struct {
								TotalCount Int
							}
						}
					}
				} `graphql:"addReaction(input:$input)"`
			}{},
			inVariables: map[string]interface{}{
				"input": AddReactionInput{
					SubjectID: "MDU6SXNzdWUyMzE1MjcyNzk=",
					Content:   ReactionContentThumbsUp,
				},
			},
			want: `mutation($input:AddReactionInput!){addReaction(input:$input){subject{reactionGroups{users{totalCount}}}}}`,
		},
	}
	for _, tc := range tests {
		got := constructMutation(tc.inV, tc.inVariables)
		if got != tc.want {
			t.Errorf("\ngot:  %q\nwant: %q\n", got, tc.want)
		}
	}
}

func TestQueryArguments(t *testing.T) {
	tests := []struct {
		name string
		in   map[string]interface{}
		want string
	}{
		{
			in:   map[string]interface{}{"a": Int(123), "b": NewBoolean(true)},
			want: "$a:Int!$b:Boolean",
		},
		{
			in:   map[string]interface{}{"states": []IssueState{IssueStateOpen, IssueStateClosed}},
			want: "$states:[IssueState!]",
		},
		{
			in:   map[string]interface{}{"states": []IssueState(nil)},
			want: "$states:[IssueState!]",
		},
		{
			in:   map[string]interface{}{"states": [...]IssueState{IssueStateOpen, IssueStateClosed}},
			want: "$states:[IssueState!]",
		},
	}
	for _, tc := range tests {
		got := queryArguments(tc.in)
		if got != tc.want {
			t.Errorf("%s: got: %q, want: %q", tc.name, got, tc.want)
		}
	}
}
