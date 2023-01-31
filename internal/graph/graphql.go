package graph

import (
	"backend/internal/models"
	"errors"
	"strings"

	"github.com/graphql-go/graphql"
)

type Graph struct {
	Movies []*models.Movie
	QueryString string
	Config graphql.SchemaConfig
	fields graphql.Fields // defines the available actions on the data (in this case movieType data)
	movieType *graphql.Object //describe the data as it exists in the database
}

// New define factory function to get the instance of variable of the type Graph
func New(movies []*models.Movie) *Graph{
	// describes the kinds of things we wanna expose from our database
	var movieType = graphql.NewObject(
		graphql.ObjectConfig{
			Name : "Movie",
			Fields : graphql.Fields{
				"id":&graphql.Field{
					Type : graphql.Int,
				},
				"title":&graphql.Field{
					Type : graphql.String,
				},
				"release_date":&graphql.Field{
					Type : graphql.DateTime,
				},
				"runtime":&graphql.Field{
					Type : graphql.Int,
				},
				"mpaa_rating":&graphql.Field{
					Type : graphql.String,
				},
				"description":&graphql.Field{
					Type : graphql.String,
				},
				"image":&graphql.Field{
					Type : graphql.String,
				},
				"created_at":&graphql.Field{
					Type : graphql.DateTime,
				},
				"updated_at":&graphql.Field{
					Type : graphql.DateTime,
				},
			},
		},
	)

	// fields gives me the kinds of actions I want to perform, and every entry in fields has a name for the action itself.
	var fields = graphql.Fields{
		"list" : &graphql.Field{
			Type : graphql.NewList(movieType),
			Description:"Get all movies",
			Resolve : func(params graphql.ResolveParams) (interface{}, error){
				return movies, nil
			},
		},

		"search": &graphql.Field{
			Type : graphql.NewList(movieType),
			Description:  "Search movies by title",
			Args: graphql.FieldConfigArgument{
				"titleContains":&graphql.ArgumentConfig{
					Type:graphql.String,
				},
			},
			Resolve : func(params graphql.ResolveParams)(interface{}, error){
				var theList []*models.Movie
				search, ok := params.Args["titleContains"].(string)
				if ok {
					for _, currentMovie := range movies {
						if strings.Contains(strings.ToLower(currentMovie.Title), strings.ToLower(search)){
							theList = append(theList, currentMovie)
						}
					}
				}
				return theList, nil								
			},
		},

		"get":&graphql.Field{
			Type: movieType,
			Description : "Get movie by id",
			Args : graphql.FieldConfigArgument{
				"id" : &graphql.ArgumentConfig{
					Type:graphql.Int,
				},
			},
			Resolve : func(params graphql.ResolveParams) (interface{}, error){
				id, ok := params.Args["id"].(int)
				if ok {
					for _, movie := range movies {
						if movie.ID == id{
							return movie, nil
						}
					}
				}
				return nil, nil
			},
		},
	}

	return &Graph {
		Movies: movies,
		fields : fields,
		movieType : movieType,
	}
}

func (g *Graph) Query() (*graphql.Result, error){
	rootQuery := graphql.ObjectConfig{Name : "RootQuery", Fields : g.fields}
	schemaConfig := graphql.SchemaConfig{Query:graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)

	if err != nil {
		return nil, err
	}

	params := graphql.Params{Schema: schema, RequestString: g.QueryString}
	resp := graphql.Do(params)
	if len(resp.Errors) > 0{
		return nil, errors.New("error executing query")
	}

	return resp, nil
}