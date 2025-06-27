package main

import (
	"bufio"
	"fmt"
	"gin-server/app/module/idiom/model"
	"gin-server/app/module/idiom/service"
	"os"
	"strings"
)

func main() {
	if err := service.InitIdiomDB(); err != nil {
		panic(fmt.Sprintf("init fail: %v", err))
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to, game start！")

	var lastIdiom string
	var lastChar string
	usedIdioms := make(map[string]bool)

	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == model.EXIT {
			fmt.Println(model.ExitMessage)
			break
		}

		if input == model.FAIL {
			fmt.Println(model.FailMessage)
			break
		}

		if !service.IsValidIdiom(input) {
			fmt.Println(model.ValidMessage)
			continue
		}

		if lastChar != "" && !strings.HasPrefix(input, lastChar) {
			fmt.Printf(model.NextMessage, lastIdiom, lastChar)
			continue
		}

		lastIdiom = input
		lastChar = service.GetLastChar(input)
		usedIdioms[input] = true

		matches := service.FindMatchingIdioms(lastChar, usedIdioms)
		if len(matches) == 0 {
			fmt.Println(model.ComFailMessage)
			break
		}

		// computer
		computerIdiom := matches[0]
		computerLastChar := service.GetLastChar(computerIdiom)
		fmt.Printf(model.ComNextMessage, computerIdiom, computerLastChar)

		lastIdiom = computerIdiom
		lastChar = computerLastChar
		usedIdioms[computerIdiom] = true
	}
}
