import { getSubtitleListAPI } from '@/api/subtitle';
import { statusCode } from '@/utils/status-code';

const LANG_CODE_TO_LABEL: Record<string, string> = {
  'zh-Hans': '简体中文',
  'zh-Hant': '繁體中文',
  'en': 'English',
  'ja': '日本語',
  'ko': '한국어',
  'vi': 'Tiếng Việt',
  'th': 'ภาษาไทย',
  'ms': 'Bahasa Melayu',
  'id': 'Bahasa Indonesia',
  'es': 'Español',
  'pt': 'Português',
  'ru': 'Русский',
};

/** 控制台过滤前缀：[Alnitak:subtitle]，失败用 [Alnitak:subtitle:warn] */
function subLog(stage: string, payload?: Record<string, unknown>) {
  if (payload && Object.keys(payload).length > 0) {
    console.info(`[Alnitak:subtitle] ${stage}`, payload);
  } else {
    console.info(`[Alnitak:subtitle] ${stage}`);
  }
}

function subWarn(stage: string, payload?: Record<string, unknown>) {
  if (payload && Object.keys(payload).length > 0) {
    console.warn(`[Alnitak:subtitle:warn] ${stage}`, payload);
  } else {
    console.warn(`[Alnitak:subtitle:warn] ${stage}`);
  }
}

/** 延迟多次查看 video 字幕轨就绪情况（blob <track> 的 load/textTracks cues 异步） */
function scheduleSubtitleDomProbe(seq: number, wplayer: WPlayerSubtitleHost | null | undefined, delaysMs: number[]) {
  if (!wplayer?.video || typeof window === 'undefined') return;
  const v = wplayer.video;

  const runOnce = (ms: number) => {
    window.setTimeout(() => {
      if (seq !== subtitleApplySeq) return;
      const elTracks = v.querySelectorAll('track[data-wplayer-subtitle]');

      const fromApi: Record<string, unknown>[] = [];
      const { textTracks } = v;
      if (textTracks && textTracks.length) {
        for (let i = 0; i < textTracks.length; i++) {
          const tt = textTracks[i];
          if (tt.kind === 'subtitles' || tt.kind === 'captions') {
            fromApi.push({
              i,
              kind: tt.kind,
              language: tt.language,
              label: tt.label,
              mode: tt.mode,
              cueCount: tt.cues?.length ?? 0,
            });
          }
        }
      }
      subLog('probe:textTracks', { msAfterApply: ms, seq, htmlTrackTags: elTracks.length, cues: fromApi });

      const tagRows: Record<string, unknown>[] = [];
      elTracks.forEach((el, idx) => {
        const te = el as HTMLTrackElement;
        const t = te.track;
        tagRows.push({
          idx,
          srclang: te.srclang,
          labelAttr: te.label,
          htmlReadyState: te.readyState,
          srcPrefix: (te.src || '').slice(0, 52),
          linkMode: t?.mode ?? '(no TextTrack)',
          linkCues: t?.cues?.length ?? '(n/a)',
        });
      });
      if (tagRows.length) subLog('probe:htmlTrackEl', { msAfterApply: ms, seq, rows: tagRows });

      try {
        const sub = (wplayer as { subtitle?: { available?: boolean; loadSuccess?: boolean; selectedTrackIndex?: number; defaultIndex?: number } }).subtitle;
        if (sub) {
          subLog('probe:wplayer.subtitle', {
            msAfterApply: ms,
            available: sub.available,
            loadSuccess: sub.loadSuccess,
            selectedTrackIndex: sub.selectedTrackIndex,
            defaultIndex: sub.defaultIndex,
          });
        }
      } catch {
        /* noop */
      }
    }, ms);
  };

  for (const ms of delaysMs) runOnce(ms);
}
const pendingSubtitleObjectUrls: string[] = [];

/** 递增后：仍在进行的上一次 fetch+hydrate 在收尾时发现序号不一致则作废并 revoke 本批 blob，避免与其它调用互相 revoke（多轨字幕尤其明显） */
let subtitleApplySeq = 0;

