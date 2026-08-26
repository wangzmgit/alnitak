<template>
  <div v-if="seasons.length > 0 || episodes.length > 0" class="collection-list">
    <div class="collection-head">
      <span class="title">正片</span>
      <span class="count">({{ episodes.length }})</span>
      <span class="view-mode-switch" @click="isGrid = !isGrid" title="切换显示模式">
        <svg v-if="!isGrid" class="icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="16" height="16">
          <path d="M3.02134 15.21437C3.02134 13.79417 4.17261 12.64294 5.59277 12.64294L8.67849 12.64294C10.09869 12.64294 11.24991 13.79417 11.24991 15.21437L11.24991 18.30009C11.24991 19.7202 10.09869 20.87151 8.67849 20.87151L5.59277 20.87151C4.17261 20.87151 3.02134 19.7202 3.02134 18.30009L3.02134 15.21437zM5.59277 14.35723C5.11939 14.35723 4.73563 14.74097 4.73563 15.21437L4.73563 18.30009C4.73563 18.7734 5.11939 19.15723 5.59277 19.15723L8.67849 19.15723C9.15189 19.15723 9.53563 18.7734 9.53563 18.30009L9.53563 15.21437C9.53563 14.74097 9.15189 14.35723 8.67849 14.35723L5.59277 14.35723z" fill="currentColor"></path>
          <path d="M3.02134 5.70002C3.02134 4.27986 4.17261 3.12859 5.59277 3.12859L8.67849 3.12859C10.09869 3.12859 11.24991 4.27986 11.24991 5.70002L11.24991 8.78571C11.24991 10.20591 10.09869 11.35714 8.67849 11.35714L5.59277 11.35714C4.17261 11.35714 3.02134 10.20591 3.02134 8.78571L3.02134 5.70002zM5.59277 4.84287C5.11939 4.84287 4.73563 5.22663 4.73563 5.70002L4.73563 8.78571C4.73563 9.25911 5.11939 9.64286 5.59277 9.64286L8.67849 9.64286C9.15189 9.64286 9.53563 9.25911 9.53563 8.78571L9.53563 5.70002C9.53563 5.22663 9.15189 4.84287 8.67849 4.84287L5.59277 4.84287z" fill="currentColor"></path>
          <path d="M12.75 5.10859C12.75 4.63521 13.15937 4.25145 13.66431 4.25145L20.06426 4.25145C20.5692 4.25145 20.97857 4.63521 20.97857 5.10859C20.97857 5.58198 20.5692 5.96573 20.06426 5.96573L13.66431 5.96573C13.15937 5.96573 12.75 5.58198 12.75 5.10859z" fill="currentColor"></path>
          <path d="M12.75 9.39429C12.75 8.92089 13.15937 8.53716 13.66431 8.53716L20.06426 8.53716C20.5692 8.53716 20.97857 8.92089 20.97857 9.39429C20.97857 9.86769 20.5692 10.25143 20.06426 10.25143L13.66431 10.25143C13.15937 10.25143 12.75 9.86769 12.75 9.39429z" fill="currentColor"></path>
          <path d="M12.75 14.57143C12.75 14.09803 13.15937 13.71429 13.66431 13.71429L20.06426 13.71429C20.5692 13.71429 20.97857 14.09803 20.97857 14.57143C20.97857 15.04483 20.5692 15.42857 20.06426 15.42857L13.66431 15.42857C13.15937 15.42857 12.75 15.04483 12.75 14.57143z" fill="currentColor"></path>
          <path d="M12.75 18.85714C12.75 18.38374 13.15937 18 13.66431 18L20.06426 18C20.5692 18 20.97857 18.38374 20.97857 18.85714C20.97857 19.33054 20.5692 19.71429 20.06426 19.71429L13.66431 19.71429C13.15937 19.71429 12.75 19.33054 12.75 18.85714z" fill="currentColor"></path>
        </svg>
        <svg v-else class="icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="16" height="16">
          <path d="M12.64286 15.21437C12.64286 13.79417 13.79409 12.64294 15.21429 12.64294L18.3 12.64294C19.7202 12.64294 20.87143 13.79417 20.87143 15.21437L20.87143 18.30009C20.87143 19.7202 19.7202 20.87151 18.3 20.87151L15.21429 20.87151C13.79409 20.87151 12.64286 19.7202 12.64286 18.30009L12.64286 15.21437zM15.21429 14.35723C14.74089 14.35723 14.35714 14.74097 14.35714 15.21437L14.35714 18.30009C14.35714 18.7734 14.74089 19.15723 15.21429 19.15723L18.3 19.15723C18.7734 19.15723 19.15714 18.7734 19.15714 18.30009L19.15714 15.21437C19.15714 14.74097 18.7734 14.35723 18.3 14.35723L15.21429 14.35723z" fill="currentColor"></path>
          <path d="M3.12849 15.21437C3.12849 13.79417 4.27976 12.64294 5.69991 12.64294L8.78563 12.64294C10.20583 12.64294 11.35706 13.79417 11.35706 15.21437L11.35706 18.30009C11.35706 19.7202 10.20583 20.87151 8.78563 20.87151L5.69991 20.87151C4.27976 20.87151 3.12849 19.7202 3.12849 18.30009L3.12849 15.21437zM5.69991 14.35723C5.22653 14.35723 4.84277 14.74097 4.84277 15.21437L4.84277 18.30009C4.84277 18.7734 5.22653 19.15723 5.69991 19.15723L8.78563 19.15723C9.25903 19.15723 9.64277 18.7734 9.64277 18.30009L9.64277 15.21437C9.64277 14.74097 9.25903 14.35723 8.78563 14.35723L5.69991 14.35723z" fill="currentColor"></path>
          <path d="M12.64286 5.70002C12.64286 4.27986 13.79409 3.12859 15.21429 3.12859L18.3 3.12859C19.7202 3.12859 20.87143 4.27986 20.87143 5.70002L20.87143 8.78571C20.87143 10.20591 19.7202 11.35714 18.3 11.35714L15.21429 11.35714C13.79409 11.35714 12.64286 10.20591 12.64286 8.78571L12.64286 5.70002zM15.21429 4.84287C14.74089 4.84287 14.35714 5.22663 14.35714 5.70002L14.35714 8.78571C14.35714 9.25911 14.74089 9.64286 15.21429 9.64286L18.3 9.64286C18.7734 9.64286 19.15714 9.25911 19.15714 8.78571L19.15714 5.70002C19.15714 5.22663 18.7734 4.84287 18.3 4.84287L15.21429 4.84287z" fill="currentColor"></path>
          <path d="M3.12849 5.70002C3.12849 4.27986 4.27976 3.12859 5.69991 3.12859L8.78563 3.12859C10.20583 3.12859 11.35706 4.27986 11.35706 5.70002L11.35706 8.78571C11.35706 10.20591 10.20583 11.35714 8.78563 11.35714L5.69991 11.35714C4.27976 11.35714 3.12849 10.20591 3.12849 8.78571L3.12849 5.70002zM5.69991 4.84287C5.22653 4.84287 4.84277 5.22663 4.84277 5.70002L4.84277 8.78571C4.84277 9.25911 5.22653 9.64286 5.69991 9.64286L8.78563 9.64286C9.25903 9.64286 9.64277 9.25911 9.64277 8.78571L9.64277 5.70002C9.64277 5.22663 9.25903 4.84287 8.78563 4.84287L5.69991 4.84287z" fill="currentColor"></path>
        </svg>
      </span>
      <span class="autoplay-switch" @click="autonext = !autonext">
        <span class="autoplay-text">自动连播</span>
        <span class="switch-button" :class="{ 'on': autonext }"></span>
      </span>
    </div>

    <div class="season-tabs" v-if="seasons.length > 0">
      <span
        v-for="(s, idx) in seasons"
        :key="s.pgc_id"
        :class="['tab', activeSeasonId === String(s.pgc_id) ? 'active' : '']"
        @click="changeSeason(String(s.pgc_id))"
      >
        第{{ idx + 1 }}季
      </span>
    </div>

    <el-scrollbar max-height="340px">
      <ul v-if="!isGrid" class="ep-list">
        <li
          v-for="ep in orderedEpisodes"
          :key="ep.id"
          :class="['ep-item', isActiveEpisode(ep) ? 'active' : '']"
          @click="goEpisode(ep)"
        >
          <span class="idx">第{{ ep.episode_number }}话</span>
          <span class="name">{{ ep.title || `第${ep.episode_number}话` }}</span>
          <span v-if="isActiveEpisode(ep)" class="playing-tag">正在播放</span>
        </li>
      </ul>
      <div v-else class="ep-grid">
        <div
          v-for="ep in orderedEpisodes"
          :key="ep.id"
          :class="['ep-no', isActiveEpisode(ep) ? 'active-number' : '']"
          @click="goEpisode(ep)"
        >
          {{ ep.episode_number }}
        </div>
      </div>
    </el-scrollbar>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue';
