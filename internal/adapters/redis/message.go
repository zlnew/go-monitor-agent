package redis

import (
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func ParseStreamMessages[T any](msgs []redis.XMessage) ([]T, []string, error) {
	var (
		items []T
		ids   []string
	)

	for _, msg := range msgs {
		raw, ok := msg.Values["data"]
		if !ok {
			return nil, nil, fmt.Errorf("missing data field in stream message %s", msg.ID)
		}

		var bytes []byte

		switch v := raw.(type) {
		case string:
			bytes = []byte(v)
		case []byte:
			bytes = v
		default:
			return nil, nil, fmt.Errorf(
				"unexpected data type %T in stream message %s",
				raw,
				msg.ID,
			)
		}

		var item T
		if err := json.Unmarshal(bytes, &item); err != nil {
			return nil, nil, fmt.Errorf(
				"failed to unmarshal stream message %s: %w",
				msg.ID,
				err,
			)
		}

		items = append(items, item)
		ids = append(ids, msg.ID)
	}

	return items, ids, nil
}
