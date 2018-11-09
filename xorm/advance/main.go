package main

import (
	"fmt"
	"log"

	"github.com/xiaozefeng/golang-third-party-lib/xorm/models"
)

var printFn = func(idx int, bean interface{}) error {
	fmt.Printf("%d, %#v\n", idx, bean.(*models.Account))
	return nil
}

func main() {
	count, err := models.CountAccount()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Account count= ", count)

	// 统计
	for i := count; i < 10; i++ {
		if err = models.NewAccount(fmt.Sprintf("joe%d", i), float64(i)*100); err != nil {
			log.Fatalf("Fail to create account: %v\n", err)
		}
	}

	fmt.Println("Query all records:")
	//如何需要添加过滤条件可以在account对象中添加属性
	// 例如 Account{Id:1}
	models.IterateAccount(&models.Account{}, printFn)

	fmt.Println()
	a := new(models.Account)
	rows, err := models.GetAccountRows(a)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(a); err != nil {
			log.Fatalf("Fail to get row: %+v\n", err)
		} else {
			fmt.Printf("%#v\n", a)
		}
	}

	fmt.Println("\nOnly query id,name")
	e := models.Engine
	// 查询指定字段
	e.Cols("id,name").Iterate(new(models.Account), printFn)

	// 排除某些字段
	fmt.Println("\nQuery all but name")
	e.Omit("name").Iterate(new(models.Account), printFn)

	fmt.Println("\nOffset 2 and limit 3")
	e.Limit(3, 2).Iterate(new(models.Account), printFn)

}
