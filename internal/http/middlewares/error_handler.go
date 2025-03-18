package middlewares

import (
	"errors"
	"net/http"
	"reflect"

	"github.com/Joe5451/go-oauth2-server/internal/adapter/handlers"
	"github.com/Joe5451/go-oauth2-server/internal/domain"
	"github.com/Joe5451/go-oauth2-server/internal/socialproviders"
	"github.com/gin-gonic/gin"
)

func ErrorHandler(errMap ...*errorMapping) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		lastErr := c.Errors.Last()
		if lastErr == nil {
			return
		}

		for _, e := range errMap {
			for _, e2 := range e.fromErrors {
				if errors.Is(lastErr.Err, e2) {
					e.toResponse(c, lastErr.Err)
					return
				}
			}
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "INTERNAL_ERROR",
			"message": "An unexpected error occurred.",
		})
	}
}

func isType(a, b interface{}) bool {
	return reflect.TypeOf(a) == reflect.TypeOf(b)
}

type errorMapping struct {
	fromErrors   []error
	toStatusCode int
	toResponse   func(ctx *gin.Context, err error)
}

func (r *errorMapping) ToStatusCode(statusCode int) *errorMapping {
	r.toStatusCode = statusCode
	r.toResponse = func(ctx *gin.Context, err error) {
		ctx.Status(statusCode)
	}
	return r
}

func (r *errorMapping) ToResponse(response func(ctx *gin.Context, err error)) *errorMapping {
	r.toResponse = response
	return r
}

func Map(err ...error) *errorMapping {
	return &errorMapping{
		fromErrors: err,
	}
}

func InitErrorHandler() gin.HandlerFunc {
	return ErrorHandler(
		Map(handlers.ErrValidation).ToResponse(func(c *gin.Context, err error) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			})
		}),
		Map(handlers.ErrUnauthorized).ToResponse(func(c *gin.Context, err error) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    "UNAUTHORIZED",
				"message": "Requires authentication.",
			})
		}),
		Map(handlers.ErrMissingFile).ToResponse(func(c *gin.Context, err error) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    "MISSING_FILE",
				"message": err.Error(),
			})
		}),
		Map(handlers.ErrInvalidFileFormat).ToResponse(func(c *gin.Context, err error) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    "INVALID_FILE_FORMAT",
				"message": err.Error(),
			})
		}),
		Map(domain.ErrUserNotFound).ToResponse(func(c *gin.Context, err error) {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    "USER_NOT_FOUND",
				"message": "The user does not exist.",
			})
		}),
		Map(domain.ErrInvalidCredentials).ToResponse(func(c *gin.Context, err error) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    "INVALID_CREDENTIALS",
				"message": "Incorrect email or password.",
			})
		}),
		Map(domain.ErrDuplicateEmail).ToResponse(func(c *gin.Context, err error) {
			c.JSON(http.StatusConflict, gin.H{
				"code":    "DUPLICATE_EMAIL",
				"message": "The email is already in use.",
			})
		}),
		Map(socialproviders.ErrInvalidProvider).ToResponse(func(c *gin.Context, err error) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "INVALID_SOCIAL_PROVIDER",
				"message": err.Error(),
			})
		}),
		Map(domain.ErrInvalidLinkToken).ToResponse(func(c *gin.Context, err error) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "INVALID_LINK_TOKEN",
				"message": err.Error(),
			})
		}),
		Map(domain.ErrMismatchedLinkedUser).ToResponse(func(c *gin.Context, err error) {
			c.JSON(http.StatusConflict, gin.H{
				"code":    "MISMATCHED_LINKED_USER",
				"message": "The linked social account belongs to a different user.",
			})
		}),
		Map(domain.ErrSocialAccountAlreadyLinked).ToResponse(func(c *gin.Context, err error) {
			c.JSON(http.StatusConflict, gin.H{
				"code":    "SOCIAL_ACCOUNT_ALREADY_LINKED",
				"message": "The social account has already been linked to another user.",
			})
		}),
		Map(domain.ErrSocialAccountAlreadyUnlinked).ToResponse(func(c *gin.Context, err error) {
			c.JSON(http.StatusConflict, gin.H{
				"code":    "SOCIAL_ACCOUNT_ALREADY_UNLINKED",
				"message": "The social account has either not been linked or has already been unlinked.",
			})
		}),
		Map(socialproviders.ErrOAuth2RetrieveError).ToResponse(func(c *gin.Context, err error) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "OAUTH2_RETRIEVE_ERROR",
				"message": err.Error(),
			})
		}),
	)
}
