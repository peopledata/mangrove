package schema

import "time"

type DemandListResp struct {
	DemandId       int64     `json:"demand_id"`
	Name           string    `json:"name"`
	ValidAt        time.Time `json:"valid_at"`
	Status         int       `json:"status"`
	Category       string    `json:"category"`
	Content        string    `json:"content"`
	NeedUsers      int       `json:"need_users"`
	UseTimes       int       `json:"use_times"`
	ExistingUsers  int       `json:"existing_users"`
	AvailableTimes int       `json:"available_times"`
	CreatedAt      time.Time `json:"created_at"`
}

type APIDemandListResp struct {
	DemandId  int64     `json:"demand_id"`
	Name      string    `json:"name"`
	Brief     string    `json:"brief"`
	ValidAt   time.Time `json:"valid_at"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
}

type APIDemandContractListResp struct {
	DemandId       int64           `json:"demand_id"`
	Name           string          `json:"name"`
	Brief          string          `json:"brief"`
	DemandContract *DemandContract `json:"contract"`
}

type APIDemandDetailResp struct {
	DemandId       int64           `json:"demand_id"`
	Name           string          `json:"name"`
	Brief          string          `json:"brief"`
	ValidAt        time.Time       `json:"valid_at"`
	Category       string          `json:"category"`
	Content        string          `json:"content"`
	Purpose        string          `json:"purpose"`
	Agreement      string          `json:"agreement"`
	CreatedAt      time.Time       `json:"created_at"`
	DemandContract *DemandContract `json:"contract"`
}

type DemandContract struct {
	Address   string `json:"address"`
	TokenName string `json:"token"`
	Symbol    string `json:"symbol,omitempty"`
	ABI       string `json:"abi,omitempty"`
}

type DemandCreateReq struct {
	Name      string    `json:"name" binding:"required"`
	ValidAt   time.Time `json:"valid_at" binding:"required"`
	Brief     string    `json:"brief" binding:"required"`
	Category  string    `json:"category" binding:"required"`
	Content   string    `json:"content" binding:"required"`
	NeedUsers int       `json:"need_users" binding:"required,gt=0"`
	UseTimes  int       `json:"use_times" binding:"required,gt=0"`
	Purpose   string    `json:"purpose" binding:"required"`
	Algorithm string    `json:"algorithm" binding:"required"`
	Agreement string    `json:"agreement" binding:"required"`
}

type DemandDetailResp struct {
	DemandId        int64     `json:"demand_id"`
	Name            string    `json:"name"`
	Brief           string    `json:"brief"`
	ValidAt         time.Time `json:"valid_at"`
	CreatedAt       time.Time `json:"created_at"`
	Category        string    `json:"category"`
	Content         string    `json:"content"`
	NeedUsers       int       `json:"need_users"`
	UseTimes        int       `json:"use_times"`
	ExistingUsers   int       `json:"existing_users"`
	AvailableTimes  int       `json:"available_times"`
	ContractAddress string    `json:"contract_address"`
	ContractSymbol  string    `json:"contract_symbol"`
	Purpose         string    `json:"purpose"`
	Algorithm       string    `json:"algorithm"`
	Agreement       string    `json:"agreement"`
}

type DemandInfoResp struct {
	DemandId  int64     `json:"demand_id"`
	Name      string    `json:"name"`
	Brief     string    `json:"brief"`
	ValidAt   time.Time `json:"valid_at"`
	Category  string    `json:"category"`
	Content   string    `json:"content"`
	NeedUsers int       `json:"need_users"`
	UseTimes  int       `json:"use_times"`
	Purpose   string    `json:"purpose"`
	Algorithm string    `json:"algorithm"`
	Agreement string    `json:"agreement"`
}

type DemandUpdateReq struct {
	DemandId  int64     `json:"demand_id" binding:"required"`
	Name      string    `json:"name"`
	Brief     string    `json:"brief"`
	ValidAt   time.Time `json:"valid_at"`
	Category  string    `json:"category"`
	Content   string    `json:"content"`
	NeedUsers int       `json:"need_users"`
	UseTimes  int       `json:"use_times"`
	Purpose   string    `json:"purpose"`
	Algorithm string    `json:"algorithm"`
	Agreement string    `json:"agreement"`
}
