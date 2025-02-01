interface UserListParam {
  page: number;
  pageSize: number;
}

interface EditUserRoleParam {
  uid: number;
  code: string;
}

interface UserInfoType {
  uid: number;
  name: string;
  avatar: string;
  spaceCover?: string;
  email?: string;
  gender?: number;
  sign?: string;
  birthday?: string;
  createdAt?: string;
  role?: string;
}

interface UserFormType {
  uid: number;
  name: string;
  avatar: string;
  spaceCover: string;
  email: string;
  sign: string;
  role: string;
}