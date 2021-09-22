package client

type MIMEType string

const (
	MIMETypeHTML    MIMEType = "text/html"
	MIMETypeCSV     MIMEType = "text/csv"
	MIMETypeJSON    MIMEType = "application/json"
	MIMETypeXML     MIMEType = "application/xml"
	MIMETypeUnknown MIMEType = ""
)

func (m MIMEType) String() string {
	return string(m)
}
