<template>
  <!-- 有分P列表或合集列表时才显示 -->
  <div v-if="mergedList.length > 1" class="collection-list">
    <div class="collection-head">
      <span class="title">{{ listType === 'parts' ? '视频分集' : displayTitle }}</span>
      <span class="count">({{ currentIndex }}/{{ displayCount }})</span>
      <span class="view-mode-switch" @click="toggleViewMode" title="切换显示模式">
        <svg v-if="isNumberMode" class="icon" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 24 24" width="16" height="16">
          <path d="M3.02134 15.21437C3.02134 13.79417 4.17261 12.64294 5.59277 12.64294L8.67849 12.64294C10.09869 12.64294 11.24991 13.79417 11.24991 15.21437L11.24991 18.30009C11.24991 19.7202 10.09869 20.87151 8.67849 20.87151L5.59277 20.87151C4.17261 20.87151 3.02134 19.7202 3.02134 18.30009L3.02134 15.21437zM5.59277 14.35723C5.11939 14.35723 4.73563 14.74097 4.73563 15.21437L4.73563 18.30009C4.73563 18.7734 5.11939 19.15723 5.59277 19.15723L8.67849 19.15723C9.15189 19.15723 9.53563 18.7734 9.53563 18.30009L9.53563 15.21437C9.53563 14.74097 9.15189 14.35723 8.67849 14.35723L5.59277 14.35723z" fill="currentColor"></path>
          <path d="M3.02134 5.70002C3.02134 4.27986 4.17261 3.12859 5.59277 3.12859L8.67849 3.12859C10.09869 3.12859 11.24991 4.27986 11.24991 5.70002L11.24991 8.78571C11.24991 10.20591 10.09869 11.35714 8.67849 11.35714L5.59277 11.35714C4.17261 11.35714 3.02134 10.20591 3.02134 8.78571L3.02134 5.70002zM5.59277 4.84287C5.11939 4.84287 4.73563 5.22663 4.73563 5.70002L4.73563 8.78571C4.73563 9.25911 5.11939 9.64286 5.59277 9.64286L8.67849 9.64286C9.15189 9.64286 9.53563 9.25911 9.53563 8.78571L9.53563 5.70002C9.53563 5.22663 9.15189 4.84287 8.67849 4.84287L5.59277 4.84287z" fill="currentColor"></path>
          <path d="M12.75 5.10859C12.75 4.63521 13.15937 4.25145 13.66431 4.25145L20.06426 4.25145C20.5692 4.25145 20.97857 4.63521 20.97857 5.10859C20.97857 5.58198 20.5692 5.96573 20.06426 5.96573L13.66431 5.96573C13.15937 5.96573 12.75 5.58198 12.75 5.10859z" fill="currentColor"></path>
          <path d="M12.75 9.39429C12.75 8.92089 13.15937 8.53716 13.66431 8.53716L20.06426 8.53716C20.5692 8.53716 20.97857 8.92089 20.97857 9.39429C20.97857 9.86769 20.5692 10.25143 20.06426 10.25143L13.66431 10.25143C13.15937 10.25143 12.75 9.86769 12.75 9.39429z" fill="currentColor"></path>
          <path d="M12.75 14.57143C12.75 14.09803 13.15937 13.71429 13.66431 13.71429L20.06426 13.71429C20.5692 13.71429 20.97857 14.09803 20.97857 14.57143C20.97857 15.04483 20.5692 15.42857 20.06426 15.42857L13.66431 15.42857C13.15937 15.42857 12.75 15.04483 12.75 14.57143z" fill="currentColor"></path>
          <path d="M12.75 18.85714C12.75 18.38374 13.15937 18 13.66431 18L20.06426 18C20.5692 18 20.97857 18.38374 20.97857 18.85714C20.97857 19.33054 20.5692 19.71429 20.06426 19.71429L13.66431 19.71429C13.15937 19.71429 12.75 19.33054 12.75 18.85714z" fill="currentColor"></path>
        </svg>
        <svg v-else class="icon" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 24 24" width="16" height="16">
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
    <el-scrollbar max-height="340px">
      <ul v-if="!isNumberMode" class="list-box">
