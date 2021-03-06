/*
 * This is gen middleware for the web apps using the middlewares that supports http handleFunc
 */
package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"

	"github.com/dshechao/gen-api/gen"
	"github.com/dshechao/gen-api/gen/models"
)

/* 32 MB in memory max */
const MaxInMemoryMultipartSize = 32000000

var reqWriteExcludeHeaderDump = map[string]bool{
	"Host":              true, // not in Header map anyway
	"Content-Length":    true,
	"Transfer-Encoding": true,
	"Trailer":           true,
	"Accept-Encoding":   false,
	"Accept-Language":   false,
	"Cache-Control":     false,
	"Connection":        false,
	"Origin":            false,
	"User-Agent":        false,
}

type YaagHandler struct {
	nextHandler http.Handler
}

func Handle(nextHandler http.Handler) http.Handler {
	return &YaagHandler{nextHandler: nextHandler}
}

func (y *YaagHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !gen.IsOn() {
		y.nextHandler.ServeHTTP(w, r)
		return
	}
	writer := NewResponseRecorder(w)
	apiCall := models.ApiCall{}
	Before(&apiCall, r)
	y.nextHandler.ServeHTTP(writer, r)
	After(&apiCall, writer, r)
}

func HandleFunc(next func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !gen.IsOn() {
			next(w, r)
			return
		}
		apiCall := models.ApiCall{}
		writer := NewResponseRecorder(w)
		Before(&apiCall, r)
		next(writer, r)
		After(&apiCall, writer, r)
	})
}

func Before(apiCall *models.ApiCall, req *http.Request) {
	apiCall.RequestHeader = ReadHeaders(req, apiCall)
	apiCall.RequestUrlParams = ReadQueryParams(req)
	val, ok := apiCall.RequestHeader["Content-Type"]
	log.Println(val)
	if ok {
		ct := strings.TrimSpace(apiCall.RequestHeader["Content-Type"])
		switch ct {
		case "application/x-www-form-urlencoded":
			fallthrough
		case "application/json, application/x-www-form-urlencoded":
			log.Println("Reading form")
			apiCall.PostForm = ReadPostForm(req)
		case "application/json":
			log.Println("Reading body")
			r := *ReadBody(req)
			r1, r2 := HandlerParamComment(req, r)
			apiCall.RequestBody = r1
			apiCall.RequestBodyComment = r2
		default:
			if strings.Contains(ct, "multipart/form-data") {
				handleMultipart(apiCall, req)
			} else {
				log.Println("Reading body")
				r := *ReadBody(req)
				r1, r2 := HandlerParamComment(req, r)
				apiCall.RequestBody = r1
				apiCall.RequestBodyComment = r2
			}
		}

	}
}

func ReadQueryParams(req *http.Request) map[string]string {
	params := map[string]string{}
	u, err := url.Parse(req.RequestURI)
	if err != nil {
		return params
	}
	for k, v := range u.Query() {
		if len(v) < 1 {
			continue
		}
		// TODO: v is a list, and we should be showing a list of values
		// rather than assuming a single value always, gotta change this
		params[k] = v[0]
	}
	return params
}

func printMap(m map[string]string) {
	for key, value := range m {
		log.Println(key, "=>", value)
	}
}

func handleMultipart(apiCall *models.ApiCall, req *http.Request) {
	apiCall.RequestHeader["Content-Type"] = "multipart/form-data"
	req.ParseMultipartForm(MaxInMemoryMultipartSize)
	apiCall.PostForm = ReadMultiPostForm(req.MultipartForm)
}

func ReadMultiPostForm(mpForm *multipart.Form) map[string]string {
	postForm := map[string]string{}
	for key, val := range mpForm.Value {
		postForm[key] = val[0]
	}
	return postForm
}

func ReadPostForm(req *http.Request) map[string]string {
	postForm := map[string]string{}
	for _, param := range strings.Split(*ReadBody(req), "&") {
		value := strings.Split(param, "=")
		postForm[value[0]] = value[1]
	}
	return postForm
}

