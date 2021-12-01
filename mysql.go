// @Title mysql.go
// @Description mysql gorm 连接操作
// @Author shigx 2021/6/10 11:28 上午
package dbable

import (
	"database/sql"
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"strings"
	"sync"
	"time"
)

var (
	configMap    map[string]*MysqlConfig
	debug        bool
	mysqlPollMap sync.Map
)

// gorm mysql配置定义
type MysqlConfig struct {
	Driver      string     `yaml:"driver"`
	Master      *DSNConfig `yaml:"master"`
	Slave       []*DSNConfig
	MaxOpenCon  int           `yaml:"maxOpenCon"`
	MaxIdleCon  int           `yaml:maxIdleCon`
	MaxLifeTime time.Duration `yaml:"maxLifeTime"`
}

// DSN配置定义
type DSNConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DbName   string `yaml:"dbName"`
	UserName string `yaml:"userName"`
	Password string `yaml:"password"`
	Charset  string `yaml:"charset"`
}

// @Description 初始化mysql组件，添加配置，支持多实例配置
// @Auth shigx
// @Date 2021/6/10 1:10 下午
// @param
// @return
func Init(config *MysqlConfig, instance string) {
	if configMap == nil {
		configMap = make(map[string]*MysqlConfig)
	}

	configMap[instance] = config
}

// @Description 从连接池获取数据库链接
// @Auth shigx
// @Date 2021/6/10 2:06 下午
// @param instance string 实例名称
// @return
func GetMysql(instance string) (db *gorm.DB, err error) {
	if mysqlPool, ok := mysqlPollMap.Load(instance); ok {
		db = mysqlPool.(*gorm.DB)
		_, err = db.DB()
		if err == nil {
			return
		}
		mysqlPollMap.Delete(instance)
		return GetMysql(instance)
	} else {
		db, err = newConnect(instance)
		if err == nil {
			mysqlPollMap.Store(instance, db)
		}
		return
	}
}

// @Description 创建mysql数据库链接
// @Auth shigx
// @Date 2021/6/10 1:57 下午
// @param instance string 数据库实例名称
// @return
func newConnect(instance string) (db *gorm.DB, err error) {
	if _, ok := configMap[instance]; !ok {
		err = errors.New("mysql newConnect Error: not found config data in configMap")
		return nil, err
	}
	config := configMap[instance]
	masterDSN := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=True&loc=Local&charset=%s",
		config.Master.UserName,
		config.Master.Password,
		config.Master.Host,
		config.Master.Port,
		config.Master.DbName,
		config.Master.Charset)

	var replicas []gorm.Dialector
	for _, conf := range config.Slave {
		slaveDSN := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?parseTime=True&loc=Local&charset=%s",
			conf.UserName,
			conf.Password,
			conf.Host,
			conf.Port,
			conf.DbName,
			conf.Charset)
		replicas = append(replicas, mysql.Open(slaveDSN))
	}
	if strings.ToUpper(config.Driver) == "MYSQL" {
		if db, err = gorm.Open(mysql.Open(masterDSN), &gorm.Config{
			SkipDefaultTransaction: true,
		}); err != nil {
			err = errors.New("gorm open failed " + masterDSN)
			return nil, err
		} else {
			err = db.Use(dbresolver.Register(dbresolver.Config{
				Sources:  []gorm.Dialector{mysql.Open(masterDSN)},
				Replicas: replicas,
			}).SetMaxOpenConns(config.MaxOpenCon).
				SetMaxIdleConns(config.MaxIdleCon).
				SetConnMaxLifetime(config.MaxLifeTime * time.Second))
			if debug {
				db = db.Debug()
			}

			var DB *sql.DB
			DB, _ = db.DB()
			if pErr := DB.Ping(); pErr != nil {
				return nil, pErr
			}
		}

		return db, err
	}

	return nil, errors.New("Driver is not mysql ")
}

// @Description 设置开启debug模式
// @Auth shigx
// @Date 2021/6/10 1:06 下午
// @param
// @return
func SetDebugOn() {
	debug = true
}
