/**
 * 字幕轨记忆：与 vendor wplayer-next `subtitle.js` 使用同一 localStorage key。
 * 「显示字幕」时按上次选择的 label/lang 匹配当前列表，不匹配则退回第一条。
 */

export const ALNITAK_SUBTITLE_PREF_LS_KEY = 'alnitak-pref-subtitle-track';

export type StoredSubtitlePreference = {
  label?: string;
  lang?: string;
};

function normLabel(s: string | undefined): string {
  return (s ?? '').trim().toLowerCase();
}

/** 写入当前选中的轨（外链 / TS 调用） */
export function writeStoredSubtitlePreference(pref: StoredSubtitlePreference): void {
  if (typeof localStorage === 'undefined') return;
  try {
    const label = pref.label != null ? String(pref.label).trim() : '';
    const lang = pref.lang != null ? String(pref.lang).trim() : '';
    localStorage.setItem(ALNITAK_SUBTITLE_PREF_LS_KEY, JSON.stringify({ label, lang }));
  } catch {
    /* quota / private mode */
  }
}

export function readStoredSubtitlePreference(): StoredSubtitlePreference | null {
  if (typeof localStorage === 'undefined') return null;
  try {
    const raw = localStorage.getItem(ALNITAK_SUBTITLE_PREF_LS_KEY);
    if (!raw) return null;
    const o = JSON.parse(raw) as unknown;
    if (!o || typeof o !== 'object') return null;
    const po = o as Record<string, unknown>;
    const label = typeof po.label === 'string' ? po.label : '';
    const lang = typeof po.lang === 'string' ? po.lang : '';
    return { label, lang };
  } catch {
    return null;
  }
}

/**
 * @param tracks 当前可用轨（顺序与播放器列表一致）
 * @returns 下标：先按 label（不区分大小写）再按 srclang/lang，均无匹配则 `0`
 */
export function pickSubtitleTrackIndexByPreference<
  T extends { label?: string; srclang?: string; lang?: string },
>(tracks: T[]): number {
  if (!tracks.length) return 0;
  const pref = readStoredSubtitlePreference();
  if (!pref || (!pref.label && !pref.lang)) return 0;
  const nl = normLabel(pref.label);
  if (nl) {
    const i = tracks.findIndex((t) => normLabel(t.label) === nl);
    if (i >= 0) return i;
  }
  const lc = normLabel(pref.lang);
  if (lc) {
    const j = tracks.findIndex(
      (t) =>
        normLabel(t.srclang) === lc ||
        normLabel(t.lang) === lc,
    );
    if (j >= 0) return j;
  }
  return 0;
}
