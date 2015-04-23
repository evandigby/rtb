package redis

import (
	"math/rand"
	"os"
	"testing"
	"time"
)

var testDataAccess NoDbDataAccess

// From: http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func AssertAllKeysGone(t *testing.T, da NoDbDataAccess) {
	keys := testDataAccess.GetKeys()
	if len(testDataAccess.GetKeys()) > 0 {
		t.Fatalf("Test did not clean itself up.", keys)
	}
}

func DeleteAllKeys(da NoDbDataAccess) {
	testDataAccess.DeleteKeys(testDataAccess.GetKeys())
}

func TestMain(m *testing.M) {
	rand.Seed(time.Now().UTC().UnixNano())

	testDataAccess = NewRedisDataAccess("tcp", "localhost:6379", randSeq(10), 1)

	result := m.Run()

	DeleteAllKeys(testDataAccess)

	os.Exit(result)
}
