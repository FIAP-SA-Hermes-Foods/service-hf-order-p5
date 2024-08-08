package strings

import (
	"encoding/json"

	l "service-hf-order-p5/external/logger"
)

func MarshalString(s interface{}) string {
	if s == nil {
		return ""
	}

	o, err := json.Marshal(s)
	if err != nil {
		l.Errorf("", "error in MarshalString client ", " | ", err)
		return ""
	}

	return string(o)
}
