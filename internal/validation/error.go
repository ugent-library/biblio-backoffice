package validation

type Errors []*Error

type Error struct {
	Pointer string
	Field   string
	Code    string
}

func (errs Errors) Error() string {
	msg := ""
	for i, e := range errs {
		msg += e.Error()
		if i < len(errs)-1 {
			msg += ", "
		}
	}
	return msg
}

func (errs Errors) At(ptr string) *Error {
	for _, e := range errs {
		if e.Pointer == ptr {
			return e
		}
	}
	return nil
}

func (e *Error) Error() string {
	msg := e.Code
	if e.Pointer != "" {
		msg += "[" + e.Pointer + "]"
	}
	return msg
}
