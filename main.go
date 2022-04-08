package main

import (
	"fmt"

	"github.com/doctorpoor/dasd/internal/util"
)

func main() {
	commands, commandsMap, getCommandsErr := util.GetCommands()
	if getCommandsErr != nil {
		fmt.Println("Error getting commands:", getCommandsErr)
		return
	}
	util.PrintCommands(commands)

	keyPressed, getKeyPressedErr := util.GetKeyPressed()
	if getKeyPressedErr != nil {
		fmt.Println("Error getting key pressed:", getKeyPressedErr)
		return
	}

	executeCommandInTerminalErr := util.ExecuteCommandInTerminal(commandsMap, keyPressed)
	if executeCommandInTerminalErr != nil {
		fmt.Println("Error executing command in terminal:", executeCommandInTerminalErr)
	}
}
