package models


type ProblemDetail struct {
 Type                   string            `json:"type"`
 Title                  string            `json:"title"`
 Status                 int               `json:"status"`
 Detail                 string            `json:"detail"`
 Instance               string            `json:"instance"`
 Timestamp              string            `json:"timestamp"`
 SchemaValidationErrors []ValidationError `json:"schemaValidationErrors"`
}


type ValidationError struct {
 Field          string `json:"field"`
 Message        string `json:"message"`
 DefaultMessage string `json:"defaultMessage"`
 RejectedValue  string `json:"rejectedValue"`
}


func (p ProblemDetail) FieldsInError() []string {
 fields := make([]string, 0, len(p.SchemaValidationErrors))
 for _, e := range p.SchemaValidationErrors {
  fields = append(fields, e.Field)
 }

 return fields
}