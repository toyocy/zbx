package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/urfave/cli"
)

type Params struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type Request struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  Params `json:"params"`
	ID      int    `json:"id"`
}

type Result struct {
	Token string `json:"result"`
}

func CreateRequestString(c *cli.Context) []byte {
	request := Request{
		Jsonrpc: "2.0",
		Method:  "user.login",
		Params: Params{
			User:     c.String("user"),
			Password: c.String("password"),
		},
		ID: 1,
	}

	values, err := json.Marshal(request)
	if err != nil {
		log.Fatalln(err)
	}

	return values
}

func GetAuthToken(request []byte, c *cli.Context) Result {
	url := "http://" + c.String("zabbix-url") + "/api_jsonrpc.php"
	res, err := http.Post(url, "application/json-rpc", bytes.NewBuffer(request))
	if err != nil {
		log.Fatalln(err)
	}

	defer res.Body.Close()

	body, error := ioutil.ReadAll(res.Body)
	if error != nil {
		log.Fatal(error)
	}

	var result Result
	json.Unmarshal(body, &result)

	if err != nil {
		fmt.Println("JSON unmarshall Error : ", err)
	}
	return result
}

func init() {
	cmdList = append(cmdList, cli.Command{
		Name:  "login",
		Usage: "Sign In to Zabbix Server",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "zabbix-url, z",
				Value: "localhost/zabbix",
				Usage: "Zabbix URL without protocol",
			},
			cli.StringFlag{
				Name:  "user, u",
				Value: "admin",
				Usage: "User name",
			},
			cli.StringFlag{
				Name:  "password, p",
				Value: "zabbix",
				Usage: "Password",
			},
		},
		Action: func(c *cli.Context) error {
			request := CreateRequestString(c)
			result := GetAuthToken(request, c)
			return action(c, &login{Token: result.Token})
		},
	})
}

type login struct {
	Token string
}

func (l *login) Run(c *cli.Context) error {
	fmt.Println(l.Token)
	return nil
}