<li
          v-for="(item, index) in mergedList"
          :key="item.vid + '-' + (item.p || index)"
          :class="['list-item', isActiveItem(item, index) ? 'active-item' : '']"
          @click="goVideo(item)"
        >
          <div class="item-content">
            <span class="item-index">{{ index + 1 }}</span>
            <span class="item-title">{{ item.partTitle || item.title }}</span>
          </div>
          <div class="item-duration">{{ formatDuration(item.duration) }}</div>
        </li>
      </ul>
      <div v-else class="number-grid">
<div
          v-for="(item, index) in mergedList"
          :key="item.vid + '-' + (item.p || index)"
          :class="['number-item', isActiveItem(item, index) ? 'active-number' : '']"
          @click="goVideo(item)"
        >
          {{ index + 1 }}
        </div>
      </div>
    </el-scrollbar>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch, readonly } from 'vue'
import { getPlaylistVideoListWithPartsAPI } from '@/api/playlist'
import { statusCode } from '@/utils/status-code'

const props = defineProps<{
  vid: number | string
  resources?: Array<ResourceType>  // 当前视频的分P列表
  currentPart?: number  // 当前播放的分P序号
}>()

interface VideoItem {
  vid: number | string
  shortId?: string
  resourceRid?: string  // 资源shortId，用于精确匹配分P
  title: string
  cover: string
  duration: number
  clicks: number
  desc: string
  p?: number  // 分P序号
  partTitle?: string  // 分P的标题
}

// 清理 vid，兼容数字和 shortId
const cleanVid = (v: number | string): string => {
  if (typeof v === 'number') return String(v)
  return v
}

// 当前视频的 vid 字符串（兼容数字和 shortId）
const currentVidStr = computed(() => cleanVid(props.vid))

interface PlaylistInfo {
  id: number
  title: string
}

const playlist = ref<PlaylistInfo | null>(null)
const videoList = ref<VideoItem[]>([])
const partList = ref<VideoItem[]>([])  // 分P列表
const autonext = ref(false)
const isNumberMode = ref(false)

// 合并后的列表：根据列表类型显示对应的列表
const mergedList = computed(() => {
  // 优先显示合集列表（包含展开的多分P视频）
  if (playlist.value && videoList.value.length > 0) {
    return videoList.value
  }
  // 没有合集但有分P，显示分P列表
  if (partList.value.length > 0) {
    return partList.value
  }
  return []
})

// 当前展示的标题
const displayTitle = computed(() => {
  if (playlist.value) {
    return playlist.value.title
  }
  if (partList.value.length > 0) {
    return '视频分集'
  }
  return ''
})

// 当前列表类型（用于判断是否有合集）
const listType = computed(() => {
  if (playlist.value && videoList.value.length > 0) {
    return 'collection'  // 合集
  }
  if (partList.value.length > 0) {
    return 'parts'  // 仅分P
  }
  return ''
})

// 当前展示的数量
const displayCount = computed(() => mergedList.value.length)

const toggleViewMode = () => {
  isNumberMode.value = !isNumberMode.value
}

