package logic

import (
	"mangrove/internal/dao/mysql"
	"mangrove/internal/schema"
	"net/url"
	"regexp"
)

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
