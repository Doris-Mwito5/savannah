package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// ---------- Generic Random Helpers ----------

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(8))
}

func RandomPhoneNumber() string {
	return fmt.Sprintf("+2547%d", RandomInt(10000000, 99999999))
}

// ---------- Customers ----------

func RandomCustomerType() string {
	types := []string{"individual", "business"}
	return types[rand.Intn(len(types))]
}

func RandomCustomer() (name, email, phone, customerType string) {
	name = RandomString(6)
	email = RandomEmail()
	phone = RandomPhoneNumber()
	customerType = RandomCustomerType()
	return
}

// ---------- Categories ----------

func RandomCategory() (name string, parentID *int64) {
	name = RandomString(8)
	if rand.Intn(2) == 0 { // 50% chance of having a parent
		id := RandomInt(1, 100) // adjust max ID based on test data
		parentID = &id
	}
	return
}

// ---------- Products ----------

func RandomProductType() string {
	types := []string{"goods", "services"}
	return types[rand.Intn(len(types))]
}

func RandomProduct(categoryID int64) (name, description string, wholesalePrice, retailPrice float64, prodType string) {
	name = RandomString(8)
	description = fmt.Sprintf("Description for %s", name)
	wholesalePrice = float64(RandomInt(100, 1000))
	retailPrice = wholesalePrice * 1.5
	prodType = RandomProductType()
	return
}

// ---------- Orders ----------

func RandomOrderStatus() string {
	statuses := []string{"pending", "paid", "cancelled", "returned"}
	return statuses[rand.Intn(len(statuses))]
}

func RandomOrderSource() string {
	sources := []string{"offline", "online"}
	return sources[rand.Intn(len(sources))]
}

func RandomPaymentMethod() string {
	methods := []string{"cash", "mpesa", "card"}
	return methods[rand.Intn(len(methods))]
}

func RandomReferenceNumber() string {
	return fmt.Sprintf("REF-%d", RandomInt(100000, 999999))
}

func RandomOrder(customerID int64, shopID string) (refNum, phone, status, source, payment string, totalItems int, totalAmount float64) {
	refNum = RandomReferenceNumber()
	phone = RandomPhoneNumber()
	status = RandomOrderStatus()
	source = RandomOrderSource()
	payment = RandomPaymentMethod()
	totalItems = int(RandomInt(1, 10))
	totalAmount = float64(RandomInt(100, 10000)) / 100.0
	return
}

// ---------- Order Items ----------

func RandomOrderItem(orderID, productID int64) (unitPrice float64, quantity int, totalAmount float64) {
	unitPrice = float64(RandomInt(100, 1000)) / 10.0
	quantity = int(RandomInt(1, 5))
	totalAmount = unitPrice * float64(quantity)
	return
}
