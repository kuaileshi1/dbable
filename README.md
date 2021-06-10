## 介绍
基于grom v2 封装的mysql连接池工具

## 安装
`go get github.com/kuaileshi1/dbable`

## 使用
```go
type DemoUser struct {
	ID   int
	Uid  int
	Name string
}

// 一主一从配置
masterDSN := &dbable.DSNConfig{
    Host:     "127.0.0.1",
    Port:     3306,
    UserName: "root",
    Password: "123456",
    DbName:   "test",
    Charset:  "utf8",
}
slaveDSN := &dbable.DSNConfig{
    Host:     "127.0.0.1",
    Port:     3306,
    UserName: "root",
    Password: "123456",
    DbName:   "test",
    Charset:  "utf8",
}
config := &dbable.MysqlConfig{
    Driver:      "mysql",
    Master:      masterDSN,
    Slave:       []*dbable.DSNConfig{slaveDSN},
    MaxOpenCon:  10,
    MaxIdleCon:  5,
    MaxLifeTime: time.Second * 10,
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
```
