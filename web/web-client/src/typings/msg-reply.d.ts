interface ReplyMessageType {
  type: number;
  article: ArticleType;
  video: BaseVideoType;
  user: UserInfoType;
  createdAt: string;
  content: string;
  targetReplyContent: string;
  rootContent: string;
  commentId: string;
}