package mongodb

type ProfileType byte

const (
	ProfileTypeBuyer  ProfileType = 0
	ProfileTypeSeller ProfileType = 1
)

type ProductStatus byte

const (
	ProductStatusAvailable ProductStatus = 0
	ProductStatusSold      ProductStatus = 1
)

type User struct {
	Username  string      `json:"username" bson:"username"`
	Password  string      `json:"password" bson:"password"`
	Profile   ProfileType `json:"profile" bson:"profile"`
	AccountId string      `json:"account" bson:"account"`
}

type Account struct {
	Id      string  `json:"id" bson:"id"`
	Balance float32 `json:"balance" bson:"balance"`
}

type Product struct {
	Id             string         `json:"id" bson:"id"`
	Name           string         `json:"name" bson:"name"`
	Price          float32        `json:"price" bson:"price"`
	TotalAvailable float32        `json:"total_available" bson:"total_available"`
	TotalSold      float32        `json:"total_sold" bson:"total_sold"`
	Stocks         []ProductStock `json:"stocks" bson:"stocks"`
}

type ReturnProduct struct {
	Id             string  `json:"id" bson:"id"`
	Name           string  `json:"name" bson:"name"`
	Price          float32 `json:"price" bson:"price"`
	TotalAvailable float32 `json:"total_available" bson:"total_available"`
	TotalSold      float32 `json:"total_sold" bson:"total_sold"`
}

type ReturnProductF struct {
	Id             string  `json:"id" bson:"id"`
	Name           string  `json:"name" bson:"name"`
	Price          float32 `json:"price" bson:"price"`
	Quantity       int     `json:"quantity" bson:"quantity"`
	TotalAvailable float32 `json:"total_available" bson:"total_available"`
	TotalSold      float32 `json:"total_sold" bson:"total_sold"`
}

type ProductStock struct {
	Id             string        `json:"id" bson:"id"`
	Name           string        `json:"name" bson:"name"`
	Price          float32       `json:"price" bson:"price"`
	TotalAvailable float32       `json:"total_available" bson:"total_available"`
	TotalSold      float32       `json:"total_sold" bson:"total_sold"`
	Status         ProductStatus `json:"status" bson:"status"`
}

type ReceiptProduct struct {
	Id       string  `json:"id" bson:"id"`
	Quantity float32 `json:"quantity" bson:"quantity"`
}

type ReceiptStatus byte

const (
	ReceiptStatusOpened ReceiptStatus = 0
	ReceiptStatusClosed ReceiptStatus = 1
)

type Currency byte

// metoda de plata
type Payment struct {
	Currency Currency `json:"currency" bson:"currency"`
	Balance  float32  `json:"balance" bson:"balance"`
}

type MyId int

type Receipt struct {
	Id         MyId             `json:"id" bson:"id"`
	Products   []ReceiptProduct `json:"products" bson:"products"`
	TotalPrice float32          `json:"total" bson:"total"`
	Status     ReceiptStatus    `json:"status" bson:"status"`
}

type IdGenerator struct {
	Id MyId `json:"id" bson:"id"`
}

type Session struct {
	Token    string      `json:"token" bson:"token"`
	Username string      `json:"username" bson:"username"`
	Profile  ProfileType `json:"profile" bson:"profile"`
}

type ResponseStatus struct {
	Status bool `json:"status"`
}
