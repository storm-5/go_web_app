package ginimp

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func init() {
	middleWares = append(middleWares, bindFromUrlencoded())
}

func bindFromUrlencoded() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Contains(c.GetHeader("Content-Type"), "application/x-www-form-urlencoded") {
			args, err := bindArgs(c, bindWithFromUrlencoded)
			if err != nil {
				c.Error(fmt.Errorf("%s", stack(1)))
				c.Error(err)
				c.Status(400)
				c.Abort()
				return
			}
			c.Set("args", args)
		}
		c.Next()
	}
}

func bindWithFromUrlencoded(c *gin.Context, ft reflect.Type) (reflect.Value, error) {
	switch ft.Kind() {
	case reflect.Int:
		return reflect.Value{}, fmt.Errorf("not supported type")
	case reflect.String:
		return reflect.Value{}, fmt.Errorf("not supported type")
	case reflect.Slice:
		return reflect.Value{}, fmt.Errorf("not supported type")
	case reflect.Struct:
		fv, err := bindStruct(c, ft, binding.Form)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("binding failed")
		}
		return fv, nil
	case reflect.Ptr:
		return reflect.Value{}, fmt.Errorf("not supported type")
	default:
		return reflect.Value{}, fmt.Errorf("not supported type")
	}
}
