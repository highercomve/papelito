package gamerepo

import (
	"context"
	"time"

	"github.com/highercomve/papelito/modules/game/gamemodels"
	"github.com/highercomve/papelito/modules/helpers"
	"github.com/highercomve/papelito/modules/helpers/helperrepo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// collectionName costumer collection name
const collectionName = "customers"

var repo *Repo

var searchableQueries = map[string]bool{
	"deleted_at": true,
	"updated_at": true,
	"created_at": true,
}

// Repo manage customer database interaction
type Repo struct {
	Repo helperrepo.Repoable
}

func SetRepo(s *Repo) {
	repo = s
}

// GetRepo Return account db singleton instance
func GetRepo() (*Repo, error) {
	if repo != nil {
		return repo, nil
	}

	st, err := helperrepo.NewStorage(helpers.Env.CollectionPrefix)
	if err != nil {
		return nil, err
	}

	collection := st.GetCollection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	t := true
	f := false

	_, err = collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: map[string]int{
				"created_at": 1,
			},
			Options: &options.IndexOptions{
				Unique:     &f,
				Background: &f,
				Sparse:     &f,
			},
		},
		{
			Keys: map[string]int{
				"deleted_at": 1,
			},
			Options: &options.IndexOptions{
				Unique:     &f,
				Background: &f,
				Sparse:     &f,
			},
		},
		{
			Keys: map[string]int{
				"updated_at": 1,
			},
			Options: &options.IndexOptions{
				Unique:     &f,
				Background: &f,
				Sparse:     &f,
			},
		},
		{
			Keys: map[string]int{
				"provider_resource_id": 1,
			},
			Options: &options.IndexOptions{
				Unique:     &t,
				Background: &f,
				Sparse:     &f,
			},
		},
		{
			Keys: bson.D{
				{Key: "provider_id", Value: 1},
				{Key: "provider_resource_id", Value: 1},
			},
			Options: &options.IndexOptions{
				Unique:     &f,
				Background: &f,
				Sparse:     &f,
			},
		},
		{
			Keys: map[string]int{
				"owner_id": 1,
			},
			Options: &options.IndexOptions{
				Unique:     &f,
				Background: &f,
				Sparse:     &f,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	repo = &Repo{
		Repo: &helperrepo.Repo{
			Collection: collection,
			Storage:    st,
		},
	}

	return repo, nil
}

// SearchBy search logs by anything
func (db *Repo) SearchBy(ctx context.Context, q, p, s, pa map[string]interface{}) ([]gamemodels.Game, error) {
	result := []gamemodels.Game{}
	query := helperrepo.Map{
		"deleted_at": nil,
	}

	for key, value := range q {
		if _, ok := searchableQueries[key]; !ok {
			continue
		}

		query[key] = value
	}

	sort := helperrepo.Map{}
	if s != nil {
		sort = s
	}

	sort["created_at"] = -1
	queryOptions := options.FindOptions{}

	if p != nil {
		queryOptions.Projection = helperrepo.MergeDefaultProjection(p)
	}

	helpers.SetMongoPagination(query, sort, pa, &queryOptions)

	err := db.Repo.Find(ctx, query, &result, &queryOptions)

	return result, err
}
