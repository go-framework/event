package inapp

import (
	"strings"
)

// Errors is error interface array, and impl error interface.
type Errors []error

func (list Errors) Error() string {
	buf := strings.Builder{}
	buf.WriteByte('[')
	for idx, item := range list {
		buf.WriteString(item.Error())
		if idx < len(list)-1 {
			buf.WriteByte(',')
		}
	}
	buf.WriteByte(']')
	return buf.String()
}

func (list Errors) Nil() error {
	if len(list) == 0 {
		return nil
	}
	return list
}
