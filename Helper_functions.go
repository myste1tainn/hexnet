package msnet

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/myste1tainn/hexlog"
)

var (
	ContextKeyLogger = "ContextKeyLogger"
)

func setLoggerToContext(l *log.Logger, ctx *gin.Context) {
	ctx.Set(ContextKeyLogger, l)
}

func getBody(res any) any {
	switch t := res.(type) {
	case error:
		if err, ok := t.(stackTracer); ok {
			lines := []string{}
			for _, f := range err.StackTrace() {
				lines = append(lines, fmt.Sprintf("%+s:%d", f, f))
			}
			return strings.Join(lines, "\n")
		} else {
			return t
		}
	default:
		return t
	}
}
