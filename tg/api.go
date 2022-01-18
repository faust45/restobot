package tg

import (
	"bytes"
	"encoding/json"
	// "errors"
	"fmt"
	// "gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
)

type APImethod string

const (
	Token  = "2146847843:AAHuZP9eXE692kk4TYzsyb_DZ4b8bdE3irs"
	ApiURL = "https://api.telegram.org/bot%s/%s"

	APItext      APImethod = "sendMessage"
	APIsendPhoto APImethod = "sendPhoto"
	APIeditPhoto APImethod = "editMessageMedia"
)

type Update struct {
	Id      int            `json:"update_id"`
	Command *Command       `json:"message"`
	Query   *CallbackQuery `json:"callback_query"`
}

type Command struct {
	Id   int    `json:"message_id"`
	Text string `json:"text"`
	Date int    `json:"date"`
	From User   `json:"from"`
	Chat Chat   `json:"chat"`
}

type QMessage struct {
	Id   int  `json:"message_id,omitempty"`
	Chat Chat `json:"chat"`
}

type CallbackQuery struct {
	Message QMessage `json:"message"`
	From    User     `json:"from"`
	Data    string   `json:"data"`
}

type User struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserName  string `json:"username"`
	IsBot     bool   `json:"is_bot"`
}

type Chat struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	UserName  string `json:"username"`
	Type      string `json:"type"`
}

type Button struct {
	Text  string `json:"text"`
	Query string `json:"callback_data,omitempty"`
}


type Text string

type Photo struct {
	Photo   string
	Caption string
}

type Keyboard [][]Button
type ReplyMarkup struct {
	Keyboard       Keyboard `json:"keyboard,omitempty"`
	InlineKeyboard Keyboard `json:"inline_keyboard,omitempty"`
	Resize         bool     `json:"resize_keyboard,omitempty"`
}

type media struct {
	Type  string `json:"type"`
	Media string `json:"media"`
	Caption string `json:"caption"`
}

type message struct {
	Id          int         `json:"message_id,omitempty"`
	ChatId      int         `json:"chat_id"`
	Text        string      `json:"text,omitempty"`
	Photo       string      `json:"photo,omitempty"`
	Caption     string      `json:"photo,omitempty"`
	Media       media       `json:"media,omitempty"` // only for edit req
	ReplyMarkup ReplyMarkup `json:"reply_markup,omitempty"`
}

type Msg interface {
	toMessage() message
	toEditMessage() message
	sendMethod() APImethod
	editMethod() APImethod
}

func (m Text) toMessage() message {
	return message{Text: string(m)}
}

func (m Photo) toMessage() message {
	return message{Photo: m.Photo, Caption: m.Caption}
}

func (m Text) toEditMessage() message {
	return message{Text: string(m)}
}

func (m Photo) toEditMessage() message {
	return message{
		Media: media{
			Type: "photo",
			Media: m.Photo,
			Caption: m.Caption,
		},
	}
}

func (m Photo) sendMethod() APImethod {
	return APIsendPhoto
}

func (m Photo) editMethod() APImethod {
	return APIsendPhoto
}

func (m Text) sendMethod() APImethod {
	return APItext
}

func (m Text) editMethod() APImethod {
	return APItext
}

func (chat Chat) Send(msg Msg, markup ReplyMarkup) {
	m := msg.toMessage()
	m.ReplyMarkup = markup
	method := msg.sendMethod()
	m.ChatId = chat.Id
	
	m.send(method)
}

func (m QMessage) Edit(msg Msg, markup ReplyMarkup) {
	edit := msg.toMessage()
	method := msg.editMethod()
	edit.ChatId = m.Chat.Id
	edit.ReplyMarkup = markup
	
	edit.send(method)
}

func (msg message) send(method APImethod) *http.Response {
	json, _ := json.Marshal(msg)
	// resp := req(method, json)
	// fmt.Printf("\n%s\n\n", resp.Status)

	url := fmt.Sprintf(ApiURL, Token, method)
	resp, _ := http.Post(url, "application/json", bytes.NewBuffer(json))

	if resp.Body != nil {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("Response \n%s\n", body)
	}

	return resp
}
