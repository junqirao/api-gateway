package ipgeo

import (
	"context"
	"net"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/oschwald/geoip2-golang"
)

var (
	db *geoip2.Reader
)

func Init(ctx context.Context) {
	path := g.Cfg().MustGet(ctx, "program.extra.ipgeo.database", "").String()
	if path == "" {
		g.Log().Infof(ctx, "ipgeo database not set")
		return
	}
	var err error
	db, err = geoip2.Open(path)
	if err != nil {
		g.Log().Errorf(ctx, "open ipgeo database failed: %v", err)
		return
	}
}

func country(addr string) string {
	if db == nil {
		return ""
	}
	record, err := db.Country(net.ParseIP(addr))
	if err != nil {
		return ""
	}
	return record.Country.IsoCode
}

func city(addr string) map[string]string {
	if db == nil {
		return map[string]string{}
	}
	record, err := db.City(net.ParseIP(addr))
	if err != nil {
		return map[string]string{}
	}
	return record.City.Names
}

func cityEN(addr string) string {
	return city(addr)["en"]
}
