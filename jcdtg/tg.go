package jcdtg

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/someview/trade-base/config"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const TGConfigKey = "tg_config"

type TGConfig struct {
	Token             string `mapstructure:"token"`
	ChatID_Warning    int64  `mapstructure:"chat_id_warning"`
	ChatID_Statistics int64  `mapstructure:"chat_id_statistics"`
	ChatID_Report     int64  `mapstructure:"chat_id_report"`
}

func (T *TGConfig) CheckValidity() error {
	return nil
}

var _ config.Configer = (*TGConfig)(nil)

// 一个 Bot 每秒不能发出超过 30 条消息，在同一个群组中每分钟不能发出超过 20 条消息,3s 发送一次
var (
	infoWarningTgC                = make(chan string, 1000)
	statisticsTgC                 = make(chan string, 1000)
	roundreportTgC                = make(chan string, 1000)
	tg_infoWarningChatID          = int64(0)
	tg_statisticsChatID           = int64(0)
	tg_reportChatID               = int64(0)
	sendMsgBot           *bot.Bot = nil
)

const ln string = "\r\n"

func sendMessage(chatid int64, msg string) {
	msglen := len(msg)
	if msglen == 0 {
		return
	}
	//一次发送的大小不能超过4096
	if msglen > 4096 {
		arr := strings.Split(msg, ln)
		seg := ""
		next := ""
		arrlen := len(arr)
		for i := 0; i < arrlen; i++ {
			next = seg
			next = next + arr[i] + ln
			if len(next) >= 4096 {
				sendMsgBot.SendMessage(context.Background(), &bot.SendMessageParams{
					ChatID: chatid,
					Text:   seg,
				})
				next = ""
				i--
			}
			seg = next
		}
		sendMsgBot.SendMessage(context.Background(), &bot.SendMessageParams{
			ChatID: chatid,
			Text:   seg,
		})

	} else {
		sendMsgBot.SendMessage(context.Background(), &bot.SendMessageParams{
			ChatID: chatid,
			Text:   msg,
		})
	}
}
func handleMessage(ctx context.Context, chatChan chan string, chatID int64, flag string) {
	msgBuilder := strings.Builder{}
	ticker := time.NewTicker(time.Millisecond * 3100) //在同一个群组中每分钟不能发出超过 20 条消息,3s 发送一次
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Println("jcdlog:", flag, "ctx.Done")
			for {
				if len(chatChan) == 0 {
					break
				}
				msg := <-chatChan
				sendMsgBot.SendMessage(context.Background(), &bot.SendMessageParams{
					ChatID: chatID,
					Text:   msg,
				})
				time.Sleep(time.Millisecond * 20)
			}
			log.Println("jcdlog:", flag, "end")
			return
		case msg := <-chatChan:
			msgBuilder.WriteString(msg)
		case <-ticker.C:
			msg := msgBuilder.String()
			msgBuilder.Reset()
			sendMessage(chatID, msg)
		}
	}
}
func TGInit(ctx context.Context, isProxy bool, ProxyURL, token string, infoChatID, statisticsChatID, reportChatID int64) {
	client := &http.Client{}
	if isProxy {
		proxy_url, _ := url.Parse(ProxyURL)
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxy_url),
		}
	}
	var err error
	//尝试3次
	for i := 0; i < 3; i++ {
		sendMsgBot, err = bot.New(token, bot.WithHTTPClient(10*time.Second, client))
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	if err != nil {
		log.Panic(err)
	}
	tg_infoWarningChatID = infoChatID
	tg_statisticsChatID = statisticsChatID
	tg_reportChatID = reportChatID

	//info
	go handleMessage(ctx, infoWarningTgC, infoChatID, "infoWarningTgC")

	//staistics
	go handleMessage(ctx, statisticsTgC, statisticsChatID, "statisticsTgC")
	//roundreport
	go handleMessage(ctx, roundreportTgC, reportChatID, "roundreportTgC")

}

func TGSyncWarning(msg string) {
	if sendMsgBot == nil {
		return
	}
	ts := msg + ln + time.Now().UTC().Format("2006-01-02T15:04:05.000000Z") + ln
	sendMsgBot.SendMessage(context.Background(), &bot.SendMessageParams{
		ChatID: tg_infoWarningChatID,
		Text:   ts,
	})
}

func TGSyncReport(msg string) {
	if sendMsgBot == nil {
		return
	}
	ts := msg + ln + time.Now().UTC().Format("2006-01-02T15:04:05.000000Z") + ln
	sendMsgBot.SendMessage(context.Background(), &bot.SendMessageParams{
		ChatID: tg_reportChatID,
		Text:   ts,
	})
}

func TGSyncStatistics(msg string) {
	if sendMsgBot == nil {
		return
	}
	ts := msg + ln + time.Now().UTC().Format("2006-01-02T15:04:05.000000Z") + ln
	sendMsgBot.SendMessage(context.Background(), &bot.SendMessageParams{
		ChatID: tg_statisticsChatID,
		Text:   ts,
	})
}

func TGWarning(msg string) {
	if sendMsgBot == nil {
		return
	}
	ts := msg + ln + time.Now().UTC().Format("2006-01-02T15:04:05.000000Z") + ln
	infoWarningTgC <- ts
}
func TGStatistics(msg string) {
	if sendMsgBot == nil {
		return
	}
	ts := msg + ln + time.Now().UTC().Format("2006-01-02T15:04:05.000000Z") + ln
	statisticsTgC <- ts
}

func TGStatisticsWithoutTime(msg string) {
	if sendMsgBot == nil {
		return
	}
	ts := msg + ln
	statisticsTgC <- ts
}

func TGRoundReport(msg string) {
	if sendMsgBot == nil {
		return
	}
	ts := msg + ln + time.Now().UTC().Format("2006-01-02T15:04:05.000000Z") + ln
	roundreportTgC <- ts
}

func TGHTMLReport(msg string) error {
	if sendMsgBot == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, err := sendMsgBot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    tg_reportChatID,
		Text:      msg, // 对特殊字符进行转义
		ParseMode: models.ParseModeHTML,
	})
	return err
}

func TGHTMLWarn(msg string) error {
	if sendMsgBot == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, err := sendMsgBot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    tg_infoWarningChatID,
		Text:      msg, // 对特殊字符进行转义
		ParseMode: models.ParseModeHTML,
	})
	return err
}
