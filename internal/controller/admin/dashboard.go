package admin

import (
	"mangrove/internal/controller"
	"mangrove/internal/schema"

	"github.com/gin-gonic/gin"
)

func DashboardHandler(c *gin.Context) {
	controller.ResponseOk(c, gin.H{
		"sales": []schema.DashboardSalesItemResp{
			{
				Name:        2008,
				Clothes:     410,
				Food:        299,
				Electronics: 511,
			},
			{
				Name:        2009,
				Clothes:     449,
				Food:        199,
				Electronics: 500,
			},
			{
				Name:        2010,
				Clothes:     412,
				Food:        216,
				Electronics: 411,
			},
			{
				Name:        2011,
				Clothes:     400,
				Food:        399,
				Electronics: 501,
			},
			{
				Name:        2012,
				Clothes:     310,
				Food:        199,
				Electronics: 531,
			},
			{
				Name:        2013,
				Clothes:     450,
				Food:        269,
				Electronics: 521,
			},
			{
				Name:        2014,
				Clothes:     350,
				Food:        319,
				Electronics: 419,
			},
			{
				Name:        2015,
				Clothes:     418,
				Food:        279,
				Electronics: 410,
			},
		},
		"numbers": []schema.DashboardNumberItemResp{
			{
				Icon:   "pay-circle-o",
				Title:  "Online Review",
				Number: 2781,
				Color:  "#64ea91",
			},
			{
				Icon:   "team",
				Title:  "New Customers",
				Number: 3241,
				Color:  "#8fc9fb",
			},
			{
				Icon:   "message",
				Title:  "Active Projects",
				Number: 253,
				Color:  "#d897eb",
			},
			{
				Icon:   "shopping-cart",
				Title:  "Referrals",
				Number: 4324,
				Color:  "#f69899",
			},
		},
	})
}

func RoutesHandler(c *gin.Context) {
	controller.ResponseOk(c, []schema.MenuItemResp{
		{
			Id:    "1",
			Icon:  "dashboard",
			Name:  "Dashboard",
			Route: "/dashboard",
			Zh: &schema.MenuLocale{
				Name: "仪表盘",
			},
		},
		{
			Id:                 "2",
			BreadcrumbParentId: "1",
			Icon:               "demand",
			Name:               "Demands",
			Route:              "/demand",
			Zh: &schema.MenuLocale{
				Name: "需求管理",
			},
		},
		{
			Id:                 "21",
			MenuParentId:       "-1",
			BreadcrumbParentId: "2",
			Name:               "Demand Detail",
			Zh: &schema.MenuLocale{
				Name: "需求详情",
			},
			Route: "/demand/:id",
		},
	})
}
