package utility

import (
	"bufio"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"log"
	"errors"
)

func GetConfigDir()(string, error){
	usr, err := os.UserHomeDir()
	if err != nil{
		return "", err
	}
	path := filepath.Join( usr, ".config/gqm")
	if _, err := os.Stat(path); !os.IsNotExist(err){
		// Path exists
		return path, nil
	}
	log.Fatalln("ERROR: Config directory does not exist!\nExiting... ")
	return "", errors.New("ERROR: Config directory does not exist!\nExiting...")
}


func GenBot(key string)  *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(key)
	if err != nil{
		log.Panic(err)
	}
	return bot
}

func ParseConfig(v bool)map[string]string {
	configPath, err := GetConfigDir()
	if err != nil{
		log.Fatalln(err)
	}

	cfgPath := filepath.Join(configPath, "gqm.ini")
	cfg,err := ini.Load(cfgPath)
	if err != nil{
		log.Fatalln(err)
	}
	apikey := cfg.Section("DEFAULT").Key("apikey").String()
	savedChatid := cfg.Section("DEFAULT").Key("saved_chat_id").String()
	if v == true{
		fmt.Printf("Api key: %s\nSaved chat id : %s\n", apikey, savedChatid)
	}
	//cfg.Section("DEFAULT").Key("saved_chat_id").SetValue("1234")
	//cfg.SaveTo("gqm.ini")
	m := make(map[string]string)
	m["apikey"] = apikey
	m["savedChatId"] = savedChatid
	return m

}

func updateIni(s string){
	configPath, err := GetConfigDir()
	if err != nil{
		log.Fatalln(err)
	}
	cfgPath := filepath.Join(configPath, "gqm.ini")
	cfg,err := ini.Load(cfgPath)
	if err != nil{
		log.Fatalln(err)
	}
	cfg.Section("DEFAULT").Key("saved_chat_id").SetValue(s)
	cfg.SaveTo(cfgPath)

}

func GetNewChatId(bot *tgbotapi.BotAPI, saved string, v bool)  (string, error){
	if v != false{
		fmt.Println("Getting updates")
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 5
	updates, err := bot.GetUpdates(u)
	if err != nil{
		fmt.Println(err)
		return saved, nil
	}

	if len(updates) < 1 && saved == "" {
		// if there's no updates, and the saved is not valid, we must throw an error
		return "", errors.New("ERROR: No update, and no chat ID valid, EXITING...")
	}
	if len(updates) < 1 && saved != "" {
		// If there's no update, and the saved is valid, try to send to the saved
		return saved, nil
	}

	newChatId := updates[len(updates)-1].Message.Chat.ID
	newChatIdstr := strconv.Itoa(int(newChatId))
	if newChatIdstr != saved{
		// If there's different chat id from the incoming message, we need to update the config file
		updateIni(newChatIdstr)
	}

	return newChatIdstr, nil
}

func escapeTelegramMsg(s string)(string, error) {
	//s = url.QueryEscape(s)
	s = strings.ReplaceAll(s, "-", `\-`)
	s = strings.ReplaceAll(s, "+", `\+`)
	s = strings.ReplaceAll(s, ".", `\.`)
	s = strings.ReplaceAll(s, "#", `\#`)

	return s, nil
}

func SendMsgFollow(bot *tgbotapi.BotAPI, chatIdInt int64 , md bool, v bool){

	reader := bufio.NewReader(os.Stdin)
	for {
		line, _, err := reader.ReadLine()
		if err != nil{
			break
		}
		fmt.Printf("Sending the following: %s\n", line)
		var str strings.Builder
		dataStr := string(line)
		if md {
			str.WriteString("`")
			str.WriteString(string(line))
			str.WriteString("`")
			dataStr = str.String()
		}
		dataEscaped, _:= escapeTelegramMsg(dataStr)
		msg := tgbotapi.NewMessage(chatIdInt, dataEscaped)
		msg.ParseMode = "MarkdownV2"
		bot.Send(msg)
	}


}

func SendMsgPredefined(bot *tgbotapi.BotAPI, chatIdInt int64, md bool,  v bool, content []byte){

	data := string(content)
	if v {
		fmt.Printf("Sending the following line: \n%s\n", data)
	}
	var str strings.Builder
	dataStr := string(data)
	if md {
		str.WriteString("```\n")
		str.WriteString(string(data))
		str.WriteString("\n```")
		dataStr = str.String()
	}
	dataEscaped, _ := escapeTelegramMsg(dataStr)
	msg := tgbotapi.NewMessage(chatIdInt, dataEscaped)
	msg.ParseMode = "MarkdownV2"
	bot.Send(msg)
}

func SendMsgQuick(bot *tgbotapi.BotAPI, chatIdInt int64, md bool,  v bool){

	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil{
		log.Fatalln(err)
	}
	if v {
		fmt.Printf("Sending the following line: \n%s\n", data)
	}
	var str strings.Builder
	dataStr := string(data)
	if md {
		str.WriteString("```\n")
		str.WriteString(string(data))
		str.WriteString("\n```")
		dataStr = str.String()
	}
	dataEscaped, _ := escapeTelegramMsg(dataStr)
	msg := tgbotapi.NewMessage(chatIdInt, dataEscaped)
	msg.ParseMode = "MarkdownV2"
	bot.Send(msg)
}

