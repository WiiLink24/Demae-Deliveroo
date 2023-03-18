package deliveroo

import (
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
