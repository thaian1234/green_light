package domain

import (
	"fmt"
	"strconv"
)

type Runtime int32

func (r Runtime) MarshalJSON() ([]byte, error) {
	var suffix string
	if r <= 1 {
		suffix = "min"
	} else {
		suffix = "mins"
	}
	jsonValue := fmt.Sprintf("%d %s", r, suffix)
	quotedJSONValue := strconv.Quote(jsonValue)
	return []byte(quotedJSONValue), nil
}
