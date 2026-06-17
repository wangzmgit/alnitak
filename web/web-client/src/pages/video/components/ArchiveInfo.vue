<template>
  <div class="archive-data">
    <!--点赞收藏-->
    <div class="archive-item">
      <el-icon :class="[likeAnimation, archive.hasLike ? 'active' : 'icon']" @click="likeClick">
        <like-icon></like-icon>
      </el-icon>
      <p>{{ stat.like }}</p>
    </div>
    <div class="archive-item">
      <el-icon :class="archive.hasCollect ? 'active' : 'icon'" @click="collectClick">
        <collect-icon></collect-icon>
      </el-icon>
      <p>{{ stat.collect }}</p>
    </div>
    <!-- 分享按钮 -->
    <div class="archive-item share-item">
      <el-icon class="icon">
        <share-icon></share-icon>
      </el-icon>
      <p>{{ stat.share }}</p>
      <div class="share-popover">
        <div class="share-popover-content">
          <el-tabs v-model="shareTab">
            <el-tab-pane label="分享链接" name="link">
              <div class="embed-box">
                <el-input v-model="shareUrl" name="share-url" readonly></el-input>
                <el-button type="primary" @click="copyUrl">复制链接</el-button>
              </div>
            </el-tab-pane>
            <el-tab-pane label="嵌入代码" name="embed">
              <div class="embed-box">
                <el-input v-model="embedCode" name="embed-code" readonly></el-input>
                <el-button type="primary" @click="copyEmbed">复制嵌入代码</el-button>
              </div>
            </el-tab-pane>
          </el-tabs>
        </div>
      </div>
    </div>
    <collection-list v-if="showCollect" :vid="shortId || vid" @close="closeCollectionCard"></collection-list>
  </div>
</template>

<script setup lang="ts">
import { ref, onBeforeMount, computed, reactive, watch } from 'vue';
import { ElIcon } from 'element-plus';
import { useRoute } from 'vue-router';
import { statusCode } from '@/utils/status-code';
import { requireLogin } from '@/utils/require-login';
import { useAuthStore } from '@/stores/auth-store';
import LikeIcon from "@/components/icons/LikeIcon.vue";
import CollectIcon from "@/components/icons/CollectIcon.vue";
import ShareIcon from '@/components/icons/ShareIcon.vue';
import { getVideoArchiveStatAPI, shareVideoAPI } from "@/api/archive";
import { getLikeVideoStatusAPI, likeVideoAPI, cancelLikeVideoAPI } from "@/api/like";
import { getCollectVideoStatusAPI } from '@/api/collect';
import CollectionList from './CollectionList.vue';

const props = defineProps<{
  vid: number | string;
  shortId?: string;
}>();

// 点赞收藏数据
const stat = ref<{ like: number, collect: number, share: number }>({
  like: 0,
  collect: 0,
  share: 0
});

const loading = ref(true);
const archive = reactive({ // 是否点赞收藏
  hasCollect: false,
  hasLike: false
})

const auth = useAuthStore();

const refreshViewerStatus = async () => {
  if (!auth.isLoggedIn) {
    archive.hasLike = false;
    archive.hasCollect = false;
    showCollect.value = false;
    return;
  }
  await getLikeStatus();
  await getCollectStatus();
};

// 分享相关
const shareTab = ref('link');

// 使用 process.client 检查是否在客户端
const shareUrl = computed(() => {
  if (process.client) {
    return window.location.href;
  }
  return '';
});

const route = useRoute();
const embedCode = computed(() => {
  if (process.client) {
    const v = props.shortId || String(props.vid);
    // 优先使用 rid 精准定位
    const rid = route.query.rid;
    if (rid) {
      const url = window.location.origin + `/embed/watch?v=${v}&rid=${rid}`;
      return `<iframe src='${url}' width='800' height='450' frameborder='0' allowfullscreen></iframe>`;
    }
    const part = Number(route.query.p) || 1;
    const url = window.location.origin + `/embed/watch?v=${v}` + (part > 1 ? `&p=${part}` : '');
    return `<iframe src='${url}' width='800' height='450' frameborder='0' allowfullscreen></iframe>`;
  }
  return '';
});

