package models

type ResponseErr struct {
	Message string
	Status  int
}

func (r *ResponseErr) Error() string {
	return r.Message
}
