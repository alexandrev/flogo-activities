package parsecsv

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// ActivityLog is the default logger for the Log Activity
var activityLog = logger.GetLogger("activity-flogo-parsecsv")

const (
	ivFieldNames = "fieldNames"
	ivCSV        = "csv"

	ovOutput = "output"
)

func init() {
	activityLog.SetLogLevel(logger.InfoLevel)
}

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
	txt := ctx.GetInput(ivCSV).(string)

	r := csv.NewReader(strings.NewReader(txt))
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

		field := make(map[string]string)

		for i := 0; i < len(record); i++ {
			field[fieldNames[i].(string)] = record[i]
		}

		obj = append(obj, field)
	}

	activityLog.Debugf("Parsed Object from CSV: %s", obj)
	ctx.SetOutput(ovOutput, obj)

	return true, nil
}
