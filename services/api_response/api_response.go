package api_response

func NewMetaData() *MetaResponse {
	return &MetaResponse{
		TraceID: "your-trace-id",
		Success: true,
	}
}

type ImageSaveResponse struct {
	Meta *MetaResponse `json:"meta"`
	Data ImageSaveURL  `json:"data"`
}

type ImageSaveURL struct {
	ImageURL string `json:"image_url"`
}
