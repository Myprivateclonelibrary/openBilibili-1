package ads

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	xsql "go-common/library/database/sql"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
)

func TestAdsVideoAds(t *testing.T) {
	convey.Convey("VideoAds", t, func(ctx convey.C) {
		ctx.Convey("When everything is correct", func(ctx convey.C) {
			res, err := d.VideoAds(context.Background())
			ctx.Convey("Error should be nil, res should not be empty(No Data)", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.SkipSo(res, convey.ShouldNotBeEmpty)
			})
		})
		ctx.Convey("When db.Query gets error", func(ctx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.db), "Query", func(_ *xsql.DB, _ context.Context, _ string, _ ...interface{}) (*xsql.Rows, error) {
				return nil, fmt.Errorf("db.Query error")
			})
			defer guard.Unpatch()
			_, err := d.VideoAds(context.Background())
			ctx.Convey("Error should not be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}
