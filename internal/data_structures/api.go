package internal

type UserId struct {
	Id int32 `json:"id" uri:"id" binding:"required"`
}

type UserResponse struct {
	Name     string `json:"name"`
	LastName string `json:"lastName"`
	Email    string `json:"email"`
}

type CustomerInsertRequest struct {
	StoreId   uint   `json:"storeId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	AddressId uint   `json:"addressId"`
	Active    int    `json:"active"`
}

func (c *CustomerInsertRequest) ToCustomerDBO() CustomerDBO {
	return CustomerDBO{
		StoreId:   c.StoreId,
		FirstName: c.FirstName,
		LastName:  c.LastName,
		Email:     c.Email,
		AddressId: c.AddressId,
		Active:    c.Active,
	}
}
