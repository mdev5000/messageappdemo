package handler

import (
	"net/http"
	"strconv"
	"strings"
)

// GetQueryParams extracts common query parameters that are used for filtering the results contained with a REST store.
// Ex. ?fields=only,some&pageSize=20&pageStartIndex=10 = fields=[]string{only,some}, limit=20, offset=200
func GetQueryParams(op string, r *http.Request) (fields map[string]struct{}, limit uint64, offset uint64, err error) {
	fields = map[string]struct{}{}

	if limitS := r.URL.Query().Get("pageSize"); limitS != "" {
		limit, err = strconv.ParseUint(limitS, 10, 64)
		if err != nil || limit < 1 {
			re := ResponseError(op)
			re.AddResponse("invalid pageSize value")
			err = &re
			return
		}
	}

	if offsetS := r.URL.Query().Get("pageStartIndex"); offsetS != "" {
		pageOffset, errPO := strconv.ParseUint(offsetS, 10, 64)
		if errPO != nil || pageOffset < 1 {
			re := ResponseError(op)
			re.AddResponse("invalid pageStartIndex value")
			err = &re
			return
		}
		if limit == 0 {
			offset = pageOffset
		} else {
			offset = (pageOffset - 1) * limit
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
