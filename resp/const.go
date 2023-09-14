package resp

import "errors"

// DataTypes
const (
	SimpleStrings DataType = "SimpleStrings"
	SimpleErrors  DataType = "SimpleErrors"
	Integers      DataType = "Integers"
	BulkStrings   DataType = "BulkStrings"
	Arrays        DataType = "Arrays"
)

// Prefix
const (
	SimpleStringPrefix byte = '+'
	SimpleErrorPrefix  byte = '-'
	IntegerPrefix      byte = ':'
	BulkStringPrefix   byte = '$'
	ArrayPrefix        byte = '*'
)

const (
	CR byte = '\r'
	LF byte = '\n'
)

var ErrUnrecognizedType = errors.New("unrecognized type")
