// Able Gao @
// ablegao@gmail.com
// description：
//
//

package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type DTContentTypeText struct {
	Content string `json:"content"`
}

type DTBody struct {
	MsgType string            `json:"msgtype"`
	Text    DTContentTypeText `json:"text"`
}

func DingTalk(msgstr string) error {
	msgBody := DTBody{MsgType: "text", Text: DTContentTypeText{msgstr}}

	msg, _ := json.Marshal(msgBody)
	req, err := http.Post("https://oapi.dingtalk.com/robot/send?access_token=d608964bf1132e6f8c1e6a3b06161a16ff17187c9965e6a3311d5d3e210d894f", "application/json", bytes.NewBuffer(msg))
	ioutil.ReadAll(req.Body)
	return err
}
