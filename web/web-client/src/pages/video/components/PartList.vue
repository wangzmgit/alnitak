<template>
  <div class="part-list">
    <!-- 修改分集头部，将自动连播开关放在标题后面 -->
    <div class="part-head">
      <span class="title">视频选集</span>
      <span class="part">({{ props.active }}/{{ resources?.length }})</span>
      <!-- 添加切换显示样式按钮 -->
      <span class="view-mode-switch" @click="toggleViewMode" title="切换显示模式">
        <!-- 网格模式时显示列表图标 -->
        <svg v-if="isNumberMode" class="icon" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 24 24" width="16" height="16">
          <path d="M3.02134 15.21437C3.02134 13.79417 4.17261 12.64294 5.59277 12.64294L8.67849 12.64294C10.09869 12.64294 11.24991 13.79417 11.24991 15.21437L11.24991 18.30009C11.24991 19.7202 10.09869 20.87151 8.67849 20.87151L5.59277 20.87151C4.17261 20.87151 3.02134 19.7202 3.02134 18.30009L3.02134 15.21437zM5.59277 14.35723C5.11939 14.35723 4.73563 14.74097 4.73563 15.21437L4.73563 18.30009C4.73563 18.7734 5.11939 19.15723 5.59277 19.15723L8.67849 19.15723C9.15189 19.15723 9.53563 18.7734 9.53563 18.30009L9.53563 15.21437C9.53563 14.74097 9.15189 14.35723 8.67849 14.35723L5.59277 14.35723z" fill="currentColor"></path>
          <path d="M3.02134 5.70002C3.02134 4.27986 4.17261 3.12859 5.59277 3.12859L8.67849 3.12859C10.09869 3.12859 11.24991 4.27986 11.24991 5.70002L11.24991 8.78571C11.24991 10.20591 10.09869 11.35714 8.67849 11.35714L5.59277 11.35714C4.17261 11.35714 3.02134 10.20591 3.02134 8.78571L3.02134 5.70002zM5.59277 4.84287C5.11939 4.84287 4.73563 5.22663 4.73563 5.70002L4.73563 8.78571C4.73563 9.25911 5.11939 9.64286 5.59277 9.64286L8.67849 9.64286C9.15189 9.64286 9.53563 9.25911 9.53563 8.78571L9.53563 5.70002C9.53563 5.22663 9.15189 4.84287 8.67849 4.84287L5.59277 4.84287z" fill="currentColor"></path>
          <path d="M12.75 5.10859C12.75 4.63521 13.15937 4.25145 13.66431 4.25145L20.06426 4.25145C20.5692 4.25145 20.97857 4.63521 20.97857 5.10859C20.97857 5.58198 20.5692 5.96573 20.06426 5.96573L13.66431 5.96573C13.15937 5.96573 12.75 5.58198 12.75 5.10859z" fill="currentColor"></path>
          <path d="M12.75 9.39429C12.75 8.92089 13.15937 8.53716 13.66431 8.53716L20.06426 8.53716C20.5692 8.53716 20.97857 8.92089 20.97857 9.39429C20.97857 9.86769 20.5692 10.25143 20.06426 10.25143L13.66431 10.25143C13.15937 10.25143 12.75 9.86769 12.75 9.39429z" fill="currentColor"></path>
          <path d="M12.75 14.57143C12.75 14.09803 13.15937 13.71429 13.66431 13.71429L20.06426 13.71429C20.5692 13.71429 20.97857 14.09803 20.97857 14.57143C20.97857 15.04483 20.5692 15.42857 20.06426 15.42857L13.66431 15.42857C13.15937 15.42857 12.75 15.04483 12.75 14.57143z" fill="currentColor"></path>
          <path d="M12.75 18.85714C12.75 18.38374 13.15937 18 13.66431 18L20.06426 18C20.5692 18 20.97857 18.38374 20.97857 18.85714C20.97857 19.33054 20.5692 19.71429 20.06426 19.71429L13.66431 19.71429C13.15937 19.71429 12.75 19.33054 12.75 18.85714z" fill="currentColor"></path>
        </svg>
        <!-- 列表模式时显示网格图标 -->
        <svg v-else class="icon" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 24 24" width="16" height="16">
          <path d="M12.64286 15.21437C12.64286 13.79417 13.79409 12.64294 15.21429 12.64294L18.3 12.64294C19.7202 12.64294 20.87143 13.79417 20.87143 15.21437L20.87143 18.30009C20.87143 19.7202 19.7202 20.87151 18.3 20.87151L15.21429 20.87151C13.79409 20.87151 12.64286 19.7202 12.64286 18.30009L12.64286 15.21437zM15.21429 14.35723C14.74089 14.35723 14.35714 14.74097 14.35714 15.21437L14.35714 18.30009C14.35714 18.7734 14.74089 19.15723 15.21429 19.15723L18.3 19.15723C18.7734 19.15723 19.15714 18.7734 19.15714 18.30009L19.15714 15.21437C19.15714 14.74097 18.7734 14.35723 18.3 14.35723L15.21429 14.35723z" fill="currentColor"></path>
          <path d="M3.12849 15.21437C3.12849 13.79417 4.27976 12.64294 5.69991 12.64294L8.78563 12.64294C10.20583 12.64294 11.35706 13.79417 11.35706 15.21437L11.35706 18.30009C11.35706 19.7202 10.20583 20.87151 8.78563 20.87151L5.69991 20.87151C4.27976 20.87151 3.12849 19.7202 3.12849 18.30009L3.12849 15.21437zM5.69991 14.35723C5.22653 14.35723 4.84277 14.74097 4.84277 15.21437L4.84277 18.30009C4.84277 18.7734 5.22653 19.15723 5.69991 19.15723L8.78563 19.15723C9.25903 19.15723 9.64277 18.7734 9.64277 18.30009L9.64277 15.21437C9.64277 14.74097 9.25903 14.35723 8.78563 14.35723L5.69991 14.35723z" fill="currentColor"></path>
          <path d="M12.64286 5.70002C12.64286 4.27986 13.79409 3.12859 15.21429 3.12859L18.3 3.12859C19.7202 3.12859 20.87143 4.27986 20.87143 5.70002L20.87143 8.78571C20.87143 10.20591 19.7202 11.35714 18.3 11.35714L15.21429 11.35714C13.79409 11.35714 12.64286 10.20591 12.64286 8.78571L12.64286 5.70002zM15.21429 4.84287C14.74089 4.84287 14.35714 5.22663 14.35714 5.70002L14.35714 8.78571C14.35714 9.25911 14.74089 9.64286 15.21429 9.64286L18.3 9.64286C18.7734 9.64286 19.15714 9.25911 19.15714 8.78571L19.15714 5.70002C19.15714 5.22663 18.7734 4.84287 18.3 4.84287L15.21429 4.84287z" fill="currentColor"></path>
          <path d="M3.12849 5.70002C3.12849 4.27986 4.27976 3.12859 5.69991 3.12859L8.78563 3.12859C10.20583 3.12859 11.35706 4.27986 11.35706 5.70002L11.35706 8.78571C11.35706 10.20591 10.20583 11.35714 8.78563 11.35714L5.69991 11.35714C4.27976 11.35714 3.12849 10.20591 3.12849 8.78571L3.12849 5.70002zM5.69991 4.84287C5.22653 4.84287 4.84277 5.22663 4.84277 5.70002L4.84277 8.78571C4.84277 9.25911 5.22653 9.64286 5.69991 9.64286L8.78563 9.64286C9.25903 9.64286 9.64277 9.25911 9.64277 8.78571L9.64277 5.70002C9.64277 5.22663 9.25903 4.84287 8.78563 4.84287L5.69991 4.84287z" fill="currentColor"></path>
        </svg>
      </span>
      <!-- 自动连播开关直接放在这里 -->
      <span class="autoplay-switch" @click="autonext = !autonext">
        <span class="autoplay-text">自动连播</span>
        <span class="switch-button" :class="{ 'on': autonext }"></span>
      </span>
    </div>
    
    <el-scrollbar max-height="340px">
      <!-- 标题模式：垂直列表 -->
      <ul v-if="!isNumberMode" class="list-box">
        <li :class="['list-item', props.active - 1 === index ? 'active-part' : '']" v-for="(item, index) in resources" :key="item.id"
          @click="changePart(index)">
          <div class="item-content">
            <span class="part-num">P{{ index + 1 }}</span>
            <span class="part-title">{{ item.title }}</span>
          </div>
          <div class="duration">{{ toDuration(item.duration) }}</div>
        </li>
      </ul>
      
      <!-- 数字模式：网格布局 -->
      <div v-else class="number-grid">
        <div 
          :class="['number-item', props.active - 1 === index ? 'active-number' : '']" 
          v-for="(item, index) in resources" :key="item.id"
          @click="changePart(index)"
        >
          {{ index + 1 }}
        </div>
      </div>
    </el-scrollbar>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, readonly } from "vue";
