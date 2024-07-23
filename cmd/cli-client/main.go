package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Ayaya-zx/mem-flow/internal/auth"
	"github.com/Ayaya-zx/mem-flow/internal/entity"
)

const URL = "http://localhost:8765"

var token string

func main() {
	var err error
	var authData auth.AuthData
	fReg := flag.Bool("reg", false, "Register instead of authenticate")
	flag.Parse()

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
		err = registrate(&authData)
	} else {
		err = authenticate(&authData)
	}
	if err != nil {
		fmt.Println(err)
		return
	}

	printHelp()

	fmt.Print("> ")
	for scanner.Scan() {
		input := scanner.Text()
		switch input {
		case "list", "l":
			list()
		case "":
		default:
			fmt.Println("Unknown command")
		}
		fmt.Print("> ")
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}

func printHelp() {
	fmt.Println("Usage:")
	fmt.Println("\thelp   (h) [topic title]   print this help")
	fmt.Println("\tlist   (l)                 print all topic titles")
	fmt.Println("\tshow   (s) [topic title]   print topic info")
	fmt.Println("\tadd    (a) [topic title]   add topic")
	fmt.Println("\trepeat (p) [topic title]   repeat topic")
	fmt.Println("\tremove (r) [topic title]   remove topic")
}

func registrate(authData *auth.AuthData) error {
	data, err := json.Marshal(&authData)
	if err != nil {
		return err
	}

	resp, err := http.Post(
		URL+"/registration",
		"application/json",
		bytes.NewReader(data),
	)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("api status code %d", resp.StatusCode)
	}

	tokenData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	resp.Body.Close()

	token = string(tokenData)
	return nil
}

func authenticate(authData *auth.AuthData) error {
	data, err := json.Marshal(&authData)
	if err != nil {
		return err
	}

	resp, err := http.Post(
		URL+"/auth",
		"application/json",
		bytes.NewReader(data),
	)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("api status code %d", resp.StatusCode)
	}

	tokenData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	resp.Body.Close()

	token = string(tokenData)
	return nil
}

func list() {
	if token == "" {
		panic("No auth token")
	}

	req, err := http.NewRequest(
		"GET",
		URL+"/topics",
		nil,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("api status code %d", resp.StatusCode)
		return
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var topics []*entity.Topic
	err = json.Unmarshal(data, &topics)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Themes list:")
	for _, t := range topics {
		fmt.Println(t.Title)
	}
}
