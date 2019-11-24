package server

import (
	"banking/mongodb"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
)

func TestHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	type tmp struct {
		Name string `json:"name"`
	}

	var query tmp
	if err = json.Unmarshal(body, &query); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	var array []interface{}
	var user mongodb.User
	var product mongodb.Product
	var receipt mongodb.Receipt
	var productStock mongodb.ProductStock
	var payment mongodb.Payment

	type Mc struct {
		Name  string      `json:"name"`
		Field interface{} `json:"example"`
	}

	array = append(array, Mc{"User", user})
	array = append(array, Mc{"Product", product})
	array = append(array, Mc{"Product", receipt})
	array = append(array, Mc{"Product", productStock})
	array = append(array, Mc{"Product", payment})

	if err := json.NewEncoder(res).Encode(array); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
}

func LoginHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	type tmp struct {
		Username string              `json:"username"`
		Password string              `json:"password"`
		Profile  mongodb.ProfileType `json:"profile"`
	}

	var query tmp
	if err = json.Unmarshal(body, &query); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	var session mongodb.Session
	if session, err = mongodb.Login(query.Username, query.Password, query.Profile); err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := json.NewEncoder(res).Encode(session); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
}

func ProductAdd(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	status := mongodb.ResponseStatus{Status: false}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(res).Encode(status); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		return
	}

	type tmp struct {
		Token        string               `json:"token"`
		ProductStock mongodb.ProductStock `json:"product_stock"`
	}

	var query tmp
	if err = json.Unmarshal(body, &query); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(res).Encode(status); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		return
	}

	if _, err := mongodb.GetSession(query.Token); err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := mongodb.AddProduct(query.Token, query.ProductStock); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(res).Encode(status)
		return
	}

	status.Status = true
	if err := json.NewEncoder(res).Encode(status); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
}

func ProductGet(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	type tmp struct {
		Token string `json:"token"`
		Id    string `json:"id"`
	}

	var query tmp
	if err = json.Unmarshal(body, &query); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := mongodb.GetSession(query.Token); err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	if product, err := mongodb.GetProductSecure(query.Token, query.Id); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	} else {
		ans := mongodb.ReturnProduct{
			Id:             product.Id,
			Name:           product.Name,
			Price:          product.Price,
			TotalAvailable: product.TotalAvailable,
			TotalSold:      product.TotalSold,
		}
		if err := json.NewEncoder(res).Encode(ans); err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func ReceiptCreate(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	type tmp struct {
		Token    string                   `json:"token" bson:"token"`
		Products []mongodb.ReceiptProduct `json:"products"`
	}

	var query tmp
	if err = json.Unmarshal(body, &query); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := mongodb.GetSession(query.Token); err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	if rec, err := mongodb.CreateReceipt(query.Token, query.Products); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	} else {
		if err := json.NewEncoder(res).Encode(rec); err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func ReceiptConfirm(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	status := mongodb.ResponseStatus{Status: false}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	type tmp struct {
		Token    string `json:"token" bson:"token"`
		Id       int    `json:"id"`
		UserFrom string `json:"user_from"`
		UserTo   string `json:"user_to"`
	}

	var query tmp
	if err = json.Unmarshal(body, &query); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(res).Encode(status)
		return
	}

	if _, err := mongodb.GetSession(query.Token); err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(res).Encode(status)
		return
	}

	if err := mongodb.ConfirmReceipt(query.UserFrom, query.UserTo, query.Id); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(res).Encode(status)
		return
	}

	status.Status = true
	_ = json.NewEncoder(res).Encode(status)
}

func ReceiptGet(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	type tmp struct {
		Token string `json:"token"`
		Id    int    `json:"id"`
	}

	var query tmp
	if err = json.Unmarshal(body, &query); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := mongodb.GetSession(query.Token); err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	var receipt mongodb.Receipt
	if receipt, err = mongodb.GetReceipt(query.Id); err != nil {
		fmt.Println(err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	type ans struct {
		Id         mongodb.MyId             `json:"id" bson:"id"`
		Products   []mongodb.ReturnProductF `json:"products" bson:"products"`
		TotalPrice float32                  `json:"total" bson:"total"`
		Status     mongodb.ReceiptStatus    `json:"status" bson:"status"`
	}

	var rsp ans
	rsp.Id = receipt.Id
	rsp.TotalPrice = receipt.TotalPrice
	rsp.Status = receipt.Status

	for _, obj := range receipt.Products {
		if prod, err := mongodb.GetProduct(obj.Id); err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			newProd := mongodb.ReturnProductF{
				Id:             prod.Id,
				Name:           prod.Name,
				Price:          prod.Price,
				Quantity:       int(math.Round(float64(obj.Quantity))),
				TotalAvailable: obj.Quantity,
				TotalSold:      prod.TotalSold,
			}
			rsp.Products = append(rsp.Products, newProd)
		}
	}

	if err := json.NewEncoder(res).Encode(rsp); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}
