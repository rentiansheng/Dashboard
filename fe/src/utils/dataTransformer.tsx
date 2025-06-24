import dayjs from 'dayjs';
import { Tag, Space, Typography } from 'antd';
import { 
  UserOutlined, 
  TeamOutlined, 
  BankOutlined, 
  ShopOutlined,
  BarChartOutlined,
  PieChartOutlined,
  LineChartOutlined,
  AreaChartOutlined,
  CalendarOutlined,
  DollarOutlined,
  PercentageOutlined,
  NumberOutlined
} from '@ant-design/icons';
import React from 'react';

const { Text } = Typography;

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
  DATE: 'date',
  CURRENCY: 'currency',
  PERCENTAGE: 'percentage',
  EMAIL: 'email',
  PHONE: 'phone',
  URL: 'url',
} as const;

// 数据转换配置接口
export interface DataTransformConfig {
  type: string;
  format?: string;
  precision?: number;
  currency?: string;
  locale?: string;
  showIcon?: boolean;
  color?: string;
  size?: 'small' | 'default' | 'large';
  customRender?: (value: any, record: any) => React.ReactNode;
}

// 基础数据转换函数
export const transformData = (value: any, config: DataTransformConfig): any => {
  if (value === null || value === undefined) {
    return 'N/A';
  }

  switch (config.type) {
    case DATA_TYPES.DATE:
      return transformDate(value, config);
    case DATA_TYPES.CURRENCY:
      return transformCurrency(value, config);
    case DATA_TYPES.PERCENTAGE:
      return transformPercentage(value, config);
    case DATA_TYPES.NUMERICAL:
      return transformNumber(value, config);
    case DATA_TYPES.EMAIL:
      return transformEmail(value, config);
    case DATA_TYPES.PHONE:
      return transformPhone(value, config);
    case DATA_TYPES.URL:
      return transformUrl(value, config);
    case DATA_TYPES.USER:
      return transformUser(value, config);
    case DATA_TYPES.DEPARTMENT:
      return transformDepartment(value, config);
    case DATA_TYPES.ORGANIZATION:
      return transformOrganization(value, config);
    case DATA_TYPES.PRODUCT:
      return transformProduct(value, config);
    case DATA_TYPES.METRIC:
      return transformMetric(value, config);
    case DATA_TYPES.FINANCIAL:
      return transformFinancial(value, config);
    default:
      return value;
  }
};

// 日期转换
const transformDate = (value: any, config: DataTransformConfig): React.ReactNode => {
  const format = config.format || 'YYYY-MM-DD HH:mm:ss';
  const date = dayjs(value);
  
  if (!date.isValid()) {
    return 'Invalid Date';
  }

  if (config.showIcon) {
    return (
      <Space>
        <CalendarOutlined style={{ color: config.color || '#1890ff' }} />
        <Text>{date.format(format)}</Text>
      </Space>
    );
  }

  return date.format(format);
};

// 货币转换
const transformCurrency = (value: any, config: DataTransformConfig): React.ReactNode => {
  const currency = config.currency || 'CNY';
  const locale = config.locale || 'zh-CN';
  const precision = config.precision || 2;
  
  const numValue = parseFloat(value);
  if (isNaN(numValue)) {
    return 'Invalid Amount';
  }

  const formattedValue = new Intl.NumberFormat(locale, {
    style: 'currency',
    currency: currency,
    minimumFractionDigits: precision,
    maximumFractionDigits: precision,
  }).format(numValue);

  if (config.showIcon) {
    return (
      <Space>
        <DollarOutlined style={{ color: config.color || '#52c41a' }} />
        <Text>{formattedValue}</Text>
      </Space>
    );
  }

  return formattedValue;
};

// 百分比转换
const transformPercentage = (value: any, config: DataTransformConfig): React.ReactNode => {
  const precision = config.precision || 2;
  const numValue = parseFloat(value);
  
  if (isNaN(numValue)) {
    return 'Invalid Percentage';
  }

  const percentage = (numValue * 100).toFixed(precision) + '%';

  if (config.showIcon) {
    return (
      <Space>
        <PercentageOutlined style={{ color: config.color || '#722ed1' }} />
        <Text>{percentage}</Text>
      </Space>
    );
  }

  return percentage;
};

