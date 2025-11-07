package http

import (
	"time"
	// nolint:revive // intentional blank import to include tzdata in the final binary
	_ "time/tzdata"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// init registers custom validation tags for Gin's validator engine.
func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("iana_tz", func(fl validator.FieldLevel) bool {
			s := fl.Field().String()
			if s == "" {
				// Let "required" handle empty values; this tag only checks validity when present.
				return true
			}
			_, err := time.LoadLocation(s)
			return err == nil
		})
	}
}
