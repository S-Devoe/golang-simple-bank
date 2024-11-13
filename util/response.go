package util

import "fmt"

// i added this part
type Response struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
	Error  interface{} `json:"error"`
}

// Define the struct for the paginated response format
type PaginatedResponse[T any] struct {
	Status int              `json:"status"`
	Data   PaginatedData[T] `json:"data"`
	Error  interface{}      `json:"error"`
}

// Define the structure of the paginated data (inside the "data" field)
type PaginatedData[T any] struct {
	Results []T   `json:"results"` // Results is a generic slice of any type
	Page    int32 `json:"page"`
	Limit   int32 `json:"limit"`
	Total   int64 `json:"total"`
}

func CreateResponse(status int, data interface{}, err interface{}) Response {

	// If there's an error, check its type and convert it to a string
	var errorField interface{}
	if err != nil {
		// If err is an error type, convert it to a string
		switch e := err.(type) {
		case error:
			errorField = e.Error() // Convert the error to string
		case string:
			errorField = e // If it's already a string, use it directly
		default:
			errorField = fmt.Sprintf("Unknown error type: %v", e) // For any other type, handle gracefully
		}
	} else {
		// No error, so set error field to null
		errorField = nil
	}
	return Response{
		Status: status,
		Data:   data,
		Error:  errorField,
	}
}

// Standardize the paginated response format
func CreatePaginatedResponse[T any](status int, results []T, page, limit int32, totalItems int64, err interface{}) PaginatedResponse[T] {
	var errorField interface{}
	if err != nil {
		switch e := err.(type) {
		case error:
			errorField = e.Error() // Convert error to a string
		case string:
			errorField = e // Directly use the string
		default:
			errorField = fmt.Sprintf("Unknown error type: %v", e)
		}
	} else {
		errorField = nil
	}

	// Create the paginated data object
	data := PaginatedData[T]{
		Results: results,
		Page:    page,
		Limit:   limit,
		Total:   totalItems,
	}

	return PaginatedResponse[T]{
		Status: status,
		Data:   data,
		Error:  errorField,
	}
}
