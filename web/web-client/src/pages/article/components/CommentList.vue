<template>
  <div class="comment-navigation" id="comment-title">
    <ul class="nav-bar">
      <li class="nav-title">
        <span class="nav-title-text">评论</span>
        <span class="total-comment">{{ commentCount }}</span>
      </li>
    </ul>
  </div>
  <!-- 评论输入框 -->
  <div class="comment-box">
    <common-avatar v-if="isLoggedIn" class="avatar" :url="userInfo?.avatar" :size="40"></common-avatar>
    <div v-else class="avatar">
      <div class="login-btn">登录</div>
    </div>
    <el-input class="comment-input" v-model="commentContent" resize="none" :rows="3" type="textarea"
      placeholder="善语结善缘，恶言伤人心" />
    <button class="comment-submit" :class="isLoggedIn ? '' : 'submit-disabled'" @click="submitComment">发表评论</button>
  </div>
  <!-- 评论列表 -->
  <div class="comment-container" v-for="(item, i) in commentList">
    <div class="comment-avatar">
      <common-avatar :url="item.author.avatar" :size="40"></common-avatar>
    </div>
    <div class="content-warp">
      <div class="user-info">
        <div class="user-name">{{ item.author.name }}</div>
      </div>
      <div class="comment-content-container">
        <span class="comment-content" v-for="content in handleMention(item.content, item.atUserIds, item.atUsernames)">
          <span v-if="!content.key">{{ content.value }}</span>
          <nuxt-link v-else class="at" :to="`/user/${content.key}`">{{ content.value }}</nuxt-link>
        </span>
      </div>
      <div class="comment-info">
        <span class="comment-time">{{ formatRelativeTime(item.createdAt) }}</span>
        <span class="reply-btn" :class="isLoggedIn ? '' : 'btn-disabled'" @click="showReplyBox(item)">回复</span>
        <client-only v-if="isLoggedIn">
          <el-popconfirm title="是否删除该条评论？" confirm-button-text="确认" cancel-button-text="取消"
            @confirm="deleteComment(i, item)">
            <template #reference>
              <span class="del-btn" v-if="item.uid === userInfo?.uid">删除</span>
            </template>
          </el-popconfirm>
        </client-only>
      </div>
    </div>
    <!-- 回复 -->
    <div class="reply-list" v-if="item.reply">
      <div class="reply-item" v-for="(reply, j) in item.reply">
        <div class="reply-user-info">
          <div class="reply-avatar">
            <common-avatar :url="reply.author.avatar" :size="24"></common-avatar>
          </div>
          <div class="reply-user-name">{{ reply.author.name }}</div>
        </div>
        <span class="reply-content-container">
          <span v-if="reply.replyUserName">
            回复<nuxt-link class="at" :to="`/user/${reply.replyUserId}`">@{{ reply.replyUserName }}</nuxt-link> :
          </span>
          <span class="reply-content" v-for="content in handleMention(reply.content, reply.atUserIds, reply.atUsernames)">
            <span v-if="!content.key">{{ content.value }}</span>
            <nuxt-link v-else class="at" :to="`/user/${content.key}`">{{ content.value }}</nuxt-link>
          </span>
        </span>
        <div class="reply-info">
          <span class="reply-time">{{ formatRelativeTime(reply.createdAt) }}</span>
          <span class="reply-btn" :class="isLoggedIn ? '' : 'btn-disabled'" @click="showReplyBox(item, reply)">回复</span>
          <client-only v-if="isLoggedIn">
            <el-popconfirm title="是否删除该条回复？" confirm-button-text="确认" cancel-button-text="取消"
              @confirm="deleteComment(j, item, reply)">
              <template #reference>
                <span class="del-btn" v-if="reply.uid === userInfo?.uid">删除</span>
              </template>
            </el-popconfirm>
          </client-only>
        </div>
      </div>
    </div>
    <!-- 加载更多 -->
    <div v-if="item.replyCount > 3 && !item.hiddenMoreBtn" class="view-more">
      <span>共{{ item.replyCount }}条回复，</span>
      <span class="view-more-btn" @click="replyPageChange(item, 1)">点击查看</span>
    </div>
    <div v-else-if="item.hiddenMoreBtn && item.replyCount > replyPageSize">
      <el-pagination v-model:current-page="item.page" :page-size="replyPageSize" small layout="prev, pager, next"
        :total="item.replyCount" @current-change="page => replyPageChange(item, page)" />
    </div>
    <!-- 回复输入框 -->
    <div v-show="item.showReplyBox" :id="`reply-box-${item.id}`" class="comment-box reply-box">
      <common-avatar class="avatar" :url="userInfo?.avatar" :size="32"></common-avatar>
      <el-input class="comment-input" v-model="commentForm.content" resize="none" :rows="3" type="textarea"
        :placeholder="replyTip" />
      <button class="comment-submit" @click="submitReply(item)">回复</button>
    </div>
    <div class="bottom-line"></div>
  </div>
