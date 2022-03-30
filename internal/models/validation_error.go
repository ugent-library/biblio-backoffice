package models

type ValidationErrors []ValidationError

type ValidationError struct {
	Pointer string
	Code    string
}

func (e ValidationError) Error() string {
	msg := e.Code
	if e.Pointer != "" {
		msg += "[" + e.Pointer + "]"
	}
	return msg
}

func (errs ValidationErrors) Error() string {
	msg := ""
	for i, e := range errs {
		msg += e.Error()
		if i < len(errs)-1 {
			msg += ", "
		}
	}
	return msg
}