import { getPGCPlayPanelByVideoAPI } from '@/api/pgc';
import { statusCode } from '@/utils/status-code';

type PGCSeason = { pgc_id: string; title?: string };
type PGCEpisode = { id: number; ep_id?: number; episode_number: number; title: string; vid: number };

const props = defineProps<{
  vid: number | string;
  initialSeasons?: any[];
  initialEpisodes?: any[];
  initialActiveSeasonId?: string;
}>();

const seasons = ref<PGCSeason[]>((props.initialSeasons || []) as PGCSeason[]);
const episodes = ref<PGCEpisode[]>((props.initialEpisodes || []) as PGCEpisode[]);
const activeSeasonId = ref(props.initialActiveSeasonId || '');
const isGrid = ref(false);
const autonext = ref(false);

const orderedEpisodes = computed(() => {
  const list = [...episodes.value];
  list.sort((a, b) => a.episode_number - b.episode_number);
  return list;
});

const toCleanVid = (raw: unknown): number => {
  const n = Number(String(raw ?? '').trim());
  return Number.isFinite(n) && n > 0 ? Math.floor(n) : 0;
};

const isActiveEpisode = (ep: PGCEpisode) => toCleanVid(ep.vid) === props.vid;

const goEpisode = (ep: PGCEpisode) => {
  const epId = Number(ep.ep_id || ep.id || 0);
  if (Number.isFinite(epId) && epId > 0) {
    navigateTo(`/watch?ep=${Math.floor(epId)}&mode=pgc`);
    return;
  }
  const vid = toCleanVid(ep.vid);
  if (vid > 0) {
    navigateTo(`/watch?v=${vid}&mode=pgc`);
  }
};

