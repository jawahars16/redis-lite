package resp

import (
	"strconv"

	"golang.org/x/exp/slog"
)

type ArrayItem struct {
	Value    any
	DataType DataType
}

func Serialize(dataType DataType, data any) ([]byte, error) {
	switch {
	case dataType == SimpleStrings:
		return serializeToSimpleString(data.(string))
	case dataType == SimpleErrors:
		return serializeToSimpleError(data.(string))
	case dataType == Integers:
		return serializeToInteger(data.(int))
	case dataType == BulkStrings:
		return serializeToBulkStrings(data)
	case dataType == Arrays:
		return serializeToArrays(data)
	}
	return []byte{}, nil
}

func serializeToSimpleString(input string) ([]byte, error) {
	result := []byte{}
	result = append(result, SimpleStringPrefix)
	result = append(result, []byte(input)...)
	result = append(result, CR, LF)
	return result, nil
}

func serializeToSimpleError(message string) ([]byte, error) {
	result := []byte{}
	result = append(result, SimpleErrorPrefix)
	result = append(result, []byte(message)...)
	result = append(result, CR, LF)
	return result, nil
}

func serializeToInteger(number int) ([]byte, error) {
	result := []byte{}
	numberStr := strconv.Itoa(int(number))
	result = append(result, IntegerPrefix)
	result = append(result, []byte(numberStr)...)
	result = append(result, CR, LF)
	return result, nil
}

func serializeToBulkStrings(input any) ([]byte, error) {
	result := []byte{}
	result = append(result, BulkStringPrefix)
	if input == nil {
		result = append(result, []byte("-1")...)
		result = append(result, CR, LF)
		return result, nil
	}
	inputStr := input.(string)
	length := strconv.Itoa(len(inputStr))
	result = append(result, []byte(length)...)
	result = append(result, CR, LF)
	result = append(result, []byte(inputStr)...)
	result = append(result, CR, LF)
	return result, nil
}

func serializeToArrays(input any) ([]byte, error) {
	result := []byte{}
	result = append(result, ArrayPrefix)
	if input == nil {
		result = append(result, []byte("-1")...)
		result = append(result, CR, LF)
		return result, nil
	}

	inputArr := input.([]ArrayItem)
	length := strconv.Itoa(len(inputArr))
	result = append(result, []byte(length)...)
	result = append(result, CR, LF)

	for _, item := range inputArr {
		var (
			bytes []byte
			err   error
		)
		switch item.DataType {
		case SimpleStrings:
			bytes, err = serializeToSimpleString(item.Value.(string))
			if err != nil {
				slog.Error(err.Error())
			}
		case BulkStrings:
			bytes, err = serializeToBulkStrings(item.Value.(string))
			if err != nil {
				slog.Error(err.Error())
			}
		case Integers:
			bytes, err = serializeToInteger(item.Value.(int))
			if err != nil {
				slog.Error(err.Error())
			}
		case SimpleErrors:
			bytes, err = serializeToSimpleError(item.Value.(string))
			if err != nil {
				slog.Error(err.Error())
			}
		}
		result = append(result, bytes...)
	}

	if len(inputArr) == 0 {
		result = append(result, CR, LF)
	}

	return result, nil
}
