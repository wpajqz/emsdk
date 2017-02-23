package emsdk

import (
	"encoding/json"
	"strings"
)

// 授权模式注册用户
func (c *Client) CreateAccount(username, password, nickname string) error {
	data := `{"username":"` + username + `","password":"` + password + `","nickname":"` + nickname + `"}`
	_, err := c.sendRequest("users", strings.NewReader(data), "POST")
	if err != nil {
		if e, ok := err.(EMError); ok {
			// 对环信错误进行处理,目前发现的是数据库中保存的昵称和环信的昵称规则不一致引起的
			switch e.Code {
			case 400:
				data := `{"username":"` + username + `","password":"` + password + `"}`
				_, err = c.sendRequest("users", strings.NewReader(data), "POST")
			}
		}
	}

	return err
}

// 从环信服务器中删除用户
func (c *Client) DeleteAccount(username string) error {
	url := "users/" + username
	_, err := c.sendRequest(url, strings.NewReader(""), "DELETE")

	return err
}

// 修改用户密码
func (c *Client) ChangePassword(username, password string) error {
	url := "users/" + username + "/password"
	data := `{"newpassword":"` + password + `"}`
	_, err := c.sendRequest(url, strings.NewReader(data), "PUT")

	return err
}

// 修改用户昵称
func (c *Client) ChangeNickname(username, nickname string) error {
	url := "users/" + username
	data := `{"nickname":"` + nickname + `"}`
	_, err := c.sendRequest(url, strings.NewReader(data), "PUT")
	return err
}

// 判断一个用户是否存在
func (c *Client) IsExists(username string) (bool, error) {
	url := "users/" + username
	_, err := c.sendRequest(url, nil, "GET")
	if err != nil {
		if e, ok := err.(EMError); ok {
			if e.Code == 404 {
				return false, nil
			}
		}

		return false, err
	}

	return true, nil
}

// 查看一个用户的在线状态
func (c *Client) IsOnline(username string) bool {
	url := "users/" + username + "/status"
	res, err := c.sendRequest(url, nil, "GET")

	var result map[string]interface{}
	err = json.Unmarshal([]byte(res), &result)
	if err != nil {
		return false
	}

	v, ok := result["data"].(map[string]interface{})
	if !ok {
		return false
	}

	if v[username].(string) != "online" {
		return false
	}

	return true
}

// 禁用某个 IM 用户的账号，禁用后该用户不可登录，下次解禁后该账户恢复正常使用
func (c *Client) Deactivate(username string) bool {
	url := "users/" + username + "/deactivate"
	_, err := c.sendRequest(url, nil, "POST")
	if err != nil {
		return false
	}

	return true
}

// 解除对某个 IM 用户账号的禁用，解禁后用户恢复正常使用
func (c *Client) Activate(username string) bool {
	url := "users/" + username + "/activate"
	_, err := c.sendRequest(url, nil, "POST")
	if err != nil {
		return false
	}

	return true
}

// 如果某个 IM 用户已经登录环信服务器，强制其退出登录
func (c *Client) Disconnect(username string) bool {
	url := "users/" + username + "/disconnect"
	_, err := c.sendRequest(url, nil, "GET")
	if err != nil {
		return false
	}

	return true
}
