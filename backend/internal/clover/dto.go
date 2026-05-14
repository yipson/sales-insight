package clover

// TokenResponse represents the OAuth token response from Clover.
type TokenResponse struct {
	AccessToken           string `json:"access_token"`
	AccessTokenExpiration int64  `json:"access_token_expiration"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiration int64 `json:"refresh_token_expiration"`
}

// Merchant represents a Clover merchant in API responses.
type MerchantResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Order represents a Clover order in API responses.
type OrderResponse struct {
	ID           string                `json:"id"`
	ModifiedTime int64                 `json:"modifiedTime"`
	Total        int64                 `json:"total"`
	LineItems    []OrderLineItemResponse `json:"lineItems"`
	Employee     *EmployeeResponse       `json:"employee"`
}

// OrderLineItem represents a line item in a Clover order.
type OrderLineItemResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Price    int64  `json:"price"`
	Item     *ItemResponse `json:"item"`
}

// Item represents a Clover item/product.
type ItemResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Categories []CategoryResponse `json:"categories"`
}

// Category represents a Clover category.
type CategoryResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Employee represents a Clover employee.
type EmployeeResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	Active bool   `json:"active"`
}

// Payment represents a Clover payment.
type PaymentResponse struct {
	ID       string `json:"id"`
	Amount   int64  `json:"amount"`
	OrderID  string `json:"orderId"`
	Employee *EmployeeResponse `json:"employee"`
}