const loadPanel = async (seasonId?: string) => {
  const res = await getPGCPlayPanelByVideoAPI(props.vid, seasonId);
  if (res.data.code !== statusCode.OK) return;
  const data = res.data.data || {};
  seasons.value = data.seasons || [];
  episodes.value = data.episodes || [];
  activeSeasonId.value = String(data.active_season_id || '');
};

const changeSeason = async (seasonId: string) => {
  if (!seasonId || seasonId === activeSeasonId.value) return;
  await loadPanel(seasonId);
};

const getNextVideo = () => {
  const idx = orderedEpisodes.value.findIndex((ep) => toCleanVid(ep.vid) === props.vid);
  if (idx >= 0 && idx < orderedEpisodes.value.length - 1) {
    const next = orderedEpisodes.value[idx + 1];
    return { vid: toCleanVid(next.vid), epId: Number(next.ep_id || next.id || 0) || undefined };
  }
  return null;
};

onMounted(() => {
  const savedAuto = localStorage.getItem('video-autonext-pgc-season');
  autonext.value = savedAuto === 'true';
  const savedView = localStorage.getItem('video-pgc-ep-view-mode');
  isGrid.value = savedView === 'grid';
});

watch(autonext, (val) => localStorage.setItem('video-autonext-pgc-season', String(val)));
watch(isGrid, (val) => localStorage.setItem('video-pgc-ep-view-mode', val ? 'grid' : 'list'));
watch(() => props.vid, () => loadPanel());
watch(() => props.initialSeasons, (val) => {
  if (val && val.length > 0) seasons.value = val as PGCSeason[];
});
watch(() => props.initialEpisodes, (val) => {
  if (val && val.length > 0) episodes.value = val as PGCEpisode[];
});
watch(() => props.initialActiveSeasonId, (val) => {
  if (val) activeSeasonId.value = val;
});

