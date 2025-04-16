// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package main

import (
	"net/http"

	"github.com/clin211/miniblog-v2/internal/pkg/errno"
	"github.com/gin-gonic/gin"
)

func main() {
	// Create and configure Gin router
	r := gin.Default()

	// Register handlers
	r.GET("/error-demo", errorDemoHandler)
	r.GET("/error-check", errorCheckHandler)

	// Start server
	r.Run(":8080")
}

// errorDemoHandler demonstrates error creation and modification
func errorDemoHandler(c *gin.Context) {
	// Create an ErrorX error, representing a database connection failure
	err := errno.New(500, "InternalError.DBConnection", "Something went wrong: %s", "DB connection failed")

	// Add metadata to enhance the error context for debugging
	err.WithMetadata(map[string]string{
		"user_id":    "12345",
		"request_id": "abc-def",
	})

	// Add more metadata using the KV method
	err.KV("trace_id", "xyz-789")

	// Update the error message
	err.WithMessage("Updated message: %s", "retry failed")

	// Return the error as JSON response
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": err.Error(),
		"details": gin.H{
			"code":     err.Code,
			"reason":   err.Reason,
			"message":  err.Message,
			"metadata": err.Metadata,
		},
	})
}

// errorCheckHandler demonstrates error type checking
func errorCheckHandler(c *gin.Context) {
	// Generate an error from the doSomething function
	someErr := doSomething()

	// Check if the error matches predefined error types
	isUsernameError := errno.ErrUsernameInvalid.Is(someErr)
	isPasswordError := errno.ErrPasswordInvalid.Is(someErr)

	c.JSON(http.StatusBadRequest, gin.H{
		"error":           someErr.Error(),
		"isUsernameError": isUsernameError,
		"isPasswordError": isPasswordError,
		"details": gin.H{
			"code":    someErr.(*errno.ErrorX).Code,    // Type assertion - be careful in production code
			"reason":  someErr.(*errno.ErrorX).Reason,  // Type assertion
			"message": someErr.(*errno.ErrorX).Message, // Type assertion
		},
	})
}

// doSomething returns a predefined error
func doSomething() error {
	// Return a predefined error type with a custom message
	return errno.ErrUsernameInvalid.WithMessage("Username is too short")
}
