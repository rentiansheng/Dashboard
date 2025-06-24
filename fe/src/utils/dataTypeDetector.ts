// 数据类型检测器
export interface DataTypeInfo {
  type: string;
  confidence: number;
  fields: string[];
  description: string;
}

// 数据类型定义
export const DATA_TYPES = {
  USER: 'user',
  DEPARTMENT: 'department',
  ORGANIZATION: 'organization',
  PRODUCT: 'product',
  METRIC: 'metric',
  FINANCIAL: 'financial',
  TIME_SERIES: 'time_series',
  GEOGRAPHIC: 'geographic',
  CATEGORICAL: 'categorical',
  NUMERICAL: 'numerical',
} as const;

// 字段模式匹配规则
const FIELD_PATTERNS = {
  [DATA_TYPES.USER]: {
    patterns: [
      /user(name|id|_name|_id)/i,
      /employee(name|id|_name|_id)/i,
      /staff(name|id|_name|_id)/i,
      /member(name|id|_name|_id)/i,
      /name/i,
      /email/i,
      /phone/i,
      /role/i,
      /position/i,
    ],
    required: ['name'],
    optional: ['email', 'phone', 'role', 'position'],
  },
  [DATA_TYPES.DEPARTMENT]: {
    patterns: [
      /dept(name|id|_name|_id)/i,
      /department(name|id|_name|_id)/i,
      /team(name|id|_name|_id)/i,
      /division(name|id|_name|_id)/i,
      /branch(name|id|_name|_id)/i,
      /level/i,
      /membercount/i,
      /headcount/i,
    ],
    required: ['name'],
    optional: ['level', 'memberCount', 'headCount'],
  },
  [DATA_TYPES.ORGANIZATION]: {
    patterns: [
      /org(name|id|_name|_id)/i,
      /organization(name|id|_name|_id)/i,
      /company(name|id|_name|_id)/i,
      /corporation(name|id|_name|_id)/i,
      /enterprise(name|id|_name|_id)/i,
      /scale/i,
      /size/i,
      /industry/i,
      /sector/i,
    ],
    required: ['name'],
    optional: ['scale', 'size', 'industry', 'sector'],
  },
  [DATA_TYPES.PRODUCT]: {
    patterns: [
      /product(name|id|_name|_id)/i,
      /item(name|id|_name|_id)/i,
      /goods(name|id|_name|_id)/i,
      /sku/i,
      /category/i,
      /brand/i,
      /model/i,
      /status/i,
      /price/i,
      /cost/i,
    ],
    required: ['name'],
    optional: ['category', 'brand', 'status', 'price'],
  },
  [DATA_TYPES.METRIC]: {
    patterns: [
      /metric(name|id|_name|_id)/i,
      /kpi(name|id|_name|_id)/i,
      /indicator(name|id|_name|_id)/i,
      /target/i,
      /goal/i,
      /threshold/i,
      /baseline/i,
      /performance/i,
    ],
    required: ['name', 'value'],
    optional: ['target', 'goal', 'threshold'],
  },
  [DATA_TYPES.FINANCIAL]: {
    patterns: [
      /revenue/i,
      /income/i,
      /profit/i,
      /loss/i,
      /expense/i,
      /cost/i,
      /budget/i,
      /amount/i,
      /currency/i,
      /price/i,
    ],
    required: ['value'],
    optional: ['currency', 'amount'],
  },
  [DATA_TYPES.TIME_SERIES]: {
    patterns: [
      /date/i,
      /time/i,
      /timestamp/i,
      /period/i,
      /month/i,
      /year/i,
      /quarter/i,
      /week/i,
      /day/i,
    ],
    required: ['date', 'value'],
    optional: ['period', 'time'],
  },
  [DATA_TYPES.GEOGRAPHIC]: {
    patterns: [
      /location/i,
      /address/i,
      /city/i,
      /state/i,
      /country/i,
      /region/i,
      /province/i,
      /district/i,
      /latitude/i,
      /longitude/i,
      /geo/i,
    ],
    required: ['location'],
    optional: ['city', 'country', 'region'],
  },
  [DATA_TYPES.CATEGORICAL]: {
    patterns: [
      /category/i,
      /type/i,
      /class/i,
      /group/i,
      /tag/i,
      /label/i,
      /status/i,
      /state/i,
    ],
    required: ['category'],
    optional: ['type', 'status'],
  },
  [DATA_TYPES.NUMERICAL]: {
    patterns: [
      /value/i,
      /number/i,
      /count/i,
      /total/i,
      /sum/i,
      /average/i,
      /mean/i,
      /median/i,
      /score/i,
      /rating/i,
    ],
    required: ['value'],
    optional: ['count', 'total'],
  },
};

// 检测单个字段是否匹配模式
const matchFieldPattern = (fieldName: string, patterns: RegExp[]): boolean => {
  return patterns.some(pattern => pattern.test(fieldName));
};

// 计算字段匹配度
const calculateFieldMatch = (fields: string[], patterns: RegExp[]): number => {
  const matchedFields = fields.filter(field => matchFieldPattern(field, patterns));
  return matchedFields.length / fields.length;
};

