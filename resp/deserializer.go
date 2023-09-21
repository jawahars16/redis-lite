package resp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strconv"
)

type DataType string

func Deserialize(reader io.Reader) (DataType, any, error) {
	r := bufio.NewReader(reader)
	line, err := readLine(r)
	if err != nil && err != io.EOF {
		return "", nil, err
	}
	dataType, data, err := extract(line, r)
	if err != nil && err != io.EOF {
		slog.Error(err.Error())
		return "", nil, err
	}
	return dataType, data, nil
}

func extract(line []byte, r *bufio.Reader) (DataType, any, error) {
	switch {
	case line[0] == SimpleStringPrefix:
		return simpleString(line)
	case line[0] == SimpleErrorPrefix:
		return simpleError(line)
	case line[0] == IntegerPrefix:
		return integer(line)
	case line[0] == BulkStringPrefix:
		return bulkStrings(line, r)
	case line[0] == ArrayPrefix:
		return arrays(line, r)
	default:
		if len(line) > 0 {
			inlineCommand := []byte{'+'}
			inlineCommand = append(inlineCommand, line...)
			return simpleString(inlineCommand)
		}
		return "", nil, ErrUnrecognizedType
	}
}

func simpleString(data []byte) (DataType, string, error) {
	if len(data) < 1 {
		return "", "", errors.New("no data available")
	}
	if len(data) < 2 {
		return SimpleStrings, "", nil
	}
	return SimpleStrings, string(data[1:]), nil
}

func simpleError(data []byte) (DataType, string, error) {
	if len(data) < 2 {
		return "", "", errors.New("no data available")
	}
	return SimpleErrors, string(data[1:]), nil
}

func integer(data []byte) (DataType, int, error) {
	n, err := strconv.Atoi(string(data[1:]))
	if err != nil {
		return "", 0, err
	}
	return Integers, n, nil
}

func bulkStrings(data []byte, reader *bufio.Reader) (DataType, any, error) {
	length, err := strconv.Atoi(string(data[1:]))
	if err != nil {
		return "", nil, err
	}

	if length < 0 {
		return BulkStrings, nil, nil
	}

	result := make([]byte, length)
	_, err = reader.Read(result)
	if err != nil && err != io.EOF {
		slog.Error(err.Error())
		return "", nil, err
	}

	// Bulk string ends with CRLF
	// So discard the last 2 bytes
	reader.Discard(2)

	return BulkStrings, string(result), nil
}

func arrays(data []byte, reader *bufio.Reader) (DataType, []ArrayItem, error) {
	length, err := strconv.Atoi(string(data[1:]))
	if err != nil {
		return "", nil, err
	}

	if length < 0 {
		return Arrays, nil, nil
	}

	resultArray := []ArrayItem{}
	if length == 0 {
		return Arrays, resultArray, nil
	}

	for {
		line, err := readLine(reader)
		if err != nil && err != io.EOF {
			slog.Error(err.Error())
			break
		}
		dataType, data, err := extract(line, reader)
		if err != nil && err != io.EOF {
			slog.Error(err.Error())
			break
		}

		resultArray = append(resultArray, ArrayItem{
			Value:    data,
			DataType: dataType,
		})
		if len(resultArray) >= length {
			return Arrays, resultArray, nil
		}
	}

	return "", nil, errors.New("cannot form array")
}

func readLine(reader *bufio.Reader) ([]byte, error) {
	line, err := reader.ReadBytes('\n')
	if len(line) < 1 {
		return nil, fmt.Errorf("empty line cannot process")
	}
	// skip the LF at end
	line = line[:len(line)-1]
	// skip if there is CR
	if line[len(line)-1] == '\r' {
		line = line[:len(line)-1]
	}
	if err != nil {
		slog.Error(err.Error())
		return line, err
	}
	return line, nil
}
