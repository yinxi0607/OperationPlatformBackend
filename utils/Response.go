package utils

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

const (
	SuccessMessage    = "success"
	SuccessCode       = 0
	DBErrorCode       = 100
	ParamsErrorCode   = 101
	InternalErrorCode = 102
	OthersErrorCode   = 103
)
