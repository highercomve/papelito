package types

import (
	"net/url"

	"github.com/highercomve/papelito/modules/helpers"
	"go.mongodb.org/mongo-driver/bson"
)

type QueryOP struct {
	Query      bson.M
	Projection bson.M
	Sort       bson.M
	Pagination bson.M
}

func GetQueryOPFromURL(url *url.URL) QueryOP {
	asp := QueryOP{}

	asp.Query = helpers.GetMongoQueryFromQuery(url.Query())
	asp.Sort = helpers.GetMongoSortingFromQuery(url.Query())
	asp.Projection = helpers.GetMongoFieldsFromQuery(url.Query())
	asp.Pagination = helpers.GetMongoPaginationFromQuery(url.Query())

	if _, ok := asp.Pagination["offset"]; !ok {
		asp.Pagination["offset"] = 0
		query := url.Query()
		query.Add("page[offset]", "0")
		url.RawQuery = query.Encode()
	}

	return asp
}