const copyText = async (text: string, msg: string) => {
  try {
    await navigator.clipboard.writeText(text);
    ElMessage.success(msg);
  } catch (e) {
    // 降级处理
    const textarea = document.createElement('textarea');
    textarea.value = text;
    textarea.style.position = 'fixed';
    textarea.style.opacity = '0';
    document.body.appendChild(textarea);
    textarea.focus();
    textarea.select();
    try {
      document.execCommand('copy');
      ElMessage.success(msg);
    } catch (err) {
      ElMessage.error('复制失败，请手动复制');
    }
    document.body.removeChild(textarea);
  }
};

const doShare = async () => {
  await shareVideoAPI(props.shortId || props.vid);
  stat.value.share++;
};
const copyUrl = () => { copyText(shareUrl.value, '播放地址已复制'); doShare(); };
const copyEmbed = () => { copyText(embedCode.value, '嵌入代码已复制'); doShare(); };

//获取点赞收藏关注信息
const getArchiveStat = async () => {
  const res = await getVideoArchiveStatAPI(props.shortId || props.vid);
  if (res.data.code === statusCode.OK) {
    stat.value = res.data.data.stat;
  }
}

// 获取是否点赞
const getLikeStatus = async () => {
  const res = await getLikeVideoStatusAPI(props.shortId || props.vid);
  if (res.data.code === statusCode.OK) {
    archive.hasLike = res.data.data.like;
  }
}

// 获取是否收藏
const getCollectStatus = async () => {
  const res = await getCollectVideoStatusAPI(props.shortId || props.vid);
  if (res.data.code === statusCode.OK) {
    archive.hasCollect = res.data.data.collect;
  }
}

const likeAnimation = ref('');
const likeClick = async () => { // 点赞点赞按钮
  if (loading.value) return;
  if (!(await requireLogin('点赞'))) return;
  const videoId = props.shortId || props.vid;
  if (!archive.hasLike) {
    //调用点赞接口
    await likeVideoAPI(videoId);
    likeAnimation.value = 'like-active';
    stat.value.like++;
  } else {
    await cancelLikeVideoAPI(videoId);
    stat.value.like--;
  }

  archive.hasLike = !archive.hasLike;
}


const showCollect = ref(false);
const collectClick = async () => {
  if (loading.value) return;
  if (!(await requireLogin('收藏'))) return;
  showCollect.value = true;
};
// 关闭收藏弹窗
const closeCollectionCard = (val: number) => {
  if (val === 1) {
    stat.value.collect++;
    archive.hasCollect = true;
  } else if (val === -1) {
    stat.value.collect--;
    archive.hasCollect = false;
  }

  showCollect.value = false;
}

onBeforeMount(async () => {
  await getArchiveStat();
  await refreshViewerStatus();

  loading.value = false;
})

watch(
  () => auth.isLoggedIn,
  async () => {
    if (loading.value) return;
    await refreshViewerStatus();
  }
);
</script>

<style lang="scss" scoped>
.archive-data {
  height: 30px;

  .archive-item {
    float: left;
    user-select: none;
    margin-right: 20px;

    i,
    .icon {
      font-size: 26px;
      width: 26px;
      height: 26px;
      line-height: 30px;
      cursor: pointer;
      vertical-align: middle;
    }

    p {
      font-size: 16px;
      float: right;
      margin: 0 6px;
      line-height: 30px;
    }

    .icon:hover {
      color: var(--primary-color);
    }

    .active {
      color: var(--primary-color);
    }

    .like-active {
      animation: scaleDraw .3s ease-in-out;
    }
  }

  .share-item {
    position: relative;

    .share-popover {
      display: none;
      position: absolute;
      bottom: 26px;
      left: 0;
      z-index: 100;

      .share-popover-content {
        margin-bottom: 10px;
        background: var(--bg-elev-1);
        border: 1px solid var(--border-color);
        border-radius: 8px;
        box-shadow: 0 2px 8px var(--shadow-weak);
        padding: 16px 20px 8px 20px;
        min-width: 320px;
        min-height: 120px;

        .embed-box {
          display: flex;
          align-items: center;
          gap: 8px;
          margin-bottom: 8px;
        }
      }
    }

    &:hover .share-popover {
      display: block;
    }
  }
}


@keyframes scaleDraw {
  0% {
    transform: scale(1);
    /*开始为原始大小*/
  }

  25% {
    transform: scale(1.2);
    /*放大1.1倍*/
  }

  100% {
    transform: scale(1);
  }
}
</style>
