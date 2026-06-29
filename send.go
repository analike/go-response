/**
 * @package go-response (2026)
 * @author Emmanuel Analike <emmanuel@analike.dev>
 * @created Feb 19, 2026; 1:04 PM
 */

package response

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

func (r *Response) Send() error {
	isNoContent := r.IsNoContent()
	isNoBody := r.noBody
	hasBody := r.body != nil
	hasFile := r.file != nil

	var payload []byte
	var err error
	if !isNoContent && !hasBody && !isNoBody {
		r.buildBody()
	}
	if !isNoContent && !isNoBody && !hasFile {
		if r.html != nil {
			payload = []byte(*r.html)
		} else if r.body != nil {
			payload, err = json.Marshal(*r.body)
		}
	}
	if err == nil {
		err = r.setFileHeaders()
	}
	r.setHeaders().sendHeaders()
	r.doSend(&payload)
	return err
}

func (r *Response) SendError(code int, message string) error {
	return r.
		SetHttpCode(code).
		SetStatus(code).
		SetMessage(message).
		Send()
}

func (r *Response) SendErrorFormatted(code int, key string, setHttp bool) error {
	req := r.reqClient
	uri := req.Uri
	method := req.Method
	var message string
	switch code {
	case Status.NotFound:
		message = "Unrecognized endpoint or resource not found"
		break
	case Status.MethodNotAllowed:
		message = "Unrecognized method"
		break
	case Status.Unauthorized:
		message = "Authorization failed"
		break
	case Status.Forbidden:
		message = "Request forbidden"
		break
	case Status.BadRequest:
		message = "Bad Request. Invalid parameters provided"
		break
	default:
		code = Status.ServerError
		message = "The server encountered an unexpected error in trying to fulfill your request"
	}
	key = strings.TrimSpace(key)
	if key != "" {
		var keyMsg string
		re := regexp.MustCompile(`\s*\[\(`)
		hasBraces := re.MatchString(key)
		if hasBraces {
			keyMsg = key
		} else {
			keyMsg = fmt.Sprintf("[%s]", key)
		}
		message = fmt.Sprintf("%s %s", message, keyMsg)
	}
	r.SetStatus(code).
		SetMessage(message).
		SetMeta("method", method).
		SetMeta("uri", uri)
	if setHttp {
		r.SetHttpCode(code)
	}

	return r.Send()
}

func (r *Response) doSend(payload *[]byte) {
	switch cl := r.client.(type) {
	case *gin.Context:
		cl.Data(r.code, r.contentType, *payload)
	case *http.ResponseWriter:
		c := *cl
		c.WriteHeader(r.code)
		_, err := c.Write(*payload)
		if err != nil {
			log.Printf("[analike/response] error writing response: %#v", err)
		}
	}
}
