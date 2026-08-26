/** 与后端 /api/v1/video/subtitle/list 返回的 tracks 项一致 */
interface SubtitleTrackItemType {
  id: number;
  shortId?: string;
  lang: string;
  label: string;
  url: string;
  backupUrl?: string;
  isDefault: boolean;
}

/** wplayer-next：video.subtitles / player.updateSubtitles 配置项（见 vendor/wplayer-next/docs/guide/subtitles.md） */
interface WPlayerSubtitleConfigItem {
  src: string;
  label?: string;
  srclang?: string;
  default?: boolean;
  kind?: 'subtitles' | 'captions';
}
