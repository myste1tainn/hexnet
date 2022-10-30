package msnet

import "fmt"

type ErrorResponse struct {
	Code           string      `json:"code"`
	Message        string      `json:"message"`
	Description    string      `json:"description"`
	Title          string      `json:"title"`
	AdditionalInfo interface{} `json:"additionalInfo"`
}

func (this ErrorResponse) Error() string {
	return fmt.Sprintf("%s: %s, %s, %s", this.Code, this.Title, this.Message, this.Description)
}
