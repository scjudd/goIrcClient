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

    gc.Echo(false)
    gc.CBreak(true)
    //gc.Raw(true)
    chatArea.ScrollOk(true)
    scr.Keypad(true)

    field, _ := gc.NewField(1, cols, rows-1, 0, 0, 0)
    defer field.Free()

    form, _ := gc.NewForm([]*gc.Field{ field })
    defer form.Free()
    
    form.Post()
    defer form.UnPost()

    for {
        ch := scr.GetChar()        
        switch ch {
        case gc.Key(27):
            return
        case gc.KEY_DOWN:
            form.Driver(gc.REQ_NEXT_FIELD)
            form.Driver(gc.REQ_END_LINE)
        case gc.KEY_UP:
            form.Driver(gc.REQ_PREV_FIELD)
            form.Driver(gc.REQ_END_LINE)
        default:
            form.Driver(ch)
        }
    }


/*    buffer := ""
    for {
        chatArea.Refresh()
        key := inputArea.GetChar()
        switch key {
            case gc.Key(27):
                return
            case gc.KEY_RETURN:
                chat <- buffer
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
    }*/
    
}