// 数字转换
const transformNumber = (value: any, config: DataTransformConfig): React.ReactNode => {
  const precision = config.precision || 0;
  const numValue = parseFloat(value);
  
  if (isNaN(numValue)) {
    return 'Invalid Number';
  }

  const formattedValue = new Intl.NumberFormat(config.locale || 'zh-CN', {
    minimumFractionDigits: precision,
    maximumFractionDigits: precision,
  }).format(numValue);

  if (config.showIcon) {
    return (
      <Space>
        <NumberOutlined style={{ color: config.color || '#1890ff' }} />
        <Text>{formattedValue}</Text>
      </Space>
    );
  }

  return formattedValue;
};

// 邮箱转换
const transformEmail = (value: any, config: DataTransformConfig): React.ReactNode => {
  const email = String(value);
  const isValidEmail = /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
  
  if (!isValidEmail) {
    return <Text type="danger">{email}</Text>;
  }

  return (
    <a href={`mailto:${email}`} style={{ color: config.color || '#1890ff' }}>
      {email}
    </a>
  );
};

// 电话转换
const transformPhone = (value: any, config: DataTransformConfig): React.ReactNode => {
  const phone = String(value);
  const formattedPhone = phone.replace(/(\d{3})(\d{4})(\d{4})/, '$1-$2-$3');
  
  return (
    <a href={`tel:${phone}`} style={{ color: config.color || '#52c41a' }}>
      {formattedPhone}
    </a>
  );
};

// URL转换
const transformUrl = (value: any, config: DataTransformConfig): React.ReactNode => {
  const url = String(value);
  const isValidUrl = /^https?:\/\/.+/.test(url);
  
  if (!isValidUrl) {
    return <Text type="danger">{url}</Text>;
  }

  return (
    <a href={url} target="_blank" rel="noopener noreferrer" style={{ color: config.color || '#1890ff' }}>
      {url}
    </a>
  );
};

// 用户数据转换
const transformUser = (value: any, config: DataTransformConfig): React.ReactNode => {
  const userData = typeof value === 'object' ? value : { name: value };
  
  return (
    <Space>
      <UserOutlined style={{ color: config.color || '#1890ff' }} />
      <Text strong>{userData.name || userData.userName || '未知用户'}</Text>
      {userData.role && (
        <Tag color="blue" size={config.size || 'small'}>{userData.role}</Tag>
      )}
    </Space>
  );
};

// 部门数据转换
const transformDepartment = (value: any, config: DataTransformConfig): React.ReactNode => {
  const deptData = typeof value === 'object' ? value : { name: value };
  
  return (
    <Space>
      <TeamOutlined style={{ color: config.color || '#52c41a' }} />
      <Text strong>{deptData.name || deptData.deptName || '未知部门'}</Text>
      {deptData.level && (
        <Tag color="green" size={config.size || 'small'}>L{deptData.level}</Tag>
      )}
    </Space>
  );
};

// 组织数据转换
const transformOrganization = (value: any, config: DataTransformConfig): React.ReactNode => {
  const orgData = typeof value === 'object' ? value : { name: value };
  
  return (
    <Space>
      <BankOutlined style={{ color: config.color || '#722ed1' }} />
      <Text strong>{orgData.name || orgData.orgName || '未知组织'}</Text>
      {orgData.type && (
        <Tag color="purple" size={config.size || 'small'}>{orgData.type}</Tag>
      )}
    </Space>
  );
};

// 产品数据转换
const transformProduct = (value: any, config: DataTransformConfig): React.ReactNode => {
  const productData = typeof value === 'object' ? value : { name: value };
  
  return (
    <Space>
      <ShopOutlined style={{ color: config.color || '#fa8c16' }} />
      <Text strong>{productData.name || productData.productName || '未知产品'}</Text>
      {productData.category && (
        <Tag color="orange" size={config.size || 'small'}>{productData.category}</Tag>
      )}
    </Space>
  );
};

