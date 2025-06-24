package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/rentiansheng/dashboard/app/metrics/repository"
	dataSourceEngine "github.com/rentiansheng/dashboard/app/metrics/repository/data_source/impl/mysql"
	groupKeyEngine "github.com/rentiansheng/dashboard/app/metrics/repository/group_key/impl/mysql"
	"github.com/rentiansheng/dashboard/pkg/config"
	"github.com/rentiansheng/ges"
)

func initDependResource() error {

	if config.Config.Engine.EngineType == "mysql" {
		var dsDB *gorm.DB
		var gkDB *gorm.DB
		if config.Config.Engine.IsAll {
			db, err := initMysqlEngine(config.Config.Engine.DefaultEngineConfig.Mysql)
			if err != nil {
				return err
			}
			dsDB = db
			gkDB = db
		} else {
			var err error
			dsDB, err = initMysqlEngine(config.Config.Engine.DataSourceEngine.Mysql)
			if err != nil {
				return err
			}
			gkDB, err = initMysqlEngine(config.Config.Engine.GroupKeyEngine.Mysql)
			if err != nil {
				return err
			}
		}
		dsEngine := dataSourceEngine.NewDataSourceRepo(dsDB)

		gkEngine := groupKeyEngine.NewGroupKeyRepo(gkDB)
		repository.SetEngine(dsEngine, gkEngine)
	}

	if err := initES(); err != nil {
		return err
	}

	return nil

}

func initMysqlEngine(dns string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func initES() error {
	cfg := config.Config.ES
	return ges.InitClientWithCfg(cfg.Addrs, cfg.Username, cfg.Password)
}
