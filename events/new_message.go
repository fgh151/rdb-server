package events

import "fmt"

type subscribers struct {
	list map[string][]string
}

var Subscribers = subscribers{}

func init() {
	Subscribers.list = make(map[string][]string)
}

func Subscribe(topic string, listener string) {

	var currentList []string

	if _, ok := Subscribers.list[topic]; ok {
	} else {
		currentList = []string{}
	}

	currentList = append(currentList, listener)

	Subscribers.list[topic] = currentList
}

func RegisterNewMessage(topic string, content interface{}) {

	if currentList, ok := Subscribers.list[topic]; ok {

		for _, listener := range currentList {

			fmt.Println(listener)
			fmt.Println(content)

		}
	}
}
