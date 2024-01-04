package network

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/soulteary/docker-quick-docs/internal/config"
)

func MuteNoise(body []byte) []byte {
	content := UpdateBody(body, []byte("https://www.googletagmanager.com/gtag/js"), []byte("/keep-quiet.js"))
	content = UpdateBody(content, []byte("https://www.google-analytics.com/analytics.js"), []byte("/keep-quiet.js"))
	content = UpdateBody(content, []byte("https://hm.baidu.com/hm.js"), []byte("/keep-quiet.js"))
	return content
}

func UpdateBody(content []byte, from []byte, to []byte) []byte {
	return bytes.ReplaceAll(content, from, to)
}

type UpdateJob struct {
	From string
	To   string
}

func Forward(port int) func(c *gin.Context) {
	target := fmt.Sprintf("http://%s:%d", config.DOCS_INTERNAL_HOST, port)
	url, _ := url.Parse(target)
	internal := httputil.NewSingleHostReverseProxy(url)
	return func(c *gin.Context) {
		internal.ModifyResponse = func(response *http.Response) error {
			if len(config.PostRules) == 0 || response.ContentLength == 0 || response.Body == nil {
				return nil
			}

			mimeType := strings.ToLower(response.Header.Get("Content-Type"))
			needUpdate := false

			var jobs []UpdateJob
			for _, rule := range config.PostRules {
				// match rule type
				if strings.HasPrefix(mimeType, rule.Type) {
					// match rule dir
					if rule.Dir == "*" || strings.HasPrefix(c.Request.URL.Path, rule.Dir) {
						needUpdate = true
						var job UpdateJob
						job.From = rule.From
						job.To = rule.To
						jobs = append(jobs, job)
					}
				}
			}

			// only allow html or need update content
			if mimeType != "text/html" && !needUpdate {
				return nil
			}

			body, err := io.ReadAll(response.Body)
			if err != nil {
				return err
			}
			response.Body.Close()
			bodyUpdated := MuteNoise(body)

			for _, job := range jobs {
				bodyUpdated = UpdateBody(bodyUpdated, []byte(job.From), []byte(job.To))
			}

			bodyLength := len(bodyUpdated)
			response.Body = io.NopCloser(bytes.NewReader(bodyUpdated))
			response.ContentLength = int64(bodyLength)
			response.Header.Set("Content-Length", strconv.Itoa(bodyLength))
			return nil
		}

		internal.ServeHTTP(c.Writer, c.Request)
	}
}
