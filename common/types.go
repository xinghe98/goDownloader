package common

import (
	"net/http"

	"github.com/vbauerster/mpb"
)

type Tasks struct {
	Index int
	Req   *http.Request
	Size  int64
	Bar   *mpb.Bar
}

type Resluts struct {
	Err     error
	Success bool
}
