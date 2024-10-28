// hrm-bio/app/response/response.go
package response

import (
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
		Message: message.(string), // Set the message directly here
	})
}

func NotFound(ctx *gin.Context, message any, payloadCode ...string) {
	ctx.JSON(http.StatusNotFound, StatusResponse{
		Code:    404,
		Message: message.(string), // Set the message directly here
	})
}

// BadRequest ส่งผลลัพธ์เมื่อมีข้อผิดพลาดจากการขอข้อมูลที่ไม่ถูกต้อง
func BadRequest(ctx *gin.Context, message any, payloadCode ...string) {
	ctx.JSON(http.StatusBadRequest, StatusResponse{
		Code:    400,
		Message: message.(string), // Set the message directly here
	})
}

func Unauthorized(ctx *gin.Context, message any, payloadCode ...string) {
	ctx.JSON(http.StatusInternalServerError, StatusResponse{
		Code:    401,
		Message: message.(string), // Set the message directly here
	})
}

type Pagination struct {
	CurrentPage int `json:"current_page"`
	PerPage     int `json:"per_page"`
	TotalPages  int `json:"total_pages"`
	Total       int `json:"total"`
}

func SuccessWithPaginate(ctx *gin.Context, data any, limit, page, count int) {
	// Calculate pagination details
	totalPages := 1
	perPage := count
	if limit > 0 {
		totalPages = (count + limit - 1) / limit
		perPage = limit
	}

	pagination := Pagination{
		CurrentPage: page,
		PerPage:     perPage,
		TotalPages:  totalPages,
		Total:       count,
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
