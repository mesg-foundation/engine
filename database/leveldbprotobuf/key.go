package leveldbprotobuf

import "strings"

const collectionSeparatorKey string = "_|EOP|_"

func collectionKey(collection string) []byte {
	return []byte(collection + collectionSeparatorKey)
}

func makeKey(collection string, key string) []byte {
	return append(collectionKey(collection), []byte(key)...)
}

func splitKey(index []byte) (collection string, key string) {
	strings := strings.Split(string(index), collectionSeparatorKey)
	if len(strings) == 2 {
		collection = strings[0]
		key = strings[1]
	}
	return
}
