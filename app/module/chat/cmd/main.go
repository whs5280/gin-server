package main

import (
	"bufio"
	"fmt"
	"gin-server/app/module/chat/model"
	"gin-server/app/module/chat/service"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(model.WelcomeMessage)

	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == model.EXIT {
			fmt.Println(model.ExitMessage)
			break
		}

		if service.ValidLength(input) == false {
			fmt.Println(model.ValidMessage)
			break
		}

		// reply := service.Reply(input, model.NoStream)
		reply := service.Reply(input, model.ChatStream)
		fmt.Printf(model.ReplyMessage, reply)
	}
}
