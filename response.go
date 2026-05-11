/**
 * @package go-response (2026)
 * @author Emmanuel Analike <emmanuel@analike.dev>
 * @created Feb 18, 2026; 9:20 PM
 */

package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.analike.dev/request"
)

var (
	Mimetypes = mimetypes{
		JSON:       "application/json",
		XML:        "application/xml",
		HTML:       "text/html",
		CSS:        "text/css",
		JAVASCRIPT: "application/javascript",
		FILE:       "application/octet-stream",
	}
	Seconds = seconds{
		SixHours: 21600,
		HalfDay:  43200,
		OneDay:   86400,
		OneWeek:  604800,
		OneMonth: 2592000,
	}
	Status = status{
		Ok:                    200,
		NoContent:             204,
		BadRequest:            400,
		Unauthorized:          401,
		Forbidden:             403,
		NotFound:              404,
		MethodNotAllowed:      405,
		Conflict:              409,
		ServerError:           500,
		ServerBadGateway:      502,
		ServiceGatewayTimeout: 504,
	}
	Redirect = redirect{
		Permanent: 301,
		Temporary: 302,
	}
)

type body struct {
	Status  int            `json:"status"`
	Meta    map[string]any `json:"meta,omitempty"`
	Message string         `json:"message,omitempty"`
	Data    any            `json:"data,omitempty"`
}

type File struct {
	Name         string
	PathRelative string
	PathAbsolute string
}

type header struct {
	name  string
	value []string
}

type Response struct {
	status  int
	code    int
	data    any
	body    *body
	message string
	noBody  bool

	cacheAge       *int
	cacheImmutable bool
	noCache        bool

	html *string
	meta map[string]any

	headers        []header
	contentType    string
	acceptRanges   bool
	exposeDuration bool

	file *File

	reqHeaders *http.Header
	reqMethod  string
	reqTime    *time.Time
	reqUri     string

	client any

	reqClient *request.Request
}

func (r *Response) SetStatus(status int) *Response {
	r.status = status
	return r
}

func (r *Response) GetStatus() int {
	return r.status
}

func (r *Response) GetHttpCode() int {
	return r.code
}

func (r *Response) SetHttpCode(code int) *Response {
	r.code = code
	return r
}

func (r *Response) SetStatusWithHttp(status int) *Response {
	return r.SetStatus(status).SetHttpCode(status)
}

func (r *Response) SetExposeDuration(s bool) *Response {
	r.exposeDuration = s
	return r
}

func (r *Response) IsExposeDuration() bool {
	return r.exposeDuration
}

func (r *Response) SetMessage(m string) *Response {
	r.message = m
	return r
}

func (r *Response) GetMessage() string {
	return r.message
}

func (r *Response) SetData(data any) *Response {
	r.data = data
	return r
}

func (r *Response) GetData() any {
	return r.data
}

func (r *Response) SetHtml(html *string) *Response {
	r.html = html
	return r
}

func (r *Response) GetHtml() *string {
	return r.html
}

func (r *Response) SetMeta(key string, value any) *Response {
	if r.meta == nil {
		r.meta = make(map[string]any)
	}
	r.meta[key] = value
	return r
}

func (r *Response) GetMeta(key string) any {
	return r.meta[key]
}

func (r *Response) AddHeader(name string, value ...string) *Response {
	r.headers = append(r.headers, header{name, value})
	return r
}

func (r *Response) SetNoBody(n bool) *Response {
	r.noBody = n
	return r
}

func (r *Response) SetNoCache(s bool) *Response {
	r.noCache = s
	return r
}

func (r *Response) IsNoCache() bool {
	return r.noCache
}

func (r *Response) SetCacheAge(age int) *Response {
	r.cacheAge = &age
	return r
}

func (r *Response) SetCacheImmutable(s bool) *Response {
	r.cacheImmutable = s
	return r
}

func (r *Response) IsCacheImmutable() bool {
	return r.cacheImmutable
}

func (r *Response) SetContentType(t string) *Response {
	r.contentType = t
	return r
}

func (r *Response) SetRedirectParams(url string, rType int) *Response {
	return r.AddHeader("Location", url).SetHttpCode(rType)
}

func (r *Response) SetFile(file *File) *Response {
	r.file = file
	if file != nil {
		r.SetContentType(Mimetypes.FILE)
	}
	return r
}

func (r *Response) SetAcceptRanges(s bool) *Response {
	r.acceptRanges = s
	return r
}

func (r *Response) IsNoContent() bool {
	return r.status == Status.NoContent
}

func (r *Response) buildBody() *Response {
	if r.body != nil {
		return r
	}
	body := body{
		Status:  r.status,
		Meta:    nil,
		Message: "",
		Data:    nil,
	}
	if r.meta != nil {
		body.Meta = r.meta
	}
	if r.message != "" {
		body.Message = r.message
	}
	if r.data != nil {
		body.Data = r.data
	}
	r.body = &body

	return r
}

func CreateNewGin(c *gin.Context, req *request.Request) Response {
	now := time.Now()
	return Response{
		status:         Status.Ok,
		code:           Status.Ok,
		cacheImmutable: false,
		noBody:         false,
		noCache:        true,
		contentType:    Mimetypes.JSON,
		acceptRanges:   false,
		exposeDuration: false,
		client:         c,

		reqHeaders: req.Headers,
		reqClient:  req,
		reqMethod:  req.Method,
		reqTime:    &now,
	}
}