</template>

<script setup lang="ts">
import { onBeforeMount, onBeforeUnmount, reactive, ref } from "vue";
import { handleMention } from '@/utils/mention';
import { statusCode } from '@/utils/status-code';
import { ElMessage } from "element-plus";
import { formatRelativeTime } from "@/utils/format";
import { asyncGetUserBaseInfoAPI, getUserInfoAPI } from "@/api/user";
import CommonAvatar from "@/components/common-avatar/index.vue";
import { addArticleCommentAPI, getArticleCommentAPI, getArticleReplyAPI, deleteArticleCommentAPI } from "@/api/comment";
import { scrollToViewCenter } from "@/utils/scroll";

const emit = defineEmits(["updateCount"])
const props = defineProps<{
  aid: number
}>();

const getUserInfo = async () => {
  const res = await getUserInfoAPI();
  if (res.data.code === statusCode.OK) {
    userInfo.value = res.data.data.userInfo;
    isLoggedIn.value = true;
  }
}

const isLoggedIn = ref(false);
const userId = useCookie('user_id');
const userInfo = ref<UserInfoType>();
const { data } = await asyncGetUserBaseInfoAPI(userId.value!);
if ((data.value as any).code === statusCode.OK) {
  userInfo.value = (data.value as any).data.userInfo;
  isLoggedIn.value = true;
}

if (!process.server) {
  if (!userInfo.value) {
    getUserInfo();
  }
}

const pagination = reactive({
  page: 1,
  pageSize: 10,
  noMore: false,
  loading: false,
})

const commentCount = ref(0);
const commentList = ref<CommentType[]>([]);
const getCommentList = async () => {
  pagination.loading = true;
  const res = await getArticleCommentAPI(props.aid, pagination.page, pagination.pageSize);
  if (res.data.code === statusCode.OK) {
    commentCount.value = res.data.data.total;
    commentList.value = commentList.value.concat(res.data.data.comments);
    if (res.data.data.comments < pagination.pageSize) {
      pagination.noMore = true;
    }
    emit("updateCount", res.data.data.total);
  }
  pagination.loading = false;
}

const replyPageSize = 10;
const getReplyList = async (comment: CommentType) => {
  const res = await getArticleReplyAPI(comment.id, comment.page || 1, replyPageSize);
  if (res.data.code === statusCode.OK) {
    comment.reply = res.data.data.replies;
  }
}

// 加载更多评论
const replyPageChange = (comment: CommentType, page: number) => {
  comment.page = page;
  comment.hiddenMoreBtn = true;
  getReplyList(comment)
}

const commentContent = ref("");
const commentForm = reactive<AddCommentType>({
  cid: props.aid,
  content: "",
  parentId: 0,
  replyUserId: 0,
  replyUserName: "",
  at: [],
  replyContent: ""
})

const submitComment = async () => {
  if (!isLoggedIn.value) return;

  commentForm.parentId = 0;
  commentForm.content = commentContent.value;
  const comment = await addComment();
  if (comment) {
    // 关闭所有回复框
    commentList.value.forEach((item) => {
      item.showReplyBox = false;
    })

    // 手动添加评论
    if (pagination.noMore) {
      commentList.value.push(comment);
    }
  }
}

