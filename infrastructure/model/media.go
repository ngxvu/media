package model

import (
	"media-service/services/api_response"
)

type OTPDataResponse struct {
	Meta *api_response.MetaResponse `json:"meta"`
	Data *api_response.DataResponse `json:"data"`
}
