package api

type Response struct {
	Status 			string `json:"status"`
	Message 		string `json:"message"`
	StatusCode 		int	   `json:"status_code"`
	TransactionId 	string `json:"transaction_id,omitempty"`
}

