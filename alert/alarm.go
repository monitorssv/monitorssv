package alert

import (
	"errors"
	"github.com/monitorssv/monitorssv/alert/discord"
	"github.com/monitorssv/monitorssv/alert/telegram"
	"strings"
)

var TestAlarmMsg = "Welcome to MonitorSSV!"

type Alarm interface {
	Send(msg string) error
	Platform() string
}

type AlarmType int

const (
	DiscordType AlarmType = iota
	TelegramType
)

func NewAlarm(alarmType int, AlarmChannel string) (Alarm, error) {
	var alarm Alarm
	switch AlarmType(alarmType) {
	case DiscordType:
		alarm = discord.NewDiscordClient(AlarmChannel)
	case TelegramType:
		channelInfos := strings.Split(AlarmChannel, ",")
		if len(channelInfos) != 2 {
			return nil, errors.New("invalid Telegram channel")
		}
		alarm = telegram.NewTelegramClient(channelInfos[0], channelInfos[1])
	default:
		return nil, errors.New("unknown alarm type")
	}

	return alarm, nil
}
