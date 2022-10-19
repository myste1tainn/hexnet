package msnet

type ErrorResponse struct {
	Code           string      `json:"code"`
	Message        string      `json:"message"`
	Description    string      `json:"description"`
	Title          string      `json:"title"`
	AdditionalInfo interface{} `json:"additionalInfo"`
}