const submitReply = async (comment: CommentType) => {
  commentForm.parentId = comment.id;

  const reply = await addComment();
  if (reply) {
    // 关闭所有回复框
    commentList.value.forEach((item) => {
      item.showReplyBox = false;
    })

    comment.replyCount++;
    if (!comment.reply) comment.reply = [];
    if (comment.reply.length < 3 || (comment.hiddenMoreBtn && comment.reply.length < replyPageSize)) {
      comment.reply.push(reply);
    }
  }
}

const replyTip = ref('善语结善缘，恶言伤人心');
const showReplyBox = async (comment: CommentType, reply?: ReplyType) => {
  // 关闭所有回复框
  commentList.value.forEach((item) => {
    item.showReplyBox = false;
  })
  // 判断是否登录
  if (!isLoggedIn.value) return;
  // 初始化数据
  comment.showReplyBox = true;
  commentForm.content = '';
  commentForm.replyContent = "";
  commentForm.parentId = comment.id;

  if (reply) {
    commentForm.replyUserId = reply.author.uid;
    commentForm.replyUserName = reply.author.name;
    commentForm.replyContent = reply.content;
    replyTip.value = `回复 @${reply.author.name}: `
  } else {
    commentForm.replyContent = comment.content;
    replyTip.value = `回复 @${comment.author.name}: `
  }

  // 滚动到指定位置
  nextTick(() => {
    const el = document.getElementById(`reply-box-${comment.id}`);
    if (el) scrollToViewCenter(el);
  })
}

const addComment = async () => {
  // 处理@
  const regex = /@(\S+)\s/g;
  const matches = commentForm.content.match(regex);
  if (matches) {
    commentForm.at = matches.map(match => match.trim().substring(1));
  }

  const res = await addArticleCommentAPI(commentForm);
  if (res.data.code === statusCode.OK) {
    const comment = res.data.data.comment;
    comment.author = userInfo.value;
    commentForm.content = "";
    commentContent.value = "";
    return comment;
  } else {
    ElMessage.error('发送失败');
    return 0;
  }
}

const deleteComment = async (index: number, comment: CommentType, reply?: ReplyType) => {
  const deleteId = reply ? reply.id : comment.id;
  const res = await deleteArticleCommentAPI(deleteId);
  if (res.data.code === statusCode.OK) {
    if (reply) {
      comment.replyCount--;
      comment.reply.splice(index, 1);
    } else {
      commentList.value.splice(index, 1);
    }
  }
}

// 加载更多评论
const lazyLoading = () => {
  const scrollTop = document.documentElement.scrollTop || document.body.scrollTop;
  const clientHeight = document.documentElement.clientHeight;
  const scrollHeight = document.documentElement.scrollHeight;
  if (scrollTop + clientHeight >= (scrollHeight - 10)) {
    if (!pagination.noMore && !pagination.loading) {
      pagination.page++;
      getCommentList();
    }
  }
}

onBeforeMount(() => {
  getCommentList();
  window.addEventListener("scroll", lazyLoading, true);
})

onBeforeUnmount(() => {
  window.removeEventListener("scroll", lazyLoading, true);
})
</script>

<style lang="scss" scoped>
.comment-navigation {
  margin-bottom: 22px;

  .nav-bar {
    display: flex;
    align-items: center;
    list-style: none;
    margin: 0;
    padding: 0;

    .nav-title {
      font-size: 20px;
      display: flex;
      align-items: center;

      .nav-title-text {
        color: #18191c;
        font-weight: 500;
      }

      .total-comment {
        font-size: 13px;
        margin: 0 36px 0 6px;
        font-weight: 400;
        color: #9499a0;
      }
    }
  }
}

.comment-box {
  display: flex;
  align-items: center;
  justify-content: space-between;

  .avatar {
    margin: 0 20px;
  }

  .login-btn {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    color: var(--primary-text-color);
    text-align: center;
    line-height: 40px;
    font-size: 14px;
    font-weight: 500;
    background-color: var(--primary-hover-color);
  }

  .comment-input {
    flex: 1;
    margin-right: 6px;
  }

  .comment-submit {
    width: 70px;
    height: 70px;
    padding: 4px 15px;
    font-size: 14px;
    color: #fff;
    border-radius: 4px;
    text-align: center;
    min-width: 60px;
    vertical-align: top;
    cursor: pointer;
    background-color: var(--primary-color);
    border: 1px solid var(--primary-color);
    transition: 0.1s;
    user-select: none;
    outline: none;

    &:hover {
      background-color: var(--primary-hover-color);
    }
  }
}

