interface AddCommentType {
  cid: number;
  content: string;
  parentId: number;
  replyUserId?: number;
  replyUserName?: string;
  at: string[];
  replyContent: string;
}


interface CommentType {
  id: number;
  uid: number;
  atUserIds: string;
  atUsernames: string;
  author: UserInfoType;
  content: string;
  createdAt: string;
  parentId: number;
  reply: ReplyType[];
  replyCount: number;

  // 本地使用的数据，接口不返回
  page?: number;
  showReplyBox?: boolean;
  hiddenMoreBtn?: boolean;
}

interface ReplyType {
  id: number;
  uid: number;
  atUserIds: string;
  atUsernames: string;
  author: UserInfoType;
  content: string;
  createdAt: string;
  parentId: number;
  replyUserId: number;
  replyUserName: string;
}

interface CommentManageType {
  video: BaseVideoType;
  author: UserInfoType;
  target: UserInfoType;
  createdAt: string;
  content: string;
  targetReplyContent: string;
  rootContent: string;
  commentId: string;
}