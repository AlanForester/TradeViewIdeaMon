package libs

import (
	"fmt"
	"github.com/AlexCollin/TradeViewIdeaMon/sql"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"strconv"
	"time"
)

const NotifyText = `
*Дата:* %s
*Символ:* %s
*Направление:* %s
*Заголовок:* %s
*Автор:* %s
Описание:
%s
`

type Telebot struct {
	Connect *tb.Bot
}

func (t *Telebot) Sender(ch chan sql.Post) {
	for post := range ch {
		i, err := strconv.ParseFloat(post.Date, 64)
		if err != nil {
			panic(err)
		}
		tm := time.Unix(int64(i), 0)
		NotifyTextMsg := fmt.Sprintf(NotifyText, tm.Format("2006-01-02"), post.Pair, post.Tp, post.Title, post.Author, post.Descr)
		var users []sql.User
		sql.DB.Find(&users)
		for _, user := range users {
			mess := &tb.Photo{File: tb.FromDisk(post.Image), Caption: NotifyTextMsg}
			_, _ = t.Connect.Send(user, mess, , tb.SendOptions{ParseMode: tb.ModeMarkdownV2})
		}
	}
}

func (t *Telebot) Start() {
	var err error
	t.Connect, err = tb.NewBot(tb.Settings{
		Token:  "1669602029:AAH20CYggKwpCbncssBSJ6gdvQn5HjfNOJA",
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	t.Connect.Handle("/start", func(m *tb.Message) {
		user := sql.User{}
		db := sql.DB.First(&user, "uid = ?", m.Sender.Recipient())
		if db.Error != nil {
			user.Uid = m.Sender.Recipient()
			user.Status = "active"
			user.IsBlocked = false
			user.Username = m.Sender.Username
			db = sql.DB.Model(user).Create(&user)
			if db.Error != nil {
				log.Printf("Error on user save: %v", db.Error)
			}
			_, _ = t.Connect.Send(user, "Вы успешно подписались", tb.SendOptions{ParseMode: tb.ModeMarkdownV2})
		}
	})

	t.Connect.Start()
}
