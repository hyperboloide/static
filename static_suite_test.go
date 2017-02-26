package static_test

import (
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestStatic(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Static Suite")
}

var _ = BeforeSuite(func() {
	gin.SetMode(gin.ReleaseMode)
})