.comment-container {
  position: relative;
  padding: 16px 0 0 80px;

  .comment-avatar {
    display: flex;
    justify-content: center;
    position: absolute;
    left: 0;
    width: 80px;
    cursor: pointer;
  }

  .content-warp {
    position: relative;

    .user-info {
      display: flex;
      align-items: center;
      margin-bottom: 4px;
      font-size: 14px;
      height: 30px;

      .user-name {
        font-weight: 500;
        margin-right: 5px;
        line-height: 20px;
        color: #61666d;
        cursor: pointer;
      }
    }

    .comment-content-container {
      display: block;
      overflow: hidden;
      width: 100%;
      font-size: 15px;
      line-height: 24px;
      position: relative;
      padding: 2px 0;
      font-weight: 400;

      .comment-content {
        color: #18191c;
        overflow: hidden;
        word-wrap: break-word;
        word-break: break-word;
        white-space: pre-wrap;
        font-size: 15px;
        line-height: 24px;
        vertical-align: baseline;
      }
    }

    .comment-info {
      display: flex;
      align-items: center;
      position: relative;
      height: 24px;
      margin-top: 2px;
      font-size: 13px;
      color: #9499a0;
      font-weight: 400;

      .comment-time {
        margin-right: 20px;
      }

      .reply-btn {
        cursor: pointer;

        &:hover {
          color: var(--primary-hover-color);
        }
      }
    }
  }
}

.reply-list {
  .reply-item {
    position: relative;
    padding: 8px 0 8px 36px;
    border-radius: 4px;
    font-size: 15px;
    line-height: 24px;

    .reply-user-info {
      display: inline-flex;
      align-items: center;
      margin-right: 9px;
      line-height: 24px;
      vertical-align: baseline;
      white-space: nowrap;

      .reply-avatar {
        position: absolute;
        left: 2px;
        cursor: pointer;
      }

      .reply-user-name {
        font-size: 13px;
        line-height: 24px;
        font-weight: 500;
        margin-right: 5px;
        color: #61666d;
      }
    }

    .reply-content-container {
      .reply-content {
        color: #18191c;
        overflow: hidden;
        word-wrap: break-word;
        word-break: break-word;
        white-space: pre-wrap;
        line-height: 24px;
        vertical-align: baseline;
      }
    }

    .reply-info {
      display: flex;
      align-items: center;
      position: relative;
      height: 24px;
      margin-top: 2px;
      font-size: 13px;
      color: #9499a0;
      font-weight: 400;

      .reply-time {
        margin-right: 20px;
      }
    }
  }
}

.view-more {
  font-size: 13px;
  color: #9499a0;

  .view-more-btn {
    cursor: pointer;

    &:hover {
      color: var(--primary-hover-color);
    }
  }
}

.reply-box {
  margin: 25px 0 10px 0;

  .avatar {
    margin: 0 12px 0 0;
  }
}

.bottom-line {
  margin-top: 14px;
  border-bottom: 1px solid #e3e5e7;
}

.at {
  cursor: pointer;
  padding: 0 2px;
  color: var(--primary-color);
  text-decoration: none;

  &:hover {
    color: var(--primary-hover-color);
  }
}

.reply-btn {
  cursor: pointer;

  &:hover {
    color: var(--primary-hover-color);
  }
}

.del-btn {
  cursor: pointer;
  margin-left: 16px;

  &:hover {
    color: var(--primary-hover-color);
  }
}

.btn-disabled {
  color: #9499a0 !important;
  cursor: not-allowed !important;

  &:hover {
    color: #9499a0 !important;
  }
}

.submit-disabled {
  background-color: var(--primary-hover-color) !important;
  cursor: not-allowed !important;

  &:hover {
    background-color: var(--primary-hover-color) !important;
  }
}
</style>