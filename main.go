package main
import (
    gc "code.google.com/p/goncurses"
    "fmt"
    "time"
)
const (
    ChatAreaColor = iota
    InputAreaColor
)
func main() {
    scr,_ := gc.Init()
    defer gc.End()

    rows, cols := scr.Maxyx()

    gc.InitPair(ChatAreaColor, gc.C_WHITE, gc.C_BLACK)
    gc.InitPair(InputAreaColor, gc.C_BLACK, gc.C_WHITE)

    chatArea := scr.Derived(rows-1, cols, 0, 0)
    chatArea.SetBackground(gc.Character(' ' | gc.ColorPair(ChatAreaColor)))

    inputArea := scr.Derived(1, cols, rows-1, 0)
    chatArea.SetBackground(gc.Character(' ' | gc.ColorPair(InputAreaColor)))

    chat := make(chan string)
    go func() {
        for msg := range(chat) {
            msg = fmt.Sprintf("%v| %v", time.Now().Format("15:04:05"), msg)
            chatArea.Scroll(1)
            chatArea.MovePrint(rows - 2, 0, msg)
            chatArea.Refresh()
        }
    }()
    defer close(chat)


    userInputChan := make(chan string)
    go func() {
        for msg := range(userInputChan) {
            if len(msg) >0 && msg[0] == '/' {
                chat <- fmt.Sprintf("Command: %s", msg[1:])
            } else {
                chat <- msg
            }
        }
    }()
    defer close(userInputChan)

    gc.Echo(false)
    gc.CBreak(true)
    gc.Raw(true)
    chatArea.ScrollOk(true)
    scr.Keypad(true)

    chat <- "Welcome to snails shitty chat thing."
    chat <- "Press esc to quit, it may or may not break stuff. "
    chat <- "If it does, do a 'reset' to fix it."
    buffer := ""
    for {
        chatArea.Refresh()
        key := inputArea.GetChar()
        switch key {
            case gc.Key(27):
                return
            case gc.KEY_RETURN:
                userInputChan <- buffer
                buffer = ""
                chatArea.Refresh()
            case gc.Key(127)://backspace
                l := len(buffer)
                if l > 0 {
                    buffer = buffer[:l-1]
                }
            default:
                buffer = fmt.Sprintf("%s%c", buffer, key)
        }
        inputArea.Clear()
        inputArea.MovePrint(0, 0, buffer)
    }

}
