import React, { useState } from 'react';
import { Routes, Route, useNavigate, useLocation } from 'react-router-dom';
import {
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  DashboardOutlined,
  SettingOutlined,
} from '@ant-design/icons';
import { Layout, Menu, Button, theme, Space } from 'antd';
import type { MenuProps } from 'antd';

// 导入页面组件
import DataAnalysis from './pages/DataAnalysis';
import { DeptSelect } from './components/Header/DeptSelect';
import { GroupKeyInfo } from './components/Header/GroupKeyInfo';
import { GlobalProvider } from './context/GlobalContext';

const { Header, Sider, Content } = Layout;

type MenuItem = Required<MenuProps>['items'][number];

const App: React.FC = () => {
  const [collapsed, setCollapsed] = useState(false);
  const navigate = useNavigate();
  const location = useLocation();
  const {
    token: { colorBgContainer, borderRadiusLG },
  } = theme.useToken();

  const menuItems: MenuItem[] = [
    {
      key: '/',
      icon: <DashboardOutlined />,
      label: 'Dashboard',
    },
 
  ];

  const handleMenuClick: MenuProps['onClick'] = (e) => {
    navigate(e.key);
  };

  return (
    <GlobalProvider>
      <Layout style={{ minHeight: '100vh' }}>
        <Sider trigger={null} collapsible collapsed={collapsed}>
          <div className="logo">
            {collapsed ? 'D' : 'Dashboard'}
          </div>
          <Menu
            theme="dark"
            mode="inline"
            selectedKeys={[location.pathname]}
            items={menuItems}
            onClick={handleMenuClick}
          />
        </Sider>
        <Layout>
          <Header style={{ padding: 0, background: colorBgContainer }}>
            <Button
              type="text"
              icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
              onClick={() => setCollapsed(!collapsed)}
              style={{
                fontSize: '16px',
                width: 64,
                height: 64,
              }}
            />

            <Space style={{ float: 'right', marginRight: 16 }}>
              <GroupKeyInfo />
              <DeptSelect />
            </Space>
            
          </Header>
          <Content
            style={{
              margin: '24px 16px',
              padding: 24,
              minHeight: 280,
              background: colorBgContainer,
              borderRadius: borderRadiusLG,
            }}
          >
            <Routes>
              <Route path="/" element={<DataAnalysis />} />
              <Route path="/data-analysis" element={<DataAnalysis />} />
            </Routes>
          </Content>
        </Layout>
      </Layout>
    </GlobalProvider>
  );
};

export default App; 