// 检测数据类型
export const detectDataType = (data: any[]): DataTypeInfo[] => {
  if (!data || data.length === 0) {
    return [];
  }

  // 获取所有字段名
  const allFields = new Set<string>();
  data.forEach(item => {
    if (typeof item === 'object' && item !== null) {
      Object.keys(item).forEach(key => allFields.add(key));
    }
  });
  
  const fields = Array.from(allFields);
  
  // 分析每个数据类型
  const typeAnalysis: DataTypeInfo[] = [];
  
  Object.entries(FIELD_PATTERNS).forEach(([type, config]) => {
    const fieldMatch = calculateFieldMatch(fields, config.patterns);
    
    // 检查必需字段
    const hasRequiredFields = config.required.every(reqField => 
      fields.some(field => field.toLowerCase().includes(reqField.toLowerCase()))
    );
    
    // 检查可选字段
    const optionalFieldCount = config.optional.filter(optField =>
      fields.some(field => field.toLowerCase().includes(optField.toLowerCase()))
    ).length;
    
    // 计算置信度
    let confidence = fieldMatch * 0.6; // 字段匹配度占60%
    
    if (hasRequiredFields) {
      confidence += 0.3; // 必需字段存在占30%
    }
    
    confidence += (optionalFieldCount / config.optional.length) * 0.1; // 可选字段占10%
    
    // 数据内容分析
    const contentAnalysis = analyzeDataContent(data, type);
    confidence = (confidence + contentAnalysis) / 2;
    
    if (confidence > 0.3) { // 只返回置信度大于30%的结果
      typeAnalysis.push({
        type,
        confidence: Math.round(confidence * 100) / 100,
        fields: fields.filter(field => matchFieldPattern(field, config.patterns)),
        description: getDataTypeDescription(type),
      });
    }
  });
  
  // 按置信度排序
  return typeAnalysis.sort((a, b) => b.confidence - a.confidence);
};

// 分析数据内容
const analyzeDataContent = (data: any[], type: string): number => {
  if (data.length === 0) return 0;
  
  const sample = data.slice(0, Math.min(10, data.length));
  let score = 0;
  
  switch (type) {
    case DATA_TYPES.USER:
      // 检查是否包含用户相关信息
      const hasUserInfo = sample.some(item => 
        item.name && typeof item.name === 'string' && item.name.length > 0
      );
      if (hasUserInfo) score += 0.5;
      
      const hasEmail = sample.some(item => 
        item.email && /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(item.email)
      );
      if (hasEmail) score += 0.3;
      
      const hasRole = sample.some(item => 
        item.role && typeof item.role === 'string'
      );
      if (hasRole) score += 0.2;
      break;
      
    case DATA_TYPES.DEPARTMENT:
      // 检查是否包含部门相关信息
      const hasDeptInfo = sample.some(item => 
        item.name && typeof item.name === 'string' && item.name.length > 0
      );
      if (hasDeptInfo) score += 0.4;
      
      const hasLevel = sample.some(item => 
        item.level && (typeof item.level === 'number' || typeof item.level === 'string')
      );
      if (hasLevel) score += 0.3;
      
      const hasMemberCount = sample.some(item => 
        item.memberCount && typeof item.memberCount === 'number'
      );
      if (hasMemberCount) score += 0.3;
      break;
      
    case DATA_TYPES.METRIC:
      // 检查是否包含指标相关信息
      const hasValue = sample.some(item => 
        item.value !== undefined && item.value !== null
      );
      if (hasValue) score += 0.4;
      
      const hasTarget = sample.some(item => 
        item.target !== undefined && item.target !== null
      );
      if (hasTarget) score += 0.3;
      
      const hasNumericValue = sample.some(item => 
        typeof item.value === 'number'
      );
      if (hasNumericValue) score += 0.3;
      break;
      
    case DATA_TYPES.FINANCIAL:
      // 检查是否包含财务相关信息
      const hasAmount = sample.some(item => 
        item.value && typeof item.value === 'number' && item.value > 0
      );
      if (hasAmount) score += 0.5;
      
      const hasCurrency = sample.some(item => 
        item.currency && typeof item.currency === 'string'
      );
      if (hasCurrency) score += 0.3;
      
      const hasPrice = sample.some(item => 
        item.price && typeof item.price === 'number'
      );
      if (hasPrice) score += 0.2;
      break;
      
    default:
      score = 0.5; // 默认分数
  }
  
  return score;
};

// 获取数据类型描述
const getDataTypeDescription = (type: string): string => {
  const descriptions: Record<string, string> = {
    [DATA_TYPES.USER]: '用户数据，包含用户基本信息',
    [DATA_TYPES.DEPARTMENT]: '部门数据，包含组织架构信息',
    [DATA_TYPES.ORGANIZATION]: '组织数据，包含公司或机构信息',
    [DATA_TYPES.PRODUCT]: '产品数据，包含商品或服务信息',
    [DATA_TYPES.METRIC]: '指标数据，包含KPI和性能指标',
    [DATA_TYPES.FINANCIAL]: '财务数据，包含收入和支出信息',
    [DATA_TYPES.TIME_SERIES]: '时间序列数据，包含时间维度的数据',
    [DATA_TYPES.GEOGRAPHIC]: '地理数据，包含位置信息',
    [DATA_TYPES.CATEGORICAL]: '分类数据，包含类别和标签',
    [DATA_TYPES.NUMERICAL]: '数值数据，包含统计和计算数据',
  };
  
  return descriptions[type] || '未知数据类型';
};

// 获取推荐的数据类型
export const getRecommendedDataType = (data: any[]): DataTypeInfo | null => {
  const analysis = detectDataType(data);
  return analysis.length > 0 ? analysis[0] : null;
};

// 验证数据类型
export const validateDataType = (data: any[], expectedType: string): boolean => {
  const analysis = detectDataType(data);
  return analysis.some(result => 
    result.type === expectedType && result.confidence > 0.7
  );
}; 