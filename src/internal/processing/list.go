package processing

import (
	"encoding/json"
	"fmt"
	"strings"

	"net/http"
	"net/url"
	"reflect"

	"github.com/emurray647/audioServer/internal/model"

	log "github.com/sirupsen/logrus"
)

func (p *RequestProcessor) CreateListHandler() http.HandlerFunc {

	// build a map from query name to string we want to use to filter
	// ex maxduration maps to "duration <= {value}"
	var temp model.WavFileDetails
	t := reflect.TypeOf(temp)
	queryMap := make(map[string]string)

	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("json")

		if t.Field(i).Type.Kind() == reflect.String {
			queryMap[tag] = fmt.Sprintf("%s = '{value}'", tag)
		} else {
			queryMap[tag] = fmt.Sprintf("%s = {value}", tag)
			queryMap[fmt.Sprintf("max%s", tag)] = fmt.Sprintf("%s <= {value}", tag)
			queryMap[fmt.Sprintf("min%s", tag)] = fmt.Sprintf("%s >= {value}", tag)
		}

	}

	return func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()

		details, err := p.list(values, queryMap)
		if err != nil {
			log.Errorf(err.Error())
			setStatus(w, http.StatusInternalServerError, "unknown err", false)
			return
		}

		err = json.NewEncoder(w).Encode(details)
		if err != nil {
			log.Errorf(err.Error())
			setStatus(w, http.StatusInternalServerError, "failed to encode wav JSON", false)
			return
		}
	}
}

func (p *RequestProcessor) list(values url.Values, queryMap map[string]string) (*model.WavFilesDetailsSlice, error) {
	// create a slice for all the filters we want to apply for this search
	filters := make([]string, 0)
	for key, value := range values {
		queryFormat, ok := queryMap[key]
		if !ok {
			log.Warnf("unknown query parameter %s", key)
			continue
		}

		replacer := strings.NewReplacer("{value}", value[0])
		filterString := replacer.Replace(queryFormat)
		filters = append(filters, filterString)
	}

	result, err := p.db.GetWavs(filters)
	if err != nil {
		return nil, fmt.Errorf("failed retrieving wav details from db: %w", err)
	}

	return result, nil
}
