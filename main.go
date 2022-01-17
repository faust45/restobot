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
		_, query, reply, _ := ParseUpdate(body)
	
		switch query.Cmd {
		case "start":
			msg = tg.Msg{
				Text: "Приветртвуем вас в кафе Буратино",
				// Keyboard: MainMenu()
			}
		case "menu":
			msg = menu.Show()
		case "group":
			msg = menu.ShowGroup(query.Value)
		case "item":
			msg = menu.ShowItem(query.Value)
		}

		reply(msg)
	}
}

type ReplyFn func(m tg.Msg) 

func ParseUpdate(bytes []byte) (*tg.Update, *Query, ReplyFn, error) {
	var update tg.Update
	err := json.Unmarshal(bytes, &update)
	if err != nil {
		fmt.Printf("JSON parse error: %v", err)
		return nil, nil, nil, nil
	}
	// fmt.Printf("\n incoming req \n%+v\n%+v\n",
		// update.Command, update.Query)

	if update.Command != nil {
		reply := func(msg tg.Msg) {
			update.Command.Chat.Send(msg)
		}
		
		return &update, &Query{update.Command.Text, ""}, reply, nil
	} else if update.Query != nil {
		q, err := ParseQuery(update.Query.Data)
		if err != nil {
			fmt.Println(err)
		}
		
		reply := func(msg tg.Msg) {
			update.Command.Chat.Send(msg)
		}
		
		return &update, &q, reply, nil

	} else {
		return nil, nil, nil, nil
	}
}

func (menu Menu) Show() tg.Msg {
	groups := func(a MenuItem) string {
		return a.Group
	}

	items := Unique(menu.Map(groups))
	// nav := []string{"<< назад", "menu"}

	return tg.Msg{
		Photo:          Photo,
		InlineKeyboard: Keyboard("group", items),
	}
}

func (menu Menu) ShowGroup(group string) tg.Msg {
	byGroup := func(item MenuItem) bool {
		return item.Group == group
	}

	titles := func(a MenuItem) string {
		return a.Title
	}

	items := menu.Search(byGroup).Map(titles)

	return tg.Msg{
		Photo:          Photo1,
		InlineKeyboard: Keyboard("item", items),
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

func Keyboard(q string, items []string) tg.Keyboard {
	var keyboard tg.Keyboard
	for _, el := range items {
		q := fmt.Sprintf("%s=%s", q, el)
		keyboard = append(keyboard,
			[]tg.Button{
				tg.Button{Text: el, Query: q}})
	}

	return keyboard
}

func (menu Menu) ShowItem(title string) tg.Msg {
	// items := menu.Search(func(item MenuItem) bool {
	// 	return item.Title == title
	// })

	return tg.Msg{}
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