function revokeSubtitleObjectUrls(): void {
  for (const u of pendingSubtitleObjectUrls) {
    try {
      URL.revokeObjectURL(u);
    } catch {
      /* noop */
    }
  }
  pendingSubtitleObjectUrls.length = 0;
}

function revokeUrls(urls: string[]): void {
  for (const u of urls) {
    try {
      URL.revokeObjectURL(u);
    } catch {
      /* noop */
    }
  }
}

/** 相对 /api/subtitle/... 拼同源；http(s) 保持原 URL */
export function resolveSubtitleSrc(url: string): string {
  if (!url) return url;
  if (typeof window === 'undefined') return url;
  const origin = window.location.origin;
  if (/^https?:\/\//i.test(url)) return url;
  if (url.startsWith('/')) return `${origin}${url}`;
  return `${origin}/${url}`;
}

function isCrossOriginResolved(resolved: string): boolean {
  if (typeof window === 'undefined') return false;
  try {
    const u = new URL(resolved, window.location.href);
    return u.origin !== window.location.origin;
  } catch {
    return true;
  }
}

/**
 * Chrome 禁止跨源直接把 https://…oss…/subtitle.x?vtt 塞进 <track src>（報 Unsafe attempt），
 * /api/subtitle 再 302 到 OSS 时跟随后仍会跨源。此处统一 fetch 正文后 blob: ，与頁面同源。
 */
async function hydrateSubtitleSrcForTrack(url: string, blobCollector: string[]): Promise<string> {
  const resolved = resolveSubtitleSrc(url);
  if (typeof window === 'undefined') {
    return resolved;
  }
  const res = await fetch(resolved, {
    mode: 'cors',
    credentials: isCrossOriginResolved(resolved) ? 'omit' : 'include',
  });
  if (!res.ok) {
    throw new Error(`subtitle fetch ${res.status}`);
  }
  const text = await res.text();
  const positioned = text.replace(
    /^(\d+(?::\d+)*[.,]\d+\s*-->\s*\d+(?::\d+)*[.,]\d+).*$/gm,
    '$1 line:98%',
  );
  const blob = new Blob([positioned], { type: 'text/vtt;charset=utf-8' });
  const objectUrl = URL.createObjectURL(blob);
  blobCollector.push(objectUrl);
  return objectUrl;
}

/** 后端 tracks → wplayer 配置（blob 先记入 blobCollector，由 fetchAndApply 提交到 pending） */
export async function hydrateTracksToWPlayerConfig(
  tracks: SubtitleTrackItemType[],
  blobCollector: string[],
): Promise<WPlayerSubtitleConfigItem[]> {
  const out: WPlayerSubtitleConfigItem[] = [];
  for (const t of tracks) {
    try {
      const src = await hydrateSubtitleSrcForTrack(t.url, blobCollector);
      out.push({
        src,
        label: t.label || LANG_CODE_TO_LABEL[t.lang] || t.lang,
        srclang: t.lang,
        default: !!t.isDefault,
        kind: 'subtitles',
      });
      subLog('hydrate:ok', {
        lang: t.lang,
        label: t.label,
        blobPrefix: src.slice(0, 48),
      });
    } catch (e) {
      subWarn('hydrate:fail', {
        lang: t.lang,
        url: resolveSubtitleSrc(t.url).slice(0, 160),
        err: e instanceof Error ? e.message : String(e),
      });
    }
  }
  return out;
}

/** 同步形态（无跨源 blob）；仅在同源场景或测试需要时使用 */
export function tracksToWPlayerConfig(tracks: SubtitleTrackItemType[]): WPlayerSubtitleConfigItem[] {
  return tracks.map((t) => ({
    src: resolveSubtitleSrc(t.url),
    label: t.label || LANG_CODE_TO_LABEL[t.lang] || t.lang,
    srclang: t.lang,
    default: !!t.isDefault,
    kind: 'subtitles',
  }));
}

/** 外链 Artplayer：与同站 WPlayer 一致拉列表 → fetch→blob→同源 URL；页面/组件卸载时调用返回的 revoke */
export async function fetchSubtitleTracksForArtplayer(resourceShortId: string): Promise<{
  tracks: Array<{ url: string; label: string; srclang: string; isDefault: boolean }>;
  revoke: () => void;
}> {
  const blobCollector: string[] = [];
  const revoke = () => revokeUrls(blobCollector);
  try {
    const res = await getSubtitleListAPI(resourceShortId);
    if (res.data.code !== statusCode.OK) {
      return { tracks: [], revoke };
    }
    const rawTracks = (res.data.data?.tracks as SubtitleTrackItemType[]) ?? [];
    if (!rawTracks.length) {
      return { tracks: [], revoke };
    }
    const config = await hydrateTracksToWPlayerConfig(rawTracks, blobCollector);
    if (!config.length) {
      revoke();
      return { tracks: [], revoke: () => {} };
    }
    const tracks = config.map((c) => ({
      url: c.src,
      label: (c.label || c.srclang || 'subtitle').trim() || 'subtitle',
      srclang: (c.srclang ?? '').trim(),
      isDefault: !!c.default,
    }));
    subLog('embed:subtitle:hydrate', {
      resourceShortId,
      trackCount: tracks.length,
      labels: tracks.map((x) => x.label),
    });
    return { tracks, revoke };
  } catch (e) {
    subWarn('embed:subtitle:hydrateFail', {
      resourceShortId,
      err: e instanceof Error ? e.message : String(e),
    });
    revoke();
    return { tracks: [], revoke: () => {} };
  }
}

/** 调用 wplayer-next 内部挂载（与 vendor/wplayer-next/src/js/player.js updateSubtitles 一致） */
export function applyWPlayerSubtitleConfig(
  wplayer: { updateSubtitles?: (c: WPlayerSubtitleConfigItem[] | WPlayerSubtitleConfigItem) => void; subtitle?: { updateSubtitleConfig?: (c: WPlayerSubtitleConfigItem[] | WPlayerSubtitleConfigItem) => void }; video?: HTMLVideoElement } | null | undefined,
  config: WPlayerSubtitleConfigItem[],
) {
  if (!wplayer) return subWarn('apply:skipped', { reason: 'no wplayer instance' });
  if (typeof wplayer.updateSubtitles === 'function') {
    subLog('apply:path', { method: 'updateSubtitles', count: config.length });
    wplayer.updateSubtitles(config);
    return;
  }
  if (wplayer.subtitle && typeof wplayer.subtitle.updateSubtitleConfig === 'function') {
    subLog('apply:path', { method: 'updateSubtitleConfig', count: config.length });
    wplayer.subtitle.updateSubtitleConfig(config);
    return;
  }
  subWarn('apply:fallback:no-api', { count: config.length });
}

/** 无 WPlayer API 时的回退：<track data-alnitak-sub> */
export function setVideoSubtitleTracksFromConfig(
  video: HTMLVideoElement | null | undefined,
  config: WPlayerSubtitleConfigItem[],
) {
  if (!video) return;
  video.querySelectorAll('track[data-alnitak-sub="1"]').forEach((el) => el.remove());
  for (const item of config) {
    const el = document.createElement('track');
    el.kind = item.kind || 'subtitles';
    el.label = item.label || item.srclang || '';
    if (item.srclang) el.srclang = item.srclang;
    el.src = item.src;
    if (item.default) {
      el.setAttribute('default', '');
    }
    el.setAttribute('data-alnitak-sub', '1');
    video.appendChild(el);
  }
}

/** @deprecated 优先用 `fetchAndApplySubtitles`；本函数为异步以支持 OSS 字幕 fetch→blob */
export async function setVideoSubtitleTracks(
  video: HTMLVideoElement | null | undefined,
  tracks: SubtitleTrackItemType[],
): Promise<void> {
  revokeSubtitleObjectUrls();
  const blobs: string[] = [];
  const config = await hydrateTracksToWPlayerConfig(tracks, blobs);
  pendingSubtitleObjectUrls.push(...blobs);
  setVideoSubtitleTracksFromConfig(video, config);
}

export type WPlayerSubtitleHost = {
  updateSubtitles?: (config: WPlayerSubtitleConfigItem[] | WPlayerSubtitleConfigItem) => void;
  subtitle?: { updateSubtitleConfig?: (c: WPlayerSubtitleConfigItem[] | WPlayerSubtitleConfigItem) => void };
  video?: HTMLVideoElement;
};

/**
 * 拉取字幕并交给 WPlayer（启用控制栏 CC 与设置-字幕菜单；依赖 data-wplayer-subtitle 与 load 成功）。
 */
export async function fetchAndApplySubtitles(
  resourceShortId: string | undefined,
  wplayer?: WPlayerSubtitleHost | null,
) {
  const seq = ++subtitleApplySeq;
  subLog('fetch:start', {
    seq,
    resourceShortId: resourceShortId != null && resourceShortId !== '' ? String(resourceShortId).slice(0, 48) : '(none)',
  });

  const apply = (config: WPlayerSubtitleConfigItem[]) => {
    if (!wplayer) return;
    applyWPlayerSubtitleConfig(wplayer, config);
    if (config.length === 0 && wplayer.video) {
      setVideoSubtitleTracksFromConfig(wplayer.video, []);
    }
  };

  if (!resourceShortId) {
    revokeSubtitleObjectUrls();
    subWarn('fetch:abort', { reason: 'missing resourceShortId', seq });
    apply([]);
    return;
  }

  try {
    const res = await getSubtitleListAPI(resourceShortId);
    if (seq !== subtitleApplySeq) {
      subLog('fetch:discard', { reason: 'seq stale after API', seq, latestSeq: subtitleApplySeq });
      return;
    }

    if (res.data.code !== statusCode.OK) {
      revokeSubtitleObjectUrls();
      subWarn('fetch:list:fail', {
        seq,
        apiCode: res.data.code,
        msg: typeof res.data.msg === 'string' ? res.data.msg : '',
      });
      apply([]);
      return;
    }

    const tracks = (res.data.data?.tracks as SubtitleTrackItemType[]) ?? [];
    subLog('fetch:list:ok', { seq, trackCount: tracks.length, langs: tracks.map((t) => t.lang) });

    const blobBatch: string[] = [];
    const config = await hydrateTracksToWPlayerConfig(tracks, blobBatch);

    if (seq !== subtitleApplySeq) {
      revokeUrls(blobBatch);
      subLog('fetch:discard', { reason: 'seq stale after hydrate', seq, latestSeq: subtitleApplySeq, blobs: blobBatch.length });
      return;
    }

    revokeSubtitleObjectUrls();
    pendingSubtitleObjectUrls.push(...blobBatch);

    subLog('fetch:apply', {
      seq,
      configLen: config.length,
      blobs: blobBatch.length,
      entries: config.map((c) => ({ lang: c.srclang, label: c.label, default: c.default, srcPrefix: c.src.slice(0, 40) })),
    });

    applyWPlayerSubtitleConfig(wplayer, config);
    scheduleSubtitleDomProbe(seq, wplayer, [80, 500, 2000]);
    if (
      wplayer &&
      typeof wplayer.updateSubtitles !== 'function' &&
      !(wplayer.subtitle && typeof wplayer.subtitle.updateSubtitleConfig === 'function') &&
      wplayer.video
    ) {
      subLog('apply:dom-fallback-tracks');
      setVideoSubtitleTracksFromConfig(wplayer.video, config);
    }
  } catch (e) {
    if (seq === subtitleApplySeq) {
      revokeSubtitleObjectUrls();
      apply([]);
    }
    subWarn('fetch:exception', { seq, latestSeq: subtitleApplySeq, err: e instanceof Error ? e.message : String(e) });
  }
}
