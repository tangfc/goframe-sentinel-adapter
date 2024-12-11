package goframe

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/stretchr/testify/assert"
)

func initSentinel(t *testing.T) {
	err := sentinel.InitDefault()
	if err != nil {
		t.Fatalf("Unexpected error: %+v", err)
	}

	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "GET:/ping",
			Threshold:              1.0,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "/api/users/:id",
			Threshold:              0.0,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			StatIntervalInMs:       1000,
		},
	})
	if err != nil {
		t.Fatalf("Unexpected error: %+v", err)
		return
	}
}

func TestSentinelMiddleware(t *testing.T) {
	type args struct {
		opts    []Option
		method  string
		path    string
		reqPath string
		handler func(r *ghttp.Request)
		body    io.Reader
	}
	type want struct {
		code int
	}
	var (
		tests = []struct {
			name string
			args args
			want want
		}{
			{
				name: "default get",
				args: args{
					opts:    []Option{},
					method:  http.MethodGet,
					path:    "/ping",
					reqPath: "/ping",
					handler: func(r *ghttp.Request) {
						r.Response.Writeln("ping")
					},
					body: nil,
				},
				want: want{
					code: http.StatusOK,
				},
			},
			{
				name: "customize resource extract",
				args: args{
					opts: []Option{
						WithResourceExtractor(func(r *ghttp.Request) string {
							return r.URL.Path
						}),
					},
					method:  http.MethodPost,
					path:    "/api/users/:id",
					reqPath: "/api/users/123",
					handler: func(r *ghttp.Request) {
						r.Response.Writeln("ping")
					},
					body: nil,
				},
				want: want{
					code: http.StatusTooManyRequests,
				},
			},
			{
				name: "customize block fallback",
				args: args{
					opts: []Option{
						WithBlockFallback(func(r *ghttp.Request) {
							r.Response.WriteHeader(http.StatusBadRequest)
							r.Response.Writeln("block")
						}),
					},
					method:  http.MethodGet,
					path:    "/ping",
					reqPath: "/ping",
					handler: func(r *ghttp.Request) {
						r.Response.Writeln("ping")
					},
					body: nil,
				},
				want: want{
					code: http.StatusBadRequest,
				},
			},
		}
	)
	initSentinel(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := g.Server()
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(SentinelMiddleware(tt.args.opts...))
				group.ALL(tt.args.path, tt.args.handler)
			})

			r := httptest.NewRequest(tt.args.method, tt.args.reqPath, tt.args.body)
			w := httptest.NewRecorder()

			s.ServeHTTP(w, r)

			assert.Equal(t, tt.want.code, w.Code)
		})
	}
}
