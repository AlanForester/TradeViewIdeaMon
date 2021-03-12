package libs

import (
	"fmt"
	"github.com/AlexCollin/TradeViewIdeaMon/sql"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"net/url"
	"strconv"
	"time"
)

const NotifyText = `
<b>Дата:</b> %s
<b>Символ:</b> %s
<b>Направление:</b> %s
<b>Заголовок:</b> %s
<b>Автор:</b> %s
<b>Описание:</b>
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
		var rmupc *tb.ReplyMarkup
		for _, user := range users {
			text := NotifyTextMsg

			rmupc = &tb.ReplyMarkup{}
			path, err := url.Parse(post.Url)
			if err != nil {
				log.Printf("Url parse error: %v", err)
			}

			more := rmupc.URL("Читать полностью", fmt.Sprintf("https://trview-bot.com%s", path.Path))
			rmupc.Inline(rmupc.Row(more))
			//pOptions := tb.SendOptions{ReplyMarkup: rmupc}

			if len(post.Descr) > 200 {
				text = fmt.Sprintf("%s...", post.Descr[0:200])
			}

			var mess tb.Sendable
			if post.Video != "" {
				mess = &tb.Video{File: tb.FromDisk(post.Video), Caption: text}
			} else {
				mess = &tb.Photo{File: tb.FromDisk(post.Image), Caption: text, ParseMode: tb.ModeHTML}
			}
			_, err = t.Connect.Send(user, mess, rmupc, tb.ModeHTML)
			if err != nil {
				log.Printf("Error on send: %v", err)
			}
		}
	}
}

func (t *Telebot) Start() {
	var err error
	t.Connect, err = tb.NewBot(tb.Settings{
		//Token:  "1669602029:AAH20CYggKwpCbncssBSJ6gdvQn5HjfNOJA",
		Token:  "1656961529:AAH9dKZXphT75tK0ulfqkSdO15dJBpfsOJQ",
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
			_, _ = t.Connect.Send(user, "Вы успешно подписались")
		}
	})

	t.Connect.Start()
}
