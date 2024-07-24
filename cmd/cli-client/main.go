package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/Ayaya-zx/mem-flow/internal/auth"
	"github.com/Ayaya-zx/mem-flow/internal/client"
	"github.com/Ayaya-zx/mem-flow/internal/entity"
)

const URL = "http://localhost:8765"

var cs *client.ClientService

func main() {
	var err error
	var authData auth.AuthData
	fReg := flag.Bool("reg", false, "Register instead of authenticate")
	flag.Parse()

	cs = client.NewClientService(URL)
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("name: ")
	scanner.Scan()
	authData.Name = scanner.Text()
	if err = scanner.Err(); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Print("password: ")
	scanner.Scan()
	authData.Password = scanner.Text()
	if err = scanner.Err(); err != nil {
		fmt.Println(err)
		return
	}

	if *fReg {
		err = cs.Register(authData)
	} else {
		err = cs.Auth(authData)
	}
	if err != nil {
		fmt.Println(err)
		return
	}

	help()

	fmt.Print("> ")
	for scanner.Scan() {
		input := scanner.Text()
		handleCommand(input)
		fmt.Print("> ")
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}

func help() {
	fmt.Println("Usage:")
	fmt.Println("\thelp    (h)                print this help")
	fmt.Println("\tlist    (l)                print all topic titles")
	fmt.Println("\tshow    (s) [topic id]     print topic info")
	fmt.Println("\tadd     (a) [topic title]  add topic")
	fmt.Println("\trepeat  (r) [topic id]     repeat topic")
	fmt.Println("\tdelete  (d) [topic id]     delete topic")
}

func shortHelp() {
	fmt.Println("Bad command format")
	fmt.Println("To print help type 'help' or 'h'")
}

func handleCommand(input string) {
	var cmd, arg string

	split := strings.Split(input, " ")
	if len(split) > 2 {
		shortHelp()
		return
	}
	cmd = split[0]
	if len(split) > 1 {
		arg = split[1]
	} else {
		arg = ""
	}
	switch cmd {
	case "list", "l":
		list()
	case "add", "a":
		if arg == "" {
			shortHelp()
			return
		}
		add(arg)
	case "show", "s":
		if arg == "" {
			shortHelp()
			return
		}
		id, err := strconv.Atoi(arg)
		if err != nil {
			fmt.Println(err)
			return
		}
		show(id)
	case "repeat", "r":
		if arg == "" {
			shortHelp()
			return
		}
		id, err := strconv.Atoi(arg)
		if err != nil {
			fmt.Println(err)
			return
		}
		repeat(id)
	case "delete", "d":
		if arg == "" {
			shortHelp()
			return
		}
		id, err := strconv.Atoi(arg)
		if err != nil {
			fmt.Println(err)
			return
		}
		remove(id)
	case "help", "h":
		help()
	case "":
	default:
		fmt.Println("Unknown command")
	}
}

func list() {
	topics, err := cs.GetAllTopics()
	if err != nil {
		fmt.Println(err)
		return
	}

	slices.SortFunc(topics, func(a, b entity.Topic) int {
		return a.Id - b.Id
	})

	fmt.Println("Themes list:")
	for _, t := range topics {
		fmt.Printf("%d: %s\n", t.Id, t.Title)
	}
}

func add(title string) {
	err := cs.AddTopic(title)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("OK")
	}
}

func show(id int) {
	topic, err := cs.GetTopicById(id)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Id:", topic.Id)
	fmt.Println("Title:", topic.Title)
	fmt.Println("Created:", topic.Created)
	fmt.Println("Last repeated:", topic.LastRepeated)
	fmt.Println("Next repeat:", topic.NextRepeat)
}

func repeat(id int) {
	err := cs.RepeatTopic(id)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("OK")
	}
}

func remove(id int) {
	err := cs.RemoveTopic(id)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("OK")
	}
}
