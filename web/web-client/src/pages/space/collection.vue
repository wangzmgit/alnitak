<template>
  <div class="collection">
    <p class="collection-title">收藏夹</p>
    <ul class="collection-list">
      <li class="collection-item" v-for="(item, index) in collectionList">
        <div class="cover">
          <img v-if="item.cover" :src="item.cover" alt="视频封面">
        </div>
        <div class="info-box">
          <nuxt-link class="info-title" :to="`/collection/${item.id}`">{{ item.name }}</nuxt-link>
          <div class="info-desc">{{ item.desc }}</div>
          <div class="info-meta">
            <span class="open">{{ item.open ? '公开' : '私密' }}</span>
            <span v-if="item.createdAt" class="time">{{ formatTime(item.createdAt) }}</span>
          </div>
        </div>
        <el-dropdown>
          <el-button class="more-btn" plain>
            <el-icon size="16">
              <more-icon></more-icon>
            </el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item @click="editCollection(item)">编辑</el-dropdown-item>
              <el-dropdown-item @click="deleteCollection(item.id, index)">删除</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </li>
    </ul>
    <!-- 编辑对话框 -->
    <client-only>
      <collection-modal ref="collectionModalRef" @edit-finish="editFinish"></collection-modal>
    </client-only>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onBeforeMount } from 'vue';
import { MoreOne as MoreIcon } from '@icon-park/vue-next';
import CollectionModal from "./components/CollectionModal.vue";
import { getCollectionListAPI, deleteCollectionAPI } from '@/api/collection';

const collectionList = ref<CollectionType[]>([]);
const getCollectionList = async () => {
  const res = await getCollectionListAPI();
  if (res.data.code === statusCode.OK) {
    collectionList.value = res.data.data.collections;
  }
}

const collectionModalRef = ref<InstanceType<typeof CollectionModal>>();
const editCollection = (collection: CollectionType) => {
  collectionModalRef.value?.open(collection);
}

const editFinish = (collection: EditCollectionType) => {
  const index = collectionList.value.findIndex(item => item.id === collection.id);
  if (index !== -1) {
    collectionList.value[index].name = collection.name || "";
    collectionList.value[index].cover = collection.cover;
    collectionList.value[index].desc = collection.desc;
    collectionList.value[index].open = collection.open;
  }
}

const deleteCollection = (id: number, index: number) => {
  ElMessageBox.confirm('此操作将会删除收藏夹和收藏夹的内容，且不可恢复，是否继续？', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  }).then(async () => {
    const res = await deleteCollectionAPI(id);
    if (res.data.code === statusCode.OK) {
      collectionList.value.splice(index, 1);
    } else {
      ElMessage.error(res.data.msg || '删除失败');
    }
  })
}

onBeforeMount(() => {
  getCollectionList();
})
</script>

<style lang="scss" scoped>
.collection {
  padding-left: 12px;

  .collection-title {
    font-size: 18px;
    margin-top: 20px;
    padding-left: 10px;
  }
}

.collection-list {
  list-style: none;
  padding: 0;

  .collection-item {
    display: flex;
    align-items: center;
    border-bottom: 1px solid #eee;
    padding: 20px 0;

    .cover {
      width: 160px;
      height: 100px;
      background-color: #f1f2f3;
      border-radius: 2px;

      img {
        width: 100%;
        height: 100%;
        border-radius: 2px;
      }
    }

    .info-box {
      width: calc(100% - 160px - 64px);
      margin: 0 20px;
      height: 100px;

      .info-title {
        display: block;
        color: #1f2225;
        font-size: 16px;
        height: 20px;
        line-height: 20px;
        margin-bottom: 10px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;

        &:hover {
          color: var(--primary-hover-color);
        }
      }

      .info-desc {
        height: 40px;
        font-size: 12px;
        margin-bottom: 10px;
        color: #777;
        line-height: 20px;
        overflow: hidden;
      }

      .info-meta {
        color: #aaa;
        font-size: 12px;
        line-height: 20px;

        .time {
          margin-left: 16px;
        }
      }
    }

    .more-btn {
      width: 32px;
      height: 32px;
      margin-right: 20px;
    }
  }
}
</style>
