# Dashboard Backend API 接口文档

## 📋 概述

Dashboard Backend 提供了一套完整的 RESTful API 接口，用于数据源管理、数据查询、图表生成和分组管理等功能。

**基础信息：**
- 基础URL: `http://localhost:8080`
- 内容类型: `application/json`
- 字符编码: `UTF-8`

## 🔗 数据源管理接口

### 1. 获取数据源列表

**接口地址：** `GET /api/data/source/list`

**功能描述：** 获取所有可用的数据源列表

**请求参数：** 无

**响应格式：**
```json
{
  "retcode": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "name": "用户行为数据",
        "tips": "用户行为分析数据源",
        "remark": "用于分析用户行为的数据源",
        "data_name": "user_behavior",
        "status": 1,
        "data_type": "mysql",
        "sort_fields": "create_time desc",
        "formatter": "json",
        "template": "",
        "group_key": "group_key",
        "enable_group_key_name": true,
        "mtime": 1640995200,
        "ctime": 1640995200
      }
    ],
    "total": 1
  }
}
```

**字段说明：**
- `id`: 数据源ID
- `name`: 数据源显示名称
- `tips`: 数据源提示信息
- `remark`: 数据源备注
- `data_name`: 数据源名称
- `status`: 状态 (1: 启用, 0: 禁用)
- `data_type`: 数据源类型 (mysql, es等)
- `sort_fields`: 排序字段
- `formatter`: 格式化器
- `template`: 模板
- `group_key`: 分组键
- `enable_group_key_name`: 是否启用分组键名称
- `mtime`: 修改时间戳
- `ctime`: 创建时间戳

---

### 2. 获取数据源元数据

**接口地址：** `POST /api/data/source/meta`

**功能描述：** 获取指定数据源的字段元数据信息

**请求参数：**
```json
{
  "data_source_id": 1
}
```

**响应格式：**
```json
{
  "retcode": 0,
  "message": "success",
  "data": {
    "fields": [
      {
        "id": 1,
        "data_source_id": 1,
        "name": "数据源名字（es:索引名）",
        "display_name": "显示名字",
        "data_type": "数据类型",
        "output_data_type": "展示类型（预留信息）",
        "field_tips": "提示信息",
        "action": {
          "sort": true, 
          "filter": true, 
          "key": true,  
          "time": false,  
          "api": false,  
          "detail": true  
        },
        "enum": {
          "api": {
            "path": "",  
            "dynamic": false
          },
          "values": {
            "type": "kv/array", 
            "values": []
          }
        },
        "formatter": "",
        "template": "",
        "mtime": 1640995200,
        "ctime": 1640995200,
        "nested_path": ""
      }
    ]
  }
}
```

**字段说明：**
- `action`: 字段操作权限
  - `sort`: 是否可排序
  - `filter`: 是否可过滤
  - `key`: 是否可作为键
  - `time`: 是否为时间字段
  - `api`: 是否支持API枚举
  - `detail`: 是否显示在详情中
- `enum`: 枚举配置
  - `api`: API枚举配置
    - `path`: API地址
    - `dynamic`: 是否动态获取
  - `values`: 
  - `type`: 枚举类型 (kv/array)
  - `values`: 枚举值列表

---

### 3. 查询数据表

**接口地址：** `POST /api/data/source/query/table`

**功能描述：** 根据查询条件获取表格数据

**请求参数：**
```json
{
  "date": {
    "start": 1640995200,
    "end": 1641081600
  },
  "type": "100",
  "indexes": [
    {
      "name": "user_behavior",
      "output": {
        "time_field": "create_time",
        "sorts": ["create_time desc"],
        "cycle": 5,
        "aggregator": {
          "v2": false,
          "is_data_histogram": false,
          "is_top_n": false,
          "agg_type": 0,
          "name": "count"
        },
        "fields": ["user_id", "action", "create_time"],
        "x_axis_field": ""
      },
      "filters": [
        {
          "name": "user_filter",
          "join_operator": "and",
          "rules": [
            {
              "field": "user_id",
              "operator": "in",
              "values": [1, 2, 3],
              "pause": false,
              "not": false
            }
          ]
        }
      ]
    }
  ],
  "group_key_id": 1,
  "page": {
    "page_num": 1,
    "page_size": 20
  },
  "data_source_id": 1
}
```

**字段说明：**

**时间范围 (date):**
- `start`: 开始时间戳 (秒)
- `end`: 结束时间戳 (秒)

**查询类型 (type):**
- `1`: 柱状图
- `2`: 折线图
- `3`: 饼图
- `100`: 详情表格

**查询索引 (indexes):**
- `name`: 索引名称
- `output`: 输出配置
  - `time_field`: 时间字段
  - `sorts`: 排序字段
  - `cycle`: 时间周期类型
  - `aggregator`: 聚合配置
  - `fields`: 查询字段
  - `x_axis_field`: X轴字段
- `filters`: 过滤条件

