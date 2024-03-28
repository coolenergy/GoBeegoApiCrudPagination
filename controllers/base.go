package controllers

import (
	"github.com/astaxie/beego"
	"strconv"
)

type BaseController struct {
	beego.Controller
}

func (this *BaseController) getPagingParams() (page, pageSize int) {
	pageSizeStr := this.GetString("page_size")
	pageStr := this.GetString("page")

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 5
	}

	page, err = strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	return
}
