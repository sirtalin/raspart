package model

import "strings"

type StringList struct {
	List []string `csv:"-"`
}

// MarshalCSV converts the internal date as CSV string
func (stringList *StringList) MarshalCSV() (string, error) {
	return strings.Join(stringList.List[:], ";"), nil
}

// UnmarshalCSV converts CSV string as the internal date
func (stringList *StringList) UnmarshalCSV(csv string) (err error) {
	stringList.List = strings.Split(csv, ";")
	return nil
}
