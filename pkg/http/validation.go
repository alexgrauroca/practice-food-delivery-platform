package http

import (
	"reflect"
	"regexp"
	"time"
	// nolint:revive // intentional blank import to include tzdata in the final binary
	_ "time/tzdata"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var (
	// Precompiled regexes for phone validators
	phonePrefRx = regexp.MustCompile(`^\+\d{1,4}$`) // '+' followed by 1..4 digits
	phoneNumRx  = regexp.MustCompile(`^\d{4,14}$`)  // 4..14 digits
)

// init registers custom validation tags for Gin's validator engine.
func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// Validates IANA time zone ID.
		// Usage: binding:"iana_tz"
		_ = v.RegisterValidation("iana_tz", func(fl validator.FieldLevel) bool {
			s := fl.Field().String()
			if s == "" {
				// Let "required" handle empty values; this tag only checks validity when present.
				return true
			}
			_, err := time.LoadLocation(s)
			return err == nil
		})

		// Validates E.164-like phone prefix: leading '+' and 1..4 digits.
		// Usage: binding:"phone_pref"
		_ = v.RegisterValidation("phone_pref", func(fl validator.FieldLevel) bool {
			f := fl.Field()
			if f.Kind() != reflect.String {
				return false
			}
			s := f.String()
			if s == "" {
				return true // let "required" enforce presence
			}
			return phonePrefRx.MatchString(s)
		})

		// Validates phone number: 4..14 digits.
		// Usage: binding:"phone_num"
		_ = v.RegisterValidation("phone_num", func(fl validator.FieldLevel) bool {
			f := fl.Field()
			if f.Kind() != reflect.String {
				return false
			}
			s := f.String()
			if s == "" {
				return true // let "required" enforce presence
			}
			return phoneNumRx.MatchString(s)
		})
	}
}
