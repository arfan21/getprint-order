package validation

import(
	"github.com/arfan21/getprint-order/models"
	validator "github.com/go-ozzo/ozzo-validation/v4"
)

func Validate(order models.Order) error{
	return validator.Errors{
		"partner_id": validator.Validate(order.PartnerID, validator.Required),
		"user_id":    validator.Validate(order.UserID, validator.Required),
		"order_details": validator.Validate(order.OrderDetails, validator.Required),
	}.Filter()
}