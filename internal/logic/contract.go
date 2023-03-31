package logic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mangrove/internal/dao/mysql"
	"mangrove/internal/schema"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/spf13/viper"
)

func GetContractRecords(demandId int64, page, pageSize int) (*schema.RespContract, error) {
	apiHost := viper.GetString("marketplace.host")
	apiKey := viper.GetString("marketplace.api_key")
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/demands/%d/contract?page=%d&pageSize=%d", apiHost, demandId, page, pageSize), nil)
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
	var rtl schema.RespContractList
	if err := json.Unmarshal(body, &rtl); err != nil {
		return nil, err
	}
	if rtl.Code == 1000 {
		return rtl.Data, nil
	}
	return nil, fmt.Errorf("get demand contract record list failed: %s", rtl.Msg)
}

func GetContractRecordsByDemandId(demandId int64, page, pageSize int) []schema.ContractRecordResp {
	contractRecords := mysql.GetContractRecordsByDemandId(demandId, page, pageSize)
	var records []schema.ContractRecordResp
	for idx := range contractRecords {
		record := contractRecords[idx]
		// dic -> did:ppld:zQmSB8zjpigGkNNjSjLmxTgut2bQoSJ9iyW8Patv3owYnBU%40http%3A%2F%2Fppldid.peopledata.org.cn:3000
		decodedStr, _ := url.QueryUnescape(record.Did)
		r := regexp.MustCompile(`did:ppld:(.*?)@`)
		match := r.FindStringSubmatch(decodedStr)
		if len(match) > 1 {
			records = append(records, schema.ContractRecordResp{
				ID:       record.ID,
				Did:      match[1],
				SignTime: record.SignTime,
			})
		}
	}
	return records
}

func TotalContractRecordsByDemandId(demandId int64) int64 {
	return mysql.GetContractRecordsByDemandIdCount(demandId)
}
