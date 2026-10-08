package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gix-coder/gix-coder/shared/errors"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

// ValidationMiddleware validates request bodies using protovalidate
func ValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// For now, we rely on Gin's ShouldBindJSON which does basic validation
		// Full protovalidate integration would require parsing the request body
		// into the protobuf message type and calling the Validate method
		c.Next()
	}
}

// ValidateProtoMessage validates a protobuf message using protovalidate
func ValidateProtoMessage(msg proto.Message) error {
	if v, ok := interface{}(msg).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return errors.WithCause(err, errors.ErrCodeValidation, "request validation failed")
		}
	}
	return nil
}

// ProtoValidationMiddleware creates middleware that validates protobuf requests
func ProtoValidationMiddleware(messageFactory func() proto.Message) gin.HandlerFunc {
	return func(c *gin.Context) {
		msg := messageFactory()
		if err := c.ShouldBindJSON(msg); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "VALIDATION_ERROR",
					"message": "invalid request body: " + err.Error(),
				},
			})
			return
		}

		if err := ValidateProtoMessage(msg); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "VALIDATION_ERROR",
					"message": err.Error(),
				},
			})
			return
		}

		c.Set("validated_request", msg)
		c.Next()
	}
}

// ConvertStructToProto converts a map to a protobuf Struct
func ConvertStructToProto(m map[string]string) (*structpb.Struct, error) {
	data := make(map[string]interface{}, len(m))
	for k, v := range m {
		data[k] = v
	}
	return structpb.NewStruct(data)
}

// ConvertProtoToStruct converts a protobuf Struct to a map
func ConvertProtoToStruct(s *structpb.Struct) map[string]string {
	if s == nil {
		return nil
	}
	result := make(map[string]string)
	for k, v := range s.AsMap() {
		if str, ok := v.(string); ok {
			result[k] = str
		}
	}
	return result
}