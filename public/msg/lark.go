package msg

import (
	"fmt"
	"github.com/astaxie/beego/httplib"
	"regexp"
)

const (
	//NODE_ERR_LARK = "05e8ef56-6e70-4b07-bc27-f60c56a3019b"
	NODE_ERR_LARK = "a5383493-cbd9-4a56-a797-45f0390a1c75"
)

type larkNotifyMessage struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

func SendLarkNotify(larkType string,title, message string) (string, error) {
	var rgx = regexp.MustCompile(`apikey=.*?:`)
	rs := rgx.ReplaceAllString(message,"apikey=*******************")
	fmt.Println("send lark,title:",title)
	req := httplib.Post("https://open.larksuite.com/open-apis/bot/hook/"+larkType)
	req.JSONBody(larkNotifyMessage{
		Title: title,
		Text:  rs,
	})
	return req.String()
}