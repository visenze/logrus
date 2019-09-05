package logrus

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

// OSSFormatter implementing the open source logrus.Formatter
// Formats messages in Visenze's logging standards
// Reference document: https://wiki.visenze.com/display/VE/Logging+System+-+Conventions
// This implements Pattern 2, with the key value pattern.
type OSSFormatter struct {
	// projectName to be logged as
	projectName string
	// serviceName that this Formatter is for
	// This is the only required field for our Visenze logging configuration
	componentName string
}

// NewOSSFormatter initializes a new formatter object
func NewOSSFormatter(projectName, componentName string, opts ...Option) (*OSSFormatter, error) {
	if projectName == "" {
		return nil, errors.New("A project name must be specified")
	}
	if componentName == "" {
		return nil, errors.New("A component name must be specified")
	}
	return &OSSFormatter{
		projectName:   projectName,
		componentName: componentName,
	}, nil
}

// Format the message
func (f *OSSFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	b.WriteString(f.projectName)
	b.WriteByte(':')
	b.WriteString(f.componentName)
	b.WriteString(fmt.Sprintf(" time=\"%s\"", entry.Time.UTC().Format(TimestampFormat)))
	b.WriteString(fmt.Sprintf(" log_level=%s", strings.ToUpper(entry.Level.String())))
	b.WriteString(fmt.Sprintf(" msg='''%s'''", entry.Message))
	b.WriteByte(' ')
	for key, value := range entry.Data {
		appendKeyValue(b, replaceInvalidChars(key), value)
	}
	b.WriteByte('\n')
	return b.Bytes(), nil
}

func replaceInvalidChars(content string) string {
	return strings.Replace(content, ".", "_", -1)
}

var _ logrus.Formatter = (*OSSFormatter)(nil)

// Option for the Visenze formatter
type Option func(*Formatter) error

// TimestampFormat for Visenze logging
const TimestampFormat = "2006-01-02T15:04:05.000Z"
