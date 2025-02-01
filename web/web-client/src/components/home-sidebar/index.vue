<template>
  <div class="sidebar" :class="menuFold ? 'sidebar-fold' : ''">
    <div class="scrollbar" v-if="!menuFold">
      <el-scrollbar height="calc(100vh - 56px)" style="width: 220px;">
        <div class="menu-group">
          <nuxt-link class="menu-item menu-item-with-icon" to="/">
            <el-icon class="menu-item-icon" size="20">
              <home-icon></home-icon>
            </el-icon>
            <span class="menu-text">首页</span>
          </nuxt-link>
          <nuxt-link v-if="globalConfig.article" class="menu-item menu-item-with-icon" to="/article/list">
            <el-icon class="menu-item-icon" size="20">
              <article-icon></article-icon>
            </el-icon>
            <span class="menu-text">专栏</span>
          </nuxt-link>
        </div>
        <div class="menu-group partition-list">
          <nuxt-link class="menu-item menu-item-only-text" v-for="item in partitionList" :to="`/partition/${item.id}`">
            {{ item.name }}
          </nuxt-link>
        </div>
        <div class="menu-footer">
          <div class="links">
            <span>
              <a :href="globalConfig.mobile">移动端</a>
            </span>
            <span>
              <a href="https://beian.miit.gov.cn/#/Integrated/index" target="_blank">
                {{ globalConfig.icp }}
              </a>
            </span>
            <span v-show="globalConfig.security">
              <img class="security" src="@/assets/images/common/filing.png" alt="公网安备图标" />
              <a href="https://www.beian.gov.cn/portal/registerSystemInfo" target="_blank">
                {{ globalConfig.security }}
              </a>
            </span>
          </div>
          <div class="copyright">
            <span>
              Powered by <a href="https://github.com/wangzmgit/alnitak" target="_blank">alnitak</a>
            </span>
          </div>
        </div>
      </el-scrollbar>
    </div>
    <div v-else class="sidebar-content-fold">
      <nuxt-link to="/">
        <el-icon class="fold-menu-icon-btn" :size="22">
          <home-icon></home-icon>
        </el-icon>
      </nuxt-link>
      <nuxt-link to="/space/history">
        <el-icon class="fold-menu-icon-btn" size="22">
          <article-icon></article-icon>
        </el-icon>
      </nuxt-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from "vue";
import { asyncGetPartition } from "@/api/partition";
import { globalConfig } from "@/utils/global-config";
import HomeIcon from "@/components/icons/HomeIcon.vue";
import ArticleIcon from "@/components/icons/ArticleIcon.vue";

const props = withDefaults(defineProps<{
  fold: boolean
}>(), {
  fold: false
})

const menuFold = ref(props.fold);
watch(() => props.fold, (newValue) => {
  menuFold.value = newValue;
});

// 获取分区
const partitionList = ref<Array<PartitionType>>([])
const { data } = await asyncGetPartition();
if ((data.value as any).code === statusCode.OK) {
  const partitions = (data.value as any).data.partitions as PartitionType[];
  if (partitions) {
    partitionList.value = partitions.filter((item) => {
      return item.parentId === 0;
    })
  }
}
</script>

<style lang="scss" scoped>
.sidebar {
  width: 220px;
  white-space: nowrap;
  overflow-x: hidden;
  background-color: #fff;

  .scrollbar {
    width: 220px;
  }
}

.sidebar-fold {
  width: 50px;
}


.menu-group {
  border-bottom: 1px solid rgba(0, 0, 0, .1);

  .menu-item {
    display: block;
    height: 30px;
    margin: 6px 10px;
    padding: 8px 0;
    line-height: 30px;
    border-radius: 6px;
    cursor: pointer;
    color: #18191c;

    &-with-icon {
      display: flex;
      align-items: center;
    }

    &-only-text {
      padding-left: 30px;
    }

    &-icon {
      color: #808080;
      padding-left: 6px;
      margin-right: 10px;
    }

    &:hover {
      background-color: rgba(0, 0, 0, .1);
    }
  }

  .menu-active {
    background-color: rgba(0, 0, 0, .1);
  }

  .menu-text {
    color: #18191c;
  }
}

.partition-list {
  min-height: calc(100vh - 326px);
}

.menu-footer {
  color: #666;
  box-sizing: border-box;
  width: 100%;
  padding: 6px 10px;

  a {
    color: #666;
    text-decoration: none;

    &:hover {
      text-decoration: underline;
    }
  }

  .links {
    span {
      display: block;
      margin-top: 6px;
      font-size: 13px;
      line-height: 2;
    }

    .security {
      width: 12px;
      padding-right: 3px;
    }
  }

  .copyright {
    font-size: 12px;
    margin-top: 20px;
    padding-bottom: 8px;
  }
}

// 侧边栏折叠
.sidebar-content-fold {
  .fold-menu-icon-btn {
    color: #808080;
    display: block;
    margin: 10px 6px;
    padding: 8px;
    border-radius: 50%;
    cursor: pointer;

    &:hover {
      background-color: rgba(0, 0, 0, .1);
    }
  }
}
</style>