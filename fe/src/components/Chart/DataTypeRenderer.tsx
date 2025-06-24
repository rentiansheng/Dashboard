import React from 'react';
import { Tag, Space, Typography, Progress, Statistic, Card, Row, Col } from 'antd';
import { 
  UserOutlined, 
  TeamOutlined, 
  BankOutlined, 
  ShopOutlined,
  BarChartOutlined,
  PieChartOutlined,
  LineChartOutlined,
  AreaChartOutlined
} from '@ant-design/icons';
import styled from 'styled-components';

const { Text, Title } = Typography;

const DataTypeCard = styled(Card)`
  margin-bottom: 16px;
  border-radius: 8px;
  
  .ant-card-head {
    background: #fafafa;
    border-bottom: 1px solid #f0f0f0;
  }
  
  .ant-card-body {
    padding: 16px;
  }
`;

const MetricValue = styled.div`
  font-size: 24px;
  font-weight: bold;
  color: #1890ff;
  margin: 8px 0;
`;

const DataTypeRenderer: React.FC<{
  data: any[];
  dataType: string;
  title?: string;
  onItemClick?: (item: any) => void;
}> = ({ data, dataType, title, onItemClick }) => {
  
  // 根据数据类型获取图标
  const getDataTypeIcon = (type: string) => {
    switch (type.toLowerCase()) {
      case 'user':
      case 'users':
        return <UserOutlined style={{ color: '#1890ff' }} />;
      case 'department':
      case 'dept':
        return <TeamOutlined style={{ color: '#52c41a' }} />;
      case 'organization':
      case 'org':
        return <BankOutlined style={{ color: '#722ed1' }} />;
      case 'product':
      case 'item':
        return <ShopOutlined style={{ color: '#fa8c16' }} />;
      case 'metric':
      case 'kpi':
        return <BarChartOutlined style={{ color: '#eb2f96' }} />;
      default:
        return <BarChartOutlined style={{ color: '#1890ff' }} />;
    }
  };

  // 根据数据类型获取颜色
  const getDataTypeColor = (type: string) => {
    switch (type.toLowerCase()) {
      case 'user':
      case 'users':
        return 'blue';
      case 'department':
      case 'dept':
        return 'green';
      case 'organization':
      case 'org':
        return 'purple';
      case 'product':
      case 'item':
        return 'orange';
      case 'metric':
      case 'kpi':
        return 'pink';
      default:
        return 'default';
    }
  };

  // 渲染用户类型数据
  const renderUserData = (data: any[]) => {
    return (
      <Row gutter={[16, 16]}>
        {data.map((item, index) => (
          <Col xs={24} sm={12} md={8} lg={6} key={index}>
            <Card 
              hoverable 
              size="small"
              onClick={() => onItemClick?.(item)}
              style={{ cursor: 'pointer' }}
            >
              <Space direction="vertical" style={{ width: '100%' }}>
                <Space>
                  <UserOutlined style={{ color: '#1890ff' }} />
                  <Text strong>{item.name || item.userName || '未知用户'}</Text>
                </Space>
                {item.value && (
                  <MetricValue>
                    {typeof item.value === 'number' 
                      ? item.value.toLocaleString() 
                      : item.value}
                  </MetricValue>
                )}
                {item.role && (
                  <Tag color="blue">{item.role}</Tag>
                )}
              </Space>
            </Card>
          </Col>
        ))}
      </Row>
    );
  };

  // 渲染部门类型数据
  const renderDepartmentData = (data: any[]) => {
    return (
      <Row gutter={[16, 16]}>
        {data.map((item, index) => (
          <Col xs={24} sm={12} md={8} lg={6} key={index}>
            <Card 
              hoverable 
              size="small"
              onClick={() => onItemClick?.(item)}
              style={{ cursor: 'pointer' }}
            >
              <Space direction="vertical" style={{ width: '100%' }}>
                <Space>
                  <TeamOutlined style={{ color: '#52c41a' }} />
                  <Text strong>{item.name || item.deptName || '未知部门'}</Text>
                </Space>
                {item.value && (
                  <MetricValue>
                    {typeof item.value === 'number' 
                      ? item.value.toLocaleString() 
                      : item.value}
                  </MetricValue>
                )}
                {item.memberCount && (
                  <Text type="secondary">
                    成员: {item.memberCount}人
                  </Text>
                )}
                {item.level && (
                  <Tag color="green">L{item.level}</Tag>
                )}
              </Space>
            </Card>
          </Col>
        ))}
      </Row>
    );
  };

  // 渲染指标类型数据
  const renderMetricData = (data: any[]) => {
    return (
      <Row gutter={[16, 16]}>
        {data.map((item, index) => (
          <Col xs={24} sm={12} md={8} lg={6} key={index}>
            <Card 
              hoverable 
              size="small"
              onClick={() => onItemClick?.(item)}
              style={{ cursor: 'pointer' }}
            >
              <Statistic
                title={item.name || item.metricName || '指标'}
                value={item.value || 0}
                precision={2}
                valueStyle={{ color: '#3f8600' }}
                prefix={<BarChartOutlined />}
              />
              {item.target && (
                <Progress 
                  percent={Math.round((item.value / item.target) * 100)} 
                  size="small" 
                  status={item.value >= item.target ? 'success' : 'active'}
                />
              )}
            </Card>
          </Col>
        ))}
      </Row>
    );
  };

  // 渲染产品类型数据
  const renderProductData = (data: any[]) => {
    return (
      <Row gutter={[16, 16]}>
        {data.map((item, index) => (
          <Col xs={24} sm={12} md={8} lg={6} key={index}>
            <Card 
              hoverable 
              size="small"
              onClick={() => onItemClick?.(item)}
              style={{ cursor: 'pointer' }}
            >
              <Space direction="vertical" style={{ width: '100%' }}>
                <Space>
                  <ShopOutlined style={{ color: '#fa8c16' }} />
                  <Text strong>{item.name || item.productName || '未知产品'}</Text>
                </Space>
                {item.value && (
                  <MetricValue>
                    {typeof item.value === 'number' 
                      ? `¥${item.value.toLocaleString()}` 
                      : item.value}
                  </MetricValue>
                )}
                {item.category && (
                  <Tag color="orange">{item.category}</Tag>
                )}
                {item.status && (
                  <Tag color={item.status === 'active' ? 'green' : 'red'}>
                    {item.status === 'active' ? '活跃' : '停用'}
                  </Tag>
                )}
              </Space>
            </Card>
          </Col>
        ))}
      </Row>
    );
  };

  // 渲染组织类型数据
  const renderOrganizationData = (data: any[]) => {
    return (
      <Row gutter={[16, 16]}>
        {data.map((item, index) => (
          <Col xs={24} sm={12} md={8} lg={6} key={index}>
            <Card 
              hoverable 
              size="small"
              onClick={() => onItemClick?.(item)}
              style={{ cursor: 'pointer' }}
            >
              <Space direction="vertical" style={{ width: '100%' }}>
                <Space>
                  <BankOutlined style={{ color: '#722ed1' }} />
                  <Text strong>{item.name || item.orgName || '未知组织'}</Text>
                </Space>
                {item.value && (
                  <MetricValue>
                    {typeof item.value === 'number' 
                      ? item.value.toLocaleString() 
                      : item.value}
                  </MetricValue>
                )}
                {item.type && (
                  <Tag color="purple">{item.type}</Tag>
                )}
                {item.scale && (
                  <Text type="secondary">
                    规模: {item.scale}
                  </Text>
                )}
              </Space>
            </Card>
          </Col>
        ))}
      </Row>
    );
  };

  // 根据数据类型选择渲染方式
  const renderDataByType = () => {
    const type = dataType.toLowerCase();
    
    if (type.includes('user')) {
      return renderUserData(data);
    } else if (type.includes('dept') || type.includes('department')) {
      return renderDepartmentData(data);
    } else if (type.includes('metric') || type.includes('kpi')) {
      return renderMetricData(data);
    } else if (type.includes('product') || type.includes('item')) {
      return renderProductData(data);
    } else if (type.includes('org') || type.includes('organization')) {
      return renderOrganizationData(data);
    } else {
      // 默认渲染方式
      return (
        <Row gutter={[16, 16]}>
          {data.map((item, index) => (
            <Col xs={24} sm={12} md={8} lg={6} key={index}>
              <Card 
                hoverable 
                size="small"
                onClick={() => onItemClick?.(item)}
                style={{ cursor: 'pointer' }}
              >
                <Space direction="vertical" style={{ width: '100%' }}>
                  <Text strong>{item.name || item.title || `项目 ${index + 1}`}</Text>
                  {item.value && (
                    <MetricValue>
                      {typeof item.value === 'number' 
                        ? item.value.toLocaleString() 
                        : item.value}
                    </MetricValue>
                  )}
                  {item.type && (
                    <Tag color={getDataTypeColor(item.type)}>{item.type}</Tag>
                  )}
                </Space>
              </Card>
            </Col>
          ))}
        </Row>
      );
    }
  };

  return (
    <DataTypeCard
      title={
        <Space>
          {getDataTypeIcon(dataType)}
          <Title level={4} style={{ margin: 0 }}>
            {title || `${dataType} 数据`}
          </Title>
          <Tag color={getDataTypeColor(dataType)}>
            {data.length} 项
          </Tag>
        </Space>
      }
    >
      {renderDataByType()}
    </DataTypeCard>
  );
};

export default DataTypeRenderer; 