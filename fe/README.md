# Dashboard - 左侧菜单管理系统

一个基于 React + TypeScript + Ant Design 构建的现代化管理后台系统，具有响应式左侧菜单导航。

## 功能特性

- 🎨 现代化 UI 设计，基于 Ant Design 组件库
- 📱 响应式布局，支持移动端和桌面端
- 🔧 TypeScript 支持，提供完整的类型检查
- 🚀 基于 Vite 的快速开发和构建
- 📊 丰富的页面模块：仪表板、用户管理、视频管理等
- 🎯 可折叠的左侧菜单导航
- 🔄 路由导航支持

## 技术栈

- **前端框架**: React 18
- **开发语言**: TypeScript
- **UI 组件库**: Ant Design 5.x
- **路由管理**: React Router DOM 6.x
- **构建工具**: Vite
- **图标库**: Ant Design Icons

## 项目结构

```
src/
├── pages/           # 页面组件
│   ├── Dashboard.tsx    # 仪表板
│   ├── Users.tsx        # 用户管理
│   ├── Videos.tsx       # 视频管理
│   ├── Upload.tsx       # 文件上传
│   ├── Documents.tsx    # 文档中心
│   └── Settings.tsx     # 系统设置
├── App.tsx          # 主应用组件
├── main.tsx         # 应用入口
└── index.css        # 全局样式
```

## 快速开始

### 安装依赖

```bash
npm install
```

### 启动开发服务器

```bash
npm run dev
```

应用将在 http://localhost:3000 启动

### 构建生产版本

```bash
npm run build
```

### 预览生产构建

```bash
npm run preview
```

## 页面功能

### 🏠 仪表板
- 数据统计卡片
- 项目进度展示
- 系统信息监控

### 👥 用户管理
- 用户列表展示
- 搜索和筛选功能
- 用户状态管理

### 🎬 视频管理
- 视频卡片展示
- 视频状态标识
- 响应式网格布局

### 📤 文件上传
- 拖拽上传支持
- 上传进度显示
- 文件状态管理

### 📚 文档中心
- 文档列表展示
- 分类标签管理
- 搜索和下载功能

### ⚙️ 系统设置
- 基本配置管理
- 通知设置
- 备份配置

## 开发说明

项目使用现代化的前端开发技术栈，具有良好的代码组织结构和类型安全保障。所有组件都采用 TypeScript 编写，提供完整的类型提示和检查。

## 许可证

MIT License 