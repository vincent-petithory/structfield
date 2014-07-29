package structfield_test

import (
	"encoding/json"
	"fmt"
	"github.com/vincent-petithory/structfield"
	"log"
)

func ExampleTransform() {
	// In the context of a REST API:
	// we want to transform the friends field in a friends_url field
	// that contains the URL to the list of friends of the user.
	type User struct {
		Id      string `json:"id"`
		Name    string `json:"name"`
		Age     int    `json:"age"`
		Friends []User `json:"friends"`
	}
	user := User{
		Id:   "4fa654a",
		Name: "Lelouch",
		Age:  22,
		Friends: []User{
			{Id: "65de67a", Name: "Ringo", Age: 25},
			{Id: "942ab70", Name: "Vivi", Age: 28},
		},
	}

	userFriendsToURL := structfield.TransformerFunc(func(field string, value interface{}) (string, interface{}) {
		return field + "_url", fmt.Sprintf("https://some.api.com/users/%s/friends", user.Id)
	})

	m := structfield.Transform(user, map[string]structfield.Transformer{
		"friends": userFriendsToURL,
	})
	_, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("%s", m["friends_url"])
	// Output: https://some.api.com/users/4fa654a/friends
}
