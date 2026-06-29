/**
 * @package go-response (2026)
 * @author Emmanuel Analike <emmanuel@analike.dev>
 * @created Feb 19, 2026; 10:24 AM
 */

package response

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func (r *Response) sendHeaders() *Response {
	switch cl := r.client.(type) {
	case *gin.Context:
		for _, h := range r.headers {
			for _, v := range h.value {
				cl.Header(h.name, v)
			}
		}
	case *http.ResponseWriter:
		for _, h := range r.headers {
			for _, v := range h.value {
				(*cl).Header().Set(h.name, v)
			}
		}
	}
	return r
}

func (r *Response) setHeaders() *Response {
	r.setCorsHeaders()

	if r.IsNoCache() {
		r.setNoCacheHeaders()
	} else {
		imt := ""
		if r.IsCacheImmutable() {
			imt = ", immutable"
		}
		r.AddHeader("Cache-Control", fmt.Sprintf("max-age=%d%s", r.cacheAge, imt))
	}

	if r.IsExposeDuration() && r.reqTime != nil {
		duration := time.Now().UnixMilli() - (*r.reqTime).UnixMilli()
		took := fmt.Sprintf("%dms", duration)
		if duration > 999 {
			durationSec := float64(duration) / float64(1000)
			took = fmt.Sprintf("%.2fs", float64(durationSec)/float64(1000))
		}
		r.AddHeader("X-Took", took)
	}

	return r
}

func (r *Response) setFileHeaders() error {
	if r.file != nil {
		f := *r.file
		_, err := os.Stat(f.PathAbsolute)
		if f.PathAbsolute != "" && err == nil {
			if r.contentType != "" {
				r.AddHeader("Content-Type", Mimetypes.FILE)
			}
			if f.Name != "" {
				r.AddHeader("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", f.Name))
			}
			if r.acceptRanges {
				r.AddHeader("Accept-Ranges", "bytes")
			}
			r.AddHeader("X-Accel-Redirect", fmt.Sprintf("/%s", f.PathRelative))
		} else if err != nil {
			return err
		}
	}
	return nil
}

func removeDuplicates(list []string) []string {
	seen := map[string]bool{}
	out := []string{}

	for _, v := range list {
		if !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	return out
}

func (r *Response) setNoCacheHeaders() *Response {
	r.
		AddHeader(
			"Cache-Control",
			"no-store, no-cache, must-revalidate, max-age=0",
			"post-check=0, pre-check=0").
		AddHeader("Pragma", "no-cache")
	return r
}

func (r *Response) setCorsHeaders() *Response {
	hd := r.reqClient.Headers
	origin := hd.Get("Origin")
	if origin != "" {
		r.
			AddHeader("Access-Control-Allow-Origin", origin).
			AddHeader("Access-Control-Allow-Credentials", "true").
			AddHeader("Access-Control-Max-Age", "86400")
	}
	reqMethod := hd.Get("Access-Control-Request-Method")
	reqHeaders := hd.Get("Access-Control-Request-Headers")

	if reqMethod != "" {
		methods := []string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
		pieces := strings.Split(reqMethod, ", ")
		var upperPieces []string
		for _, one := range pieces {
			upperPieces = append(upperPieces, strings.ToUpper(one))
		}
		reqMethods := append(methods, upperPieces...)
		r.AddHeader("Access-Control-Allow-Methods", strings.Join(removeDuplicates(reqMethods), ", "))
	}

	if reqHeaders != "" {
		r.AddHeader("Access-Control-Allow-Headers", reqHeaders)
	}

	return r
}
