# prefixed-api-key

TODO: badge for godoc

Go implementation of [seamapi/prefixed-api-key](https://github.com/seamapi/prefixed-api-key).

## Usage

```
go get github.com/joemiller/prefixed-api-key
```

```go
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
	log.Println(key.Token())

	ok, err := apikey.CheckAPIKey(key.Token(), key.LongTokenHash())
	if err != nil {
		log.Fatal(err)
	}
	if ok {
		log.Println("token matches hash")
	}

	key, err = apikey.ParseAPIKey("mycompany_5TJMbnP3thd_DjzvCr9MQLaKcaMisJuyUntS7Jpk61ZMp")
	log.Println(key.LongTokenHash())
}
```