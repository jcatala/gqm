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
	"os"
	"strconv"
)



func genBot(key string)  *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(key)
	if err != nil{
		log.Panic(err)
	}
	return bot
}

func parseConfig(v bool)map[string]string {
	cfg,err := ini.Load("gqm.ini")
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
	cfg, err := ini.Load("gqm.ini")
	if err != nil{
		log.Fatalln(err)
	}
	cfg.Section("").Key("saved_chat_id").SetValue(s)
	cfg.SaveTo("gqm.ini")

}

func getNewChatId(bot *tgbotapi.BotAPI, saved string, v bool)  string{
	if v != false{
		fmt.Println("Getting updates")
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 10
	updates, err := bot.GetUpdates(u)
	if err != nil{
		fmt.Println(err)
		return saved
	}

	newChatId := updates[len(updates)-1].Message.Chat.ID
	newChatIdstr := strconv.Itoa(int(newChatId))
	if newChatIdstr != saved{
		// If there's different chat id from the incoming message, we need to update the config file
		updateIni(newChatIdstr)
	}

	return newChatIdstr
}

func sendMsgFollow(bot *tgbotapi.BotAPI, chatIdInt int64 ,v bool){

	reader := bufio.NewReader(os.Stdin)
	for {
		line, _, err := reader.ReadLine()
		if err != nil{
			break
		}
		fmt.Printf("Sending the following: %s\n", line)
		msg := tgbotapi.NewMessage(chatIdInt, string(line))
		msg.ParseMode = "MarkdownV2"
		bot.Send(msg)
	}


}

func sendMsgQuick(bot *tgbotapi.BotAPI, chatIdInt int64, v bool){
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil{
		log.Fatalln(err)
	}
	if v {
		fmt.Printf("Sending the following line: \n%s\n", data)
	}
	msg := tgbotapi.NewMessage(chatIdInt, string(data))
	msg.ParseMode = "MarkdownV2"
	bot.Send(msg)
}

func main() {
	// They're just pointers
	verbose := flag.Bool("verbose",false, "-verbose to be verbose")
	follow := flag.Bool("follow", false, "-follow to keep the stdin open ")
	flag.Parse()
	fmt.Printf("The value of test is %s\n", *verbose)


	m := parseConfig(*verbose)
	bot := genBot(m["apikey"])
	bot.Debug = true
	if *verbose != false{
		log.Printf("Authorizing account %s\n", bot.Self.UserName)
	}

	// We ask for a new chat id, and see if its need to be updated or not
	chatId := getNewChatId(bot,m["savedChatId"] , *verbose)
	chatIdInt, err := strconv.ParseInt(chatId, 10 , 64)
	if err != nil{
		log.Fatalln(err)
	}

	// Here we craft a new msg
	if *follow == false{
		sendMsgQuick(bot, chatIdInt, *verbose)
	}else {
		sendMsgFollow(bot, chatIdInt,  *verbose)
	}
	//msg := tgbotapi.NewMessage(chatIdInt, "`test from gou`")
	//msg.ParseMode = "MarkdownV2"
	//bot.Send(msg)


	//msg := tgbotapi.NewMessage(-454559526, "test from go")
	//bot.Send(msg)





}