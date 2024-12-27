/*
 * Kafka Basics Pet Project API
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type Customer struct {

	ID string `json:"ID,omitempty"`

	FirstName string `json:"FirstName,omitempty"`

	LastName string `json:"LastName,omitempty"`

	EmailAddress string `json:"EmailAddress,omitempty"`

	ShippingAddress *Address `json:"ShippingAddress,omitempty"`
}
