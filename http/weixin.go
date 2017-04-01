package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"weixin-provider/config"

	"github.com/toolkits/web/param"
)

type AccessTokenResponse struct {
	AccessToken string  `json:"access_token"`
	ExpiressIn  float64 `json:"expires_in"`
}

type AccessTokenErrorRespinse struct {
	Errcode float64
	Errmsg  string
}

type CustomServiceMsg struct {
	ToUser  string         `json:"touser"`
	AgentId int            `json:"agentid"`
	MsgType string         `json:"msgtype"`
	Text    TextMsgContent `json:"text"`
}

type TextMsgContent struct {
	Content string `json:"content"`
}

func configProRoutes() {
	http.HandleFunc("/weixin", func(w http.ResponseWriter, r *http.Request) {
		cfg := config.Config()
		token := param.String(r, "token", "")
		if cfg.Http.Token != token {
			http.Error(w, "no privilege", http.StatusForbidden)
			return
		}

		tos := param.MustString(r, "tos")
		content := param.MustString(r, "content")

		accessToken, expireln, err := getToken()
		if err != nil {
			log.Println("Get access_token error:", err)
			return
		}
		fmt.Println(accessToken, expireln)

		tos = strings.Replace(tos, ",", "|", -1)
		//fmt.Println(tos, content)

		//s := strings.Split(tos, ",")
		//for _, v := range s {
		//fmt.Println(v)
		//fmt.Println(content)

		err = pushCustomMsg(accessToken, config.Config().Weixin.Agentid, tos, content)
		if err != nil {
			log.Fatalln("Push Msg err:", err)
			return
		}

		//}
	})
}

func getToken() (string, float64, error) {
	tokenurl := config.Config().Weixin.Tokenurl
	corpid := config.Config().Weixin.Corpid
	corpsecret := config.Config().Weixin.Corpsecret
	requesturl := strings.Join([]string{tokenurl, "?corpid=", corpid, "&corpsecret=", corpsecret}, "")

	fmt.Println(requesturl)

	resp, err := http.Get(requesturl)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "", 0.0, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", 0.0, err
	}

	if bytes.Contains(body, []byte("access_token")) {
		atr := AccessTokenResponse{}
		err = json.Unmarshal(body, &atr)
		if err != nil {
			return "", 0.0, err
		}
		return atr.AccessToken, atr.ExpiressIn, nil
	} else {
		fmt.Println("return err")
		ater := AccessTokenErrorRespinse{}
		err = json.Unmarshal(body, &ater)
		if err != nil {
			return "", 0.0, err
		}
		return "", 0.0, fmt.Errorf("%s", ater.Errmsg)
	}
}

func pushCustomMsg(accessToken string, agentId int, toUser, msg string) error {
	csMsg := &CustomServiceMsg{
		ToUser:  toUser,
		AgentId: agentId,
		MsgType: "text",
		Text:    TextMsgContent{Content: msg},
	}

	//fmt.Println(msg)
	body, err := json.MarshalIndent(csMsg, " ", " ")
	if err != nil {
		return err
	}

	fmt.Println(string(body))

	body = bytes.Replace(body, []byte("\\u0026"), []byte("&"), -1)
	body = bytes.Replace(body, []byte("\\u003c"), []byte("<"), -1)
	body = bytes.Replace(body, []byte("\\u003e"), []byte(">"), -1)
	body = bytes.Replace(body, []byte("\\u003d"), []byte("="), -1)

	//fmt.Println(string(body))

	postReq, err := http.NewRequest("POST",
		strings.Join([]string{config.Config().Weixin.Sendurl, "?access_token=", accessToken}, ""),
		bytes.NewReader(body))

	if err != nil {
		return err
	}

	postReq.Header.Set("Content-Type", "application/json; encoding=utf-8")

	client := &http.Client{}
	resp, err := client.Do(postReq)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil

}
