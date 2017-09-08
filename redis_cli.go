package main

import (
	"fmt"
	"github.com/holys/goredis"
	"strings"
	"errors"
	"strconv"
)

// Read-Eval-Print Loop
func repl(server RedisServer, commandLine string) string {
	cmds, err := parseEditorCommand(commandLine)
	if err != nil {
		return err.Error()
	}

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
		client := cliConnect(server)
		defer client.Close()
		return cliSendCommand(client, cmds)
	}
}

// Returns the executable path and arguments
func parseEditorCommand(editorCmd string) ([]string, error) {
	var args []string
	state := "start"
	current := ""
	quote := "\""
	for i := 0; i < len(editorCmd); i++ {
		c := editorCmd[i]

		if state == "quotes" {
			if string(c) != quote {
				if c == '\\' {
					i++
					if i >= len(editorCmd) {
						return []string{}, errors.New(fmt.Sprintf("nothing escape in command line: %s", editorCmd))
					}
					current += editorCmd[i-1:i+1]
				} else {
					current += string(c)
				}
			} else {
				current = "\"" + current + "\""
				unquoted, _ := strconv.Unquote(current)
				args = append(args, unquoted)
				current = ""
				state = "start"
			}
			continue
		}

		if c == '"' || c == '\'' {
			state = "quotes"
			quote = string(c)
			continue
		}

		if state == "arg" {
			if c == ' ' || c == '\t' {
				args = append(args, current)
				current = ""
				state = "start"
			} else {
				current += string(c)
			}
			continue
		}

		if c != ' ' && c != '\t' {
			state = "arg"
			current += string(c)
		}
	}

	if state == "quotes" {
		return []string{}, errors.New(fmt.Sprintf("Unclosed quote in command line: %s", editorCmd))
	}

	if current != "" {
		args = append(args, current)
	}

	if len(args) <= 0 {
		return []string{}, errors.New("Empty command line")
	}

	return args, nil
}

func cliSendCommand(client *goredis.Client, cmds []string) string {
	args := make([]interface{}, len(cmds[1:]))
	for i := range args {
		args[i] = string(cmds[1+i])
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
