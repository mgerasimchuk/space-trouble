package util

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func RequestLogger(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		buf, _ := io.ReadAll(c.Request.Body)
		rdr1 := io.NopCloser(bytes.NewBuffer(buf))
		rdr2 := io.NopCloser(bytes.NewBuffer(buf)) //We have to create a new Buffer, because rdr1 will be read.

		msg := fmt.Sprintf("[%s] %s %s %s", time.Now().Format("02/Jan/2006:15:04:05 -0700"), c.Request.Method, c.Request.URL.Path, readBody(rdr1))
		logger.Debug(msg)

		c.Request.Body = rdr2
		c.Next()
	}
}

func readBody(reader io.Reader) string {
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(reader)

	s := buf.String()
	return s
}