const currentIndex = computed(() => {
  const currentPart = props.currentPart || 1
  
  // 分P列表类型：直接用分P序号
  if (listType.value === 'parts') {
    return currentPart
  }
  
  // 合集类型：用p字段匹配（兼容 shortId）
  const currentVideoParts = mergedList.value.filter(v => {
    const itemVidStr = cleanVid(v.vid)
    const itemShortId = (v as any).shortId
    return itemVidStr === currentVidStr.value || itemShortId === currentVidStr.value
  })
  
  if (currentVideoParts.length > 1) {
    const idx = currentVideoParts.findIndex(v => (v.p || 1) === currentPart)
    if (idx >= 0) {
      const firstIndex = mergedList.value.findIndex(v => {
        const itemVidStr = cleanVid(v.vid)
        const itemShortId = (v as any).shortId
        return itemVidStr === currentVidStr.value || itemShortId === currentVidStr.value
      })
      return firstIndex + idx + 1
    }
  }
  // 单分P或找不到，返回第一个
  const firstIdx = mergedList.value.findIndex(v => {
    const itemVidStr = cleanVid(v.vid)
    const itemShortId = (v as any).shortId
    return itemVidStr === currentVidStr.value || itemShortId === currentVidStr.value
  })
  return firstIdx >= 0 ? firstIdx + 1 : 0
})

const formatDuration = (seconds: number) => {
  if (!seconds) return ''
  const m = Math.floor(seconds / 60)
  const s = Math.floor(seconds % 60)
  return `${m}:${s.toString().padStart(2, '0')}`
}

// 判断列表项是否为当前播放项
const isActiveItem = (item: VideoItem, index: number) => {
  const currentPart = props.currentPart || 1
  
  // 分P列表类型：精确匹配分P序号
  if (listType.value === 'parts') {
    return (index + 1) === currentPart
  }
  
  // 合集类型：需要同时匹配主稿件ID + 资源rid
  const itemShortId = (item as any).shortId;
  const itemResourceRid = (item as any).resourceRid;
  const itemP = (item as any).p || 1;
  const itemVidStr = cleanVid(item.vid);
  
  // 获取当前播放的资源 shortId (rid)，从 partList 中获取
  const currentRid = partList.value[currentPart - 1]?.shortId;
  
  // 主稿件ID匹配（兼容 shortId）
  const mainVideoMatch = itemVidStr === currentVidStr.value || itemShortId === currentVidStr.value;
  // 资源rid匹配
  const resourceMatch = itemResourceRid && currentRid && itemResourceRid === currentRid;
  
  // 合集类型：精确匹配主稿件 + 资源rid
  if (mainVideoMatch && resourceMatch) {
    return true;
  }
  
  // 同一主稿件的不同分P，用 p 匹配
  if (mainVideoMatch && itemP) {
    return itemP === currentPart;
  }
  
  return false
}

const emits = defineEmits(['changePart'])

const goVideo = (item: VideoItem) => {
  const itemP = (item as any).p || 1
  const itemShortId = (item as any).shortId;
  const itemResourceRid = (item as any).resourceRid;
  const isMultiPart = itemP > 1;
  const currentPart = props.currentPart || 1;
  const itemVidStr = cleanVid(item.vid);
  
  // 获取当前播放的资源 shortId (rid)
  const currentRid = partList.value[currentPart - 1]?.shortId;
  
  // 判断是否同一分P（兼容 shortId）
  const mainVideoMatch = itemVidStr === currentVidStr.value || itemShortId === currentVidStr.value;
  const resourceMatch = itemResourceRid && currentRid && itemResourceRid === currentRid;
  
  if (mainVideoMatch && resourceMatch) {
    return;
  }
  
  if (mainVideoMatch) {
    if (isMultiPart && itemResourceRid) {
      navigateTo(`/watch?v=${itemShortId || itemVidStr}&rid=${itemResourceRid}`);
    } else {
      emits('changePart', itemP)
    }
    return;
  }
  
  // 切换到其他视频
  const v = itemShortId || itemVidStr;
  if (isMultiPart && itemResourceRid) {
    navigateTo(`/watch?v=${v}&rid=${itemResourceRid}`)
  } else if (itemP > 1) {
    navigateTo(`/watch?v=${v}&p=${itemP}`)
  } else {
    navigateTo(`/watch?v=${v}`)
  }
}

