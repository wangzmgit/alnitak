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
      <div class="danmaku-item" v-for="item in danmakuList" :key="`${item.time}-${item.text}`">
        <div class="time">{{ formatDanmakuTime(item.time) }}</div>
        <div class="text">{{ item.text }}</div>
        <div class="send-time">{{ moment(item.createdAt).format('MM-DD HH:mm') }}</div>
      </div>
    </el-scrollbar>
  </div>
</template>


<script setup lang="ts">
import { ref } from 'vue';
import { Down as DownIcon } from "@icon-park/vue-next";
import moment from 'moment';

const props = withDefaults(defineProps<{
  height: number;
}>(), {
  height: 300,
})

// 添加弹幕列表相关的代码
const showDanmakuList = ref(false);
const danmakuList = ref<DanmakuType[]>([]);
// 格式化弹幕时间
const formatDanmakuTime = (seconds: number) => {
  const minutes = Math.floor(seconds / 60);
  const remainingSeconds = Math.floor(seconds % 60);
  return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`;
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