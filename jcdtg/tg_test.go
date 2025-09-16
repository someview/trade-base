package jcdtg

import (
	"testing"
	"time"

	"golang.org/x/net/context"
)

// token: "7684416419:AAHEvVAD6n_Gd03D0L7bkPVfy8HrMSK9lRs"
// chat_id_statistics: -4798746963
// chat_id_warning:  -4642910006
// chat_id_report:  -4504294083

func TestTgSendMessage(t *testing.T) {
	conf := &TGConfig{
		Token:             "7684416419:AAHEvVAD6n_Gd03D0L7bkPVfy8HrMSK9lRs",
		ChatID_Statistics: -4642910006,
		ChatID_Warning:    -4798746963,
		ChatID_Report:     -4504294083,
	}
	// 添加测试代码
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	TGInit(ctx, false, "", conf.Token, conf.ChatID_Statistics, conf.ChatID_Warning, conf.ChatID_Report)
	TGWarning("hello world Warning")
	TGStatistics("hello world statistics")
	TGRoundReport("1234567879 report")
	time.Sleep(time.Second * 5)
}
