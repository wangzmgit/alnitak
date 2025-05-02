<template>
  <div class="danmaku-list-container" :style="showDanmakuList ? `height: ${props.height}px` : 'height: 44px'">
    <div class="danmaku-header" @click="toggleDanmakuList">
      <div class="header-left">
        <span class="title">弹幕列表</span>
        <span class="count">({{ danmakuList.length }})</span>
      </div>
      <div class="header-right">
        <el-icon :class="{ 'is-fold': !showDanmakuList }">
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

defineExpose({
  setDanmaku,
});
</script>

<style lang="scss" scoped>
.danmaku-list-container {
  overflow: hidden;
  transition: height 0.3s;


  .danmaku-header {
    height: 44px;
    border-radius: 6px;
    background-color: #f1f2f3;
    padding: 0 10px 0 16px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    cursor: pointer;
    user-select: none;

    &:hover {
      background-color: #f6f7f8;
    }

    .header-left {
      display: flex;
      align-items: center;

      .title {
        font-size: 14px;
        font-weight: 500;
        color: #18191c;
      }

      .count {
        margin-left: 8px;
        color: #9499a0;
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
    background-color: #fff;
    font-size: 12px;
    color: #61666d;
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
    color: #61666d;
    cursor: pointer;

    .time {
      width: 45px;
      flex-shrink: 0;
    }

    .text {
      flex: 1;
      margin: 0 12px;
      color: #18191c;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }

    .send-time {
      width: 85px;
      flex-shrink: 0;
      text-align: right;
    }

    &:hover {
      background-color: #f6f7f8;
    }
  }
}
</style>