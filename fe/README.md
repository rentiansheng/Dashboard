# Dashboard 数据可视化平台

一个基于 React + TypeScript + Ant Design 构建的现代化数据可视化仪表板平台，提供强大的数据分析和图表展示功能。

## 🚀 功能特性

### 核心功能
- **数据可视化**: 支持多种图表类型，包括折线图、柱状图、饼图等
- **数据分析**: 强大的数据查询和过滤功能
- **部门管理**: 完整的部门层级结构管理
- **实时数据**: 支持实时数据更新和展示
- **响应式设计**: 适配各种屏幕尺寸

### 技术特性
- **TypeScript**: 完整的类型安全支持
- **组件化**: 高度模块化的组件设计
- **状态管理**: 基于 Context API 的全局状态管理
- **路由管理**: React Router 实现单页应用路由
- **样式系统**: 支持 SCSS 和 Styled Components
- **代码规范**: ESLint 代码质量检查

## 🛠 技术栈

### 前端框架
- **React 18.2.0** - 用户界面库
- **TypeScript 5.2.2** - 类型安全的 JavaScript
- **Vite 5.0.8** - 快速构建工具

### UI 组件库
- **Ant Design 5.12.8** - 企业级 UI 设计语言
- **Ant Design Pro Components 2.8.9** - 高级组件库
- **Ant Design Charts 2.2.7** - 数据可视化图表库
- **Ant Design Icons 5.2.6** - 图标库

### 状态管理与数据获取
- **React Context API** - 全局状态管理
- **TanStack React Query 5.80.7** - 数据获取和缓存
- **Axios 1.9.0** - HTTP 客户端

### 工具库
- **Day.js 1.11.13** - 日期处理
- **Lodash 4.17.21** - 实用工具库
- **Data Forge 1.10.4** - 数据处理
- **Styled Components 6.1.18** - CSS-in-JS

### 开发工具
- **ESLint** - 代码质量检查
- **Sass Embedded 1.89.2** - CSS 预处理器

## 📦 安装与运行

### 环境要求
- Node.js >= 16.0.0
- npm >= 8.0.0

### 安装依赖
```bash
npm install
```

### 开发环境运行
```bash
npm run dev
```
访问 http://localhost:5173

### 生产环境构建
```bash
npm run build
```

### 预览生产构建
```bash
npm run preview
```

### 代码检查
```bash
npm run lint
```

## 📁 项目结构

```
fe/
├── docs/                    # 项目文档
│   ├── api-loop-fix.md     # API 循环请求修复文档
│   ├── data-transformation.md # 数据转换文档
│   ├── global-state-management.md # 全局状态管理文档
│   ├── multi-value-input.md # 多值输入文档
│   └── smart-data-renderer.md # 智能数据渲染文档
├── src/
│   ├── components/         # 组件目录
│   │   ├── Antd/          # Ant Design 相关组件
│   │   ├── Chart/         # 图表组件
│   │   ├── Header/        # 头部组件
│   │   └── Table/         # 表格组件
│   ├── context/           # React Context
│   │   ├── AuthContext.tsx # 认证上下文
│   │   ├── GlobalContext.tsx # 全局状态上下文
│   │   └── index.ts
│   ├── hooks/             # 自定义 Hooks
│   │   ├── useChartStore.ts # 图表存储 Hook
│   │   ├── useDepartmentQuery.ts # 部门查询 Hook
│   │   ├── useGroupKeyTree.ts # 部门树 Hook
│   │   ├── useSelectedGroupKey.ts # 选中部门 Hook
│   │   └── useTableColumnState.ts # 表格列状态 Hook
│   ├── pages/             # 页面组件
│   │   └── DataAnalysis/  # 数据分析页面
│   ├── services/          # API 服务
│   │   ├── dataSourceService.ts # 数据源服务
│   │   ├── groupKeyService.ts # 部门服务
│   │   ├── index.ts
│   │   └── types.ts
│   ├── types/             # TypeScript 类型定义
│   │   ├── chart.ts       # 图表类型
│   │   ├── datasource.ts  # 数据源类型
│   │   ├── groupKey.ts    # 部门类型
│   │   └── index.ts
│   ├── utils/             # 工具函数
│   │   ├── dataTransformer.tsx # 数据转换器
│   │   ├── dataTypeDetector.ts # 数据类型检测
│   │   ├── dayjs.tsx      # 日期工具
│   │   ├── enum.ts        # 枚举定义
│   │   ├── http.ts        # HTTP 工具
│   │   ├── index.ts
│   │   └── tree.ts        # 树形数据处理
│   ├── App.tsx            # 主应用组件
│   ├── main.tsx           # 应用入口
│   └── index.css          # 全局样式
├── index.html             # HTML 模板
├── package.json           # 项目配置
├── tsconfig.json          # TypeScript 配置
├── vite.config.ts         # Vite 配置
└── README.md              # 项目说明
```

## 🎯 核心功能说明

### 1. 部门管理系统

#### 部门选择器 (`DeptSelect`)
- 支持多选部门
- 树形结构展示
- 搜索功能
- 选中状态持久化

#### 部门信息展示 (`GroupKeyInfo`)
- 显示当前选中部门信息
- 支持显示完整部门路径
- 多选时显示汇总信息
- 子部门数量统计

### 2. 数据分析功能

#### 数据源管理
- 支持多种数据源类型
- 动态字段配置
- 数据类型自动检测

#### 图表配置
- 多种图表类型支持
- 动态配置选项
- 实时数据更新

#### 条件过滤
- 灵活的过滤条件设置
- 多条件组合
- 动态条件生成

### 3. 全局状态管理

使用 React Context API 实现全局状态管理，主要管理：
- 选中的部门信息
- 用户认证状态
- 应用配置信息

## 🔧 开发指南

### 添加新组件

1. 在 `src/components/` 下创建组件目录
2. 使用 TypeScript 编写组件
3. 添加必要的类型定义
4. 编写组件文档

### 添加新页面

1. 在 `src/pages/` 下创建页面组件
2. 在 `src/App.tsx` 中添加路由配置
3. 更新菜单配置

### 添加新 API 服务

1. 在 `src/services/` 下创建服务文件
2. 定义相关的 TypeScript 类型
3. 使用 `src/utils/http.ts` 进行 HTTP 请求

### 状态管理

使用 `useGlobalContext` Hook 访问全局状态：

```tsx
import { useGlobalContext } from '@/context/GlobalContext';

const MyComponent = () => {
  const { selectedGroupKey, setSelectedGroupKey } = useGlobalContext();
  // 使用全局状态
};
```


## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

 
**注意**: 这是一个企业内部项目，请确保遵循公司的开发规范和保密要求。 