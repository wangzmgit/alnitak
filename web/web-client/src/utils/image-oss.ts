/**
 * 从已拼接域名的完整 URL 推导备用 OSS URL。
 * 仅对经过后端代理的路径有效（/api/image/、/api/subtitle/）。
 * 直接 OSS 直链无法从客户端推导备用 URL。
 */
export function getBackupUrl(fullUrl: string): string | null {
  if (!fullUrl) return null;
  // 已经在用备用 OSS，不再重试
  if (fullUrl.includes('backup=true')) return null;

  // 只处理经过后端代理的路径
  if (/\/api\/(image|subtitle)\//.test(fullUrl)) {
    const sep = fullUrl.includes('?') ? '&' : '?';
    return fullUrl + sep + 'backup=true';
  }

  return null;
}
