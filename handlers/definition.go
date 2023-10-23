package handlers

type Data map[string]interface{}

type JSONResponse struct {
	Status  int                    `json:"status"`
	Data    map[string]interface{} `json:"data"`
	Message string                 `json:"message"`
}



type UpdateStatusBody struct {

	Remarks string `json:"remarks"`
}