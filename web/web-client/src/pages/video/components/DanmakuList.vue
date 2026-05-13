<template>
  <div class="danmaku-list-container" :style="showDanmakuList ? `height: ${props.height}px` : 'height: 44px'">
    <div class="danmaku-header" @click="toggleDanmakuList">
      <div class="header-left">
        <span class="title">弹幕列表</span>
        <span class="count">({{ danmakuList.length }})</span>
      </div>
      <div class="header-right">
        <el-icon :class="{ 'is-fold': showDanmakuList }">
          <down-icon />
        </el-icon>
      </div>
    </div>
    <!-- 表头 -->
    <div class="danmaku-header-row">
      <div class="time">时间</div>
      <div class="text">弹幕内容</div>
      <div class="send-time">发送时间</div>
    </div>
    <!-- 弹幕列表 -->
    <el-scrollbar :height="props.height - 76">
      <el-tooltip
        v-for="item in danmakuList"
        :key="`${item.time}-${item.text}`"
        content="点击跳转到该弹幕出现的进度"
        placement="left"
        :show-after="400"
        :hide-after="0"
        popper-class="danmaku-item-tip"
      >
        <div class="danmaku-item" @click="handleItemClick(item)">
          <div class="time">{{ formatDanmakuTime(item.time) }}</div>
          <div class="text">{{ item.text }}</div>
          <div class="send-time">{{ formatDate(item.createdAt) }}</div>
        </div>
      </el-tooltip>
    </el-scrollbar>
  </div>
</template>


<script setup lang="ts">
import { ref } from 'vue';
import { Down as DownIcon } from "@icon-park/vue-next";

const props = withDefaults(defineProps<{
  height: number;
}>(), {
  height: 300,
})

const emit = defineEmits<{
  (e: 'seek-time', seconds: number): void;
}>();

const handleItemClick = (item: DanmakuType) => {
  if (typeof item.time === 'number' && isFinite(item.time) && item.time >= 0) {
    emit('seek-time', item.time);
  }
}

// 添加弹幕列表相关的代码
const showDanmakuList = ref(false);
const danmakuList = ref<DanmakuType[]>([]);
// 格式化弹幕时间
const formatDanmakuTime = (seconds: number) => {
  const minutes = Math.floor(seconds / 60);
  const remainingSeconds = Math.floor(seconds % 60);
  return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`;
};

const formatDate = (dateStr: string) => {
  const d = new Date(dateStr);
  const mm = String(d.getMonth() + 1).padStart(2, '0');
  const dd = String(d.getDate()).padStart(2, '0');
  const hh = String(d.getHours()).padStart(2, '0');
  const min = String(d.getMinutes()).padStart(2, '0');
  return `${mm}-${dd} ${hh}:${min}`;
};

// 切换弹幕列表显示状态
const toggleDanmakuList = () => {
  showDanmakuList.value = !showDanmakuList.value;
};

const setDanmaku = (data: DanmakuType[]) => {
  danmakuList.value = data;
}

// 列表按发送时间（createdAt）升序排列：新弹幕直接追加到末尾
// ws 回广播天然按发送顺序到达，与后端初始查询的插入顺序一致
const addDanmaku = (item: DanmakuType) => {
  danmakuList.value.push(item);
}

defineExpose({
  setDanmaku,
  addDanmaku,
});
</script>

<style lang="scss">
/* 非 scoped：el-tooltip popper 被 teleport 到 body，需全局样式才能命中 */
.danmaku-item-tip.el-popper {
  background: var(--bg-elev-1) !important;
  color: var(--font-primary-1) !important;
  border: 1px solid var(--border-color) !important;
  box-shadow: 0 4px 12px var(--shadow-weak);
  font-size: 12px;

  .el-popper__arrow::before {
    background: var(--bg-elev-1) !important;
    border: 1px solid var(--border-color) !important;
  }
}
</style>

<style lang="scss" scoped>
.danmaku-list-container {
  overflow: hidden;
  transition: height 0.3s;


  .danmaku-header {
    height: 44px;
    border-radius: 6px;
    background-color: var(--hover-bg);
    padding: 0 10px 0 16px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    cursor: pointer;
    user-select: none;

    &:hover { background-color: var(--hover-bg); }

    .header-left {
      display: flex;
      align-items: center;

      .title {
        font-size: 14px;
        font-weight: 500;
        color: var(--font-primary-1);
      }

      .count {
        margin-left: 8px;
        color: var(--font-primary-3);
        font-size: 12px;
      }
    }

    .header-right {
      .el-icon {
        transition: transform 0.3s;

        &.is-fold {
          color: #61666d;
          transform: rotate(-180deg);
        }
      }
    }
  }

  .danmaku-header-row {
    padding: 8px 16px;
    display: flex;
    align-items: center;
    background-color: var(--bg-elev-1);
    font-size: 12px;
    color: var(--font-primary-2);
    font-weight: 500;

    .time {
      width: 45px;
    }

    .text {
      flex: 1;
      margin: 0 12px;
    }

    .send-time {
      width: 85px;
      text-align: right;
    }
  }

  .danmaku-item {
    padding: 0 8px 0 16px;
    display: flex;
    font-size: 12px;
    align-items: center;
    height: 24px;
    color: var(--font-primary-2);
    cursor: pointer;

    .time {
      width: 45px;
      flex-shrink: 0;
    }

    .text {
      flex: 1;
      margin: 0 12px;
      color: var(--font-primary-1);
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }

    .send-time {
      width: 85px;
      flex-shrink: 0;
      text-align: right;
    }

    &:hover { background-color: var(--hover-bg); }
  }
}
</style>