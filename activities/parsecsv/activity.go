package parsecsv

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// ActivityLog is the default logger for the Log Activity
var activityLog = logger.GetLogger("activity-flogo-parsecsv")

const (
	ivFieldNames = "fieldNames"
	ivCSV        = "csv"
	ivFile       = "file"
	delimiter    = "delimiter"
	ovOutput     = "output"
)

// MyActivity is a stub for your Activity implementation
type ParseCSVActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &ParseCSVActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *ParseCSVActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *ParseCSVActivity) Eval(ctx activity.Context) (done bool, err error) {
	fieldNames := ctx.GetInput(ivFieldNames).([]interface{})

	var reader io.Reader


	if txt, ok := ctx.GetInput(ivCSV).(string); ok && len(txt) > 0 {
		reader = strings.NewReader(txt)
	} else if file, ok := ctx.GetInput(ivFile).(string); ok {
		osFile, err := os.Open(file)
		if err != nil {
			return false, fmt.Errorf("error opening the specified file: %v", err)
		}
		reader = bufio.NewReader(osFile)
	} else {
		return false, fmt.Errorf("either a filename or a string containing the CSV must be supplied")
	}

	dDelimiter, ok := ctx.GetInput(delimiter).(string)
	if ok == false || dDelimiter == "" {
		dDelimiter = ","
	}

	r := csv.NewReader(reader)
	rDelimiter := []rune(dDelimiter)
	r.Comma = rDelimiter[0]

	obj := make([]interface{}, 0)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			activityLog.Errorf("Failed to read csv string: %s", err)
			return false, err
		}

		if len(record) != len(fieldNames) {
			activityLog.Error("Mismatch between number of fields and field names specified")
			return false, fmt.Errorf("Fields supplied do not match total fields in csv. Expected %d but got %d", len(fieldNames), len(record))
		}

		field := make(map[string]interface{})

		for i := 0; i < len(record); i++ {
			field[fieldNames[i].(string)] = record[i]
		}

		obj = append(obj, field)
	}

	activityLog.Debugf("Parsed Object from CSV: %s", obj)
	ctx.SetOutput(ovOutput, obj)

	return true, nil
}
