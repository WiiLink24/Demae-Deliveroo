package deliveroo

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"sync"
	"testing"
)

func TestDeliveroo_GetBareShops(t *testing.T) {
	d := Deliveroo{mutex: sync.RWMutex{}}
	_, err := d.GetBareShops()
	if err != nil {
		panic(err)
	}
}

func TestDeliveroo_GetShops(t *testing.T) {
	d := Deliveroo{mutex: sync.RWMutex{}}
	_, err := d.GetShops("01")
	if err != nil {
		panic(err)
	}
}

func TestDeliveroo_GetStore(t *testing.T) {
	d := Deliveroo{mutex: sync.RWMutex{}}
	_, err := d.GetStore("78615")
	if err != nil {
		panic(err)
	}
}

func TestDeliveroo_GetMenuCategories(t *testing.T) {
	d := Deliveroo{mutex: sync.RWMutex{}}
	_, err := d.GetMenuCategories("78615")
	if err != nil {
		panic(err)
	}
}

func TestDeliveroo_GetItems(t *testing.T) {
	d := Deliveroo{mutex: sync.RWMutex{}}
	_, err := d.GetItems("78615", "135990441")
	if err != nil {
		panic(err)
	}
}

func TestDecryption(t *testing.T) {
	encryptedAuth := "fb51b87381e6af7aa0530f1373be5d2b2a752054abe776d24abf2c6dda1e45516482cad9425d592954fac0075335a0622f783f5834f7c7c9d0ad9cf6c3e9f9845749a14c9e110b3253bfc1351bacc62e6747a9a7b101131ea80b43e335e7308a4fd6af42eba2b4f3415796f8e6f731003a868a6bc0fe1438cc3e3b704d5eb2dd454f4ee3e82b37bf2edf3099a2179860e4c7907a2dfc39bbfb5d68a6fad9204091f4624f1d14af38071d06421ba3e448e2453b4dc726754b6165a9a38bbc786979ee8c64d9a171d1a6a62c6230bf38c2be6934e236cbfa471f2bc2b4aaad9ba41bdd8dbe947cb98357f37930fe2a650269445dcb620631ae6c874847da34c66fd0b91409aabd99b42ac3afbcdbd1cb1cf1dcee9b611c22bb4601223223d3c31451b5310808662d37c3bcbf6451b8293279e108dce7a84780d1db982232e49a5578f956fade41601dead88579266b97f34d2954b765e9a49813bb6c71ef64240885a849233bb390418b02016f0cb0a2833f750f375d74ded0c6507674105435a699aaece014b4344b5d676672cf7c5073d7cfd912c0c825248c573b0da1b121d83db56d86341a615455bf7edddfaef91b4183f648d2b5b8c78be4681617c086fd917b7804580393e8a23074af7edeed0d204bde93a3a6534e781388011ad4ed5c6408a4261297901adafbe250a6256896d4840776013ea608685ab3d217a464137a2e913a9f76482beffce9336366b46beb92a8e76d26cf73285a03a61ffa4f93"
	block, err := aes.NewCipher([]byte(AuthKey))
	if err != nil {
		panic(err)
	}

	auth := make([]byte, len(encryptedAuth))
	cipherText, err := hex.DecodeString(encryptedAuth)
	if err != nil {
		panic(err)
	}

	mode := cipher.NewCBCDecrypter(block, []byte(AuthIV))
	mode.CryptBlocks(auth, cipherText)
	auth = bytes.Trim(auth, "\x00")
	auth, err = pkcs7Unpad(auth, aes.BlockSize)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(auth))
}

func TestDeliveroo_SendBasket(t *testing.T) {
	/*dbString := fmt.Sprintf("postgres://%s:%s@%s/%s", "noahpistilli", "2006", "127.0.0.1", "deliveroo")
	dbConf, _ := pgxpool.ParseConfig(dbString)
	pool, err := pgxpool.ConnectConfig(context.Background(), dbConf)
	if err != nil {
		panic(err)
	}

	var basketStr string
	row := pool.QueryRow(context.Background(), `SELECT "user".basket, "user".auth_key FROM "user" WHERE "user".wii_id = $1 LIMIT 1`, "610515068")
	err = row.Scan(&basketStr, nil)
	if err != nil {
		panic(err)
	}

	var mapBasket []map[string]any
	err = json.Unmarshal([]byte(basketStr), &mapBasket)
	if err != nil {
		panic(err)
	}

	d := Deliveroo{}
	total, subtotal, charge, _ := d.SendBasket("60819", mapBasket)

	fmt.Println(subtotal)
	fmt.Println(charge)
	fmt.Println(total)*/
}

func TestDeliveroo_CreatePaymentPlan(t *testing.T) {
	d := Deliveroo{
		mutex:     sync.RWMutex{},
		auth:      "Basic Njk0MjIxMTI6b3JkZXJhcHBfYW5kcm9pZCxleUpoYkdjaU9pSkZVekkxTmlJc0ltcHJkU0k2SW1oMGRIQnpPaTh2WkdWc2FYWmxjbTl2TG1OdkxuVnJMMmxrWlc1MGFYUjVMV3RsZVhNdk1TNXFkMnNpZlEuZXlKbGVIQWlPakUyTnpreU56SXdPVEVzSW1OMWMzUWlPalk1TkRJeU1URXlMQ0prY201ZmFXUWlPaUprT0RNeU5tWmxNUzAyWVdabExUUXpNelV0T0RCa1lpMDFNR1UzTVRaaE1tWmxOV0lpTENKelpYTnpJam9pYjNKa1pYSmhjSEJmWVc1a2NtOXBaQ3hsTTJZME5XUmpZVE5tT0dNME1UQm1PVGczTldJMlltUmxabVEwWkdabE1DSjkuUUJ6dm1rV1VqT21qQ3BhSkVpYkpPOTcyY1gwc0pTMFFRRGoxaFJsVzd0alN5X2Y2SFMyX2hKX1I0QmNMZ2dTbUFTQ2xyNWxBYVdKSmxwdlhhV3U1RHc=",
		longitude: 0,
		latitude:  0,
		userId:    "",
	}
	_, err := d.CreatePaymentPlan()
	if err != nil {
		panic(err)
	}
}
