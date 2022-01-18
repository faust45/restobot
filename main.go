package main

import (
	// "menubot/coll"
	"errors"
	"menubot/tg"
	"strings"
	// "bytes"
	"encoding/json"
	"fmt"
	// "strings"
	// "gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
)

type Menu []MenuItem
type MenuItem struct {
	Title string `json:"title"`
	Price int    `json:"price"`
	Group string `json:"group"`
}

type Query struct {
	Cmd   string
	Value string
}

var (
	Photo1 = "https://www.onceuponachef.com/images/2019/07/Big-Italian-Salad.jpg"
	Photo  = "https://c8.alamy.com/comp/2C4DMYH/vintage-poster-menu-collection-retro-template-design-menu-restaurant-or-diner-2C4DMYH.jpg"
	menu   Menu
)

func init() {
	menuFile, _ := os.Open("data/bot_conf.json")
	defer menuFile.Close()

	bytesData, _ := ioutil.ReadAll(menuFile)
	_ = json.Unmarshal([]byte(bytesData), &menu)

	// fmt.Printf("Json menu err: %s\n menu: %+v", err, menu)
}

func main() {
	// a := http.HandleFunc("/", parseJsonBody)
	mux := http.NewServeMux()
	handler := http.HandlerFunc(handleReq)
	mux.Handle("/", handler)

	port := "3000"
	fmt.Printf("Listening on port %s", port)

	http.ListenAndServe(":"+port, mux)
}

func handleReq(w http.ResponseWriter, r *http.Request) {
	if r.Body != nil {
		body, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()

		var msg tg.Msg
		var markup tg.ReplyMarkup
		query, reply, _ := ParseUpdate(body)

		switch query.Cmd {
		case "start":
			msg = tg.Text("Приветртвуем вас в кафе Буратино")
			markup = Keyboard(
				)
		case "menu":
			msg, markup = menu.Show()
		case "group":
			// msg, markup = menu.ShowGroup(query.Value)
		case "item":
			// msg, markup = menu.ShowItem(query.Value)
		}

		reply(msg, markup)
	}
}

type ReplyFn func(tg.Msg, tg.ReplyMarkup)

func Keyboard(buttons ...tg.Button) tg.ReplyMarkup {
	var keyboard tg.Keyboard

	for _, b := range buttons {
		keyboard = append(keyboard, []tg.Button{b})
	}
	
	return tg.ReplyMarkup{Keyboard: keyboard}
}

func ParseUpdate(bytes []byte) (*Query, ReplyFn, error) {
	var update tg.Update
	err := json.Unmarshal(bytes, &update)
	if err != nil {
		fmt.Printf("JSON parse error: %v", err)
		return nil, nil, err
	}
	// fmt.Printf("\n incoming req \n%+v\n%+v\n",
	// update.Command, update.Query)

	if update.Command != nil {
		reply := func(msg tg.Msg, markup tg.ReplyMarkup) {
			update.Command.Chat.Send(msg, markup)
		}

		return &Query{update.Command.Text, ""}, reply, nil
	} else if update.Query != nil {
		q, err := ParseQuery(update.Query.Data)
		if err != nil {
			return nil, nil, err
		}

		reply := func(msg tg.Msg, markup tg.ReplyMarkup) {
			update.Query.Message.Edit(msg, markup)
		}

		return &q, reply, nil

	} else {
		return nil, nil, errors.New("Update parse error Command & Query both are nil")
	}
}

func (menu Menu) Show() (tg.Photo, tg.ReplyMarkup) {
	groups := func(a MenuItem) string {
		return a.Group
	}

	items := Unique(menu.Map(groups))
	key := Keyboard1("group", items)
	return tg.Photo{Photo: Photo}, InlineKeyboard(key)
}

func InlineKeyboard(keyboard tg.Keyboard) tg.ReplyMarkup {
	return tg.ReplyMarkup{}
}

func (menu Menu) ShowGroup(group string) tg.Photo {
	byGroup := func(item MenuItem) bool {
		return item.Group == group
	}

	titles := func(a MenuItem) string {
		return a.Title
	}

	_ = menu.Search(byGroup).Map(titles)
	// back := tg.Button{"<< назад", "menu"}
	// keyboard := append(Keyboard("item", items), []tg.Button{back})

	return tg.Photo{
		Photo: Photo1,
	}
}

func (items Menu) Map(f func(item MenuItem) string) []string {
	var arr []string
	for _, el := range items {
		arr = append(arr, f(el))
	}

	return arr
}

func Unique(items []string) []string {
	var keys []string
	set := make(map[string]bool)

	for _, el := range items {
		if !set[el] {
			set[el] = true
			keys = append(keys, el)
		}
	}

	return keys

}

func Keyboard1(q string, items []string) tg.Keyboard {
	var keyboard tg.Keyboard
	for _, el := range items {
		q := fmt.Sprintf("%s=%s", q, el)
		keyboard = append(keyboard,
			[]tg.Button{
				tg.Button{el, q}})
	}

	return keyboard
}

func (menu Menu) ShowItem(title string) tg.Photo {
	// items := menu.Search(func(item MenuItem) bool {
	// 	return item.Title == title
	// })

	return tg.Photo{}
}

func (menu Menu) Search(f func(item MenuItem) bool) Menu {
	var arr Menu

	for _, el := range menu {
		if f(el) {
			arr = append(arr, el)
		}
	}

	return arr
}

func ParseQuery(q string) (Query, error) {
	var query Query
	arr := strings.Split(q, "=")

	if 0 < len(arr) {
		key := arr[0]
		value := arr[1]
		query := Query{key, value}

		return query, nil
	} else {
		return query, errors.New("could't parse query")
	}
}

// func ParseQuery(q string) (Query, error) {
// 	var query Query
// 	re, _ := regexp.Compile(`^(\w+)\?(\w+)=(\w+)`)
// 	matched := re.FindAllStringSubmatch(q, -1)

// 	if 0 < len(matched) {
// 		cmd := matched[0][1]
// 		params := matched[0][2:]
// 		query := Query{Command: cmd, Params: make(map[string]string)}

// 		for i := 0; i < len(params); i = i + 2 {
// 			k := params[i]
// 			v := params[i+1]
// 			query.Params[k] = v
// 		}

// 		return query, nil
// 	} else {
// 		return query, errors.New("could't parse query")
// 	}
// }
