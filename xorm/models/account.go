package models

import (
	"errors"
	"log"
	"os"

	"github.com/go-xorm/xorm"

	_ "github.com/go-sql-driver/mysql"
)

type Account struct {
	Id      int64
	Name    string `xorm:"unique"`
	Balance float64
	Version int `xorm:"version"`
}

var Engine *xorm.Engine

func init() {
	var err error
	//Engine, err = xorm.NewEngine("mysql", "root:1qaz@WSX123@tcp(192.168.1.96:3306)/foo?charset=utf8")
	Engine, err = xorm.NewEngine("mysql", "root:root@tcp(localhost:3306)/foo?charset=utf8")
	if err != nil {
		log.Fatal(err)
	}

	err = Engine.Sync(new(Account))
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Create("sql.log")
	if err != nil {
		log.Fatal(err)
	}
	Engine.SetLogger(xorm.NewSimpleLogger(f))
	Engine.ShowSQL(true)

	cacher := xorm.NewLRUCacher(xorm.NewMemoryStore(), 1000)
	Engine.SetDefaultCacher(cacher)

}

// 插入数据
func NewAccount(name string, balance float64) error {
	_, err := Engine.Insert(&Account{Name: name, Balance: balance})
	if err != nil {
		return err
	}
	return nil
}

// 查询数据
func GetAccount(id int64) (*Account, error) {
	account := &Account{}
	has, err := Engine.Id(id).Get(account)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, errors.New("Account not found")
	}
	return account, nil
}

// 存款
func MakeDeposit(id int64, deposit float64) (*Account, error) {
	a, err := GetAccount(id)
	if err != nil {
		return nil, err
	}
	a.Balance += deposit
	_, err = Engine.Update(a)
	return a, err
}

// 取款
func MakeWithdraw(id int64, withdraw float64) (*Account, error) {
	a, err := GetAccount(id)
	if err != nil {
		return nil, err
	}
	if a.Balance <= withdraw {
		return nil, errors.New("Not enough balance")
	}
	a.Balance -= withdraw
	_, err = Engine.Update(a)
	return a, err
}

// 转账
func MakeTransfer(from, to int64, balance float64) error {
	fromAccount, err := GetAccount(from)
	if err != nil {
		return err
	}

	toAccount, err := GetAccount(to)
	if err != nil {
		return err
	}

	if fromAccount.Balance <= balance {
		return errors.New("Not enought balance")
	}

	fromAccount.Balance -= balance
	toAccount.Balance += balance
	// 添加事务支持
	sess := Engine.NewSession()
	defer sess.Close()
	if err := sess.Begin(); err != nil {
		return err
	}
	if _, err = sess.Update(fromAccount); err != nil {
		sess.Rollback()
		return err
	} else if _, err = sess.Update(toAccount); err != nil {
		sess.Rollback()
		return err
	}
	return sess.Commit()
}

// 根据id正序获取数据
func ListAccountsAscId() (as []*Account, err error) {
	err = Engine.Asc("id").Find(&as)
	return as, err
}

// 根据balance倒叙获取
func ListAccountsDescBalance() (as []*Account, err error) {
	err = Engine.Desc("balance").Find(&as)
	return as, err
}

// 根据id删除账户
func DeleteAccount(id int64) error {
	_, err := Engine.Delete(&Account{Id: id})
	return err
}

// count
func CountAccount() (int64, error) {
	return Engine.Count(&Account{})
}

// 迭代所有行
func IterateAccount(account *Account, iterateFn func(idx int, bean interface{}) error) {
	Engine.Iterate(account, iterateFn)
}

// 获取所有行
func GetAccountRows(account *Account) (*xorm.Rows, error) {
	return Engine.Rows(account)
}
