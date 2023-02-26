package schema

type MenuItemResp struct {
	Id                 string      `json:"id"`
	Name               string      `json:"name"`
	Icon               string      `json:"icon,omitempty"`
	Zh                 *MenuLocale `json:"zh"`
	Route              string      `json:"route"`
	MenuParentId       string      `json:"menuParentId,omitempty"`
	BreadcrumbParentId string      `json:"breadcrumbParentId,omitempty"`
}

type MenuLocale struct {
	Name string `json:"name"`
}

type DashboardSalesItemResp struct {
	Name        int `json:"name"`
	Clothes     int `json:"Clothes"`
	Food        int `json:"Food"`
	Electronics int `json:"Electronics"`
}

type DashboardNumberItemResp struct {
	Icon   string `json:"icon"`
	Color  string `json:"color"`
	Title  string `json:"title"`
	Number int    `json:"number"`
}

type UserInfoResp struct {
	Id          int64     `json:"id"`
	Username    string    `json:"username"`
	Permissions *UserRole `json:"permissions"`
	Avatar      string    `json:"avatar"`
}

type UserRole struct {
	Role string `json:"role"`
}
