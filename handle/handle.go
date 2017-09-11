package handle

import "encoding/json"

// Response : Structure of the JSON
type Response struct {
	Error    string `json:"error,omitempty"`
	Message  string `json:"message,omitempty"`
	QueryURL string `json:"query,omitempty"`
	ShortURL string `json:"shorturl,omitempty"`
}

// SetMessage : Sets a response message
func (r *Response) SetMessage(message string) {
	r.Message = message
}

func jsonifyError(err error) []byte {
	response := Response{Error: err.Error()}
	jsonData, _ := json.Marshal(response)
	return jsonData
}
