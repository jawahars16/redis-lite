package core

import "github.com/jawahars16/redis-lite/resp"

func HandlePing(args ...any) ([]byte, error) {
	return resp.Serialize(resp.SimpleStrings, "PONG")
}

func HandleCommand(args ...any) ([]byte, error) {
	if args[0] == "DOCS" {
		return resp.Serialize(resp.Arrays, []resp.ArrayItem{
			{
				DataType: resp.SimpleStrings,
				Value:    "PING",
			},
		})
	}
	return []byte{}, nil
}