import { toDuration } from "@/utils/format";

const emits = defineEmits(['change']);
const props = withDefaults(defineProps<{
  active?: number
  resources: Array<ResourceType>
}>(), {
  active: 1
})

const changePart = (part: number) => {
  emits('change', part + 1)
}

// 添加自动连播状态
const autonext = ref(false);

// 添加显示模式状态
const isNumberMode = ref(false);

// 从本地存储恢复状态
onMounted(() => {
  const saved = localStorage.getItem('video-autonext-parts');
  autonext.value = saved === 'true';
  
  // 恢复显示模式状态
  const savedViewMode = localStorage.getItem('video-part-view-mode');
  isNumberMode.value = savedViewMode === 'number';
});

// 监听状态变化并保存
watch(autonext, (val) => {
  localStorage.setItem('video-autonext-parts', val.toString());
  console.log('分集自动连播状态变更:', val);
});

// 监听显示模式变化并保存
watch(isNumberMode, (val) => {
  localStorage.setItem('video-part-view-mode', val ? 'number' : 'title');
  console.log('分集显示模式变更:', val ? '数字模式' : '标题模式');
});

// 切换显示模式
const toggleViewMode = () => {
  isNumberMode.value = !isNumberMode.value;
};

// 获取下一个分集
const getNextPart = () => {
  const currentIndex = props.resources.findIndex((_, index) => index + 1 === props.active);
  if (currentIndex >= 0 && currentIndex < props.resources.length - 1) {
    return currentIndex + 2; // 返回下一个分集号（从1开始）
  }
  return null;
};

