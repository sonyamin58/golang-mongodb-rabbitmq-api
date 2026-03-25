package response

import (
	"encoding/json"
	"math"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type PaginationResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

type Pagination struct {
	Limit     int   `json:"limit"`
	Offset    int   `json:"offset"`
	Total     int64 `json:"total"`
	TotalPage int   `json:"total_page"`
}

func Success(data interface{}) Response {
	return Response{
		Success: true,
		Data:    data,
	}
}

func Error(message string) Response {
	return Response{
		Success: false,
		Message: message,
	}
}

func Paginated(data interface{}, total int64, limit, offset int) PaginationResponse {
	totalPage := int(math.Ceil(float64(total) / float64(limit)))
	
	return PaginationResponse{
		Success: true,
		Data:    data,
		Pagination: Pagination{
			Limit:     limit,
			Offset:    offset,
			Total:     total,
			TotalPage: totalPage,
		},
	}
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
EOF; __hermes_rc=$?; printf '__HERMES_FENCE_a9f7b3__'; exit $__hermes_rc
