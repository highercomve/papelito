package types

import (
	"net/url"

	"github.com/highercomve/papelito/utils"
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

	asp.Query = utils.GetMongoQueryFromQuery(url.Query())
	asp.Sort = utils.GetMongoSortingFromQuery(url.Query())
	asp.Projection = utils.GetMongoFieldsFromQuery(url.Query())
	asp.Pagination = utils.GetMongoPaginationFromQuery(url.Query())

	if _, ok := asp.Pagination["offset"]; !ok {
		asp.Pagination["offset"] = 0
		query := url.Query()
		query.Add("page[offset]", "0")
		url.RawQuery = query.Encode()
	}

	return asp
}
