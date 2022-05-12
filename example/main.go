package main

import (
	"log"

	apikey "github.com/joemiller/prefixed-api-key"
)

func main() {
	generator, err := apikey.NewGenerator("mycompany")
	if err != nil {
		log.Fatal(err)
	}

	key, err := generator.GenerateAPIKey()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(key.Token())

	ok, err := apikey.CheckAPIKey(key.Token(), key.LongTokenHash())
	if err != nil {
		log.Fatal(err)
	}
	if ok {
		log.Println("token matches hash")
	}

	key, err = apikey.ParseAPIKey("mycompany_5TJMbnP3thd_DjzvCr9MQLaKcaMisJuyUntS7Jpk61ZMp")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(key.LongTokenHash())
}
