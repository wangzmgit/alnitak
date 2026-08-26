import { globalConfig } from "./src/utils/global-config";
import fs from "node:fs";
import path from "node:path";
import https from "node:https";

// 当 globalConfig.https 为 true 时，开发服务器以 HTTPS 启动（Nuxt 需要文件路径字符串）
const certDir = path.resolve(process.cwd(), "../../../../server/alnitak/server/conf");
const keyPath = path.join(certDir, "private.key");
const certPath = path.join(certDir, "public.crt");
// 环境变量 NUXT_DEV_HTTP=true 强制 HTTP，覆盖 globalConfig.https
const useHttps = process.env.NUXT_DEV_HTTP !== 'true' && globalConfig.https;
const httpsConfig = useHttps && fs.existsSync(keyPath)
  ? { key: keyPath, cert: certPath }
  : undefined;

// 代理到后端 HTTPS 时跳过自签证书验证
const proxyAgent = new https.Agent({ rejectUnauthorized: false });

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
  routeRules: {
    '/embed/**': { ssr: false },
  },
  devtools: false,
  srcDir: 'src/',
  vite: {
    define: {
      __VUE_OPTIONS_API__: true,
      __VUE_PROD_DEVTOOLS__: false,
      __VUE_PROD_HYDRATION_MISMATCH_DETAILS__: false,
    },
    optimizeDeps: {
      include: [
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
        'artplayer-plugin-danmuku',
        'dayjs',
        'dayjs/plugin/*.js',
        'lodash-unified',
        'vuedraggable',
        'js-base64',
      ]
    },
    build: {
      rollupOptions: {
        output: {
          manualChunks(id: string) {
            if (id.includes('/node_modules/artplayer') ||
                id.includes('/node_modules/hls.js') ||
                id.includes('/node_modules/dashjs') ||
                id.includes('/node_modules/artplayer-plugin-danmuku') ||
                id.includes('/node_modules/wplayer-next')) {
              return 'player-vendor';
            }
            if (id.includes('/node_modules/element-plus')) {
              return 'ui-vendor';
            }
            if (id.includes('/node_modules/@wangeditor')) {
              return 'editor-vendor';
            }
          },
        },
      },
    },
    server: {
      warmup: {
        clientFiles: [
          './src/pages/**/*.vue',
          './src/layouts/**/*.vue',
          './src/components/**/*.vue',
        ],
      },
      headers: {
        // 避免切回标签页触发强刷时，SFC 子模块（如 ?vue&type=style）被错误缓存复用
        'Cache-Control': 'no-store',
      },
      hmr: {
        overlay: false
      },
      proxy: {
        '/api': {
          target: `https://${globalConfig.domain}`,
          changeOrigin: true,
          secure: false,
          agent: proxyAgent,
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
