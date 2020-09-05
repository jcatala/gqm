package main

/*

usage: sio-notification.py [-h] [-c CONFIG] [-i INFILE] [-v] [-F FILTER] [-f] [-m]

Simple Shout-it-out telegram notificator

optional arguments:
  -h, --help            show this help message and exit
  -c CONFIG, --config CONFIG
                        Full path to config file (default is ~/.SIO.conf
  -i INFILE, --infile INFILE
                        Send a text file (default is stdin)
  -v, --verbose         Turn on the verbose mode
  -F FILTER, --filter FILTER
                        Add a filter before sending the message (string: default: None)
  -f, --follow          Send one line at a time
  -m, --markdown        Force markdown on the entire message, if is not, do it by yourself adding backquotes


*/

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
	"gopkg.in/ini.v1"
	"net/url"
	"os"
	"strconv"
	"strings"
	"path/filepath"
	"errors"
)


func getConfigDir()(string, error){
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


func genBot(key string)  *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(key)
	if err != nil{
		log.Panic(err)
	}
	return bot
}

func parseConfig(v bool)map[string]string {
	configPath, err := getConfigDir()
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
	configPath, err := getConfigDir()
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

func getNewChatId(bot *tgbotapi.BotAPI, saved string, v bool)  (string, error){
	if v != false{
		fmt.Println("Getting updates")
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 10
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

func sendMsgFollow(bot *tgbotapi.BotAPI, chatIdInt int64 , md bool, v bool){

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
		dataEscaped := url.QueryEscape(dataStr)
		msg := tgbotapi.NewMessage(chatIdInt, dataEscaped)
		msg.ParseMode = "MarkdownV2"
		bot.Send(msg)
	}


}

func sendMsgQuick(bot *tgbotapi.BotAPI, chatIdInt int64, md bool,  v bool){
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

	msg := tgbotapi.NewMessage(chatIdInt, dataStr)
	msg.ParseMode = "MarkdownV2"
	bot.Send(msg)
}

func main() {
	// They're just pointers
	verbose := flag.Bool("verbose",false, "To be verbose")
	follow := flag.Bool("follow", false, "To keep the stdin open ")
	md := flag.Bool("markdown", false, "Force markdown on the entire message, if is not, do it by yourself adding backquotes")
	debugInfo := flag.Bool("debugInfo", false, "To get debug information")
	flag.Parse()



	m := parseConfig(*verbose)
	bot := genBot(m["apikey"])
	if *debugInfo{
		bot.Debug = true
	}
	if *verbose != false{
		log.Printf("Authorizing account %s\n", bot.Self.UserName)
	}

	// We ask for a new chat id, and see if its need to be updated or not
	chatId, err := getNewChatId(bot,m["savedChatId"] , *verbose)
	if err != nil{
		log.Fatalln(err)
	}
	chatIdInt, err := strconv.ParseInt(chatId, 10 , 64)
	if err != nil{
		log.Fatalln(err)
	}

	// Here we craft a new msg
	if *follow == false{
		sendMsgQuick(bot, chatIdInt, *md, *verbose)
	}else {
		sendMsgFollow(bot, chatIdInt, *md,   *verbose)
	}
	//msg := tgbotapi.NewMessage(chatIdInt, "`test from gou`")
	//msg.ParseMode = "MarkdownV2"
	//bot.Send(msg)


	//msg := tgbotapi.NewMessage(-454559526, "test from go")
	//bot.Send(msg)





}