package helpers

import (
	"net/url"
	"reflect"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func TestGetMongoQueryFromQuery(t *testing.T) {
	var time, _ = time.Parse(time.RFC3339, "2021-01-04T11:29:36.705Z")

	type args struct {
		querystring url.Values
	}
	tests := []struct {
		name string
		args args
		want bson.M
	}{
		{
			name: "Don't process sort_by",
			args: args{
				querystring: url.Values{
					"sort_by": []string{"sortby_blah1", "sortby_blah2"},
				},
			},
			want: bson.M{},
		},
		{
			name: "Don't process sort_by or fields",
			args: args{
				querystring: url.Values{
					"sort_by": []string{"sortby_blah1", "sortby_blah2"},
					"fields":  []string{"field_blah1,field_blah2"},
				},
			},
			want: bson.M{},
		},
		{
			name: "Process without operators",
			args: args{
				querystring: url.Values{
					"sort_by": []string{"sortby_blah1", "sortby_blah2"},
					"fields":  []string{"field_blah1,field_blah2"},
					"query1":  []string{"aaabbbccc"},
					"query2":  []string{"aaaaa", "bbbbbb"},
				},
			},
			want: bson.M{
				"query1": "aaabbbccc",
				"query2": bson.M{
					"$all": []string{"aaaaa", "bbbbbb"},
				},
			},
		},
		{
			name: "Process times",
			args: args{
				querystring: url.Values{
					"sort_by": []string{"sortby_blah1", "sortby_blah2"},
					"fields":  []string{"field_blah1,field_blah2"},
					"query1":  []string{"2021-01-04T11:29:36.705Z"},
				},
			},
			want: bson.M{
				"query1": time,
			},
		},
		{
			name: "Process with operators",
			args: args{
				querystring: url.Values{
					"sort_by": []string{"sortby_blah1", "sortby_blah2"},
					"fields":  []string{"field_blah1,field_blah2"},
					"query1":  []string{"ne:aaabbbccc"},
					"query2":  []string{"eq:aaabbbccc"},
					"query3":  []string{"gt:1"},
					"query4":  []string{"gte:1"},
					"query5":  []string{"lt:1"},
					"query6":  []string{"lte:1"},
					"query7":  []string{"all:1,2,3,4"},
				},
			},
			want: bson.M{
				"query1": bson.M{
					"$ne": "aaabbbccc",
				},
				"query2": bson.M{
					"$eq": "aaabbbccc",
				},
				"query3": bson.M{
					"$gt": 1,
				},
				"query4": bson.M{
					"$gte": 1,
				},
				"query5": bson.M{
					"$lt": 1,
				},
				"query6": bson.M{
					"$lte": 1,
				},
				"query7": bson.M{
					"$all": []string{"1", "2", "3", "4"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetMongoQueryFromQuery(tt.args.querystring)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMongoQueryFromQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMongoSortingFromQuery(t *testing.T) {
	type args struct {
		querystring url.Values
	}
	tests := []struct {
		name string
		args args
		want bson.M
	}{
		{
			name: "Don't process filters or fields",
			args: args{
				querystring: url.Values{
					"query1": []string{"sortby_blah1", "sortby_blah2"},
					"fields": []string{"sortby_blah1", "sortby_blah2"},
				},
			},
			want: bson.M{},
		},
		{
			name: "Process sorting without direction",
			args: args{
				querystring: url.Values{
					"query1":  []string{"sortby_blah1", "sortby_blah2"},
					"fields":  []string{"fields_blah1", "fields_blah2"},
					"sort_by": []string{"email", "updated_at"},
				},
			},
			want: bson.M{
				"email":      1,
				"updated_at": 1,
			},
		},
		{
			name: "Process sorting with direction",
			args: args{
				querystring: url.Values{
					"query1":  []string{"sortby_blah1", "sortby_blah2"},
					"fields":  []string{"fields_blah1", "fields_blah2"},
					"sort_by": []string{"asc:email", "desc:updated_at"},
				},
			},
			want: bson.M{
				"email":      1,
				"updated_at": -1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMongoSortingFromQuery(tt.args.querystring); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMongoSortingFromQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMongoFieldsFromQuery(t *testing.T) {
	type args struct {
		querystring url.Values
	}
	tests := []struct {
		name string
		args args
		want bson.M
	}{
		{
			name: "Don't process sorting or fields",
			args: args{
				querystring: url.Values{
					"sort_by": []string{"sortby_blah1", "sortby_blah2"},
					"query":   []string{"query_blah1", "query_blah2"},
				},
			},
			want: bson.M{},
		},
		{
			name: "Process fields without operation",
			args: args{
				querystring: url.Values{
					"sort_by": []string{"sortby_blah1", "sortby_blah2"},
					"fields":  []string{"fields_blah1", "fields_blah2"},
				},
			},
			want: bson.M{
				"fields_blah1": 1,
				"fields_blah2": 1,
			},
		},
		{
			name: "Process fields without operation and merged",
			args: args{
				querystring: url.Values{
					"sort_by": []string{"sortby_blah1", "sortby_blah2"},
					"fields":  []string{"fields_blah1,fields_blah2"},
				},
			},
			want: bson.M{
				"fields_blah1": 1,
				"fields_blah2": 1,
			},
		},
		{
			name: "Process fields with operation",
			args: args{
				querystring: url.Values{
					"sort_by": []string{"sortby_blah1", "sortby_blah2"},
					"fields":  []string{"+fields_blah1", "-fields_blah2"},
				},
			},
			want: bson.M{
				"fields_blah1": 1,
				"fields_blah2": 0,
			},
		},
		{
			name: "Process fields with operation merged",
			args: args{
				querystring: url.Values{
					"sort_by": []string{"sortby_blah1", "sortby_blah2"},
					"fields":  []string{"+fields_blah1,-fields_blah2"},
				},
			},
			want: bson.M{
				"fields_blah1": 1,
				"fields_blah2": 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMongoFieldsFromQuery(tt.args.querystring); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMongoFieldsFromQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMongoPaginationFromQuery(t *testing.T) {
	beforeString := "2021-01-10T00:22:59.095Z"
	afterString := "2020-12-20T00:22:59.095Z"
	before, _ := time.Parse(time.RFC3339, beforeString)
	after, _ := time.Parse(time.RFC3339, afterString)

	type args struct {
		querystring url.Values
	}
	tests := []struct {
		name string
		args args
		want bson.M
	}{
		{
			name: "Don't process other fields as pagination pagination",
			args: args{
				querystring: url.Values{
					"query1":  []string{"sortby_blah1", "sortby_blah2"},
					"fields":  []string{"fields_blah1", "fields_blah2"},
					"sort_by": []string{"asc:email", "desc:updated_at"},
				},
			},
			want: bson.M{
				"limit": 100,
			},
		},
		{
			name: "Process pagination size and before",
			args: args{
				querystring: url.Values{
					"query1":       []string{"sortby_blah1", "sortby_blah2"},
					"fields":       []string{"fields_blah1", "fields_blah2"},
					"sort_by":      []string{"asc:email", "desc:updated_at"},
					"page[size]":   []string{"50"},
					"page[before]": []string{beforeString},
				},
			},
			want: bson.M{
				"limit":  50,
				"before": before,
			},
		},
		{
			name: "Process pagination size and after",
			args: args{
				querystring: url.Values{
					"query1":      []string{"sortby_blah1", "sortby_blah2"},
					"fields":      []string{"fields_blah1", "fields_blah2"},
					"sort_by":     []string{"asc:email", "desc:updated_at"},
					"page[size]":  []string{"50"},
					"page[after]": []string{afterString},
				},
			},
			want: bson.M{
				"limit": 50,
				"after": after,
			},
		},
		{
			name: "Process pagination size and after",
			args: args{
				querystring: url.Values{
					"query1":       []string{"sortby_blah1", "sortby_blah2"},
					"fields":       []string{"fields_blah1", "fields_blah2"},
					"sort_by":      []string{"asc:email", "desc:updated_at"},
					"page[size]":   []string{"50"},
					"page[after]":  []string{afterString},
					"page[before]": []string{beforeString},
				},
			},
			want: bson.M{
				"limit":  50,
				"after":  after,
				"before": before,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMongoPaginationFromQuery(tt.args.querystring); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMongoPaginationFromQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
