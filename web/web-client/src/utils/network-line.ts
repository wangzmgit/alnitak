/**
 * 全局网络线路探测模块。
 *
 * 策略：
 * 1. 启动时并行测主备延迟（Range: bytes=0-0），选快的
 * 2. 每 30 秒周期性重探，适应网络变化
 * 3. 各组件（oss-image, video-player）独立做本地容灾，不依赖全局翻线
 */
import { ref } from 'vue';
import { measureLatency } from './line-select';

export type Line = 'primary' | 'backup';

// ===== 常量 =====
const PROBE_URL = '/api/image/probe';
const RECHECK_INTERVAL_MS = 30 * 1000; // 30 秒重探一次，快速感知线路故障

// ===== 全局响应式状态 =====
const selectedLine = ref<Line | null>(null);
const isChecking = ref(false);
const lastCheckedAt = ref(0);
const latencyResult = ref<{ primary: number; backup: number } | null>(null);

// ===== 周期性重探 =====
let recheckTimer: ReturnType<typeof setTimeout> | null = null;

function clearRecheck(): void {
  if (recheckTimer !== null) {
    clearTimeout(recheckTimer);
    recheckTimer = null;
  }
}

function scheduleRecheck(): void {
  clearRecheck();
  recheckTimer = setTimeout(() => {
    checkNetwork();
  }, RECHECK_INTERVAL_MS);
}

// ===== 核心检测 =====

/**
 * 执行一次线路检测（只做 Range 延迟测量），更新全局状态。
 */
export async function checkNetwork(): Promise<Line> {
  isChecking.value = true;
  try {
    const backupUrl = PROBE_URL + '?backup=true';

    const [primaryLatency, backupLatency] = await Promise.all([
      measureLatency(PROBE_URL),
      measureLatency(backupUrl),
    ]);

    latencyResult.value = { primary: primaryLatency, backup: backupLatency };

    // 延迟越低越好
    selectedLine.value = backupLatency < primaryLatency ? 'backup' : 'primary';
    lastCheckedAt.value = Date.now();

    if (import.meta.dev) {
      console.log('[network-line] 线路检测:', {
        主线路_ms: primaryLatency,
        备用线路_ms: backupLatency,
        选择: selectedLine.value === 'backup' ? '备用' : '主',
      });
    }

    return selectedLine.value;
  } finally {
    isChecking.value = false;
    scheduleRecheck();
  }
}

// ===== 惰性初始化（全局只执行一次） =====
let initPromise: Promise<Line> | null = null;

/**
 * 全局只调一次，应用启动时调用。后续取缓存结果。
 */
export function initNetworkLine(): Promise<Line> {
  if (selectedLine.value) return Promise.resolve(selectedLine.value);
  if (initPromise) return initPromise;
  initPromise = checkNetwork();
  return initPromise;
}

// ===== Composable =====
export function useNetworkLine() {
  return {
    selectedLine: selectedLine as Readonly<typeof selectedLine>,
    isChecking: isChecking as Readonly<typeof isChecking>,
    lastCheckedAt,
    latencyResult,
    checkNetwork,
  };
}
