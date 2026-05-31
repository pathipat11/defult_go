// app/response/response.go
package response

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StatusResponse struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

type Response struct {
	Status StatusResponse `json:"status"`
	Data   any            `json:"data,omitempty"`
}

// toMessage safely converts an arbitrary message value into a string,
// falling back to a default when nil or of an unexpected type. This avoids
// the panic that a naked message.(string) assertion would cause on nil.
func toMessage(message any, fallback string) string {
	switch v := message.(type) {
	case nil:
		return fallback
	case string:
		if v == "" {
			return fallback
		}
		return v
	case error:
		return v.Error()
	default:
		return fmt.Sprint(v)
	}
}

type ResponsePaginate struct {
	Status     StatusResponse `json:"status"`
	Data       any            `json:"data,omitempty"`
	Pagination Pagination     `json:"pagination"`
}

type ResponsePaginate0 struct {
	Status     StatusResponse `json:"status"`
	Data       any            `json:"data,omitempty"`
	Pagination any            `json:"pagination"`
}

// Success ส่งผลลัพธ์เมื่อสำเร็จ
func Success(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, Response{StatusResponse{
		Code:    200,
		Message: "Success",
	}, data})
}

// InternalError ส่งผลลัพธ์เมื่อมีข้อผิดพลาดภายใน
func InternalError(ctx *gin.Context, message any, payloadCode ...string) {
	ctx.JSON(http.StatusInternalServerError, StatusResponse{
		Code:    500,
		Message: toMessage(message, "Internal Server Error"),
	})
}

func NotFound(ctx *gin.Context, message any, payloadCode ...string) {
	ctx.JSON(http.StatusNotFound, StatusResponse{
		Code:    404,
		Message: toMessage(message, "Not Found"),
	})
}

// BadRequest ส่งผลลัพธ์เมื่อมีข้อผิดพลาดจากการขอข้อมูลที่ไม่ถูกต้อง
func BadRequest(ctx *gin.Context, message any) {
	ctx.JSON(http.StatusBadRequest, StatusResponse{
		Code:    400,
		Message: toMessage(message, "Bad Request"),
	})
}

func Unauthorized(ctx *gin.Context, message any, payloadCode ...string) {
	ctx.JSON(http.StatusUnauthorized, StatusResponse{
		Code:    401,
		Message: toMessage(message, "Unauthorized"),
	})
}

type Pagination struct {
	Page  int `json:"page"`
	Size  int `json:"size"`
	Total int `json:"total"`
}

func SuccessWithPaginate(ctx *gin.Context, data any, size, page, count int) {

	pagination := Pagination{
		Page:  page,
		Size:  size,
		Total: count,
	}

	if pagination.Total == 0 {
		ctx.JSON(http.StatusOK, ResponsePaginate0{
			Status: StatusResponse{
				Code:    200,
				Message: "Success",
			},
			Data:       []any{},
			Pagination: gin.H{},
		})
		return
	} else {
		ctx.JSON(http.StatusOK, ResponsePaginate{
			Status: StatusResponse{
				Code:    200,
				Message: "Success",
			},
			Data:       data,
			Pagination: pagination,
		})
	}
}

func Forbidden(ctx *gin.Context, message any, payloadCode ...string) {
	ctx.JSON(http.StatusForbidden, StatusResponse{
		Code:    403,
		Message: toMessage(message, "Forbidden"),
	})
}
