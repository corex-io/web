package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/corex-io/log"
)

var zeroTime = time.Time{}

// Context context
type Context struct {
	http.ResponseWriter
	*http.Request
	statusCode int
	Timestamp  time.Time
	log.Logger
}

func (ctx *Context) reset() {
	ctx.ResponseWriter = nil
	ctx.Request = nil
	ctx.statusCode = 0
	ctx.Timestamp = zeroTime
	ctx.Logger = nil
}

// IsFinish return handle is closed or not
func (ctx *Context) IsFinish() bool {
	return ctx.statusCode != 0
}

// GetQuery get query
func (ctx *Context) GetQuery() map[string]string {
	res := make(map[string]string, len(ctx.Form))
	for key := range ctx.Form {
		res[key] = ctx.Form.Get(key)
	}
	return res
}

// GetJSONBody get json body args
func (ctx *Context) GetJSONBody(v interface{}) error {
	if ctx.Body == nil {
		return fmt.Errorf("body is nil")
	}
	defer ctx.Body.Close()
	dec := json.NewDecoder(ctx.Body)
	if err := dec.Decode(v); err != nil && err != io.EOF {
		return err
	}
	return nil
}

// Remote return request addr,  copied from net/url.stripPort
func (ctx *Context) Remote() string {
	colon := strings.IndexByte(ctx.RemoteAddr, ':')
	if colon == -1 {
		return ctx.RemoteAddr
	}
	if i := strings.IndexByte(ctx.RemoteAddr, ']'); i != -1 {
		return strings.TrimPrefix(ctx.RemoteAddr[:i], "[")
	}
	return ctx.RemoteAddr[:colon]
}

// GetCookies get cookies
func (ctx *Context) GetCookies() []*http.Cookie {
	return ctx.Request.Cookies()
}

// GetCookie get cookie
func (ctx *Context) GetCookie(name string) string {
	cookie, err := ctx.Cookie(name)
	if err != nil {
		return ""
	}
	return cookie.Value
}

// GetForm formdata, Content-Type must be multipart/form-data.
// TODO: RemoveAll removes any temporary files associated with a Form.
func (ctx *Context) GetForm() (map[string]string, map[string]*multipart.FileHeader, error) {
	reader, err := ctx.MultipartReader()
	if err != nil {
		return nil, nil, err
	}
	form, err := reader.ReadForm(10000)
	if err != nil {
		return nil, nil, err
	}
	values := make(map[string]string)
	for k, v := range form.Value {
		if len(v) > 0 {
			values[k] = v[0]
		}
	}
	files := make(map[string]*multipart.FileHeader)
	for k, v := range form.File {
		if len(v) > 0 {
			files[k] = v[0]
		}
	}
	return values, files, nil
}

// Redirect response redirect
func (ctx *Context) Redirect(url string, statusCode int) {
	ctx.statusCode = statusCode
	http.Redirect(ctx.ResponseWriter, ctx.Request, url, ctx.statusCode)
}

// Error response error
func (ctx *Context) Error(statusCode int) {
	ctx.statusCode = statusCode
	http.Error(ctx.ResponseWriter, http.StatusText(ctx.statusCode), ctx.statusCode)
}

// Download support download file
func (ctx *Context) Download(filename string) {
	info, err := os.Stat(filename)
	if err != nil {
		ctx.Error(toHTTPError(err))
		return
	}
	f, err := os.Open(filename)
	if err != nil {
		ctx.Error(toHTTPError(err))
		return
	}
	defer f.Close()
	ctx.ResponseWriter.Header().Set("Content-Length", fmt.Sprintf("%d", info.Size()))
	ctx.ResponseWriter.Header().Set("Content-Type", "multipart/form-data")
	ctx.ResponseWriter.Header().Set("Content-Disposition:", fmt.Sprintf("attachment;filename=%s", filename))
	ctx.statusCode = http.StatusOK
	io.Copy(ctx.ResponseWriter, f)
}

