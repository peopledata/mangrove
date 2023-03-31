package logic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mangrove/internal/models"
	"mangrove/internal/schema"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

const (
	StatusTaskInit = 1 // 初始状态
)

// TODO：需要将访问基础设置的API封装成SDK

// CreateTask 创建任务，直接调用基础设施的API接口
func CreateTask(demandId int64, user *models.User) error {
	apiHost := viper.GetString("marketplace.host")
	apiKey := viper.GetString("marketplace.api_key")
	task := &schema.CreateTaskReq{
		UserId:   user.ID,
		Username: user.Username,
	}
	jsonData, err := json.Marshal(task)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/demands/%d/task", apiHost, demandId), bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	header := http.Header{}
	header.Add("Content-Type", "application/json")
	header.Add("X-API-KEY", apiKey)
	req.Header = header
	httpclient := http.Client{
		Timeout: time.Second * 10, // 设置超时时间为10s
	}
	resp, err := httpclient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var rd schema.RespData
	if err := json.Unmarshal(body, &rd); err != nil {
		return err
	}
	if rd.Code == 1000 {
		return nil
	}
	return fmt.Errorf("create task failed: %s", rd.Msg)
}

func ListTasks(demandId int64) (*schema.TaskListResp, error) {
	apiHost := viper.GetString("marketplace.host")
	apiKey := viper.GetString("marketplace.api_key")
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/demands/%d/task", apiHost, demandId), nil)
	if err != nil {
		return nil, err
	}
	header := http.Header{}
	header.Add("Content-Type", "application/json")
	header.Add("X-API-KEY", apiKey)
	req.Header = header
	httpclient := http.Client{
		Timeout: time.Second * 10, // 设置超时时间为10s
	}
	resp, err := httpclient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var rtl schema.RespTaskList
	if err := json.Unmarshal(body, &rtl); err != nil {
		return nil, err
	}
	if rtl.Code == 1000 {
		return rtl.Data, nil
	}
	return nil, fmt.Errorf("get task list failed: %s", rtl.Msg)
}

func TaskAlgoRecords(taskId int64) ([]schema.AlgoRecordItem, error) {
	apiHost := viper.GetString("marketplace.host")
	apiKey := viper.GetString("marketplace.api_key")
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/tasks/%d/algo", apiHost, taskId), nil)
	if err != nil {
		return nil, err
	}
	header := http.Header{}
	header.Add("Content-Type", "application/json")
	header.Add("X-API-KEY", apiKey)
	req.Header = header
	httpclient := http.Client{
		Timeout: time.Second * 10, // 设置超时时间为10s
	}
	resp, err := httpclient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var rtl schema.RespAlgoRecordList
	if err := json.Unmarshal(body, &rtl); err != nil {
		return nil, err
	}
	if rtl.Code == 1000 {
		return rtl.Data, nil
	}
	return nil, fmt.Errorf("get task algorecord list failed: %s", rtl.Msg)
}
