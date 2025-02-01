interface PartitionType {
  id: number;
  name: string;
  parentId: number;
  type: number;
  children?: PartitionType[];
}

interface AddPartitionType {
  name: string;
  parentId: number;
  type: number;
}