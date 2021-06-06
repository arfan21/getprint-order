package validation

import (
	models2 "github.com/arfan21/getprint-order/app/models"
	validator "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func Validate(order models2.Order) error{
	return validator.Errors{
		"user_id":    validator.Validate(order.UserID, is.UUIDv4),
		"partner_id": validator.Validate(order.PartnerID, validator.Required),
		"order_details": validator.Validate(order.OrderDetails, validator.Required),
	}.Filter()
}