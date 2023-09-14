package resp_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jawahars16/redis-lite/resp"
)

type SerializerTestCase struct {
	input         any
	inputDataType resp.DataType
	expectedData  string
	expectedError error
}

func Test_Serializer(t *testing.T) {
	t.Run("Simple string", testSerialize(SerializerTestCase{
		input:         "Hello",
		inputDataType: resp.SimpleStrings,
		expectedData:  "+Hello\r\n",
	}))

	t.Run("Simple error", testSerialize(SerializerTestCase{
		input:         "Error message",
		inputDataType: resp.SimpleErrors,
		expectedData:  "-Error message\r\n",
	}))

	t.Run("Integer", testSerialize(SerializerTestCase{
		input:         1024,
		inputDataType: resp.Integers,
		expectedData:  ":1024\r\n",
	}))

	t.Run("Bulk strings", testSerialize(SerializerTestCase{
		input:         "Hello World",
		inputDataType: resp.BulkStrings,
		expectedData:  "$11\r\nHello World\r\n",
	}))

	t.Run("Empty bulk string", testSerialize(SerializerTestCase{
		input:         "",
		inputDataType: resp.BulkStrings,
		expectedData:  "$0\r\n\r\n",
	}))

	t.Run("Nil bulk string", testSerialize(SerializerTestCase{
		input:         nil,
		inputDataType: resp.BulkStrings,
		expectedData:  "$-1\r\n",
	}))

	t.Run("Arrays", testSerialize(SerializerTestCase{
		input: []resp.ArrayItem{
			resp.ArrayItem{
				Value:    "Hi",
				DataType: resp.BulkStrings,
			},
			resp.ArrayItem{
				Value:    1024,
				DataType: resp.Integers,
			},
			resp.ArrayItem{
				Value:    "Hi",
				DataType: resp.SimpleStrings,
			},
		},
		inputDataType: resp.Arrays,
		expectedData:  "*3\r\n$2\r\nHi\r\n:1024\r\n+Hi\r\n",
	}))

	t.Run("Empty arrays", testSerialize(SerializerTestCase{
		input:         []resp.ArrayItem{},
		inputDataType: resp.Arrays,
		expectedData:  "*0\r\n\r\n",
	}))

	t.Run("Nil arrays", testSerialize(SerializerTestCase{
		input:         nil,
		inputDataType: resp.Arrays,
		expectedData:  "*-1\r\n",
	}))
}

func testSerialize(testCase SerializerTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		data, err := resp.Serialize(testCase.inputDataType, testCase.input)
		if err != nil && testCase.expectedError == nil {
			t.Error(err)
		}
		if testCase.expectedError != nil && testCase.expectedError != err {
			t.Error("Unexpected", err)
		}
		if diff := cmp.Diff(string(data), testCase.expectedData); diff != "" {
			t.Error(diff)
		}
	}
}
