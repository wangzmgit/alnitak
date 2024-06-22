import path from 'path';
import { defineConfig } from 'vitepress'

export default defineConfig({
  title: "Alnitak弹幕视频网站文档",
  description: "基于GO + Nuxt的弹幕视频网站",
  head: [
    ['meta', { name: 'keywords', content: 'go, vue, nuxt, 弹幕网站' }],
    ['meta', {
      name: 'viewport',
      content: 'width=device-width,initial-scale=1,minimum-scale=1.0,maximum-scale=1.0,user-scalable=no'
    }],
    ['link', { rel: 'icon', href: '/favicon.ico' }]
  ],
  srcDir: `${path.resolve(process.cwd())}/src`,
  themeConfig: {
    editLink: {
      text: '为此页提供修改建议',
      pattern: 'https://github.com/wangzmgit/alnitak/tree/doc/docs/src/:path'
    },
    socialLinks: [
      { icon: 'github', link: 'https://github.com/wangzmgit/alnitak' },
    ],
    footer: {
      message: '根据 MIT 许可证发布',
      copyright: 'Copyright © 2020-2024'
    },
    nav: [
      { text: '项目指南', link: '/guide/introduce', activeMatch: '/guide/' },
      { text: '接口文档', link: '/api/auth', activeMatch: '/api' },
      { text: '赞助', link: '/other/donate' }
    ],
    sidebar: {
      '/guide/': [
        {
          text: '项目介绍',
          collapsed: false,
          items: [
            { text: '简介', link: '/guide/introduce' },
          ]
        },
        {
          text: '快速开始',
          collapsed: true,
          items: [
            { text: '环境配置', link: '/guide/start/env' },
            { text: '项目配置', link: '/guide/start/config' },
            { text: '项目启动', link: '/guide/start/start' },
          ]
        },
        {
          text: '部署指南',
          collapsed: true,
          items: [
            { text: 'Docker部署', link: '/guide/deploy/docker' },
            { text: '手动部署', link: '/guide/deploy/manual' },
            { text: '域名配置', link: '/guide/deploy/domain' },
          ]
        },
        {
          text: '项目指南',
          collapsed: true,
          items: [
            { text: '状态码', link: '/guide/guide/code' },
            { text: '登录流程', link: '/guide/guide/login' },
            { text: '视频处理流程', link: '/guide/guide/transcoding' },
            { text: '人机验证说明', link: '/guide/guide/captcha' },
            { text: '弹幕播放器说明', link: '/guide/guide/player' },
            // { text: '对象存储和CDN说明', link: '/guide/' },
          ]
        },
        {
          text: '其他',
          collapsed: true,
          items: [
            { text: '贡献指南', link: '/guide/other/contribution' },
            { text: '更新说明', link: '/guide/other/update' },
            { text: '常见问题解答', link: '/guide/other/qa' },
            { text: '相关截图', link: '/guide/screenshot' }
          ]
        }
      ],
      '/api/': [
        {
          text: '接口文档',
          items: [
            { text: '用户认证相关接口', link: '/api/auth' },
            { text: '用户相关接口', link: '/api/user' },
            { text: '视频相关接口', link: '/api/video' },
            { text: '视频资源接口', link: '/api/resource' },
            { text: '文章相关接口', link: '/api/article' },
            { text: '验证相关接口', link: '/api/verify' },
            { text: 'API管理相关接口', link: '/api/api' },
            { text: '视频点赞收藏接口', link: '/api/archive_video' },
            { text: '文章点赞收藏接口', link: '/api/archive_article' },
            { text: '轮播图相关接口', link: '/api/carousel' },
            { text: '收藏夹相关接口', link: '/api/collection' },
            { text: '视频评论回复接口', link: '/api/comment_video' },
            { text: '文章评论回复接口', link: '/api/comment_article' },
            { text: '配置相关接口', link: '/api/config' },
            { text: '弹幕相关接口', link: '/api/danmaku' },
            { text: '历史记录相关接口', link: '/api/history' },
            { text: '菜单管理相关接口', link: '/api/menu' },
            { text: '消息相关接口', link: '/api/message' },
            { text: '在线人数相关接口', link: '/api/online' },
            { text: '分区相关接口', link: '/api/partition' },
            { text: '粉丝关注相关接口', link: '/api/relation' },
            { text: '视频审核相关接口', link: '/api/review_video' },
            { text: '文章审核相关接口', link: '/api/review_article' },
            { text: '角色管理相关接口', link: '/api/role' },
            { text: '文件上传相关接口', link: '/api/upload' },
          ]
        }
      ]
    }
  },
})
