package api

type Error struct {
	Errors map[string]interface{} `json:"errors"`
}

func NewValidationError(err error) *Error {
	e := Error{}
	e.Errors = make(map[string]interface{})
	e.Errors["body"] = err.Error()
	return &e
}

func NewError(err error) *Error {
	e := Error{}
	e.Errors = make(map[string]interface{})
	e.Errors["message"] = err.Error()
	return &e
}
