<template>
  <div class="collection-bg">
    <div class="collection-card">
      <div class="card-title">
        <span>添加到收藏夹</span>
        <el-icon class="close-icon" @click="closeCard">
          <close-icon></close-icon>
        </el-icon>
      </div>
      <div class="card-content">
        <el-scrollbar height="300px">
          <ul v-if="collectionList.length" class="group-list">
            <li class="collection-item" v-for="item in collectionList" @click="collectionClick(item)">
              <div class="item-left"> <span class="name">{{ item.name }}</span>
                <span class="open">[{{ item.open ? '公开' : '私有' }}]</span>
              </div>
              <div class="item-right">
                <el-checkbox v-model="item.checked" @click="checkboxClick" />
              </div>
            </li>
          </ul>
          <div class="add-collection">
            <div v-if="!isAdd" class="add-btn" @click="changeAdd">
              <el-icon class="add-icon" size="16">
                <plus-icon style="height: 16px;"></plus-icon>
              </el-icon>
              <span>新建收藏夹</span>
            </div>
            <div v-else class="input-group">
              <!-- 使用mousedown触发而不是click触发，因为input的blur要早于click -->
              <input class="input" v-model="collectionName" placeholder="最多可输入20个字" maxlength="20" @blur="changeAdd">
              <button class="create-btn" @mousedown="addCollection">新建</button>
            </div>
          </div>
        </el-scrollbar>
      </div>
      <div class="card-bottom">
        <button class="submit-btn" @click="submitCollect">确定</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onBeforeMount, nextTick, ref } from 'vue';
import { ElMessage } from 'element-plus';
import { Close as CloseIcon, Plus as PlusIcon } from '@icon-park/vue-next';
import { getCollectVideoInfoAPI, collectVideoAPI } from "@/api/collect";
import { addCollectionAPI, getCollectionListAPI } from '@/api/collection';


const emits = defineEmits(['close']);
const props = defineProps<{
  vid: number
}>();

const loading = ref(true);// 加载中

const collectionList = ref<CollectionType[]>([]);
const getCollectionList = async () => {    //获取收藏夹列表
  const res = await getCollectionListAPI();
  if (res.data.code === statusCode.OK) {
    if (res.data.data.collections) {
      collectionList.value = res.data.data.collections;
    }
  }
}

const defaultChecked = ref<number[]>([]);// 默认选中
const getCollectInfo = async () => {// 获取收藏信息
  const res = await getCollectVideoInfoAPI(props.vid);
  if (res.data.code === statusCode.OK) {
    if (res.data.data.collectionIds) {
      defaultChecked.value = res.data.data.collectionIds;
      collectionList.value.forEach((item) => {
        if (defaultChecked.value.includes(item.id)) {
          item.checked = true;
        }
      })
    }
  }
}

// 点击收藏夹
const collectionClick = (collection: CollectionType) => {
  collection.checked = !collection.checked;
}

const checkboxClick = (e: Event) => {
  e.stopPropagation();
}

// 新建收藏夹
const isAdd = ref(false);
const collectionName = ref("");
const addInput = ref<HTMLElement | null>(null);
const changeAdd = () => {
  isAdd.value = !isAdd.value;
  if (isAdd.value) {
    //此时input不存在，无法修改焦点
    nextTick(() => {
      addInput.value?.focus();
    })
  }
}

// 新建收藏夹
const addCollection = async () => {
  if (!collectionName.value) {
    ElMessage.error('收藏夹名不能为空');
    return;
  }

  const res = await addCollectionAPI(collectionName.value);
  if (res.data.code === statusCode.OK) {
    collectionList.value.push({
      id: res.data.data.id as number,
      cover: "",
      name: collectionName.value,
      desc: "",
      open: false,
      checked: false,
    })
    collectionName.value = "";
  } else {
    ElMessage.error('收藏夹创建失败');
  }
}



