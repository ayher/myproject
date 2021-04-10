package fmt

import (
	"runtime"
	realfmt "fmt"
	"strings"
)

func formatLog(f interface{}, v ...interface{}) string {
	var msg string
	switch f.(type) {
	case string:
		msg = f.(string)
		if len(v) == 0 {
			return msg
		}
		if strings.Contains(msg, "%") && !strings.Contains(msg, "%%") {
			//format string
		} else {
			//do not contain format char
			msg += strings.Repeat(" %v", len(v))
		}
	default:
		msg = realfmt.Sprint(f)
		if len(v) == 0 {
			return msg
		}
		msg += strings.Repeat(" %v", len(v))
	}
	return realfmt.Sprintf(msg, v...)
}

func Println(f interface{}, v ...interface{})  {
	pc,_,line,_ := runtime.Caller(1)
	ff := runtime.FuncForPC(pc)
	me:=realfmt.Sprintf("[%s:%d] %s",ff.Name(),line,formatLog(f, v...))
	realfmt.Printf("\033[1;32;40m%s\033[0m\n",me)
}

func Error(f interface{}, v ...interface{})  {
	pc,_,line,_ := runtime.Caller(1)
	ff := runtime.FuncForPC(pc)
	me:=realfmt.Sprintf("[%s:%d] %s",ff.Name(),line,formatLog(f, v...))
	realfmt.Printf("\033[1;31;40m%s\033[0m\n",me)
}

func Debug(f interface{}, v ...interface{})  {
	pc,_,line,_ := runtime.Caller(1)
	ff := runtime.FuncForPC(pc)
	me:=realfmt.Sprintf("[%s:%d] %s",ff.Name(),line,formatLog(f, v...))
	realfmt.Printf("\033[1;33;40m%s\033[0m\n",me)
}
