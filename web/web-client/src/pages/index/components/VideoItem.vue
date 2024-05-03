<template>
  <div class="video-item">
    <div class="img">
      <img :src="getResourceUrl(info.cover)" alt="封面" />
    </div>
    <div class="video-info">
      <span v-if="!props.keywords" class="title">{{ info.title }}</span>
      <span v-else class="title">{{ info.title }}</span>
      <div class="author">
        <div class="avatar">
          <common-avatar :url="info.author.avatar" :size="26" :iconsize="16"></common-avatar>
        </div>
        <div class="name-date">
          <span class="name">{{ info.author.name }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  info: VideoType,
  keywords?: string,
}>()

const router = useRouter();

const goUser = (uid: number) => {
  // let videoUrl = router.resolve({ name: "User", params: { uid: uid } });
  // window.open(videoUrl.href, '_blank');
}

//前往视频详情
const goVideo = (vid: number) => {
  // let videoUrl = router.resolve({ name: "Video", params: { vid: vid } });
  // window.open(videoUrl.href, '_blank');
}
</script>

<style lang="scss" scoped>
.video-item {
  width: 100%;
  height: 260px;

  .img {
    width: 100%;
    height: 160px;
    border-radius: 10px;
    overflow: hidden;
    background-color: rgba(0, 0, 0, .2);

    img {
      width: 100%;
      height: 100%;
    }
  }

  .video-info {
    margin-top: 8px;

    .title {
      height: 44px;
      color: #18191c;
      padding-right: 12px;
      font-size: 15px;
      line-height: 22px;
      overflow: hidden;
      text-overflow: ellipsis;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      font-weight: 500;
      cursor: pointer;

      &:hover {
        color: var(--primary-color);
      }
    }

    .author {
      display: flex;
      align-items: center;
      font-size: 13px;
      color: #9499a0;
      margin-top: 5px;

      .avatar {
        width: 26px;
        height: 26px;
        border-radius: 50%;
        background-color: rgba(0, 0, 0, .2);
      }

      .name-date {
        margin-left: 10px;

        .name {
          cursor: pointer;

          &:hover {
            color: var(--primary-color);
          }
        }
      }

    }
  }
}
</style>