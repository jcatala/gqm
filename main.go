package main

import (
	"flag"
	"github.com/jcatala/gqm/utility"
	"log"
	"strconv"
)



func main() {
	const version = "0.0.3"
	// rewrite the usage function
	flag.Usage = func() {
		utility.Banner(version)
		flag.PrintDefaults()
	}
	// They're just pointers
	verbose := flag.Bool("verbose",false, "To be verbose")
	follow := flag.Bool("follow", false, "To keep the stdin open ")
	md := flag.Bool("markdown", false, "Force markdown on the entire message, if is not, do it by yourself adding backquotes")
	debugInfo := flag.Bool("debugInfo", false, "To get debug information")
	message := flag.String("message", "", "To send a message instead of using the stdin")
	flag.Parse()

	// If verbose, print the banner
	if *verbose{
		utility.Banner(version)
	}


	m := utility.ParseConfig(*verbose)
	bot := utility.GenBot(m["apikey"])
	if *debugInfo{
		bot.Debug = true
	}
	if *verbose != false{
		log.Printf("Authorizing account %s\n", bot.Self.UserName)
	}

	// We ask for a new chat id, and see if its need to be updated or not
	chatId, err := utility.GetNewChatId(bot,m["savedChatId"] , *verbose)
	if err != nil{
		log.Fatalln(err)
	}
	chatIdInt, err := strconv.ParseInt(chatId, 10 , 64)
	if err != nil{
		log.Fatalln(err)
	}

	// Here we craft a new msg
	if *message != "" {
		messageBytes := []byte(*message)
		utility.SendMsgPredefined(bot, chatIdInt, *md, *verbose, messageBytes)
		return
	}
	if *follow == false{
		utility.SendMsgQuick(bot, chatIdInt, *md, *verbose)
	} else {
		// We need to generate a new bot object since the "following" options, ends messed up after a couple of lines
		//utility.SendMsgFollow(bot, chatIdInt, *md,   *verbose)
		utility.SendMsgFollow2(m["apikey"], chatIdInt, *md, *verbose)
	}
	//msg := tgbotapi.NewMessage(chatIdInt, "`test from gou`")
	//msg.ParseMode = "MarkdownV2"
	//bot.Send(msg)


	//msg := tgbotapi.NewMessage(-454559526, "test from go")
	//bot.Send(msg)





}