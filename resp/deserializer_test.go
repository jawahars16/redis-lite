package resp_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jawahars16/redis-lite/resp"
)

type DeserializerTestCase struct {
	input            string
	expectedDataType resp.DataType
	expectedData     any
	expectedError    error
}

func Test_Deserialize(t *testing.T) {
	t.Run("Simple strings", testDeserialize(DeserializerTestCase{
		input:            "+Hello\r\n",
		expectedDataType: resp.SimpleStrings,
		expectedData:     "Hello",
	}))

	t.Run("Simple errors", testDeserialize(DeserializerTestCase{
		input:            "-Error\r\n",
		expectedDataType: resp.SimpleErrors,
		expectedData:     "Error",
	}))

	t.Run("Integers", testDeserialize(DeserializerTestCase{
		input:            ":1024\r\n",
		expectedDataType: resp.Integers,
		expectedData:     1024,
	}))

	t.Run("Bulk Strings", testDeserialize(DeserializerTestCase{
		input:            "$11\r\nHello World\r\n",
		expectedDataType: resp.BulkStrings,
		expectedData:     "Hello World",
	}))

	t.Run("Bulk String with new line", testDeserialize(DeserializerTestCase{
		input:            "$11\r\nHello\nWorld\r\n",
		expectedDataType: resp.BulkStrings,
		expectedData:     "Hello\nWorld",
	}))

	t.Run("Empty bulk string", testDeserialize(DeserializerTestCase{
		input:            "$0\r\n\r\n",
		expectedDataType: resp.BulkStrings,
		expectedData:     "",
	}))

	t.Run("Nil bulk string", testDeserialize(DeserializerTestCase{
		input:            "$-1\r\n",
		expectedDataType: resp.BulkStrings,
		expectedData:     nil,
	}))

	t.Run("Arrays", testDeserialize(DeserializerTestCase{
		input:            "*2\r\n$2\r\nHi\r\n$5\r\nHello\r\n",
		expectedDataType: resp.Arrays,
		expectedData: []resp.ArrayItem{
			{Value: "Hi", DataType: resp.BulkStrings},
			{Value: "Hello", DataType: resp.BulkStrings},
		},
	}))

	t.Run("Arrays with different types of data", testDeserialize(DeserializerTestCase{
		input:            "*2\r\n$2\r\nHi\r\n:20\r\n",
		expectedDataType: resp.Arrays,
		expectedData: []resp.ArrayItem{
			{Value: "Hi", DataType: resp.BulkStrings},
			{Value: 20, DataType: resp.Integers},
		},
	}))

	t.Run("Empty arrays", testDeserialize(DeserializerTestCase{
		input:            "*0\r\n\r\n",
		expectedDataType: resp.Arrays,
		expectedData:     []resp.ArrayItem{},
	}))

	// t.Run("Nil arrays", testDeserialize(DeserializerTestCase{
	// 	input:            "*-1\r\n",
	// 	expectedDataType: resp.Arrays,
	// 	expectedData:     nil,
	// }))

	t.Run("Invalid", testDeserialize(DeserializerTestCase{
		input:            "^Hello\r\n",
		expectedDataType: "",
		expectedData:     nil,
		expectedError:    resp.ErrUnrecognizedType,
	}))
}

func testDeserialize(testCase DeserializerTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		dataType, data, err := resp.Deserialize(strings.NewReader(testCase.input))
		if err != nil && testCase.expectedError == nil {
			t.Error(err)
		}
		if testCase.expectedError != nil {
			if err != testCase.expectedError {
				t.Error("Unexpected:", err)
			}
		}
		if diff := cmp.Diff(dataType, testCase.expectedDataType); diff != "" {
			t.Error(diff)
		}
		if diff := cmp.Diff(data, testCase.expectedData); diff != "" {
			t.Error(diff)
		}
	}
}
