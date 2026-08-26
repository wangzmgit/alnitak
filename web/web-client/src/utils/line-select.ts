/** 单次延迟检测超时（ms） */
const LATENCY_TIMEOUT_MS = 3000;

/**
 * 向目标 URL 发 Range 请求测延迟，失败返回 Infinity。
 * 仅请求第 0 字节，避免传输实际数据。
 *
 * 如果响应中包含 `X-Oss-Latency` 头（后端代为测量的 OSS 延迟），
 * 优先使用该值而非客户端测到的全程耗时（避免重定向引入额外延迟）。
 */
export async function measureLatency(url: string): Promise<number> {
  const controller = new AbortController();
  const timerId = setTimeout(() => controller.abort(), LATENCY_TIMEOUT_MS);
  const start = performance.now();
  try {
    const res = await fetch(url, {
      method: 'GET',
      headers: { Range: 'bytes=0-0' },
      signal: controller.signal,
    });
    const elapsed = performance.now() - start;
    // 后端代为测量的 OSS 延迟（备用线路探针专用）
    const ossLatency = res.headers.get('X-Oss-Latency');
    if (ossLatency !== null) {
      const parsed = parseInt(ossLatency, 10);
      if (!isNaN(parsed)) {
        // -1 是后端标记的"OSS 不可达"（连接失败/超时），视为 Infinity
        if (parsed < 0) return Infinity;
        return parsed;
      }
    }
    return elapsed;
  } catch {
    return Infinity;
  } finally {
    clearTimeout(timerId);
  }
}

/**
 * 并行检测主/备两条线路的延迟，返回更快的那条线路名。
 * @param primaryUrl - 主线路探测 URL
 * @param backupUrl  - 备用线路探测 URL
 */
export async function selectBestLine(
  primaryUrl: string,
  backupUrl: string,
  dev?: boolean,
): Promise<'primary' | 'backup'> {
  const [primaryLatency, backupLatency] = await Promise.all([
    measureLatency(primaryUrl),
    measureLatency(backupUrl),
  ]);

  const line = backupLatency < primaryLatency ? 'backup' : 'primary';

  if (dev) {
    const label = (ms: number) => (ms === Infinity ? '超时/失败' : `${ms.toFixed(0)}ms`);
    console.log(
      `[line-select] 线路选择: 主=${label(primaryLatency)} 备=${label(backupLatency)} → ${line === 'backup' ? '备用' : '主'}线路`,
    );
  }

  return line;
}

/** 给每个清晰度对象的 url 追加查询参数（绕过 TypeScript 宽松检查） */
export function appendParamToQualities(
  qualities: { url: string }[],
  param: string,
): void {
  for (const q of qualities) {
    q.url += (q.url.includes('?') ? '&' : '?') + param;
  }
}
