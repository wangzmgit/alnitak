interface CollectionType {
  id: number;
  name: string;
  cover?: string;
  desc?: string;
  open?: boolean;
  createdAt?: string;

  // 本地用的数据
  checked: boolean;
}

interface CollectionInfoType {
  id: number,
  name: string,
  cover: string,
  desc: string,
  open: boolean,
  createdAt: string;
  author: UserInfoType;
}



interface EditCollectionType {
  id: number,
  name?: string,
  cover?: string,
  desc?: string,
  open?: boolean,
}
