import { defineConfig, loadEnv } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'

// https://vitejs.dev/config/
export default defineConfig(({ command, mode }) => {
  // 根据当前工作目录中的 `mode` 加载 .env 文件
  const env = loadEnv(mode, process.cwd())
  
  // 使用 API 基础路径
  const baseURL = env.VITE_API_BASE_URL;

  // 使用应用标题
  const appTitle = env.VITE_APP_TITLE;

  // 使用环境标识
  const isDev = env.VITE_APP_ENV === 'development';
  
  return {
    plugins: [react()],
    resolve: {
      alias: {
        '@': path.resolve(__dirname, './src')
      }
    },
    server: {
      port: 3000,
      open: true,
      proxy: {
        '/api': {
          target: 'http://127.0.0.1:8080',
          changeOrigin: true,
        }
      }
    },
    // 环境变量配置
    envPrefix: 'VITE_',
    define: {
      __APP_ENV__: JSON.stringify(env.VITE_APP_ENV)
    }
  }
}) 