// 获取下一个分P序号
const getNextPart = () => {
  const currentPart = props.currentPart || 1
  
  if (listType.value === 'parts') {
    // 分P类型：下一个分P序号
    if (currentPart < mergedList.value.length) {
      return currentPart + 1
    }
} else {
    // 合集类型：检查当前视频是否有下一个分P
    const currentPartItems = mergedList.value.filter(v => cleanVid(v.vid) === cleanVid(props.vid))
    if (currentPartItems.length > 1) {
      // 找到当前分P的位置
      const currentPartIdx = currentPartItems.findIndex((item, idx) => idx + 1 === currentPart)
      if (currentPartIdx >= 0 && currentPartIdx < currentPartItems.length - 1) {
        return currentPart + 1
      }
    }
  }
  return null
}

// 获取下一个视频
const getNextVideo = () => {
  const currentPart = props.currentPart || 1
  
  if (listType.value === 'parts') {
    const currentIdx = currentPart - 1
    if (currentIdx >= 0 && currentIdx < mergedList.value.length - 1) {
      return mergedList.value[currentIdx + 1]
    }
  } else {
    // 合集类型：匹配当前视频的下一个分P（兼容 shortId）
    const currentVideoParts = mergedList.value.filter(v => {
      const itemVidStr = cleanVid(v.vid)
      const itemShortId = (v as any).shortId
      return itemVidStr === currentVidStr.value || itemShortId === currentVidStr.value
    })
    if (currentVideoParts.length > 1) {
      const currentIdx = currentPart - 1
      if (currentIdx >= 0 && currentIdx < currentVideoParts.length - 1) {
        return currentVideoParts[currentIdx + 1]
      }
    }
    // 当前视频没有更多分P，找下一个视频
    const vidIdx = mergedList.value.findIndex(v => {
      const itemVidStr = cleanVid(v.vid)
      const itemShortId = (v as any).shortId
      return itemVidStr === currentVidStr.value || itemShortId === currentVidStr.value
    })
    if (vidIdx >= 0 && vidIdx < mergedList.value.length - 1) {
      return mergedList.value[vidIdx + 1]
    }
  }
  return null
}

// 处理分P列表
const loadPartList = () => {
  // 优先使用API返回的当前视频分P列表
if (props.resources && props.resources.length > 0) {
    partList.value = props.resources.map((resource, index) => ({
      vid: vidAsNum.value,
      shortId: resource.shortId, // 添加资源的 shortId (rid)
      title: '',
      cover: '',
      duration: resource.duration || 0,
      clicks: 0,
      desc: '',
      p: index + 1,
      partTitle: resource.title || `P${index + 1}`
    }))
  } else {
    partList.value = []
  }
}

// 使用新API加载合集数据
const vidAsNum = computed(() => Number(props.vid));
const loadPlaylist = async (vid: number | string) => {
  playlist.value = null
  videoList.value = []
  partList.value = []

  // 优先使用 props.resources 设置分P列表（页面已有数据）
  if (props.resources && props.resources.length > 0) {
    partList.value = props.resources.map((resource, index) => ({
      vid: vidAsNum.value,
      shortId: resource.shortId, // 添加资源的 shortId (rid)
      title: '',
      cover: '',
      duration: resource.duration || 0,
      clicks: 0,
      desc: '',
      p: index + 1,
      partTitle: resource.title || `P${index + 1}`
    }))
  }

  try {
    const res = await getPlaylistVideoListWithPartsAPI(vid)
    if (res.data.code !== statusCode.OK) return
    
    const data = res.data.data
    
    // 设置是否有合集
    const hasCollection = data.hasCollection === true
    
if (hasCollection && data.playlist) {
      // 有合集
      playlist.value = { id: data.playlist.id, title: data.playlist.title }
videoList.value = (data.videos || []).map((v: any) => ({
        vid: v.vid,
        shortId: v.shortId,
        resourceRid: v.resourceRid,  // 资源shortId
        title: v.title,
        cover: v.cover,
        duration: v.duration,
        clicks: v.clicks,
        desc: v.desc,
        p: v.p,
        partTitle: v.partTitle
      }))
    } else if (!hasCollection && !partList.value.length && data.currentVideoParts && data.currentVideoParts.length > 0) {
      // 没有合集、没有props资源、但API返回了分P
partList.value = data.currentVideoParts.map((r: any, index: number) => ({
        vid: vid,
        title: '',
        cover: '',
        duration: r.duration || 0,
        clicks: 0,
        desc: '',
        p: r.p || index + 1,
        partTitle: r.title || `P${index + 1}`
      }))
    }
  } catch {
    // 静默处理
  }
}

