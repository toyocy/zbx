package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/urfave/cli"
)

type params struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type request struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  params `json:"params"`
	ID      int    `json:"id"`
}

type result struct {
	Token string `json:"result"`
}

type login struct {
	Token string
}

func createRequestString(c *cli.Context) []byte {
	request := request{
		Jsonrpc: "2.0",
		Method:  "user.login",
		Params: params{
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

func getAuthToken(request []byte, c *cli.Context) result {
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

	var result result
	json.Unmarshal(body, &result)

	if err != nil {
		fmt.Println("JSON unmarshall Error : ", err)
	}
	return result
}

func Login() cli.Command {
	return cli.Command{
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
			request := createRequestString(c)
			result := getAuthToken(request, c)
			fmt.Println(result.Token)
			return nil
		},
	}
}
