package main

import (
	"fmt"
	"github.com/holys/goredis"
	"regexp"
	"strings"
)

// Read-Eval-Print Loop
func repl(server RedisServer, commandLine string) string {
	reg, _ := regexp.Compile(`'.*?'|".*?"|\S+`)

	client := cliConnect(server)
	defer client.Close()

	cmds := reg.FindAllString(commandLine, -1)
	if len(cmds) == 0 {
		return ""
	}

	cmd := strings.ToLower(cmds[0])
	if cmd == "help" || cmd == "?" {
		return ""
	} else if cmd == "quit" || cmd == "exit" {
		return ""
	} else if cmd == "clear" {
		return ""
	} else if cmd == "connect" {
		return ""
	} else {
		return cliSendCommand(client, cmds)
	}
}

func cliSendCommand(client *goredis.Client, cmds []string) string {
	args := make([]interface{}, len(cmds[1:]))
	for i := range args {
		args[i] = strings.Trim(string(cmds[1+i]), "\"'")
	}

	cmd := strings.ToLower(cmds[0])

	r, err := client.Do(cmd, args...)
	if err != nil {
		return fmt.Sprintf("(error) %s", err.Error())
	}

	if cmd == "info" {
		return printInfo(r)
	} else {
		return printReply(0, r)
	}

}

func cliConnect(server RedisServer) *goredis.Client {
	client := goredis.NewClient(server.Addr, server.Password)
	client.SetMaxIdleConns(1)
	return client
}

func printInfo(reply interface{}) string {
	switch reply := reply.(type) {
	case []byte:
		return fmt.Sprintf("%s", reply)
		//some redis proxies don't support this command.
	case goredis.Error:
		return fmt.Sprintf("(error) %s", string(reply))
	}

	return "unknown reply"
}

func printReply(level int, reply interface{}) string {
	switch reply := reply.(type) {
	case int64:
		return fmt.Sprintf("(integer) %d", reply)
	case string:
		return fmt.Sprintf("%s", reply)
	case []byte:
		return fmt.Sprintf("%q", reply)
	case nil:
		return fmt.Sprintf("(nil)")
	case goredis.Error:
		return fmt.Sprintf("(error) %s", string(reply))
	case []interface{}:
		resp := ""
		for i, v := range reply {
			if i != 0 {
				resp += fmt.Sprintf("%s", strings.Repeat(" ", level*4))
			}

			s := fmt.Sprintf("%d) ", i+1)
			resp += fmt.Sprintf("%-4s", s)

			resp += printReply(level+1, v)
			if i != len(reply)-1 {
				resp += "\n"
			}
		}
		return resp
	default:
		return fmt.Sprintf("Unknown reply type: %+v", reply)
	}
}