// 指标数据转换
const transformMetric = (value: any, config: DataTransformConfig): React.ReactNode => {
  const metricData = typeof value === 'object' ? value : { value };
  
  return (
    <Space>
      <BarChartOutlined style={{ color: config.color || '#eb2f96' }} />
      <Text strong>{metricData.name || metricData.metricName || '指标'}</Text>
      <Text type="success">{metricData.value || value}</Text>
    </Space>
  );
};

// 财务数据转换
const transformFinancial = (value: any, config: DataTransformConfig): React.ReactNode => {
  const financialData = typeof value === 'object' ? value : { value };
  
  return (
    <Space>
      <DollarOutlined style={{ color: config.color || '#52c41a' }} />
      <Text strong>{financialData.name || '金额'}</Text>
      <Text type="success">
        {transformCurrency(financialData.value || value, {
          ...config,
          type: DATA_TYPES.CURRENCY
        })}
      </Text>
    </Space>
  );
};

// 自动检测数据类型并转换
export const autoTransformData = (value: any, fieldName: string): React.ReactNode => {
  if (value === null || value === undefined) {
    return 'N/A';
  }

  const fieldNameLower = fieldName.toLowerCase();
  
  // 根据字段名自动检测类型
  if (fieldNameLower.includes('date') || fieldNameLower.includes('time')) {
    return transformData(value, { type: DATA_TYPES.DATE, showIcon: true });
  }
  
  if (fieldNameLower.includes('price') || fieldNameLower.includes('amount') || fieldNameLower.includes('cost')) {
    return transformData(value, { type: DATA_TYPES.CURRENCY, showIcon: true });
  }
  
  if (fieldNameLower.includes('rate') || fieldNameLower.includes('percent')) {
    return transformData(value, { type: DATA_TYPES.PERCENTAGE, showIcon: true });
  }
  
  if (fieldNameLower.includes('email')) {
    return transformData(value, { type: DATA_TYPES.EMAIL });
  }
  
  if (fieldNameLower.includes('phone') || fieldNameLower.includes('tel')) {
    return transformData(value, { type: DATA_TYPES.PHONE });
  }
  
  if (fieldNameLower.includes('url') || fieldNameLower.includes('link')) {
    return transformData(value, { type: DATA_TYPES.URL });
  }
  
  if (fieldNameLower.includes('user') || fieldNameLower.includes('name')) {
    return transformData(value, { type: DATA_TYPES.USER, showIcon: true });
  }
  
  if (fieldNameLower.includes('dept') || fieldNameLower.includes('department')) {
    return transformData(value, { type: DATA_TYPES.DEPARTMENT, showIcon: true });
  }
  
  if (fieldNameLower.includes('org') || fieldNameLower.includes('organization')) {
    return transformData(value, { type: DATA_TYPES.ORGANIZATION, showIcon: true });
  }
  
  if (fieldNameLower.includes('product') || fieldNameLower.includes('item')) {
    return transformData(value, { type: DATA_TYPES.PRODUCT, showIcon: true });
  }
  
  if (fieldNameLower.includes('metric') || fieldNameLower.includes('kpi')) {
    return transformData(value, { type: DATA_TYPES.METRIC, showIcon: true });
  }
  
  // 数字类型检测
  if (typeof value === 'number' || !isNaN(parseFloat(value))) {
    return transformData(value, { type: DATA_TYPES.NUMERICAL, showIcon: true });
  }
  
  return value;
};

// 批量转换表格数据
export const transformTableData = (data: any[], columnConfigs: Record<string, DataTransformConfig>): any[] => {
  return data.map(record => {
    const transformedRecord = { ...record };
    
    Object.keys(columnConfigs).forEach(fieldName => {
      if (record.hasOwnProperty(fieldName)) {
        transformedRecord[fieldName] = transformData(record[fieldName], columnConfigs[fieldName]);
      }
    });
    
    return transformedRecord;
  });
};

// 创建表格列配置
export const createTableColumnConfig = (
  fieldName: string, 
  config: DataTransformConfig
): any => {
  return {
    title: config.title || fieldName,
    dataIndex: fieldName,
    key: fieldName,
    render: (value: any, record: any) => {
      if (config.customRender) {
        return config.customRender(value, record);
      }
      return transformData(value, config);
    },
    ...config
  };
}; 