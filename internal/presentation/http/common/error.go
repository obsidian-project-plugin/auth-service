package common

type ErrorDto struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
}

func CreateError(
	message string,
	error error,
) *ErrorDto {
	return &ErrorDto{
		Message: message,
		Error:   error.Error(),
	}
}
