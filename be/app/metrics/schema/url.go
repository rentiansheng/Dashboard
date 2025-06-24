package schema

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/rentiansheng/dashboard/app/metrics/service"
	"github.com/rentiansheng/dashboard/middleware"
)

type DataSourceSchema struct {
	dataSourceSvc service.DataSourceService
	groupKeySvc   service.GroupKeyService
}

type GroupKeySchema struct {
	groupKeySvc service.GroupKeyService
}

func NewDataSourceSchema(dataSourceSvc service.DataSourceService, groupKeySvc service.GroupKeyService) *DataSourceSchema {
	return &DataSourceSchema{
		dataSourceSvc: dataSourceSvc,
		groupKeySvc:   groupKeySvc,
	}
}

func (m *DataSourceSchema) Routes() *restful.WebService {

	ws := new(restful.WebService)
	ws.Path("/api/data/source")
	ws.Route(ws.GET("/list").To( middleware.Wrapper(m.List)))
	ws.Route(ws.POST("/meta").To(middleware.Wrapper(m.DataSourceMeta)))
	ws.Route(ws.POST("/query/table").To(middleware.Wrapper(m.DataSourceQuery)))
	ws.Route(ws.POST("/query/chart").To(middleware.Wrapper(m.DataSourceChart)))
	ws.Route(ws.POST("/meta/enum").To(middleware.Wrapper(m.DataSourceEnum)))

	return ws
}

func NewGroupKeySchema(groupKeySvc service.GroupKeyService) *GroupKeySchema {
	return &GroupKeySchema{
		groupKeySvc: groupKeySvc,
	}
}

func (g *GroupKeySchema) Routes() *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/api/group/key")
	ws.Route(ws.GET("/roots").To(middleware.Wrapper(g.Roots)))
	ws.Route(ws.GET("/tree").To(middleware.Wrapper(g.Tree)))

	return ws
}
