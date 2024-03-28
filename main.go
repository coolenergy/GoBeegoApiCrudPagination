package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/melardev/GoBeegoApiCrudPagination/infrastructure"
	_ "github.com/melardev/GoBeegoApiCrudPagination/routers"
	"github.com/melardev/GoBeegoApiCrudPagination/seeds"
	"os"
)

func main() {

	e := godotenv.Load() // Load .env file
	if e != nil {
		fmt.Print(e)
		os.Exit(0)
	}

	// if in develop mode
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	if err := infrastructure.Initialize(); err != nil {
		panic(err.Error())
	}

	seeds.Seed()

	// Print SQL statements executed
	orm.Debug = true
	beego.Run()
}
