package internal

import (
	"time"

	"github.com/uptrace/bun"
)

type CustomerDBO struct {
	bun.BaseModel `bun:"customer,select:customer,alias:customer"`
	Id            uint       `bun:"customer_id,pk,autoincrement"`
	StoreId       uint       `bun:"store_id"`
	FirstName     string     `bun:"first_name,notnull"`
	LastName      string     `bun:"last_name,notnull"`
	Email         string     `bun:"email"`
	AddressId     uint       `bun:"address_id,notnull"`
	ActiveBool    bool       `bun:"activebool,notnull"`
	CreateDate    *time.Time `bun:"create_date,notnull"`
	LastUpdate    *time.Time `bun:"last_update"`
	Active        int        `bun:"active"`
}

func (c CustomerDBO) Name() string {
	return c.FirstName
}

func (c *CustomerDBO) ToResponse() UserResponse {
	return UserResponse{
		c.FirstName, c.LastName, c.Email,
	}
}
