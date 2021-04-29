package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"ytzabbixalert/logging"
	"ytzabbixalert/model"
	"ytzabbixalert/setting"
)

const (
	WX_API_OK_CODE            = 0
	WX_API_OK_MSG             = "ok"
	WX_API_TOKEN_EXPIRED_CODE = 42001
)

var (
	TOKEN string = ``
)

type TokenResponse struct {
	ErrCode     int32  `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	Expires     int32  `json:"expires_in"`
}

type MsgRequest struct {
	ToUser  string     `json:"touser"`
	MsgType string     `json:"msgtype"`
	AgentId string     `json:"agentid"`
	Text    MsgContent `json:"text"`
	Safe    int8       `json:"safe"`
}

type MsgContent struct {
	Content string `json:"content"`
}

type MsgResponse struct {
	ErrCode      int32  `json:"errcode"`
	ErrMsg       string `json:"errmsg"`
	InvalidUser  string `json:"invaliduser"`
	InvalidParty string `json:"invalidparty"`
	InvalidTag   string `json:"invalidtag"`
}

/*
 * 获取企业微信 access_token
 */
func getAccessToken() error {
	fullUrl := fmt.Sprintf(setting.AppSetting.WeChatTokenUrl+"?corpid=%s&corpsecret=%s", setting.AppSetting.CompanyId, setting.AppSetting.CompanySecret)
	tokenResponse := TokenResponse{}
	response, err := http.Get(fullUrl)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return err
	}

	if tokenResponse.ErrMsg == WX_API_OK_MSG && tokenResponse.ErrCode == WX_API_OK_CODE {
		TOKEN = tokenResponse.AccessToken
		logging.Info("获取access_token成功")
		return nil
	} else {
		return fmt.Errorf(tokenResponse.ErrMsg)
	}
}

/*
 * token重试机制
 */
func tokenExpiredRetry(content, toUser string) (bool, error) {
	logging.Info("token已经过期，重新获取")
	err := getAccessToken()
	if err != nil {
		return false, err
	}
	_, err = sendMessage2WeChat(content, toUser)
	if err != nil {
		return false, err
	}
	return true, nil
}

func sendMessage2WeChat(content, toUser string) (bool, error) {
	if TOKEN != `` {
		fullUrl := fmt.Sprintf(setting.AppSetting.WeChatMessageUrl+"?access_token=%s", TOKEN)
		requestBody := MsgRequest{}
		requestBody.ToUser = toUser
		requestBody.MsgType = "text"
		requestBody.Safe = int8(setting.AppSetting.Safe)
		requestBody.AgentId = setting.AppSetting.AgentId
		requestBody.Text.Content = content
		requestJson, err := json.Marshal(requestBody)
		if err != nil {
			return false, err
		}

		request, err := http.NewRequest(http.MethodPost, fullUrl, strings.NewReader(string(requestJson)))
		if err != nil {
			return false, err
		}

		client := http.Client{}
		response, err := client.Do(request)
		if err != nil {
			return false, err
		}

		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return false, err
		}
		responseBody := MsgResponse{}
		if err = json.Unmarshal(body, &responseBody); err != nil {
			return false, err
		}

		if responseBody.ErrCode == WX_API_OK_CODE && responseBody.ErrMsg == WX_API_OK_MSG {
			logging.Info(fmt.Sprintf("发送微信消息成功 ====> result[%s]", string(body)))
			return true, nil
		} else if responseBody.ErrCode == WX_API_TOKEN_EXPIRED_CODE {
			return tokenExpiredRetry(content, toUser)
		} else {
			return false, fmt.Errorf(string(body))
		}
	}
	return false, fmt.Errorf("access_token 为空")
}

func main() {
	setting.Setup()
	logging.Setup()

	if len(os.Args) != 3 {
		logging.Error("入参数量不正确")
		return
	}

	// 第一个参数是告警机器类型 机器网络层面问题/系统内部服务进程问题
	// 第二个参数是告警消息内容
	for index, arg := range os.Args {
		logging.Info("arg[", index, "]: [", arg, "]")
	}

	if TOKEN == `` {
		if err := getAccessToken(); err != nil {
			logging.Error(err.Error())
		}
	}

	// 网络层面
	if strings.Contains(os.Args[1], model.NAT_PLATFORM) || strings.Contains(os.Args[1], model.ROUTER_PLATFORM) {
		_, err := sendMessage2WeChat(os.Args[2], setting.AppSetting.ToAllUser)
		if err != nil {
			logging.Error(err.Error())
		}
		// 服务层面
	} else {
		_, err := sendMessage2WeChat(os.Args[2], setting.AppSetting.ToUser)
		if err != nil {
			logging.Error(err.Error())
		}
	}

}