onMounted(() => {
  const saved = localStorage.getItem('video-autonext-collection')
  autonext.value = saved === 'true'
  const savedViewMode = localStorage.getItem('video-collection-view-mode')
  isNumberMode.value = savedViewMode === 'number'
  // 加载数据（包含合集和当前视频分P）
  if (props.vid) loadPlaylist(props.vid)
})

watch(autonext, (val) => {
  localStorage.setItem('video-autonext-collection', val.toString())
})

watch(isNumberMode, (val) => {
  localStorage.setItem('video-collection-view-mode', val ? 'number' : 'title')
})

watch(() => props.vid, (newVid) => {
  if (newVid) loadPlaylist(newVid)
})

// 监听resources变化（备用，如果API没返回分P则使用props）
watch(() => props.resources, () => {
  // 只有当partList为空且props有资源时才使用
  if (partList.value.length === 0 && props.resources && props.resources.length > 0) {
    loadPartList()
  }
})

const hasPlaylist = computed(() => !!playlist.value)

defineExpose({
  autonext: computed(() => autonext.value),
  getNextVideo,
  getNextPart,
  listType,
  hasPlaylist
})
</script>

<style lang="scss" scoped>
.collection-list {
  border-radius: 6px;
  background-color: var(--bg-elev-1);
  border: 1px solid var(--border-color);
  box-shadow: 0 2px 8px var(--shadow-weak);
  margin-bottom: 18px;

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
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
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
}

.list-box {
  list-style: none;
  margin: 0;
  padding: 0 10px 0 6px;
  color: var(--font-primary-1);

  .list-item {
    box-sizing: border-box;
    display: flex;
    justify-content: space-between;
    width: 100%;
    padding: 0 10px;
    height: 30px;
    line-height: 30px;
    color: var(--font-primary-1);
    margin: 5px 0;
    transition: all 0.3s;
    cursor: pointer;
    font-size: 13px;

    .item-content {
      display: flex;
      align-items: center;
      flex: 1;
      overflow: hidden;

      .item-index {
        width: 24px;
        flex-shrink: 0;
        text-align: center;
        color: var(--font-primary-3);
      }

      .item-title {
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
    }

    .item-duration {
      color: var(--font-primary-3);
      flex-shrink: 0;
      margin-left: 8px;
    }

    &:hover {
      background-color: var(--hover-bg);
    }
  }

  .active-item {
    color: var(--primary-hover-color) !important;
    background-color: var(--hover-bg);
    border-radius: 2px;
  }
}

.number-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(40px, 1fr));
  gap: 8px;
  padding: 16px;
  margin: 0;
  color: var(--font-primary-1);

  .number-item {
    box-sizing: border-box;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 40px;
    height: 40px;
    border-radius: 4px;
    background-color: var(--panel-bg);
    border: 1px solid var(--border-color);
    transition: all 0.2s;
    cursor: pointer;
    font-size: 14px;
    font-weight: 500;
    color: var(--font-primary-2);

    &:hover {
      background-color: var(--hover-bg);
      border-color: var(--primary-hover-color);
      color: var(--primary-hover-color);
    }
  }

  .active-number {
    background-color: var(--primary-color);
    border-color: var(--primary-color);
    color: var(--primary-text-color);
  }
}
</style>
