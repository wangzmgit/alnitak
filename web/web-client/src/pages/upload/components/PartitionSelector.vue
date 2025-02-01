<template>
  <div class="select-partition">
    <el-select v-model="partition" class="select" placeholder="选择分区" @change="partitionChange">
      <el-option v-for="item in partitions" :label="item.name" :value="item.id"></el-option>
    </el-select>
    <el-select v-if="showSubpartition" v-model="subpartition" class="select" placeholder="选择二级分区"
      @change="selectedPartition">
      <el-option v-for="item in subpartitions" :label="item.name" :value="item.id"></el-option>
    </el-select>
  </div>
</template>

<script setup lang="ts">
import { onBeforeMount, ref } from "vue";
import { getPartitionAPI } from '@/api/partition';
import { statusCode } from "@/utils/status-code";
import { ElSelect, ElOption } from "element-plus";

const emits = defineEmits(["selected"]);
const props = defineProps<{
  partitions: PartitionType[]
}>()

const partition = ref();
const subpartition = ref();
const partitionList = ref<Array<PartitionType>>(props.partitions);//所有分区
const partitions = ref<Array<PartitionType>>([]);//分区
const subpartitions = ref<Array<PartitionType>>([]);//二级分区
const showSubpartition = ref(false);//是否显示二级分区

//改变分区
const partitionChange = (parentId: number) => {
  subpartitions.value = partitionList.value.filter((item) => {
    return item.parentId === parentId;
  })
  showSubpartition.value = true;
}

const selectedPartition = (value: number) => {
  emits("selected", value);
}

onBeforeMount(async () => {
  partitions.value = partitionList.value.filter((item) => {
    return item.parentId === 0;
  })
})
</script>

<style lang="scss" scoped>
.select-partition {
  width: 330px;
  display: flex;
  justify-content: space-between;

  .select {
    width: 150px;
  }
}
</style>