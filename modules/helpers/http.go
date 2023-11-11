package helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/highercomve/papelito/modules/helpers/helpermodels"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// DefaultPageSize default page size
var DefaultPageSize int = 100
var availableFormats []string = []string{"json", "html", "css"}

type ApiSearchPagination struct {
	Filters    bson.M
	Sort       bson.M
	Fields     bson.M
	Pagination bson.M
}

func GetFormat(c string) string {
	format := c
	for _, f := range availableFormats {
		if strings.Contains(c, f) {
			format = f
		}
	}
	return format
}

// Fetch do an http request
func Fetch(ctx context.Context, method, url string, content []byte, rp interface{}, headers map[string]string, queries map[string]string) ([]byte, error) {
	request, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(content))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %s", err.Error())
	}

	request.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		request.Header.Set(key, value)
	}

	q := request.URL.Query()
	for key, value := range queries {
		q.Add(key, value)
	}
	request.URL.RawQuery = q.Encode()

	client := http.DefaultClient
	if os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") != "" {
		client = &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()

	rs, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("%s", string(rs))
	}

	err = json.Unmarshal(rs, rp)
	return rs, err
}

// GetMongoSortingFromQuery get mongo sorting from query
func GetMongoSortingFromQuery(querystring url.Values) bson.M {
	sortBy := bson.M{}

	for queryKey, value := range querystring {
		if value == nil {
			continue
		}
		if !strings.Contains(queryKey, "sort_by") {
			continue
		}

		for _, key := range value {
			match := strings.SplitN(key, ":", 2)
			switch match[0] {
			case "asc":
				sortBy[match[1]] = 1
			case "desc":
				sortBy[match[1]] = -1
			default:
				sortBy[key] = 1
			}
		}
	}
	return sortBy
}

// GetMongoPaginationFromQuery get mongo pagination from query
func GetMongoPaginationFromQuery(querystring url.Values) bson.M {
	pagination := bson.M{
		"limit": DefaultPageSize,
	}

	for queryKey, value := range querystring {
		if value == nil {
			continue
		}
		if !strings.Contains(queryKey, "page") {
			continue
		}

		if queryKey == "page[size]" {
			pagination["limit"] = processValue(value[0])
		}

		if queryKey == "page[after]" {
			pagination["after"] = processValue(value[0])
		}

		if queryKey == "page[offset]" {
			pagination["offset"] = processValue(value[0])
		}

		if queryKey == "page[before]" {
			pagination["before"] = processValue(value[0])
		}
	}
	return pagination
}

// GetMongoFieldsFromQuery get mongo fields from query
func GetMongoFieldsFromQuery(querystring url.Values) bson.M {
	selectionFields := bson.M{}
	re := regexp.MustCompile(`([+-])(.*)`)

	for key, value := range querystring {
		if value == nil {
			continue
		}
		if !strings.Contains(key, "fields") {
			continue
		}

		for _, v := range value {
			fields := strings.Split(v, ",")
			for _, field := range fields {
				match := re.FindStringSubmatch(field)
				if len(match) == 3 {
					switch match[1] {
					case "-":
						selectionFields[match[2]] = 0
					case "+":
						selectionFields[match[2]] = 1
					default:
						selectionFields[match[2]] = 1
					}
				} else {
					selectionFields[field] = 1
				}
			}
		}
	}
	return selectionFields
}

// GetMongoQueryFromQuery get mongo query from url query
func GetMongoQueryFromQuery(querystring url.Values) bson.M {
	query := bson.M{}

	for key, value := range querystring {
		if value == nil {
			continue
		}

		isNotQuery :=
			strings.Contains(key, "fields") ||
				strings.Contains(key, "sort_by") ||
				strings.Contains(key, "page")

		if isNotQuery {
			continue
		}

		if len(value) > 1 {
			query[key] = bson.M{
				"$all": value,
			}
			continue
		}

		field := value[0]
		match := strings.SplitN(field, ":", 2)

		switch match[0] {
		case "in":
			values := strings.Split(match[1], ",")
			arr := make([]interface{}, len(values))
			for index, v := range values {
				arr[index] = processValue(v)
			}
			query[key] = bson.M{
				"$in": arr,
			}
		case "nin":
			values := strings.Split(match[1], ",")
			arr := make([]interface{}, len(values))
			for index, v := range values {
				arr[index] = processValue(v)
			}
			query[key] = bson.M{
				"$nin": arr,
			}
		case "exists":
			query[key] = bson.M{
				"$exists": processValue(match[1]),
			}
		case "eq":
			query[key] = bson.M{
				"$eq": processValue(match[1]),
			}
		case "ne":
			query[key] = bson.M{
				"$ne": processValue(match[1]),
			}
		case "lt":
			query[key] = bson.M{
				"$lt": processValue(match[1]),
			}
		case "lte":
			query[key] = bson.M{
				"$lte": processValue(match[1]),
			}
		case "gt":
			query[key] = bson.M{
				"$gt": processValue(match[1]),
			}
		case "gte":
			query[key] = bson.M{
				"$gte": processValue(match[1]),
			}
		case "all":
			query[key] = bson.M{
				"$all": strings.Split(match[1], ","),
			}
		case "empty":
			query[key] = bson.M{
				"$eq": "",
			}
		default:
			query[key] = processValue(field)
		}
	}

	return query
}

