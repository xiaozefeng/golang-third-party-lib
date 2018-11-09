package main

import (
	"fmt"

	"github.com/xiaozefeng/golang-third-party-lib/xorm/models"
)

const prompt = `Plearse enter number of operation:
1. Create new account
2. Show detail of account
3. Deposit //存钱
4. Withdraw  // 取钱
5. Make transfer 
6. List account by id
7. List account by balance
8. Delete account
9. Exit
`

func main() {
Exit:
	for {
		fmt.Println()
		fmt.Println(prompt)

		var num int
		fmt.Scanf("%d\n", &num)
		switch num {
		case 1:
			fmt.Println("Please input <name> <balance>:")
			var name string
			var balance float64
			fmt.Scanf("%s %f\n", &name, &balance)
			if err := models.NewAccount(name, balance); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Create Account Successful")
			}
		case 2:
			fmt.Println("Please input <id>:")
			var id int64
			fmt.Scanf("%d\n", &id)
			if account, err := models.GetAccount(id); err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("account = %#+v\n", account)
			}
		case 3:
			fmt.Println("Please input <id> <deposit>:")
			var id int64
			var deposit float64
			fmt.Scanf("%d %f\n", &id, &deposit)
			if account, err := models.MakeDeposit(id, deposit); err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("account = %+v\n", account)
			}

		case 4:
			fmt.Println("Please input <id> <withdraw>:")
			var id int64
			var withdraw float64
			fmt.Scanf("%d %f\n", &id, &withdraw)
			if account, err := models.MakeWithdraw(id, withdraw); err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("account = %+v\n", account)
			}
		case 5:
			fmt.Println("<fromAccount id> transfer <balance> to <toAccount id>")
			var from, to int64
			var balance float64
			fmt.Scanf("%d %f %d\n", &from, &balance, &to)
			if err := models.MakeTransfer(from, to, balance); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("转账成功")
			}
		case 6:
			if as, err := models.ListAccountsAscId(); err != nil {
				fmt.Println(err)
			} else {
				for i, a := range as {
					fmt.Printf("id:%d , %#v\n", i, a)
				}
			}

		case 7:
			if as, err := models.ListAccountsDescBalance(); err != nil {
				fmt.Println(err)
			} else {
				for i, a := range as {
					fmt.Printf("%d ,%#v\n", i, a)
				}
			}
		case 8:
			fmt.Println("Please input <id>:")
			var id int64
			fmt.Scanf("%d\n", &id)
			if err := models.DeleteAccount(id); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("删除成功")
			}
		case 9:
			break Exit
		}
	}
	fmt.Println("Exited")
}
