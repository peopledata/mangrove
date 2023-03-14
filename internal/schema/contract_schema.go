package schema

import "time"

type ContractRecordResp struct {
	ID       uint      `json:"id"`
	Did      string    `json:"did"`
	SignTime time.Time `json:"sign_time"`
}
