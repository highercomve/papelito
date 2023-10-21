package types

import (
	"time"

	"github.com/highercomve/papelito/modules/helper/helpermodels"
)

func Time(t time.Time) *time.Time {
	return &t
}

func Int64(t int64) *int64 {
	return &t
}

func String(t string) *string {
	return &t
}

func Bool(t bool) *bool {
	return &t
}

func MetadataFromHub(data interface{}) map[string]string {
	result := map[string]string{}

	for key, value := range data.(map[string]interface{}) {
		result[key] = value.(string)
	}

	return result
}

func Identification(i helpermodels.Identification) *helpermodels.Identification {
	return &i
}

func Ownership(i helpermodels.Ownership) *helpermodels.Ownership {
	return &i
}

func Timestamp(i helpermodels.Timestamp) *helpermodels.Timestamp {
	return &i
}
