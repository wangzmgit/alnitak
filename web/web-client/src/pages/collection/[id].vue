<template>
  <div class="collection">
    <header-bar></header-bar>
    <div class="collection-main">
      <div class="collection-left">
        <div class="video-box">
          <el-scrollbar>
            <ul class="video-list">
              <li class="video-item" v-for="(item, index) in videoList" :key="index">
                <div class="item-left">
                  <div class="cover">
                    <img v-if="item.cover" :src="getResourceUrl(item.cover)" alt="收藏夹封面">
                  </div>
                </div>
                <div class="item-center">
                  <nuxt-link class="title" :to="`/video/${item.vid}`">{{ item.title }}</nuxt-link>
                  <span class="desc">简介：{{ item.desc }}</span>
                  <nuxt-link class="desc" :to="`/user/${item.author.uid}`">UP主：{{ item.author.name }}</nuxt-link>
                </div>
                <div class="item-right" v-if="userInfo?.uid === collection?.author.uid">
                  <el-dropdown>
                    <el-button link class="more-btn" plain>
                      <el-icon size="16">
                        <more-icon></more-icon>
                      </el-icon>
                    </el-button>
                    <template #dropdown>
                      <el-dropdown-menu>
                        <el-dropdown-item @click="removeVideo(item.vid)">取消收藏</el-dropdown-item>
                      </el-dropdown-menu>
                    </template>
                  </el-dropdown>
                </div>
              </li>
            </ul>
          </el-scrollbar>
        </div>
        <div class="page-box">
          <el-pagination small background layout="prev, pager, next" :page="page" :page-size="8" :total="total"
            @current-change="pageChange" />
        </div>
      </div>
      <div class="collection-right">
        <div class="cover">
          <img v-if="collection?.cover" :src="getResourceUrl(collection.cover)" alt="视频封面">
        </div>
        <div class="info">
          <span class="title">{{ collection?.name }}</span>
          <span class="desc" v-if="collection">简介：{{ collection.desc }}</span>
          <div class="desc" v-if="collection">
            <span v-if="collection.createdAt">{{ formatDate(collection.createdAt) }}</span>
            <span>・</span>
            <span class="open">{{ collection.open ? "公开" : "私密" }}</span>
            <nuxt-link class="author" :to="`/user/${collection.author.uid}`">{{ collection.author.name }}</nuxt-link>
          </div>
        </div>
      </div>

    </div>
  </div>
</template>

<script setup lang="ts">
import { onBeforeMount, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import HeaderBar from '@/components/header-bar/index.vue';
import { getCollectionInfoAPI, getCollectVideoAPI } from '@/api/collection';
import { MoreOne as MoreIcon } from '@icon-park/vue-next';
import { getUserInfoAPI } from '@/api/user';
import { collectVideoAPI } from '@/api/collect';

const route = useRoute();
const collectionId = Number(route.params.id);
const userInfo = ref<UserInfoType>();
const getUserInfo = async () => {
  const res = await getUserInfoAPI();
  if (res.data.code === statusCode.OK) {
    userInfo.value = res.data.data.userInfo;
  }
}


const collection = ref<CollectionInfoType>();
const getCollectionInfo = async (id: number) => {
  const res = await getCollectionInfoAPI(id);
  if (res.data.code === statusCode.OK) {
    collection.value = res.data.data.collection;
  }
}

const page = ref(1);
const total = ref(0);
const videoList = ref<Array<VideoType>>([]);
//获取收藏视频内容
const getCollectVideo = async () => {
  const res = await getCollectVideoAPI(collectionId, page.value, 8);
  if (res.data.code === statusCode.OK) {
    total.value = res.data.data.total;
    videoList.value = res.data.data.videos;
  }
}

//页码改变
const pageChange = (target: number) => {
  page.value = target;
  getCollectVideo();
}

//移除视频
const removeVideo = async (vid: number) => {
  const res = await collectVideoAPI({ vid, addList: [], cancelList: [collectionId] });
  if (res.data.code === statusCode.OK) {
    getCollectVideo();
  } else {
    ElMessage.error(res.data.msg || "移除失败");
  }
}

onBeforeMount(() => {
  getUserInfo();
  getCollectionInfo(collectionId);
  getCollectVideo();
})
</script>

<style lang="scss" scoped>
.collection {
  height: 100%;
  width: 100%;
  top: 0;
  bottom: 0;
  position: fixed;
  overflow-y: scroll;
  background-color: #f5f6f7;
}

.collection-main {
  width: calc(100% - 200px);
  box-sizing: border-box;
  margin: 20px auto 0 auto;
  display: flex;
  flex-wrap: nowrap;
  justify-content: space-between;

  .collection-left {
    width: 100%;
    margin-right: 12px;
    background-color: #fff;

    .video-box {
      height: calc(100% - 40px);
    }

    .video-list {
      list-style: none;
      box-sizing: border-box;
      width: 100%;
      margin: 0;
      padding: 16px 16px 10px;

      .video-item {
        display: flex;
        padding: 16px 0;
        width: 100%;
        height: 80px;
        margin-bottom: 12px;
        border-bottom: 1px solid #e6e6e6;
        padding-bottom: 12px;

        .item-left {
          width: 120px;
          height: 80px;
          margin-right: 10px;

          .cover {
            border-radius: 2px;
            width: 100%;
            height: 100%;
            background-color: #f1f2f3;

            img {
              width: 100%;
              height: 100%;
              border-radius: 2px;
            }
          }
        }

        .item-center {
          flex: 1;

          .title {
            font-size: 14px;
            color: #212121;
            line-height: 18px;
            margin: 0 0 26px;
            cursor: pointer;
            overflow: hidden;
            text-overflow: ellipsis;
            display: -webkit-box;
            -webkit-line-clamp: 1;
            -webkit-box-orient: vertical;

            &:hover {
              color: var(--primary-hover-color);
            }
          }

          .desc {
            font-size: 12px;
            color: #999;
            overflow: hidden;
            text-overflow: ellipsis;
            display: -webkit-box;
            -webkit-line-clamp: 1;
            -webkit-box-orient: vertical;
          }
        }

        .item-right {
          width: 90px;
          height: 100%;
          display: flex;
          align-items: center;
          justify-content: center;
        }
      }
    }

    .page-box {
      padding: 8px 16px;
    }
  }

  /**右半部分 */
  .collection-right {
    width: 320px;
    height: 600px;
    min-width: 320px;
    padding: 10px;
    background: #fff;

    .cover {
      width: 100%;
      height: 180px;
      margin-bottom: 12px;
      border-radius: 2px;
      background-color: #f1f2f3;

      img {
        width: 100%;
        height: 100%;
        border-radius: 2px;
      }
    }

    .info {
      .title {
        display: block;
        height: auto;
        font-size: 16px;
        color: #212121;
        line-height: 18px;
        margin-bottom: 20px;
      }

      .desc {
        display: block;
        font-size: 13px;
        color: #666;
        margin: 6px 0;

        .author {
          float: right;
          cursor: pointer;
          color: #666;

          &:hover {
            color: var(--primary-hover-color);
          }
        }
      }


    }
  }
}
</style>