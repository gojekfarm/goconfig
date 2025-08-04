package goconfig

import "time"

type DBConfig struct {
	driver          string
	url             string
	slaveUrl        string
	idleConn        int
	maxConn         int
	connMaxLifetime time.Duration
}

func (self *DBConfig) Driver() string {
	return self.driver
}

func (self *DBConfig) Url() string {
	return self.url
}

func (self *DBConfig) SlaveUrl() string {
	return self.slaveUrl
}

func (self *DBConfig) MaxConn() int {
	return self.maxConn
}

func (self *DBConfig) IdleConn() int {
	return self.idleConn
}

func (self *DBConfig) ConnMaxLifetime() time.Duration {
	return self.connMaxLifetime
}

func LoadDbConf(accessor ConfigAccessor) *DBConfig {
	return &DBConfig{
		driver:          getStringOrPanic(accessor, "DB_DRIVER"),
		url:             getStringOrPanic(accessor, "DB_URL"),
		slaveUrl:        getStringOrPanic(accessor, "DB_SLAVE_URL"),
		maxConn:         getIntOrPanic(accessor, "DB_MAX_CONN"),
		idleConn:        getIntOrPanic(accessor, "DB_IDLE_CONN"),
		connMaxLifetime: time.Duration(getIntOrPanic(accessor, "DB_CONN_MAX_LIFETIME")) * time.Second,
	}
}