//SaveFile save file to disk
func SaveFile(fh *multipart.FileHeader, path string, name ...string) (string, int64, error) {
	file, err := fh.Open()
	if err != nil {
		return "", 0, err
	}
	defer file.Close()

	var filename string
	if len(name) == 0 {
		filename = fh.Filename
	} else {
		filename = name[0]
	}
	filepath := filepath.Join(path, filename)
	f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return "", 0, err
	}
	defer f.Close()
	cnt, err := io.Copy(f, file)
	return filepath, cnt, err
}

// RecvFormFile recv form file
func (ctx *Context) RecvFormFile(path string) error {
	var err error
	const _24K = (1 << 10) * 24 // 24 MB
	if err := ctx.Request.ParseMultipartForm(_24K); err != nil {
		return err
	}
	for _, fheaders := range ctx.Request.MultipartForm.File {
		for _, hdr := range fheaders {
			// open uploaded
			var infile multipart.File
			if infile, err = hdr.Open(); err != nil {
				return err
			}
			defer infile.Close()
			// open destination
			var outfile *os.File //https://www.jianshu.com/p/3e0c0609d419
			if outfile, err = os.Create(path + hdr.Filename); err != nil {
				return err
			}
			defer outfile.Close()
			// 32K buffer copy
			var written int64
			if written, err = io.Copy(outfile, infile); err != nil {
				return err
			}
			fmt.Println("uploaded file:" + hdr.Filename + ";length:" + strconv.Itoa(int(written)))
		}
	}
	return nil
}

// RecvFile2 recvfile2
func (ctx *Context) RecvFile2(name string, path string) error {
	disposition := ctx.Request.Header.Get("Content-Disposition")
	filename := strings.Trim(strings.SplitN(disposition, "filename=", 2)[1], `"`)

	tmpfile, err := ioutil.TempFile(".", "tmp")
	if err != nil {
		return err
	}
	if _, err = io.Copy(tmpfile, ctx.Body); err != nil {
		return err
	}
	tmpfile.Close()
	ctx.Body.Close()
	return os.Rename(tmpfile.Name(), filename)
}

// RecvFile recv file
func (ctx *Context) RecvFile(name string, path string) (string, int64, error) {
	file, head, err := ctx.FormFile(name)
	if err != nil {
		return "", 0, err
	}
	defer file.Close()
	filepath := filepath.Join(path, head.Filename)
	f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return "", 0, err
	}
	defer f.Close()
	cnt, err := io.Copy(f, file)
	return filepath, cnt, err
}

// Render render template no cache
func (ctx *Context) Render(tpl string, data interface{}) {
	// path := filepath.Join(ctx.Config.WebPath, tpl)
	t, err := template.ParseFiles(tpl)
	if err != nil {
		ctx.Error(toHTTPError(err))
		return
	}
	t.Execute(ctx.ResponseWriter, data)
}

// Text return resp with text format
func (ctx *Context) Text(response []byte) {
	ctx.ResponseWriter.Write(response)
}

// JSON json api
func (ctx *Context) JSON(v interface{}, code int, err error) {
	var msg string
	if err != nil {
		msg = err.Error()
	}
	var data string
	switch v.(type) {
	case []byte:
		b := v.([]byte)
		if ok := json.Valid(b); ok {
			data = string(b)
		} else {
			data = fmt.Sprintf(`"%s"`, string(b))
		}

	case string:
		b := v.(string)
		if ok := json.Valid([]byte(b)); ok {
			data = b
		} else {
			data = fmt.Sprintf(`"%s"`, b)
		}

	default:
		b, err := json.Marshal(v)
		if err != nil {
			code = 101
			msg += fmt.Sprintf("\n  json marshal data: %v", err)
		} else {
			data = string(b)
		}
	}

	if msg == "" {
		msg = "success"
	}

	result := fmt.Sprintf(`{"code": %d, "msg": "%s", "data": %s}`, code, msg, data)
	ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=UTF-8")
	ctx.ResponseWriter.Write([]byte(result))
}
