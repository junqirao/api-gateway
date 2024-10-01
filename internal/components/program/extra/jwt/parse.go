package jwt

import (
	"context"
	"errors"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/golang-jwt/jwt/v5"
)

const (
	HeaderKeyID = "kid"
)

var (
	key         string
	keysMapping map[string]string // kid:key
	keyFunc     jwt.Keyfunc       = func(token *jwt.Token) (interface{}, error) {
		if kid, ok := token.Header[HeaderKeyID].(string); ok && kid != "" {
			return []byte(keysMapping[kid]), nil
		}
		return []byte(key), nil
	}
	emptyClaims = map[string]interface{}{}
)

func Init(ctx context.Context) {
	key = g.Cfg().MustGet(ctx, "program.extra.jwt.key", "").String()
	keysMapping = g.Cfg().MustGet(ctx, "program.extra.jwt.keys_mapping", map[string]string{}).MapStrStr()
}

func parseTokenMap(s string) (map[string]interface{}, error) {
	s = strings.TrimPrefix(s, "Bearer ")
	claims, err := parseToken(s, keyFunc)
	if err != nil {
		return emptyClaims, err
	}
	return gconv.Map(claims), nil
}

func parseToken(s string, keyFunc jwt.Keyfunc) (claims *jwt.MapClaims, err error) {
	claims = new(jwt.MapClaims)
	token, err := jwt.ParseWithClaims(s, claims, keyFunc)
	if err != nil {
		return
	}
	if !token.Valid {
		err = errors.New("invalid token")
		return
	}
	return
}
