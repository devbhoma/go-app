package httpserver

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetQueryContextInt(c *gin.Context, key string) int {
	_str := GetQueryContextString(c, key)
	if _val, err := strconv.Atoi(_str); err == nil && _val > 0 {
		return _val
	}
	return 0
}

func GetQueryContextString(c *gin.Context, key string) string {
	if _str, ok := c.GetQuery(key); ok && _str != "" {
		return _str
	}
	return ""
}

func GetQueryContextArray(c *gin.Context, key string) []string {
	if arStr, ok := c.GetQueryArray(key); ok {
		return arStr
	}
	return nil
}

func GetQueryContextMap(c *gin.Context, key string) map[string]string {
	if arStr, ok := c.GetQueryMap(key); ok {
		return arStr
	}
	return nil
}
