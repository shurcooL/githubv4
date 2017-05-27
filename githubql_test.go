package githubql

import "fmt"

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