func ReadHeaders(req *http.Request, apiCall *models.ApiCall) map[string]string {
	b := bytes.NewBuffer([]byte(""))
	err := req.Header.WriteSubset(b, reqWriteExcludeHeaderDump)
	if err != nil {
		return map[string]string{}
	}
	headers := map[string]string{}
	for _, header := range strings.Split(b.String(), "\n") {
		values := strings.Split(header, ":")
		if strings.EqualFold(values[0], "") {
			continue
		}
		if strings.EqualFold(values[0], "CurrentApiName") {
			apiCall.CurrentApiName = values[1]
			continue
		}
		if strings.EqualFold(values[0], "CurrentApiComment") {
			apiCall.CurrentApiComment = values[1]
			continue
		}
		headers[values[0]] = values[1]
	}
	return headers
}

func ReadHeadersFromResponse(header http.Header) map[string]string {
	headers := map[string]string{}
	for k, v := range header {
		log.Println(k, v)
		headers[k] = strings.Join(v, " ")
	}
	return headers
}

func ReadBody(req *http.Request) *string {
	save := req.Body
	var err error
	if req.Body == nil {
		req.Body = nil
	} else {
		save, req.Body, err = drainBody(req.Body)
		if err != nil {
			return nil
		}
	}
	b := bytes.NewBuffer([]byte(""))
	chunked := len(req.TransferEncoding) > 0 && req.TransferEncoding[0] == "chunked"
	if req.Body == nil {
		return nil
	}
	var dest io.Writer = b
	if chunked {
		dest = httputil.NewChunkedWriter(dest)
	}
	_, err = io.Copy(dest, req.Body)
	if chunked {
		dest.(io.Closer).Close()
	}
	req.Body = save
	body := b.String()
	return &body
}

func After(apiCall *models.ApiCall, record *responseRecorder, r *http.Request) {
	if strings.Contains(r.RequestURI, ".ico") || !gen.IsOn() {
		return
	}
	apiCall.MethodType = r.Method
	apiCall.CurrentPath = r.URL.Path
	apiCall.ResponseBody = record.Body.String()
	apiCall.ResponseCode = record.Status
	apiCall.ResponseHeader = ReadHeadersFromResponse(record.Header())
	if gen.IsStatusCodeValid(record.Status) {
		go gen.GenerateHtml(apiCall)
	}
}

// One of the copies, say from b to r2, could be avoided by using a more
// elaborate trick where the other copy is made during Request/Response.Write.
// This would complicate things too much, given that these functions are for
// debugging only.
func drainBody(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, nil, err
	}
	if err = b.Close(); err != nil {
		return nil, nil, err
	}
	return ioutil.NopCloser(&buf), ioutil.NopCloser(bytes.NewReader(buf.Bytes())), nil
}

//解析参数中的备注
func HandlerParamComment(req *http.Request, body string) (string, map[string]string) {
	var (
		e = make(map[string]string)
		f = make(map[string]interface{})
		h = make(map[string]interface{})
		m string
		n interface{}
	)
	_ = json.Unmarshal([]byte(body), &h)
	for q, w := range h {
		//对参数进行类型转换
		switch w.(type) {
		case string:
			m = w.(string)
		default:
			f[q] = w
			continue
		}
		slice := strings.Split(m, "`")
		if len(slice) > 1 {
			sliceComment := strings.Split(slice[1], ":")
			//分割备注里面的参数属性
			if len(sliceComment) > 1 {
				//按照备注进行类型转换
				switch sliceComment[1] {
				case "int":
					n, _ = strconv.Atoi(slice[0])
				case "float":
					n, _ = strconv.ParseFloat(slice[0], 64)
				default:
					//未定义处理方式,不做处理
					n = slice[0]
				}
			} else {
				//没有类型,不做处理
				n = slice[0]
			}
			//记录备注
			e[q] = slice[1]
			//重新对参数赋值
			f[q] = n
		} else {
			//没有备注直接赋值
			f[q] = w
		}
	}
	g, _ := json.Marshal(f)
	c := ioutil.NopCloser(bytes.NewReader(g))
	req.Body = c

	return string(g), e
}
