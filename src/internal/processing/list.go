package processing

import (
	"encoding/json"
	"fmt"
	"strings"

	"net/http"
	"net/url"
	"reflect"

	"github.com/emurray647/audioServer/internal/dbconnector"
	"github.com/emurray647/audioServer/internal/model"

	log "github.com/sirupsen/logrus"
)

func CreateListHandler() http.HandlerFunc {

	// build a map from query name to string we want to use to filter
	// ex maxduration maps to "duration <= {value}"
	var temp model.WavFileDetails
	t := reflect.TypeOf(temp)
	queryMap := make(map[string]string)

	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("json")
		if t.Field(i).Type.Kind() == reflect.String {
			tag = fmt.Sprintf("'%s'", tag)
		}
		queryMap[tag] = fmt.Sprintf("%s = {value}", tag)
		queryMap[fmt.Sprintf("max%s", tag)] = fmt.Sprintf("%s <= {value}", tag)
		queryMap[fmt.Sprintf("min%s", tag)] = fmt.Sprintf("%s >= {value}", tag)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()

		details, err := list(values, queryMap)
		if err != nil {
			log.Errorf(err.Error())
			setStatus(w, http.StatusInternalServerError, "unknown err", false)
			return
		}

		err = json.NewEncoder(w).Encode(details)
		if err != nil {
			log.Errorf(err.Error())
			// logError(w, http.StatusInternalServerError, fmt.Errorf("failed to encode wav JSON: %w", err))
			setStatus(w, http.StatusInternalServerError, "failed to encode wav JSON", false)
			return
		}
	}
}

func list(values url.Values, queryMap map[string]string) (*model.WavFilesDetailsSlice, error) {
	dbConnection, err := dbconnector.OpenDBConnection()
	if err != nil {
		return nil, fmt.Errorf("could not open database connection: %w", err)

	}
	defer dbConnection.Close()

	// create a slice for all the filters we want to apply for this search
	filters := make([]string, 0)
	for key, value := range values {
		queryFormat, ok := queryMap[key]
		if !ok {
			log.Warnf("unknown query paramter %s", key)
			continue
		}

		replacer := strings.NewReplacer("{value}", value[0])
		filterString := replacer.Replace(queryFormat)
		filters = append(filters, filterString)
	}

	result, err := dbConnection.GetWavs(filters)
	if err != nil {
		return nil, fmt.Errorf("failed retrieving wav details from db: %w", err)
	}

	return result, nil
}
