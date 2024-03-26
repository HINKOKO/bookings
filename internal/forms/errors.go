package forms

// 'errors' is a map of string with content being a slice of string
type errors map[string][]string

// Add - adds an error message for a given field
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Get - retunrs the first error message
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}
