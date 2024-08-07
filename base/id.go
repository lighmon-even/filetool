package base

import (
	"math/rand"
	"sync"
)

var (
	shellIDs = make(map[string]struct{})
	mu       sync.Mutex
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// 生成随机字符串
func randomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// 生成唯一ID
func generateID() string {
	mu.Lock()
	defer mu.Unlock()

	var _id string
	for {
		_id = randomString(6)
		if _, exists := shellIDs[_id]; !exists {
			shellIDs[_id] = struct{}{}
			break
		}
	}
	return _id
}
