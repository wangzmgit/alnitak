interface CollectionType {
  id: number,
  name?: string,
  cover?: string,
  desc?: string,
  open?: boolean,
  created_at?: string

  // 本地用的数据
  checked: boolean
}
