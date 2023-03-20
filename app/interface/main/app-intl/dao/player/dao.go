package player

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"go-common/app/interface/main/app-intl/conf"
	"go-common/app/interface/main/app-intl/model/player"
	httpx "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

// Dao is
type Dao struct {
	client *httpx.Client
}

// New elec dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		client: httpx.NewClient(c.HTTPClient),
	}
	return
}

// Playurl is
func (d *Dao) Playurl(c context.Context, mid, aid, cid, qn int64, npcybs, fnver, fnval, forceHost int, otype, mobiApp, buvid, fp, session, reqURL string) (playurl *player.Playurl, code int, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	params := url.Values{}
	params.Set("otype", otype)
	params.Set("buvid", buvid)
	params.Set("platform", mobiApp)
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("cid", strconv.FormatInt(cid, 10))
	params.Set("session", session)
	params.Set("force_host", strconv.Itoa(forceHost))
	if aid > 0 {
		params.Set("avid", strconv.FormatInt(aid, 10))
	}
	params.Set("fnver", strconv.Itoa(fnver))
	params.Set("fnval", strconv.Itoa(fnval))
	if qn != 0 {
		params.Set("qn", strconv.FormatInt(qn, 10))
	}
	if npcybs != 0 {
		params.Set("npcybs", strconv.Itoa(npcybs))
	}
	var res struct {
		Code int `json:"code"`
		*player.Playurl
	}
	var req *http.Request
	if req, err = d.client.NewRequest(http.MethodGet, reqURL, ip, params); err != nil {
		err = errors.Wrap(err, "d.client.NewRequest error")
		return
	}
	if fp != "" {
		req.Header.Set("X-BVC-FINGERPRINT", fp)
	}
	if err = d.client.Do(c, req, &res); err != nil {
		return
	}
	playurl = res.Playurl
	playurl.FormatDash()
	code = res.Code
	return
}
