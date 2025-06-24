package repository

import (
	dataSourceEngine "github.com/rentiansheng/dashboard/app/metrics/repository/data_source/engine"
	groupKeyEngine "github.com/rentiansheng/dashboard/app/metrics/repository/group_key/engine"
)

var (
	dsEngine dataSourceEngine.Engine
	gkEngine groupKeyEngine.Engine
)

func SetEngine(ds dataSourceEngine.Engine, gk groupKeyEngine.Engine) {
	dsEngine = ds
	gkEngine = gk
}

func DataSource() dataSourceEngine.Engine {
	return dsEngine
}

func GroupKey() groupKeyEngine.Engine {
	return gkEngine
}
