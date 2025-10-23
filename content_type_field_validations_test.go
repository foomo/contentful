package contentful

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFieldValidationLink(t *testing.T) {
	var err error

	validation := &FieldValidationLink{
		LinkContentType: []string{"test", "test2"},
	}

	data, err := json.Marshal(validation)
	require.NoError(t, err)
	assert.JSONEq(t, "{\"linkContentType\":[\"test\",\"test2\"]}", string(data))
}

func TestFieldValidationUnique(t *testing.T) {
	var err error

	validation := &FieldValidationUnique{
		Unique: false,
	}

	data, err := json.Marshal(validation)
	require.NoError(t, err)
	assert.Equal(t, "{\"unique\":false}", string(data))
}

func TestFieldValidationPredefinedValues(t *testing.T) {
	var err error

	validation := &FieldValidationPredefinedValues{
		In:           []interface{}{5, 10, "string", 6.4},
		ErrorMessage: "error message",
	}

	data, err := json.Marshal(validation)
	require.NoError(t, err)
	assert.JSONEq(t, "{\"in\":[5,10,\"string\",6.4],\"message\":\"error message\"}", string(data))
}

func TestFieldValidationRange(t *testing.T) {
	var err error

	// between
	validation := &FieldValidationRange{
		Range: &MinMax{
			Min: 60,
			Max: 100,
		},
		ErrorMessage: "error message",
	}
	data, err := json.Marshal(validation)
	require.NoError(t, err)
	assert.JSONEq(t, "{\"range\":{\"min\":60,\"max\":100},\"message\":\"error message\"}", string(data))

	var validationCheck FieldValidationRange
	err = json.NewDecoder(bytes.NewReader(data)).Decode(&validationCheck)
	require.NoError(t, err)
	assert.InDelta(t, float64(60), validationCheck.Range.Min, 0)
	assert.InDelta(t, float64(100), validationCheck.Range.Max, 0)
	assert.Equal(t, "error message", validationCheck.ErrorMessage)

	// greater than equal to
	validation = &FieldValidationRange{
		Range: &MinMax{
			Min: 10,
		},
		ErrorMessage: "error message",
	}
	data, err = json.Marshal(validation)
	require.NoError(t, err)
	assert.JSONEq(t, "{\"range\":{\"min\":10},\"message\":\"error message\"}", string(data))
	validationCheck = FieldValidationRange{}
	err = json.NewDecoder(bytes.NewReader(data)).Decode(&validationCheck)
	require.NoError(t, err)
	assert.InDelta(t, float64(10), validationCheck.Range.Min, 0)
	assert.InDelta(t, float64(0), validationCheck.Range.Max, 0)
	assert.Equal(t, "error message", validationCheck.ErrorMessage)

	// less than equal to
	validation = &FieldValidationRange{
		Range: &MinMax{
			Max: 90,
		},
		ErrorMessage: "error message",
	}
	data, err = json.Marshal(validation)
	require.NoError(t, err)
	assert.JSONEq(t, "{\"range\":{\"max\":90},\"message\":\"error message\"}", string(data))
	validationCheck = FieldValidationRange{}
	err = json.NewDecoder(bytes.NewReader(data)).Decode(&validationCheck)
	require.NoError(t, err)
	assert.InDelta(t, float64(90), validationCheck.Range.Max, 0)
	assert.InDelta(t, float64(0), validationCheck.Range.Min, 0)
	assert.Equal(t, "error message", validationCheck.ErrorMessage)
}

func TestFieldValidationSize(t *testing.T) {
	var err error

	// between
	validation := &FieldValidationSize{
		Size: &MinMax{
			Min: 4,
			Max: 6,
		},
		ErrorMessage: "error message",
	}
	data, err := json.Marshal(validation)
	require.NoError(t, err)
	assert.JSONEq(t, "{\"size\":{\"min\":4,\"max\":6},\"message\":\"error message\"}", string(data))

	var validationCheck FieldValidationSize
	err = json.NewDecoder(bytes.NewReader(data)).Decode(&validationCheck)
	require.NoError(t, err)
	assert.InDelta(t, float64(4), validationCheck.Size.Min, 0)
	assert.InDelta(t, float64(6), validationCheck.Size.Max, 0)
	assert.Equal(t, "error message", validationCheck.ErrorMessage)
}

func TestFieldValidationDate(t *testing.T) {
	var err error

	layout := "2006-01-02T03:04:05"
	minimum := time.Now()
	maximum := time.Now()

	minStr := minimum.Format(layout)
	maxStr := maximum.Format(layout)

	validation := &FieldValidationDate{
		Range: &DateMinMax{
			Min: minimum,
			Max: maximum,
		},
		ErrorMessage: "error message",
	}
	data, err := json.Marshal(validation)
	require.NoError(t, err)
	assert.Equal(t, "{\"dateRange\":{\"min\":\""+minStr+"\",\"max\":\""+maxStr+"\"},\"message\":\"error message\"}", string(data))

	var validationCheck FieldValidationDate
	err = json.NewDecoder(bytes.NewReader(data)).Decode(&validationCheck)
	require.NoError(t, err)
	assert.Equal(t, minStr, validationCheck.Range.Min.Format(layout))
	assert.Equal(t, maxStr, validationCheck.Range.Max.Format(layout))
	assert.Equal(t, "error message", validationCheck.ErrorMessage)
}
