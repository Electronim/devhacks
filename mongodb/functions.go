package mongodb

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os/exec"
	"time"
)

var Client *mongo.Client

type MongoDb struct {
	Url         string
	DbName      string
	Users       string
	Products    string
	Receipts    string
	Sessions    string
	IdGenerator string
	Accounts    string
}

var MyDb = MongoDb{
	Url:         "mongodb://localhost:27017",
	DbName:      "banking",
	Users:       "users",
	Products:    "products",
	Receipts:    "receipts",
	Sessions:    "sessions",
	IdGenerator: "id_generator",
	Accounts:    "accounts",
}

func Init() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI(MyDb.Url)
	var err error
	Client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to mongodb was successful")
}

func GetAccount(id string) (Account, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := Client.Database(MyDb.DbName).Collection(MyDb.Accounts)

	var account Account
	if err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&account); err != nil {
		return Account{}, err
	}

	return account, nil
}

func GetUser(username string) (User, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := Client.Database(MyDb.DbName).Collection(MyDb.Users)

	var user User
	if err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user); err != nil {
		return User{}, err
	}

	return user, nil
}

func GetSession(token string) (Session, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := Client.Database(MyDb.DbName).Collection(MyDb.Sessions)

	var session Session
	if err := collection.FindOne(ctx, bson.M{"token": token}).Decode(&session); err != nil {
		return Session{}, err
	}

	return session, nil
}

func GetProduct(id string) (Product, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := Client.Database(MyDb.DbName).Collection(MyDb.Products)

	var product Product
	if err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&product); err != nil {
		return Product{}, err
	}

	return product, nil
}

func GetReceipt(id int) (Receipt, error) {
	var receipt Receipt
	var err error

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := Client.Database(MyDb.DbName).Collection(MyDb.Receipts)

	if err = collection.FindOne(ctx, bson.M{"id": id}).Decode(&receipt); err != nil {
		return Receipt{}, err
	}

	return receipt, nil
}

func GetProductSecure(token string, id string) (Product, error) {
	if _, err := GetSession(token); err != nil {
		return Product{}, err
	}

	if product, err := GetProduct(id); err != nil {
		return Product{}, err
	} else {
		return product, err
	}
}

func InsertProduct(product Product) error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := Client.Database(MyDb.DbName).Collection(MyDb.Products)

	if _, err := collection.InsertOne(ctx, product); err != nil {
		return err
	}

	return nil
}

func InsertSession(session Session) error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := Client.Database(MyDb.DbName).Collection(MyDb.Sessions)

	if _, err := collection.InsertOne(ctx, session); err != nil {
		return err
	}

	return nil
}

func Login(username string, password string, profile ProfileType) (Session, error) {
	var session Session
	var user User

	var err error
	if user, err = GetUser(username); err != nil {
		return session, err
	}

	if user.Password != password {
		return session, errors.New("invalid password")
	}

	token, err := exec.Command("uuidgen").Output()
	if err != nil {
		return session, err
	}

	session = Session{
		Token:    string(token[0 : len(token)-2]),
		Username: username,
		Profile:  profile,
	}

	if err := InsertSession(session); err != nil {
		return Session{}, err
	}

	return session, nil
}