defineExpose({
  autonext: computed(() => autonext.value),
  getNextVideo,
  orderedEpisodes,
  isActiveEpisode,
});
</script>

<style lang="scss" scoped>
.collection-list {
  border-radius: 6px;
  background-color: var(--bg-elev-1);
  border: 1px solid var(--border-color);
  box-shadow: 0 2px 8px var(--shadow-weak);
  margin-bottom: 18px;
}

.collection-head {
  display: flex;
  align-items: center;
  padding: 12px 16px 8px;
  border-bottom: 1px solid var(--border-color);
  background-color: var(--panel-bg);

  .title {
    font-size: 16px;
    color: var(--font-primary-1);
    font-weight: 500;
  }

  .count {
    color: var(--font-primary-3);
    font-size: 14px;
    margin-left: 4px;
  }

  .view-mode-switch {
    display: flex;
    align-items: center;
    cursor: pointer;
    user-select: none;
    margin-left: 12px;
    padding: 4px;
    border-radius: 4px;
    transition: background-color 0.2s;

    &:hover {
      background-color: var(--hover-bg);
    }

    .icon {
      color: var(--font-primary-2);
      transition: color 0.2s;
    }

    &:hover .icon {
      color: var(--font-primary-1);
    }
  }

  .autoplay-switch {
    display: flex;
    align-items: center;
    cursor: pointer;
    user-select: none;
    margin-left: auto;

    .autoplay-text {
      font-size: 13px;
      color: var(--font-primary-2);
      margin-right: 8px;
    }

    .switch-button {
      width: 32px;
      height: 18px;
      background: var(--border-color);
      border-radius: 9px;
      position: relative;
      transition: background-color 0.3s;

      &::after {
        content: '';
        position: absolute;
        top: 2px;
        left: 2px;
        width: 14px;
        height: 14px;
        background: var(--bg-elev-1);
        border-radius: 50%;
        transition: transform 0.3s;
      }

      &.on {
        background: var(--primary-color);

        &::after {
          transform: translateX(14px);
        }
      }
    }
  }
}

.season-tabs { padding: 10px 12px 2px; display: flex; gap: 10px; overflow-x: auto; }
.tab { font-size: 13px; color: var(--font-primary-2); cursor: pointer; white-space: nowrap; }
.tab.active { color: var(--primary-hover-color); }
.ep-list { margin: 0; padding: 8px 10px 10px; list-style: none; }
.ep-item { height: 32px; line-height: 32px; border-radius: 4px; padding: 0 10px; cursor: pointer; display: flex; gap: 8px; }
.ep-item:hover { background: var(--hover-bg); }
.ep-item.active { color: var(--primary-hover-color); background: var(--hover-bg); }
.idx { width: 62px; flex-shrink: 0; font-size: 12px; color: var(--font-primary-3); }
.name { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; font-size: 13px; color: var(--font-primary-1); }
.playing-tag {
  margin-left: auto;
  font-size: 12px;
  color: #fff;
  background: var(--primary-color);
  border-radius: 10px;
  padding: 0 8px;
  line-height: 20px;
  height: 20px;
  align-self: center;
}
.ep-grid { padding: 12px; display: grid; grid-template-columns: repeat(auto-fill, minmax(44px, 1fr)); gap: 8px; }
.ep-no { height: 36px; border: 1px solid var(--border-color); border-radius: 4px; display: flex; align-items: center; justify-content: center; cursor: pointer; color: var(--font-primary-2); }
.active-number { background-color: var(--primary-color); border-color: var(--primary-color); color: var(--primary-text-color); }
</style>

