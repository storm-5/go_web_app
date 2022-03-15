package ginimp

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/storm-5/go-app/pkg/myerr"
	"github.com/storm-5/go-app/pkg/webfw"

	"github.com/gin-gonic/gin"
)

func responseHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		path := c.Request.URL.Path
		method := c.Request.Method

		args := make([]reflect.Value, 0)
		reqArgs, exist := c.Get("args")
		if exist {
			args = reqArgs.([]reflect.Value)
		}
		rsps, exist := c.Get("responses")
		if !exist {
			log.Errorf("method: %s, path: %s, requestArgs: %v, msg: %s",
				method,
				path,
				reflectValueToString(args),
				c.Errors.Errors(),
			)
			c.JSON(200, webfw.BaseRspDto{Code: c.Writer.Status(), Msg: c.Errors.Last().Error()})
			return
		}
		responses := rsps.([]reflect.Value)
		err := responses[1].Interface()
		if err != nil {
			var detailErrMsg string
			switch err := err.(type) {
			case *myerr.MyError:
				detailErrMsg = err.Marshal()
				c.JSON(200, webfw.BaseRspDto{Code: err.Code, Msg: err.Error()})
			default:
				detailErrMsg = fmt.Sprintf("%s", err)
				c.JSON(200, webfw.BaseRspDto{Code: http.StatusInternalServerError, Msg: http.StatusText(http.StatusInternalServerError)})
			}
			log.Errorf("method: %s, path: %s, requestArgs: %v, msg: %s",
				method,
				path,
				reflectValueToString(args),
				detailErrMsg,
			)
		} else {
			switch rt := responses[0].Interface().(type) {
			case http.ResponseWriter:
				// fmt.Println("")
			case webfw.ByteStreamRspDto:
				c.DataFromReader(http.StatusOK, rt.ContentLength, rt.ContentType, rt.Reader, rt.ExtraHeaders)
			case webfw.FileStreamRspDto:
				c.FileAttachment(rt.Path, rt.FileName)
			default:
				c.JSON(200, webfw.BaseRspDto{Code: 0, Msg: "ok", Data: rt})
			}
		}
	}
}
