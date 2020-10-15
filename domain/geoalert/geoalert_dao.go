package geoalerts

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/micro-gis/geo-alert-api/client/elasticsearch"
	"github.com/micro-gis/geo-alert-api/domain/queries"
	"github.com/micro-gis/utils/rest_errors"
	"strings"
)

const (
	indexgeoalert = "geoalerts"
	typegeoalert  = "_doc"
)

func (i *GeoAlert) Save() rest_errors.RestErr {
	result, err := elasticsearch.Client.Index(indexgeoalert, typegeoalert, i)
	if err != nil {
		return rest_errors.NewInternalServerError("error when trying to save geoalert", errors.New("database error"))
	}
	i.Id = result.Id
	return nil
}

func (i *GeoAlert) Get() rest_errors.RestErr {
	geoalertId := i.Id
	result, err := elasticsearch.Client.Get(indexgeoalert, typegeoalert, i.Id)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return rest_errors.NewNotFoundError(fmt.Sprintf("no geoalerts found with id : %s", i.Id))
		}
		return rest_errors.NewInternalServerError(fmt.Sprintf("error when trying get id %s", i.Id), errors.New("database error"))
	}
	bytes, err := result.Source.MarshalJSON()
	if err != nil {
		return rest_errors.NewInternalServerError("error when trying to parse database error", errors.New("database error"))
	}
	if err := json.Unmarshal(bytes, i); err != nil {
		return rest_errors.NewInternalServerError("error when trying to parse database error", errors.New("database error"))
	}
	i.Id = geoalertId
	return nil
}

func (i *GeoAlert) Search(query queries.EsQuery) ([]GeoAlert, rest_errors.RestErr) {
	result, err := elasticsearch.Client.Search(indexgeoalert, query.Build())
	if err != nil {
		return nil, rest_errors.NewInternalServerError("error when trying to search documents", errors.New("database error"))
	}

	geoalerts := make([]GeoAlert, result.TotalHits())
	for i, hit := range result.Hits.Hits {
		bytes, _ := hit.Source.MarshalJSON()
		geoalert := GeoAlert{}
		if err := json.Unmarshal(bytes, &geoalert); err != nil {
			return nil, rest_errors.NewInternalServerError("error when trying to parse response", errors.New("database error"))
		}
		geoalert.Id = hit.Id
		geoalerts[i] = geoalert
	}

	if len(geoalerts) == 0 {
		return nil, rest_errors.NewNotFoundError("not geoalerts found matching given criteria")
	}
	return geoalerts, nil
}

func (i *GeoAlert) Delete() rest_errors.RestErr {
	_, err := elasticsearch.Client.Delete(indexgeoalert, typegeoalert, i.Id)
	if err != nil {
		return rest_errors.NewNotFoundError("geoalert with given id was not found")
	}
	i = nil
	return nil
}

func (i *GeoAlert) Upsert(id string) rest_errors.RestErr {
	result, err := elasticsearch.Client.Upsert(indexgeoalert, typegeoalert, i, id)
	if err != nil {
		return rest_errors.NewInternalServerError("error when trying to save / update geoalert", errors.New("database error"))
	}
	i.Id = result.Id
	return nil
}
