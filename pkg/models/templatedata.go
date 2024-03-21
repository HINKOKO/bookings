package models

// TemplateData:  holds data sent from handlers to templates
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{} // interface can hold value of any type
	CSRFToken string
	Flash     string
	Warning   string
	Error     string
}