**聚合类型 (agg_type):**
- `0`: 无聚合 (默认)
- `1`: Top N
- `2`: Terms
- `3`: 时间直方图 + Top N
- `4`: 时间直方图 + Terms
- `5`: X轴字段模式

**过滤操作符 (operator):**
- `eq`: 等于
- `ne`: 不等于
- `in`: 包含
- `nin`: 不包含
- `gt`: 大于
- `gte`: 大于等于
- `lt`: 小于
- `lte`: 小于等于
- `like`: 模糊匹配
- `nlike`: 不模糊匹配
- `between`: 范围查询
- `nbetween`: 不在范围内

**响应格式：**
```json
{
  "retcode": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "user_id": 1,
        "action": "login",
        "create_time": "2024-01-01 10:00:00"
      }
    ],
    "total": 100
  }
}
```

---

### 4. 生成图表数据

**接口地址：** `POST /api/data/source/query/chart`

**功能描述：** 根据查询条件生成图表数据

**请求参数：** 与查询数据表接口相同

**响应格式：**
```json
{
  "retcode": 0,
  "message": "success",
  "data": {
    "x_axis": ["2024-01-01", "2024-01-02", "2024-01-03"],
    "series": [
      {
        "name": "登录次数",
        "data": [10, 15, 20]
      },
      {
        "name": "注册次数", 
        "data": [5, 8, 12]
      }
    ]
  }
}
```

---

### 5. 获取字段枚举

**接口地址：** `POST /api/data/source/meta/enum`

**功能描述：** 获取指定字段的枚举值

**请求参数：**
```json
{
  "department_id": 1,
  "data_source_id": 1,
  "date": {
    "start": 1640995200,
    "end": 1641081600
  },
  "field_name": "status",
  "field_value": [],
  "field_value_by_key": false,
  "relations": [
    {
      "field": "user_type",
      "value": "vip"
    }
  ]
}
```

**响应格式：**
```json
{
  "retcode": 0,
  "message": "success",
  "data": {
    "type": "static",
    "values": [
     
    ]
  }
}
```

---

## 🔗 分组键管理接口

### 1. 获取根分组键

**接口地址：** `GET /api/group/key/roots`

**功能描述：** 获取所有根级别的分组键

**请求参数：** 无

**响应格式：**
```json
{
  "retcode": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "分组",
      "parent_id": 0,
      "order_index": 1,
      "remark": "相关分组",
      "mtime": 1640995200,
      "ctime": 1640995200
    }
  ]
}
```

---

### 2. 获取分组键树

**接口地址：** `GET /api/group/key/tree`

**功能描述：** 获取指定根分组键的完整树形结构

**请求参数：**
- `root_id`: 根分组键ID (查询参数)

**响应格式：**
```json
{
  "retcode": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "分组",
      "order_index": 1,
      "children": [
        {
          "id": 2,
          "name": "分组1",
          "order_index": 1,
          "children": []
        },
        {
          "id": 3,
          "name": "分组2",
          "order_index": 2,
          "children": []
        }
      ]
    }
  ]
}
```

---

## 📊 时间周期类型

系统支持以下时间周期类型：

| 类型值 | 名称 | 描述 |
|--------|------|------|
| 1 | 年 (Year) | 按年统计 |
| 2 | 季度 (Quarter) | 按季度统计 |
| 3 | 月 (Month) | 按月统计 |
| 4 | 周 (Week) | 按周统计 |

## 🔧 错误码说明

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
 
## 📝 使用示例

### 示例1：查询用户行为数据

```bash
curl -X POST http://localhost:8080/api/data/source/query/table \
  -H "Content-Type: application/json" \
  -d '{
    "date": {
      "start": 1640995200,
      "end": 1641081600
    },
    "type": "100",
    "indexes": [{
      "name": "user_behavior",
      "output": {
        "time_field": "create_time",
        "sorts": ["create_time desc"],
        "cycle": 5,
        "aggregator": {"agg_type": 0},
        "fields": ["user_id", "action", "create_time"]
      }
    }],
    "group_key_id": 1,
    "page": {"page_num": 1, "page_size": 20},
    "data_source_id": 1
  }'
```

### 示例2：生成用户活跃度图表

```bash
curl -X POST http://localhost:8080/api/data/source/query/chart \
  -H "Content-Type: application/json" \
  -d '{
    "date": {
      "start": 1640995200,
      "end": 1641081600
    },
    "type": "1",
    "indexes": [{
      "name": "user_behavior",
      "output": {
        "time_field": "create_time",
        "sorts": ["create_time desc"],
        "cycle": 5,
        "aggregator": {"agg_type": 0},
        "fields": ["action", "create_time"]
      }
    }],
    "group_key_id": 1,
    "data_source_id": 1
  }'
```

## 🔒 安全说明

1. **认证授权**: 建议在生产环境中添加适当的认证和授权机制
2. **参数验证**: 所有输入参数都会进行验证，请确保参数格式正确
3. **SQL注入防护**: 系统已内置SQL注入防护机制
4. **访问控制**: 建议配置适当的访问控制策略

## 📞 技术支持

如有问题，请联系开发团队或提交Issue。 
