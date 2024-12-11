package goframe

import (
	"github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/gogf/gf/v2/net/ghttp"
	"net/http"
)

// SentinelMiddleware 返回新的 ghttp.HandlerFunc
func SentinelMiddleware(opts ...Option) ghttp.HandlerFunc {
	options := evaluateOptions(opts)
	return func(r *ghttp.Request) {
		resourceName := r.Method + ":" + r.URL.Path

		if options.resourceExtract != nil {
			resourceName = options.resourceExtract(r)
		}

		entry, err := api.Entry(
			resourceName,
			api.WithResourceType(base.ResTypeWeb),
			api.WithTrafficType(base.Inbound),
		)

		if err != nil {
			if options.blockFallback != nil {
				options.blockFallback(r)
			} else {
				r.Response.WriteHeader(http.StatusTooManyRequests)
			}
			return
		}

		defer entry.Exit()
		r.Middleware.Next()
	}
}
