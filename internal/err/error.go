package err

type Error struct {
	Code    int
	Message string
	Err     error
}

func (e *Error) Error() string {
	return e.Message
}