func processValue(v string) interface{} {
	var r interface{} = v

	if time, err := time.Parse(time.RFC3339, v); err == nil {
		return time
	}

	if i, err := strconv.Atoi(v); err == nil {
		return i
	}

	if b, err := strconv.ParseBool(v); err == nil {
		return b
	}

	if v == "null" {
		return nil
	}

	return r
}

// GetPaginationLink get pagination link
func GetPaginationLink(u url.URL, total int64, last, first helpermodels.Datable) helpermodels.Pagination {
	result := helpermodels.Pagination{
		Total:     int(total),
		PageSizes: []int{10, 20, 30, 50, 100},
	}

	finishTimestamp := last.GetCreatedAt().Format(time.RFC3339)
	startTimetamp := first.GetCreatedAt().Format(time.RFC3339)
	size := u.Query().Get("page[size]")
	if size == "" {
		size = strconv.Itoa(DefaultPageSize)
	}

	sizeInt, err := strconv.Atoi(size)
	if err != nil {
		sizeInt = DefaultPageSize
	}
	result.PageSize = sizeInt

	newURL, err := url.Parse(Env.HostURL)
	if err != nil {
		newURL = &u
	}
	newURL.Path = u.Path

	result.ServiceURL = newURL.String()

	newURL.RawQuery = u.Query().Encode()

	if u.Query().Get("page[offset]") == "" {
		prevURL := *newURL
		prevQuery := prevURL.Query()
		prevQuery.Set("page[size]", size)
		prevQuery.Set("page[before]", startTimetamp)
		prevURL.RawQuery = prevQuery.Encode()
		result.Prev = prevURL.String()

		if int(result.Total) >= sizeInt {
			nextURL := *newURL
			nextQuery := nextURL.Query()
			nextQuery.Set("page[size]", size)
			nextQuery.Set("page[after]", finishTimestamp)
			nextURL.RawQuery = nextQuery.Encode()
			result.Next = nextURL.String()
		}
	} else {
		offset, err := strconv.Atoi(u.Query().Get("page[offset]"))
		if err != nil {
			offset = 0
		}

		current := offset / sizeInt

		result.PageOffset = offset
		result.CurrentPage = current + 1

		prevURL := *newURL
		prevQuery := prevURL.Query()

		prevOffsetInt := offset - sizeInt
		if prevOffsetInt >= 0 {
			prevQuery.Set("page[offset]", strconv.Itoa(prevOffsetInt))
		}
		prevQuery.Set("page[size]", size)
		prevURL.RawQuery = prevQuery.Encode()
		result.Prev = prevURL.String()

		if int(result.Total) >= sizeInt*current {
			nextOffset := strconv.Itoa(offset + sizeInt)

			nextURL := *newURL
			nextQuery := nextURL.Query()
			nextQuery.Set("page[size]", size)
			nextQuery.Set("page[offset]", nextOffset)
			nextURL.RawQuery = nextQuery.Encode()

			result.Next = nextURL.String()
		}
	}

	return result
}

// SetMongoPagination set pagination
func SetMongoPagination(q, s bson.M, pa map[string]interface{}, queryOptions *options.FindOptions) {
	limit := int64(DefaultPageSize)
	if pa != nil {
		if l, ok := pa["limit"]; ok {
			limit = int64(l.(int))
		}
		if after, ok := pa["after"]; ok {
			s["created_at"] = -1
			if _, ok := q["created_at"]; ok {
				q["created_at"].(bson.M)["$lt"] = after
			} else {
				q["created_at"] = bson.M{
					"$lt": after,
				}
			}
		}
		if before, ok := pa["before"]; ok {
			s["created_at"] = -1
			if _, ok := q["created_at"]; ok {
				q["created_at"].(bson.M)["$gt"] = before
			} else {
				q["created_at"] = bson.M{
					"$gt": before,
				}
			}
		}
		if offset, ok := pa["offset"]; ok {
			queryOptions.SetSkip(int64(offset.(int)))
		} else {
			queryOptions.SetSkip(0)
		}
	}
	queryOptions.SetLimit(limit)

	mongoSort := bson.D{}
	for key, value := range s {
		mongoSort = append(mongoSort, bson.E{Key: key, Value: value})
	}
	queryOptions.SetSort(mongoSort)
}
