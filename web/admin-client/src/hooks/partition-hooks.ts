import { ref } from 'vue';
import { getPartitionAPI } from '@/api/partition';
import { statusCode } from '@/utils/status-code';

export default function usePartition(partitionType: "video" | "article") {
  // 获取分区列表
  const partitionList = ref<Array<PartitionType>>([]);//所有分区
  const getPartition = async () => {
    const res = await getPartitionAPI(partitionType);
    if (res.data.code === statusCode.OK) {
      partitionList.value = res.data.data.partitions;
    }
  }

  // 获取分区名
  const getPartitionName = (id: number) => {
    const subpartition = partitionList.value.find((item) => {
      return item.id === id;
    })

    const partition = partitionList.value.find((item) => {
      return item.id === subpartition?.parentId;
    })

    return `${partition?.name}/${subpartition?.name}`;
  }


  return {
    partitionList,
    getPartition,
    getPartitionName
  };
}
