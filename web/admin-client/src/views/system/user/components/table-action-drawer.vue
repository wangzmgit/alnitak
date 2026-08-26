<template>
  <n-drawer v-model:show="drawerVisible" :width="600">
    <n-drawer-content title="用户详情">
      <n-form v-if="data" label-placement="top">
        <n-grid :cols="24" :x-gap="18">
          <n-form-item-grid-item :span="12" label="用户ID">{{ data.uid }}</n-form-item-grid-item>
          <n-form-item-grid-item :span="12" label="用户名">{{ data.name }}</n-form-item-grid-item>
          <n-form-item-grid-item :span="12" label="邮箱">{{ data.email }}</n-form-item-grid-item>
          <n-form-item-grid-item :span="12" label="性别">{{ toGender(data.gender) }}</n-form-item-grid-item>
          <n-form-item-grid-item :span="12" label="生日">{{ formatDate(data.birthday || "") }}</n-form-item-grid-item>
          <n-form-item-grid-item :span="12" label="注册时间">{{ formatTime(data.createdAt || "") }}</n-form-item-grid-item>
          <n-form-item-grid-item :span="12" label="角色">{{ getRoleName(data.role) }}</n-form-item-grid-item>
          <n-form-item-grid-item :span="24" label="个性签名">{{ data.sign }}</n-form-item-grid-item>
        </n-grid>
      </n-form>
      <div class="ban-box" v-if="banList.length">
        <span class="ban-record-title">封禁记录</span>
        <n-table>
          <thead>
            <tr>
              <th>封禁时间</th>
              <th>解封时间</th>
              <th>封禁原因</th>
              <th>封禁状态</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in banList" :key="item.id">
              <td>{{ formatTime(item.createdAt) }}</td>
              <td>{{ formatDate(item.endTime) }}</td>
              <td>{{ item.reason }}</td>
              <td>
                <n-button v-if="item.status === 0" @click="unBan(item.id)">{{ toBanStatus(item.status) }}</n-button>
                <span v-else>{{ toBanStatus(item.status) }}</span>
              </td>
            </tr>
          </tbody>
        </n-table>
      </div>
      <template #footer>
        <n-button class="btn" @click="visibleBanModal = true">封禁</n-button>
        <n-button class="btn" type="primary" @click="closeDrawer">完成</n-button>
      </template>
    </n-drawer-content>
    <ban-modal v-if="props.data" v-model:visible="visibleBanModal" :uid="props.data.uid"></ban-modal>
  </n-drawer>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { formatDate, formatTime } from '@/utils/format';
import { statusCode } from '@/utils/status-code';
import { getUserBanRecordAPI, unbanUserAPI } from '@/api/user';
import { getAllRoleListAPI } from '@/api/role';
import banModal from './ban-modal.vue';
import { NButton, NTable, NDrawer, NDrawerContent, NScrollbar, NForm, NGrid, NFormItemGridItem, useDialog, useMessage } from "naive-ui";

const emit = defineEmits(['update:visible']);
const props = withDefaults(defineProps<{
  visible: boolean; //弹窗可见性
  data: UserInfoType;
}>(), {
  visible: false,
})

const dialog = useDialog();
const message = useMessage();
const visibleBanModal = ref(false);

const roleList = ref<Array<{ code: string, name: string }>>([])
const getAllRoleList = async () => {
  const res = await getAllRoleListAPI();
  if (res.data.code === statusCode.OK) {
    roleList.value = res.data.data.roles;
  }
}
getAllRoleList();

const getRoleName = (code: string | undefined) => {
  if (!code) return "未知";
  const role = roleList.value.find(r => r.code === code);
  return role ? role.name : "未知";
}

const drawerVisible = computed({
  get() {
    return props.visible;
  },
  set(visible) {
    emit('update:visible', visible);
  }
});

const closeDrawer = () => {
  drawerVisible.value = false;
}

const toGender = (gender: number | undefined) => {
  if (gender == 1) return "男";
  else if (gender == 2) return "女";
  else return "未知";
}

const toBanStatus = (status: number) => {
  switch (status) {
    case 0:
      return "封禁中";
    case 1:
      return "管理员解封";
    case 2:
      return "自动解封";
    case 3:
      return "永久封禁";
    case 4:
      return "封禁撤销";
    default:
      return "未知";
  }
}

const banList = ref<BanRecordType[]>([]);
const getBanResourceList = async (uid: number) => {
  const res = await getUserBanRecordAPI(uid);
  if (res.data.code === statusCode.OK) {
    if (res.data.data.list) {
      banList.value = res.data.data.list;
    } else {
      banList.value = [];
    }
  }
}

const unBan = (id: number) => {
  dialog.warning({
    title: '提示',
    content: '此操作将取消封禁，确定吗？',
    positiveText: '确定',
    negativeText: '不确定',
    onPositiveClick: async () => {
      const res = await unbanUserAPI(id);
      if (res.data.code === statusCode.OK) {
        message.success('取消封禁');
        getBanResourceList(props.data.uid)
      } else {
        message.error(res.data.msg);
      }
    },
    onNegativeClick: () => { }
  })
}

watch(() => props.visible, (newVal) => {
  if (newVal) {
    if (props.data) {
      getBanResourceList(props.data.uid);
    }
  }
})

watch(() => visibleBanModal.value, (newVal) => {
  if (!newVal && props.data) {
    getBanResourceList(props.data.uid);
  }
})
</script>

<style lang="scss" scoped>
.tag {
  margin-right: 10px;
}

.ban-box {
  width: 100%;
  padding-bottom: 30px;

  .ban-record-title {
    display: block;
    padding-bottom: 16px;
  }
}

.btn {
  width: 100px;
  margin-left: 10px;
}
</style>
