package main

import (
	"errors"
	"fmt"
	"justbeyourselfandenjoy/service_order/swagger"
	"log"
	"net/mail"

	"github.com/google/uuid"
)

func validateOrderPayload(order swagger.Order) error {
	log.Println("validateOrderPayload has been called")

	//validating order
	if err := uuid.Validate(order.ID); err != nil {
		return fmt.Errorf("invalid order.ID [%s]", err)
	}

	if len(order.Products) == 0 {
		return errors.New("products can not be empty")
	}

	if order.Customer == nil {
		return errors.New("customer data can not be empty")
	}

	//validating products
	for i, product := range order.Products {
		if err := uuid.Validate(product.ID); err != nil {
			return fmt.Errorf("invalid product.ID [%s] for product [%d]", err, i)
		}

		if product.ProductCode == "" {
			return fmt.Errorf("invalid product code for product [%s]: product code can not be empty", product.ID)
		}

		if product.Quantity < 1 {
			return fmt.Errorf("invalid product quantity for product [%s]: product quantity can not be empty less than 1", product.ID)
		}

		if product.Quantity > 100 {
			return fmt.Errorf("invalid product quantity for product [%s]: product quantity can not be more than than 100", product.ID)
		}
	}

	//validating customer
	if err := uuid.Validate(order.Customer.ID); err != nil {
		return fmt.Errorf("invalid customer.ID [%s]", err)
	}

	if order.Customer.FirstName == "" {
		return fmt.Errorf("invalid customer data for [%s]: first name can not be empty", order.Customer.ID)
	}

	if order.Customer.EmailAddress == "" {
		return fmt.Errorf("invalid customer data for [%s]: email can not be empty", order.Customer.ID)
	}

	if _, err := mail.ParseAddress(order.Customer.EmailAddress); err != nil {
		return fmt.Errorf("invalid customer data for [%s]: wrong email format", order.Customer.ID)
	}

	if order.Customer.ShippingAddress == nil {
		return fmt.Errorf("invalid customer data for [%s]: shipping addrest must be set", order.Customer.ID)
	}

	if err := uuid.Validate(order.Customer.ShippingAddress.ID); err != nil {
		return fmt.Errorf("invalid ShippingAddress.ID [%s]", err)
	}

	if order.Customer.ShippingAddress.Line1 == "" {
		return fmt.Errorf("invalid address data for [%s]: Line1 can not be empty", order.Customer.ShippingAddress.ID)
	}

	if order.Customer.ShippingAddress.City == "" {
		return fmt.Errorf("invalid address data for [%s]: City can not be empty", order.Customer.ShippingAddress.ID)
	}

	if order.Customer.ShippingAddress.State == "" {
		return fmt.Errorf("invalid address data for [%s]: State can not be empty", order.Customer.ShippingAddress.ID)
	}

	if order.Customer.ShippingAddress.PostalCode == "" {
		return fmt.Errorf("invalid address data for [%s]: PostalCode can not be empty", order.Customer.ShippingAddress.ID)
	}

	return nil
}