func AddProduct(token string, stock ProductStock) error {
	if _, err := GetSession(token); err != nil {
		return err
	}

	var product Product
	var err error

	if product, err = GetProduct(stock.Id); err != nil {
		// product does not exist

		var newProduct Product
		newProduct.Id = stock.Id
		newProduct.Name = stock.Name
		newProduct.Price = stock.Price
		newProduct.TotalAvailable = stock.TotalAvailable
		newProduct.TotalSold = 0
		newProduct.Stocks = append(newProduct.Stocks, stock)

		if err := InsertProduct(newProduct); err != nil {
			return err
		}
		return nil
	}

	product.TotalAvailable += stock.TotalAvailable
	product.Stocks = append(product.Stocks, stock)

	filter := bson.M{"id": product.Id}
	update := bson.M{"$set": product}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := Client.Database(MyDb.DbName).Collection(MyDb.Products)

	if _, err := collection.UpdateOne(ctx, filter, update); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func GenerateId() (MyId, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := Client.Database(MyDb.DbName).Collection(MyDb.IdGenerator)

	var id IdGenerator
	if err := collection.FindOne(ctx, bson.M{}).Decode(&id); err != nil {
		return -1, err
	}

	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	if _, err := collection.UpdateOne(ctx, bson.M{}, bson.M{"$set": bson.M{"id": id.Id + 1}}); err != nil {
		return -1, err
	}

	return id.Id, nil
}

func UpdateStock(id string, quantity float32) error {
	var product Product
	var err error
	if product, err = GetProduct(id); err != nil {
		return err
	}

	product.TotalAvailable -= quantity
	product.TotalSold += quantity

	var newStocks []ProductStock
	for _, stock := range product.Stocks {
		newStock := stock
		if quantity > newStock.TotalAvailable {
			quantity -= newStock.TotalAvailable
			newStock.TotalSold += newStock.TotalAvailable
			newStock.TotalAvailable = 0
		} else {
			newStock.TotalAvailable -= quantity
			newStock.TotalSold += quantity
			quantity = 0
		}

		if newStock.TotalAvailable == 0 {
			newStock.Status = ProductStatusSold
		}

		newStocks = append(newStocks, newStock)
	}

	product.Stocks = newStocks
	filter := bson.M{"id": id}
	update := bson.M{"$set": product}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := Client.Database(MyDb.DbName).Collection(MyDb.Products)

	if _, err := collection.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func CalculateTotalPrice(products []ReceiptProduct) (float32, error) {
	var total float32
	var err error

	for _, product := range products {
		var mock Product
		var err error
		if mock, err = GetProduct(product.Id); err != nil {
			return -1, err
		}

		price := mock.Price
		quantity := product.Quantity
		total += price * quantity
	}

	return total, err
}

func CreateReceipt(token string, recProducts []ReceiptProduct) (Receipt, error) {
	if _, err := GetSession(token); err != nil {
		return Receipt{}, err
	}

	var receipt Receipt
	var err error
	var id MyId
	if id, err = GenerateId(); err != nil {
		return Receipt{}, err
	}

	receipt.Id = id
	receipt.Products = recProducts
	if receipt.TotalPrice, err = CalculateTotalPrice(recProducts); err != nil {
		return Receipt{}, err
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := Client.Database(MyDb.DbName).Collection(MyDb.Receipts)

	if _, err := collection.InsertOne(ctx, receipt); err != nil {
		return Receipt{}, err
	}

	return receipt, nil
}

func UpdateReceipt(receipt Receipt) error {
	for _, product := range receipt.Products {
		if err := UpdateStock(product.Id, product.Quantity); err != nil {
			return err
		}
	}

	receipt.Status = ReceiptStatusClosed
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := Client.Database(MyDb.DbName).Collection(MyDb.Receipts)

	if _, err := collection.UpdateOne(ctx, bson.M{"id": receipt.Id}, bson.M{"$set": receipt}); err != nil {
		return err
	}

	return nil
}

func UpdateAccount(id string, balance float32) error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := Client.Database(MyDb.DbName).Collection(MyDb.Accounts)

	var account Account
	var err error
	if account, err = GetAccount(id); err != nil {
		return err
	}

	account.Balance += balance
	filter := bson.M{"id": id}
	update := bson.M{"$set": account}

	if _, err = collection.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func ConfirmReceipt(usernameFrom string, usernameTo string, id int) error {
	var userFrom User
	var userTo User
	var accountFrom Account
	var accountTo Account
	var receipt Receipt
	var err error

	if userFrom, err = GetUser(usernameFrom); err != nil {
		return err
	}
	if userTo, err = GetUser(usernameTo); err != nil {
		return err
	}
	if accountFrom, err = GetAccount(userFrom.AccountId); err != nil {
		return err
	}
	if accountTo, err = GetAccount(userTo.AccountId); err != nil {
		return err
	}
	if receipt, err = GetReceipt(id); err != nil {
		return err
	}

	// updating balances
	if err := UpdateAccount(accountFrom.Id, -receipt.TotalPrice); err != nil {
		return err
	}
	if err := UpdateAccount(accountTo.Id, receipt.TotalPrice); err != nil {
		return err
	}

	// updating products
	if err := UpdateReceipt(receipt); err != nil {
		return err
	}

	return nil
}