// 暴露给父组件
defineExpose({
  autonext: readonly(autonext),
  getNextPart
});
</script>

<style lang="scss" scoped>
.part-list {
  position: relative;
  border-radius: 6px;
  /* Video PartList：容器底色使用提层背景，增强与页面的层次区分 */
  background-color: var(--bg-elev-1);
  border: 1px solid var(--border-color);
  box-shadow: 0 2px 8px var(--shadow-weak);
  margin-bottom: 18px;

  .part-head {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 12px 16px 8px;
    border-bottom: 1px solid var(--border-color);
    /* Video PartList：头部采用更浅的面板背景以形成层次 */
    background-color: var(--panel-bg);

    .title {
      font-size: 16px;
      color: var(--font-primary-1);
      font-weight: 500;
    }

    .part {
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

  /* 移除重复的 .part-head 定义，以下为微调（保留一处定义即可） */
}

.list-box {
  padding: 0 6px;
  list-style: none;
  color: var(--font-primary-1);
  margin: 0;
  padding: 0 10px 0 6px;

  .list-item {
    box-sizing: border-box;
    display: flex;
    justify-content: space-between;
    overflow: hidden;
    width: 100%;
    padding: 0 10px;
    height: 30px;
    line-height: 30px;
    color: var(--font-primary-1);
    margin: 5px 0;
    transition: all 0.3s;
    cursor: pointer;
    white-space: nowrap;
    overflow: hidden;
    font-size: 13px;

    .item-content {
      display: flex;
      align-items: center;
      flex-shrink: 1;
      overflow: hidden;

      .part-num {
        margin-right: 10px;
      }

      .part-title {
        display: block;
        overflow: hidden;
        text-overflow: ellipsis;
        flex-shrink: 1;
      }
    }

    .duration {
      color: var(--font-primary-3);
    }
  }
}

  .number-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(40px, 1fr));
  gap: 8px;
  padding: 16px;
  list-style: none;
    color: var(--font-primary-1);
  margin: 0;

    .number-item {
    box-sizing: border-box;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 40px;
    height: 40px;
    border-radius: 4px;
      /* Video PartList：数字格子使用面板背景，与容器形成轻微对比 */
      background-color: var(--panel-bg);
      border: 1px solid var(--border-color);
    transition: all 0.2s;
    cursor: pointer;
    text-align: center;
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

  .active-part {
  color: var(--primary-hover-color) !important;
  background-color: var(--hover-bg);
  border-radius: 2px;
}

.list-box .list-item:hover {
  background-color: var(--hover-bg);
}
</style>
