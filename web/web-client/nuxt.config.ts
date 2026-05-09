import { globalConfig } from "./src/utils/global-config";
import fs from "node:fs";
import path from "node:path";

// 当 globalConfig.https 为 true 时，开发服务器以 HTTPS 启动（Nuxt 需要文件路径字符串）
const certDir = path.resolve(process.cwd(), "../../../../server/alnitak/server/conf");
const keyPath = path.join(certDir, "private.key");
const certPath = path.join(certDir, "public.crt");
// 环境变量 NUXT_DEV_HTTP=true 强制 HTTP，覆盖 globalConfig.https
const useHttps = process.env.NUXT_DEV_HTTP !== 'true' && globalConfig.https;
const httpsConfig = useHttps && fs.existsSync(keyPath)
  ? { key: keyPath, cert: certPath }
  : undefined;

export default defineNuxtConfig({
  compatibilityDate: '2026-05-09',
  devServer: {
    https: httpsConfig,
  },

  modules: [
    '@element-plus/nuxt',
    '@pinia/nuxt',
  ],
  app: {
    head: {
      title: globalConfig.title,
      meta: [
        {
          "name": "viewport",
          "content": "width=device-width, initial-scale=1"
        },
        {
          "charset": "utf-8"
        },
        {
          "name": "keywords",
          "content": globalConfig.keywords
        },
        {
          "name": "description",
          "content": globalConfig.description
        },
      ],
      link: [
        { rel: "icon", type: "image/x-icon", href: "/favicon.ico" }
      ]
    }
  },
  plugins: [
    {
      src: '@/plugins/wang-editor',
      mode: 'client',
    },
    {
      src: '@/plugins/auth-sync.client',
      mode: 'client',
    },
    {
      src: '@/plugins/theme-init.client',
      mode: 'client',
    },
    {
      src: '@/plugins/auth-init.server',
      mode: 'server',
    },
    {
      src: '@/plugins/error-handler.server',
      mode: 'server',
    },
  ],
  css: [
    'element-plus/dist/index.css',
    'element-plus/theme-chalk/dark/css-vars.css',
    '~/assets/styles/element.scss'
  ],
  devtools: { enabled: true },
  srcDir: 'src/',
  vite: {
    define: {
      __VUE_OPTIONS_API__: true,
      __VUE_PROD_DEVTOOLS__: false,
      __VUE_PROD_HYDRATION_MISMATCH_DETAILS__: false,
    },
    optimizeDeps: {
      include: [
        // Element Plus 组件及样式预构建
        'element-plus/es/components/form/index',
        'element-plus/es/components/form/style/css',
        'element-plus/es/components/form-item/style/css',
        'element-plus/es/components/input/index',
        'element-plus/es/components/input/style/css',
        'element-plus/es/components/radio/index',
        'element-plus/es/components/radio/style/css',
        'element-plus/es/components/radio-button/style/css',
        'element-plus/es/components/radio-group/style/css',
        'element-plus/es/components/date-picker/index',
        'element-plus/es/components/date-picker/style/css',
        'element-plus/es/components/button/index',
        'element-plus/es/components/button/style/css',
        'element-plus/es/components/pagination/index',
        'element-plus/es/components/pagination/style/css',
        'element-plus/es/components/icon/index',
        'element-plus/es/components/icon/style/css',
        'element-plus/es/components/scrollbar/index',
        'element-plus/es/components/scrollbar/style/css',
        'element-plus/es/components/dropdown/index',
        'element-plus/es/components/dropdown/style/css',
        'element-plus/es/components/dropdown-item/style/css',
        'element-plus/es/components/dropdown-menu/style/css',
        'element-plus/es/components/switch/index',
        'element-plus/es/components/switch/style/css',
        'element-plus/es/components/dialog/index',
        'element-plus/es/components/dialog/style/css',
        'element-plus/es/components/progress/index',
        'element-plus/es/components/progress/style/css',
        'element-plus/es/components/upload/index',
        'element-plus/es/components/upload/style/css',
        'element-plus/es/components/tag/index',
        'element-plus/es/components/tag/style/css',
        'element-plus/es/components/tabs/index',
        'element-plus/es/components/tabs/style/css',
        'element-plus/es/components/tab-pane/style/css',
        'element-plus/es/components/popconfirm/index',
        'element-plus/es/components/popconfirm/style/css',
        'element-plus/es/components/checkbox/index',
        'element-plus/es/components/checkbox/style/css',
        'element-plus/es/components/loading/index',
        'element-plus/es/components/loading/style/css',
        'element-plus/es/components/message/index',
        'element-plus/es/components/message/style/css',
        'element-plus/es/components/message-box/index',
        'element-plus/es/components/message-box/style/css',
        'element-plus/es/components/skeleton/index',
        'element-plus/es/components/skeleton/style/css',
        'element-plus/es/components/skeleton-item/style/css',
        'element-plus/es/components/select/index',
        'element-plus/es/components/select/style/css',
        'element-plus/es/components/option/style/css',
        'element-plus/es/components/tooltip/index',
        'element-plus/es/components/tooltip/style/css',
        'element-plus/es/components/result/index',
        'element-plus/es/components/result/style/css',
        // 其他常用依赖
        'moment',
        'hls.js',
        'wplayer-next',
        'vue-picture-cropper',
        'spark-md5',
        '@icon-park/vue-next',
        'axios',
        'js-cookie',
        '@wangeditor/editor-for-vue',
        'artplayer',
        'dashjs',
        'artplayer-plugin-danmuku'
      ]
    },
    server: {
      headers: {
        // 避免切回标签页触发强刷时，SFC 子模块（如 ?vue&type=style）被错误缓存复用
        'Cache-Control': 'no-store',
      },
      hmr: {
        overlay: false
      },
      proxy: {
        '/api': {
          target: `http${globalConfig.https ? 's' : ''}://${globalConfig.domain}`,
          changeOrigin: true,
          ws: true,
          // 可选：需要去掉 /api 前缀时，设置 API_PROXY_REWRITE=remove 并解开下一行
          // rewrite: process.env.API_PROXY_REWRITE === 'remove' ? (path) => path.replace(/^\/api/, '') : undefined,
          // 网络不稳定时适当拉长代理超时
          timeout: 30000,
          proxyTimeout: 30000,
          configure: (proxy, _options) => {
            // 捕获代理错误，防止 ECONNRESET 导致 unhandledRejection
            proxy.on('error', (err, _req, res) => {
              // 静默处理网络错误，减少日志噪音
              if (err.message?.includes('ECONNRESET') || err.message?.includes('ETIMEDOUT')) {
                return;
              }
              console.warn('[proxy error]', err.message);
              // 如果响应还没有发送，返回 502
              if (res && 'writeHead' in res && !res.headersSent) {
                res.writeHead(502, { 'Content-Type': 'application/json' });
                res.end(JSON.stringify({ code: 502, msg: 'Proxy error: ' + err.message }));
              }
            });
            // WebSocket 代理错误处理
            proxy.on('proxyReqWs', (proxyReq, _req, socket) => {
              socket.on('error', () => {
                // WebSocket 连接错误，静默忽略
              });
            });
            proxy.on('proxyReq', (proxyReq, req, _res) => {
              // 仅对非 WebSocket 请求设置 keep-alive
              // WebSocket 需要 Connection: Upgrade，不能覆盖
              if (!req.headers.upgrade) {
                proxyReq.setHeader('Connection', 'keep-alive');
              }
            });
          },
        }
      }
    }
  }
})
