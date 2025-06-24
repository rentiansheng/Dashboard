import http, { ApiResponse } from '../utils/http';

// ==================== 公共类型定义 ====================

/** 时间范围 */
export interface DateRange {
  start: number;
  end: number;
}

/** 分页参数 */
export interface PageParams {
  page_num: number;
  page_size: number;
}

/** 过滤规则 */
export interface FilterRule {
  field: string;
  operator: string;
  values: string[];
  /** 查询条件是否启用 */
  pause: boolean;
}

/** 过滤器组 */
export interface FilterGroup {
  rules: FilterRule[];
  /** 图表是作为series 名字 */
  name: string;
}

/** 基础查询参数 */
export interface BaseQueryRequest {
  date: DateRange;
  data_source_id: number;
  group_key_id: number;
}

// ==================== Table 查询相关 ====================

/** Table 查询输出配置 */
export interface TableOutput {
  time_field: string;
  /** 需要展示字段 */
  fields?: string[];
  sorts?: string[];
}

/** Table 查询索引配置 */
export interface TableIndex {
  /** 数据源 */
  name: string;
  output: TableOutput;
  filters: FilterGroup[];
}

// ==================== Metric 查询相关 ====================

/** 聚合查询项配置 */
export interface AggregatorTerms {
  /** 字段名字 */
  field?: string;
  value?: string[];
}

/** 聚合查询范围配置 */
export interface AggregatorRange {
  /** 字段筛选范围，from-to 前闭后开 */
  field?: string;
  /** 任意类型， >= */
  from?: any;
  /** 任意类型， < */
  to?: any;
  /** 任意类型， >= 后端未支持 */
  gte: any;
  /** 任意类型， <=后端未支持 */
  lte: any;
}

/** 聚合过滤器 */
export interface AggregatorFilter {
  /** series 名字 */
  name?: string;
  /** in 查询参数 */
  terms?: AggregatorTerms[];
  range?: AggregatorRange[];
  /** 聚合方法， 0,1:count, 2:avg, 3: sum */
  agg_type?: number;
  /** agg_type=(1,2) 计算的字段 */
  extra_agg_field?: string;
}

/** Top N 配置 */
export interface TopNConfig {
  /** top n 使用的field */
  field: string;
  /** 默认值10 */
  num?: number;
}

/** 聚合器配置 */
export interface Aggregator {
  /** 新版本参数为了保证兼容，看到参数的情况下，需值为true */
  v2?: boolean;
  /** 是否未x时间轴 */
  is_date_histogram?: boolean;
  /** 是否未top k 模式 */
  is_top_n?: boolean;
  /** 用来是sum,avg 或者是多个series聚合 */
  filters?: AggregatorFilter[];
  /** is_top_n=true 有效 */
  top_n?: TopNConfig;
}

/** Metric 查询输出配置 */
export interface MetricOutput {
  time_field: string;
  /** cycle 与aggregator互斥， 限定x轴 */
  cycle?: number;
  /** cycle 与aggregator互斥， 限定x轴 */
  aggregator?: Aggregator;
}

/** Metric 查询索引配置 */
export interface MetricIndex {
  /** series name */
  name: string;
  output: MetricOutput;
  filters: FilterGroup[];
}

// ==================== API 请求接口 ====================

export interface ApiDataSourceMetaRequest {
  /** 数据源 */
  data_source_id: number;
}

export interface ApiDataSourceQueryTableRequest extends BaseQueryRequest {
  indexes: TableIndex[];
  page: PageParams;
}

export interface ApiDataSourceQueryChartRequest extends BaseQueryRequest {
  /** bar stacked bar,stacked line,pie(agg_type=1,2,仅在该模式可用） */
  type: string;
  indexes: MetricIndex[];
}

export interface ApiDataSourceMetaResponse  {
  fields: ApiDataSourceMetaField[];
}

export interface ApiDataSourceMetaFieldAction {
  /**
       * 是否可以作为排序字段
       */
  sort?: boolean;
  /**
   * 是否可以作为查询字段
   */
  filter?: boolean;
  /**
   * 是否可以作为分组key
   */
  key?: boolean;
  /**
   * 是否未时间字段，时间字段可以做获取数据字段，或者分组key
   */
  time?: boolean;
  /**
   * 是否未枚举字段。枚举字段需要用api_name 来获取展示字段,  hooks 中enum 对应需要调用的api
   */
  api?: boolean;
  /**
   * 是否可以作为展示字段
   */
  detail: boolean;
}
export interface ApiDataSourceMetaField {

    name: string;
    display_name: string;
    data_type: string;
    /**
     * 详情时候字段值
     */
    output_data_type: string;
    tips: string;
    action: ApiDataSourceMetaFieldAction;
    enum: {
      /**
       * api 地址，用来获取枚举值
       */
      api: {
        path: string;
        /**
         * 是否需要动态获取数据
         */
        dynamic: boolean;
      };
      /**
       * 字段的枚举，
       */
      values: {
        /**
         * 可能数据类型，kv|array, kv 标识是key(display name)/value(field value)， array 用来做suggest(key,value相同),	是array 或者
         */
        type: string;
        Values: string;
      };
    };
    /**
     * 非空的时候，不能作为agg filter
     */
    nested: string;
}

export interface ApiDataSourceMetaEnumRequest  extends BaseQueryRequest{
  group_key_id: number;
  data_source_id: number;
  date: {
    start: number;
    end: number;
  };
  field_name: string;
  /**
   * 用来做模糊查询或者按值查询
   */
  field_value?: {}[];
  relations?: {
    field?: string;
    value?: string;
  }[];
}

export interface ApiDataSourceMetaEnumResponse {
  /**
   *
   */
  values: {};
  /**
   * 可能数据类型，kv|array, kv 标识是key(display name)/value(field value)， array 用来做suggest(key,value相同),	是array 或者
   */
  type: string;
}

export interface ApiDataSourceResponse {
  list?: {
    /**
     * 数据源id
     */
    id: number;
    /**
     * 展示的名字
     */
    name: string;
    /**
     * 数据索引名或者表名
     */
    value: string;
    tips: string;
    /**
     * 详情展示的时候，字段的顺序， 只默认值，多个字段用逗号分割，用户可以在页面，自定义需要展示的字段
     */
    sort_fields: string;
  }[];
}


// ==================== 数据源服务类 ====================

export class DataSourceService {
  /** 获取数据源列表 */
  static async getDataSourceList(): Promise<ApiResponse<any>> {
    return http.get<ApiDataSourceResponse[]>('/api/data/source/list');
  }

  /** 获取数据源元数据 */
  static async getDataSourceMeta(id: number): Promise<ApiResponse<any>> {
    const request = {
      data_source_id: id,
    } as ApiDataSourceMetaRequest;
    return http.post<ApiDataSourceMetaResponse[]>('/api/data/source/meta', request);
  }

    /** 获取数据源元数据 */
  static async getDataSourceMetaEnum(request: ApiDataSourceMetaEnumRequest): Promise<ApiResponse<any>> {
    return http.post<ApiDataSourceMetaEnumResponse>('/api/data/source/meta/enum', request);
  }
  

  /** 查询数据源表格数据 */
  static async getDataSourceQueryTable(request: ApiDataSourceQueryTableRequest): Promise<ApiResponse<any>> {
    return http.post<any>('/api/data/source/query/table', request);
  }

  /** 查询数据源指标数据 */
  static async getDataSourceQueryChart(request: ApiDataSourceQueryChartRequest): Promise<ApiResponse<any>> {
    return http.post<any>('/api/data/source/query/chart', request);
  }
}

export default DataSourceService;