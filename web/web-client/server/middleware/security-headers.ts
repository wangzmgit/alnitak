export default defineEventHandler((event) => {
  const url = getRequestURL(event);
  const path = url.pathname;
  const isEmbed = path.startsWith('/embed/');

  // 通用安全响应头
  setHeader(event, 'X-Content-Type-Options', 'nosniff');
  setHeader(event, 'Referrer-Policy', 'strict-origin-when-cross-origin');
  setHeader(event, 'Permissions-Policy', 'camera=(), microphone=(), geolocation=(), interest-cohort=()');

  // HSTS：仅在 HTTPS 时设置
  if (url.protocol === 'https:') {
    setHeader(event, 'Strict-Transport-Security', 'max-age=31536000; includeSubDomains; preload');
  }

  // 点击劫持防护
  if (isEmbed) {
    // embed 页面允许被外部站点 iframe 嵌套（嵌入播放器）
    setHeader(event, 'Content-Security-Policy', [
      "default-src 'self'",
      // unsafe-eval 已移除：Nuxt 3 预编译模板，运行时无需 eval
      // unsafe-inline on style 保留：Vue scoped style / CSS-in-JS 必需
      "script-src 'self' 'unsafe-inline'",
      "style-src 'self' 'unsafe-inline'",
      "img-src 'self' data: blob: *",
      "media-src 'self' blob: *",
      "font-src 'self' data:",
      "connect-src 'self' *",
      "frame-ancestors *",
    ].join('; '));
  } else {
    // 非 embed 页面禁止 iframe 嵌套
    setHeader(event, 'X-Frame-Options', 'DENY');
    setHeader(event, 'Content-Security-Policy', [
      "default-src 'self'",
      // unsafe-eval 已移除：Nuxt 3 预编译模板，运行时无需 eval
      // unsafe-inline on style 保留：Vue scoped style / CSS-in-JS 必需
      "script-src 'self' 'unsafe-inline'",
      "style-src 'self' 'unsafe-inline'",
      "img-src 'self' data: blob: *",
      "media-src 'self' blob: *",
      "font-src 'self' data:",
      "connect-src 'self' *",
      "frame-ancestors 'self'",
    ].join('; '));
  }
});
