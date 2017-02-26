package static

import (
	"crypto/md5"
	"fmt"
	"mime"
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	// DefaultMaxAge is 60 days.
	DefaultMaxAge = 86400
)

// ContentFunc is an alias for the Asset method of the generated bindata.go file.
type ContentFunc func(string) ([]byte, error)

// NamesFunc is an alias for the AssetNames method of the generated bindata.go file.
type NamesFunc func() []string

// Handler is the base struct used to serve the bindata.
type Handler struct {
	contentFn ContentFunc
	namesFn   NamesFunc

	basePath string
	sums     map[string]string
	mimes    map[string]string
	maxAge   uint
	indexes  []string
}

// NewHandler create a new Handler. contentFn is the Asset method
// and namesFn is the AssetNames method from the generated bindata.go.
func NewHandler(contentFn ContentFunc, namesFn NamesFunc) *Handler {
	res := &Handler{
		contentFn: contentFn,
		namesFn:   namesFn,
		sums:      map[string]string{},
		mimes:     map[string]string{},
		maxAge:    DefaultMaxAge,
		indexes:   []string{},
	}

	for _, name := range namesFn() {
		content, err := contentFn(name)
		if err != nil {
			continue

		}

		i := md5.Sum(content)
		res.sums[name] = fmt.Sprintf("%x", i)

		l := strings.Split(name, ".")
		suffix := "." + l[len(l)-1]
		if m := mime.TypeByExtension(suffix); m != "" {
			res.mimes[name] = m

		}
	}

	return res
}

// AddIndexes adds possible indexes. For example if there are one or more directories
// that contains an index.html file that you want to serve when the directory
// name is requested then add "index.html".
func (h *Handler) AddIndexes(names ...string) {
	h.indexes = append(h.indexes, names...)
}

// Register adds the static handlers to the gin.RouterGroup.
func (h *Handler) Register(r *gin.RouterGroup) {
	h.basePath = r.BasePath()
	r.GET("*any", h.get)
}

// SetMaxAge sets the max-age header (if set to 0 then caching is disabled).
func (h *Handler) SetMaxAge(ma uint) {
	h.maxAge = ma
}

func (h *Handler) getMime(path string) string {
	if mime, exists := h.mimes[path]; exists {
		return mime

	} else {
		return "text/plain"

	}
}

func (h *Handler) fileExists(pth string) bool {
	_, exists := h.sums[pth]
	return exists
}

func (h *Handler) getPath(c *gin.Context) string {
	return strings.TrimPrefix(c.Request.URL.Path, h.basePath)
}

func (h *Handler) matchSum(pth, sum string) bool {
	v, exists := h.sums[pth]
	if !exists || sum == "" {
		return false

	}
	return v == sum

}

func (h *Handler) getIndex(pth string) *string {
	for _, idx := range h.indexes {
		idxPth := path.Join(pth, idx)
		if h.fileExists(idxPth) {
			return &idxPth

		}
	}
	return nil
}

func (h *Handler) get(c *gin.Context) {
	pth := h.getPath(c)

	if !h.fileExists(pth) {
		if idxPth := h.getIndex(pth); idxPth == nil {
			c.String(http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return

		} else {
			pth = *idxPth

		}
	}

	if b, err := h.contentFn(pth); err != nil {
		c.String(http.StatusNotFound, http.StatusText(http.StatusNotFound))

	} else {
		if h.maxAge > 0 {
			c.Writer.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", h.maxAge))
			c.Writer.Header().Set("Etag", h.sums[pth])

			if h.matchSum(pth, c.Request.Header.Get("If-None-Match")) {
				c.Status(http.StatusNotModified)
				return
			}
		}

		c.Writer.Header().Set("Content-Type", h.getMime(pth))
		if _, err := c.Writer.Write(b); err != nil {
			c.Status(http.StatusInternalServerError)

		}
	}
}
