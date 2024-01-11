package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/pakojabi/simplebank/util"
)


// validCurrency gets registered as a struct validator in server.go
var validCurrency validator.Func = func (fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
	}
	return false
}

