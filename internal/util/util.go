package util

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/eiannone/keyboard"
)

const TAIL_LENGTH = 256
const KEYS_LENGTH = 36

var KEYSTROKE_MAP = [KEYS_LENGTH]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}

/////////////////////////////////////
// PUBLIC
/////////////////////////////////////

func ExecuteCommandInTerminal(commandsMap map[string]string, keyPressed string) error {
	command := strings.Replace(commandsMap[keyPressed], "\"", "\\\"", -1)
	appleScript := []string{"osascript", "-e"}
	_, errFindIterm := executeCommand([]string{"sh", "-c", "ps aux | grep iTerm2 | grep -v grep"})
	if errFindIterm != nil {
		_, errFindTerminal := executeCommand([]string{"sh", "-c", "ps aux | grep Terminal | grep -v grep"})
		if errFindTerminal != nil {
			return errFindTerminal
		}
		appleScript = append(appleScript, "tell application \"Terminal\" to do script \""+command+"\" in window 1", ">/dev/null")
	} else {
		appleScript = append(appleScript, "tell application \"iTerm2\" to tell current window to tell current session to write text \""+command+"\"")
	}

	_, executeCommandErr := executeCommand(appleScript)
	if executeCommandErr != nil {
		return executeCommandErr
	}

	return nil
}

func GetCommands() ([]string, map[string]string, error) {
	homeDir, errUserHomeDir := os.UserHomeDir()
	if errUserHomeDir != nil {
		return nil, nil, errUserHomeDir
	}
	var outTail string
	var getCommand func(input string) string
	outTailZsh, errTailZsh := executeCommand([]string{"tail", "-" + strconv.Itoa(TAIL_LENGTH), homeDir + "/.zsh_history"})
	if errTailZsh != nil {
		outTailBash, errTailBash := executeCommand([]string{"tail", "-" + strconv.Itoa(TAIL_LENGTH), homeDir + "/.bash_history"})
		if errTailBash != nil {
			return nil, nil, errTailBash
		}
		outTail = outTailBash
		getCommand = func(input string) string {
			return input
		}
	} else {
		outTail = outTailZsh
		getCommand = func(input string) string {
			split := strings.Split(input, ";")
			return strings.Join(split[1:], ";")
		}
	}

	historyLines := strings.Split(outTail, "\n")
	historyLinesMaxIndex := len(historyLines) - 1
	counter := 0
	commands := []string{}
	commandsMap := map[string]string{}
	for len(commands) < KEYS_LENGTH {
		historyLine := historyLines[historyLinesMaxIndex-counter]
		if historyLine != "" {
			command := getCommand(historyLine)
			if command != "dasd" {
				commandsMap[KEYSTROKE_MAP[len(commands)]] = command
				commands = append(commands, command)
			}
		}
		counter += 1
	}
	return commands, commandsMap, nil
}

func GetKeyPressed() (string, error) {
	fmt.Print("Press key to execute: ")
	char, _, errGetSingleKey := keyboard.GetSingleKey()
	if errGetSingleKey != nil {
		return "", errGetSingleKey
	}
	charString := string(char)
	fmt.Print(charString)
	for !includes(KEYSTROKE_MAP, charString) {
		fmt.Print("\nInvalid key, enter again: ")
		char, _, errGetSingleKeyInvalid := keyboard.GetSingleKey()
		if errGetSingleKeyInvalid != nil {
			return "", errGetSingleKeyInvalid
		}
		charString = string(char)
		fmt.Print(charString)
	}
	fmt.Print(" -> ")
	return charString, nil
}

func PrintCommands(commands []string) {
	for i := len(commands) - 1; i >= 0; i-- {
		fmt.Println(" "+KEYSTROKE_MAP[i], "—>", commands[i])
	}
}

/////////////////////////////////////
// PRIVATE
/////////////////////////////////////

func executeCommand(command []string) (string, error) {
	cmd := exec.Command(command[0], command[1:]...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

func includes(elements [KEYS_LENGTH]string, target string) bool {
	for _, element := range elements {
		if element == target {
			return true
		}
	}
	return false
}
