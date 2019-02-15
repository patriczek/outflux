package ingestion

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/timescale/outflux/internal/idrf"
)

type IdrfConverter interface {
	ConvertValues(row idrf.Row) ([]interface{}, error)
}

func newIdrfConverter(dataSet *idrf.DataSetInfo) IdrfConverter {
	return &defaultIdrfConverter{dataSet}
}

type defaultIdrfConverter struct {
	dataSet *idrf.DataSetInfo
}

func (conv *defaultIdrfConverter) ConvertValues(row idrf.Row) ([]interface{}, error) {
	if len(row) != len(conv.dataSet.Columns) {
		return nil, fmt.Errorf(
			"could not convert extracted row, number of extracted values is %d, expected %d values",
			len(row), len(conv.dataSet.Columns))
	}

	converted := make([]interface{}, len(row))
	for i, item := range row {
		converted[i] = convertByType(item, conv.dataSet.Columns[i].DataType)
	}

	return converted, nil
}

func convertByType(rawValue interface{}, expected idrf.DataType) interface{} {
	if rawValue == nil {
		return nil
	}

	switch {
	case expected == idrf.IDRFInteger32:
		valAsInt64, _ := rawValue.(json.Number).Int64()
		return int32(valAsInt64)
	case expected == idrf.IDRFInteger64:
		valAsInt64, _ := rawValue.(json.Number).Int64()
		return valAsInt64
	case expected == idrf.IDRFDouble:
		valAsFloat64, _ := rawValue.(json.Number).Float64()
		return valAsFloat64
	case expected == idrf.IDRFSingle:
		valAsFloat64, _ := rawValue.(json.Number).Float64()
		return float32(valAsFloat64)
	case expected == idrf.IDRFTimestamptz || expected == idrf.IDRFTimestamp:
		ts, _ := time.Parse(time.RFC3339, rawValue.(string))
		return ts
	default:
		return rawValue
	}
}