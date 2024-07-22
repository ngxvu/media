package api_response

type MetaResponse struct {
	TraceID string `json:"trace_id"`
	Success bool   `json:"success"`
}

type DataResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
