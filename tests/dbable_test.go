// @Title dbable_test.go
// @Description 代码测试
// @Author shigx 2021/6/10 2:28 下午
package tests

import (
	"fmt"
	"github.com/kuaileshi1/dbable"
	"testing"
	"time"
)

type DemoUser struct {
	ID   int
	Uid  int
	Name string
}

func TestSelect(t *testing.T) {
	// 一主一从配置
	masterDSN := &dbable.DSNConfig{
		Host:     "127.0.0.1",
		Port:     3306,
		UserName: "root",
		Password: "12345678",
		DbName:   "test",
		Charset:  "utf8",
	}
	slaveDSN := &dbable.DSNConfig{
		Host:     "127.0.0.1",
		Port:     3306,
		UserName: "root",
		Password: "12345678",
		DbName:   "test",
		Charset:  "utf8",
	}
	logConfig := &dbable.LogConfig{
		LogPath:                   "./logs",
		MaxAge:                    24 * 31,
		RotationTime:              24 * 7,
		LogLevel:                  "info",
		SlowThreshold:             200,
		IgnoreRecordNotFoundError: false,
	}
	config := &dbable.MysqlConfig{
		Driver:      "mysql",
		Master:      masterDSN,
		Slave:       []*dbable.DSNConfig{slaveDSN},
		MaxOpenCon:  10,
		MaxIdleCon:  5,
		MaxLifeTime: time.Second * 10,
		Logger:      logConfig,
	}
	dbable.Init(config, "test")
	db, _ := dbable.GetMysql("test")
	var result DemoUser
	db.Raw("SELECT id, uid, name FROM demo_user WHERE id = ?", 1).Scan(&result)
	fmt.Println(result)
	fmt.Printf("db:%v\n", db)

	// model
	var model DemoUser
	model.ID = 1
	db1, _ := dbable.GetMysql("test")
	db1.Table("demo_user").First(&model)

	fmt.Println(model)
	fmt.Printf("db:%v\n", db1)
}
