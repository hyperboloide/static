package static_test

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	. "github.com/hyperboloide/static"
	"github.com/hyperboloide/static/demo/files"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Static", func() {

	var ts *httptest.Server
	var staticHandler *Handler

	BeforeEach(func() {
		r := gin.Default()
		rootGroup := r.Group("/")
		{
			staticHandler = NewHandler(files.Asset, files.AssetNames)
			staticHandler.Register(rootGroup)
		}
		ts = httptest.NewServer(r)

	})

	AfterEach(func() {
		ts.Close()
	})

	checkFile := func(pth, contentType string, maxAge int) {
		resp, err := http.Get(ts.URL + "/" + pth)
		Ω(err).ToNot(HaveOccurred())
		body, err := ioutil.ReadAll(resp.Body)
		Ω(err).ToNot(HaveOccurred())
		err = resp.Body.Close()
		Ω(err).ToNot(HaveOccurred())

		Ω(resp.StatusCode).To(Equal(200))
		ma := fmt.Sprintf("max-age=%d", maxAge)
		Ω(resp.Header.Get("Cache-Control")).To(Equal(ma))
		Ω(resp.Header.Get("Content-Type")).To(Equal(contentType))

		eq := bytes.Equal(body, files.MustAsset(pth))
		Ω(eq).To(BeTrue())

		sum := md5.Sum(body)
		Ω(resp.Header.Get("Etag")).To(Equal(fmt.Sprintf("%x", sum)))
	}

	It("should get index.html", func() {
		checkFile("index.html", "text/html; charset=utf-8", 86400)
	})

	It("should get test.jpg", func() {
		checkFile("test.jpg", "image/jpeg", 86400)
	})

	It("should get bootstrap/css/bootstrap.min.css", func() {
		checkFile("bootstrap/css/bootstrap.min.css", "text/css; charset=utf-8", 86400)
	})

	It("should return 404 for /", func() {
		resp, err := http.Get(ts.URL + "/")
		Ω(err).ToNot(HaveOccurred())
		Ω(resp.StatusCode).To(Equal(404))
	})

	It("should not return 404 for / if 'index.html' is added to indexes", func() {
		staticHandler.AddIndexes("index.html", "index.htm")
		resp, err := http.Get(ts.URL + "/")
		Ω(err).ToNot(HaveOccurred())
		Ω(resp.StatusCode).To(Equal(200))
	})

})