//保存收藏
const submitCollect = async () => {
  const checkedValue = collectionList.value.filter(v => v.checked).map(v => v.id);
  //原数组不存在新数组存在表示添加
  const addList = checkedValue.filter((v) => {
    return defaultChecked.value.indexOf(v) == -1
  })

  //原数组存在新数组不存在表示删除
  const cancelList = defaultChecked.value.filter((v) => {
    return checkedValue.indexOf(v) == -1
  })

  const res = await collectVideoAPI({ vid: props.vid, addList, cancelList });
  if (res.data.code === statusCode.OK) {
    var count = 0;
    //否则收藏量不变
    if (defaultChecked.value.length === 0 && checkedValue.length !== 0)
      count = 1; //之前没有收藏，之后收藏了，收藏量+1
    else if (defaultChecked.value.length !== 0 && checkedValue.length === 0)
      count = -1;//之前收藏了，之后没有收藏，收藏量-1
    emits('close', count);
  }
}

const closeCard = () => {
  emits('close', 0);
}

onBeforeMount(async () => {
  await getCollectionList();
  await getCollectInfo();
  loading.value = false;
})
</script>

<style lang="scss" scoped>
.collection-bg {
  top: 0;
  left: 0;
  z-index: 999;
  position: fixed;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.45);

  .collection-card {
    position: absolute;
    top: 50%;
    left: 50%;
    margin: -210px 0 0 -210px;
    width: 420px;
    height: 420px;
    padding: 6px 12px;
    background-color: #fff;
    border-radius: 10px;
    animation: fadein .3s ease-in;
    box-shadow: 16px 16px 50px -12px rgba(0, 0, 0, 0.8);
  }
}

.card-title {
  position: relative;
  padding: 0 20px;
  height: 50px;
  line-height: 50px;
  font-size: 16px;
  color: #18191C;
  text-align: center;

  .close-icon {
    position: absolute;
    right: 20px;
    line-height: 50px;
    color: #adadad;
    width: 12px;
    height: 50px;
    cursor: pointer;
  }
}

.card-content {
  padding: 0 20px;

  .group-list {
    max-height: 300px;
    padding: 0 0 14px;
    margin: 12px 0 0;

    .collection-item {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding-bottom: 24px;
      font-size: 14px;
      color: #18191C;
      cursor: pointer;
      list-style: none;

      .item-left {
        display: flex;
        align-items: center;

        .name {
          max-width: 220px;
          display: inline-block;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
          vertical-align: middle;
        }

        .open {
          color: #9499A0;
        }
      }

      &:hover {
        color: var(--primary-hover-color);
      }

      &:last-child {
        padding-bottom: 0;

      }
    }
  }

  .add-collection {
    width: 100%;

    .add-btn {
      display: flex;
      align-items: center;
      height: 34px;
      line-height: 34px;
      padding: 0 12px;
      border: 1px solid #9499A0;
      border-radius: 4px;
      cursor: pointer;
      font-size: 12px;
      color: #61666D;

      .add-icon {
        margin-right: 3px;
      }

      &:hover {
        border-color: var(--primary-hover-color);
      }
    }
  }

  .input-group {
    display: flex;
    height: 34px;
    border-radius: 4px;
    border: 1px solid var(--primary-hover-color);

    .input {
      width: calc(100% - 110px);
      outline: none;
      box-shadow: none;
      border: none;
      font-size: 12px;
      margin: 0 10px;
    }

    .create-btn {
      border: none;
      cursor: pointer;
      color: #fff;
      font-size: 14px;
      width: 90px;
      height: 100%;
      background-color: var(--primary-color);
      border-radius: 0 4px 4px 0;
      border-left: 1px solid var(--primary-hover-color);

      &:hover {
        background-color: var(--primary-hover-color);

      }
    }
  }
}

.card-bottom {
  height: 76px;
  text-align: center;
  margin: 0 36px;

  .submit-btn {
    font-size: 14px;
    width: 160px;
    height: 40px;
    margin-top: 18px;
    border: none;
    border-radius: 4px;
    background-color: var(--primary-color);
    color: #fff;
    cursor: pointer;

    &:hover {
      background-color: var(--primary-hover-color);

    }
  }

  .submit-btn-disabled {
    background-color: #E3E5E7;
    color: #9499A0;

    &:hover {
      background-color: #E3E5E7;
    }
  }
}
</style>