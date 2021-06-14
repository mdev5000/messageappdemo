package handler

import (
	"net/http"
	"strconv"
	"strings"
)

// GetQueryParams extracts common query parameters that are used for filtering the results contained with a REST store.
// Ex. ?fields=only,some&pageSize=20&pageStartIndex=10 = fields=[]string{only,some}, limit=20, offset=10
func GetQueryParams(op string, r *http.Request) (fields map[string]struct{}, limit uint64, offset uint64, err error) {
	fields = map[string]struct{}{}

	if limitS := r.URL.Query().Get("pageSize"); limitS != "" {
		limit, err = strconv.ParseUint(limitS, 10, 64)
		if err != nil {
			re := ResponseError(op)
			re.AddResponse("invalid pageSize value")
			err = &re
			return
		}
	}

	if offsetS := r.URL.Query().Get("pageStartIndex"); offsetS != "" {
		offset, err = strconv.ParseUint(offsetS, 10, 64)
		if err != nil {
			re := ResponseError(op)
			re.AddResponse("invalid pageStartIndex value")
			err = &re
			return
		}
	}

	if fieldsS := r.URL.Query().Get("fields"); fieldsS != "" {
		fieldsRaw := strings.Split(fieldsS, ",")
		for _, fRaw := range fieldsRaw {
			fRaw = strings.TrimSpace(fRaw)
			if fRaw == "" {
				continue
			}
			fields[fRaw] = struct{}{}
		}
	}

	return
}
