package jwt

import (
	"testing"

	"github.com/gogf/gf/v2/util/gconv"
	"github.com/golang-jwt/jwt/v5"
)

var (
	testKey = []byte("test_key_123456")
)

func generateTestToken() (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.MapClaims{
		"foo": "bar",
		"baz": "qux",
	}).SignedString(testKey)
}

func TestParseToken(t *testing.T) {
	token, err := generateTestToken()
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("token: %s", token)
	claims, err := parseToken(token, func(token *jwt.Token) (interface{}, error) {
		return testKey, nil
	})
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("claims: %v", claims)
	t.Logf("converted claims: %v", gconv.Map(claims))
}
