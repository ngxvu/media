package common

import (
	"strconv"
	"time"
)

func GenerateIdentityID() string {
	return "M" + strconv.FormatInt(time.Now().UnixNano(), 10)
}
