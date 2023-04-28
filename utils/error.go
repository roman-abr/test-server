package utils

type HttpError struct {
	Code    int    `json:"-"`
	Message string `json:"error"`
}

func (e *HttpError) Error() string {
	return e.Message
}

func BadRequestError(text string) error {
	return &HttpError{Code: 400, Message: text}
}

func NotFoundError(text string) error {
	return &HttpError{Code: 404, Message: text}
}
