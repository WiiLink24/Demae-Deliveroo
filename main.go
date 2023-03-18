package main

import (
	"context"
	crypto "crypto/rand"
	"encoding/xml"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"math"
	"math/big"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var pool *pgxpool.Pool
var config *Config

func checkError(err error) {
	if err != nil {
		log.Fatalf("Demae Dominos server has encountered an error! Reason: %v\n", err)
	}
}

func main() {
	rawConfig, err := os.ReadFile("./config.xml")
	checkError(err)

	config = &Config{}
	err = xml.Unmarshal(rawConfig, config)
	checkError(err)

	// Before we do anything, init Sentry to capture all errors.
	err = sentry.Init(sentry.ClientOptions{
		Dsn:              config.SentryDSN,
		Debug:            true,
		TracesSampleRate: 1.0,
	})
	checkError(err)
	defer sentry.Flush(2 * time.Second)

	// Initialize database
	dbString := fmt.Sprintf("postgres://%s:%s@%s/%s", config.SQLUser, config.SQLPass, config.SQLAddress, config.SQLDB)
	dbConf, err := pgxpool.ParseConfig(dbString)
	checkError(err)
	pool, err = pgxpool.ConnectConfig(context.Background(), dbConf)
	checkError(err)

	// Ensure this Postgresql connection is valid.
	defer pool.Close()

	// Seed random number generator
	seed, err := crypto.Int(crypto.Reader, big.NewInt(math.MaxInt64))
	checkError(err)
	rand.Seed(seed.Int64())

	// Finally, initialize the HTTP server
	fmt.Printf("Starting HTTP connection (%s)...\nNot using the usual port for HTTP?\nBe sure to use a proxy, otherwise the Wii can't connect!\n", config.Address)
	r := NewRoute()
	nwapi := r.HandleGroup("nwapi.php")
	{
		nwapi.NormalResponse("webApi_document_template", documentTemplate)
		nwapi.MultipleRootNodes("webApi_category_list", categoryList)
		nwapi.NormalResponse("webApi_area_shopinfo", func(r *Response) {})
		nwapi.NormalResponse("webApi_shop_list", shopList)
		nwapi.MultipleRootNodes("webApi_shop_one", shopOne)
		nwapi.MultipleRootNodes("webApi_menu_list", menuList)
		nwapi.MultipleRootNodes("webApi_item_list", itemList)
		nwapi.MultipleRootNodes("webApi_item_one", itemOne)
		nwapi.MultipleRootNodes("webApi_Authkey", authKey)
		nwapi.MultipleRootNodes("webApi_basket_add", basketAdd)
		nwapi.MultipleRootNodes("webApi_basket_list", basketList)
	}

	log.Fatal(http.ListenAndServe(config.Address, r.Handle()))
}
