<template>
  <div class="comment-navigation" id="comment-title">
    <ul class="nav-bar">
      <li class="nav-title">
        <span class="nav-title-text">评论</span>
        <span class="total-comment">{{ commentCount }}</span>
      </li>
      <span class="comment-btn" :class="isLoggedIn ? '' : 'btn-disabled'"
        @click="showCommentBox = !showCommentBox">{{ showCommentBox?'收起':'发表评论' }}</span>
    </ul>
  </div>
  <!-- 评论输入框 -->
  <div v-show="showCommentBox" class="comment-box">
    <n-input class="comment-input" v-model="commentContent" resize="none" :rows="3" type="textarea"
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
          <span v-else class="at" :to="`/user/${content.key}`">{{ content.value }}</span>
        </span>
      </div>
      <div class="comment-info">
        <span class="comment-time">{{ formatRelativeTime(item.createdAt) }}</span>
        <span class="reply-btn" :class="isLoggedIn ? '' : 'btn-disabled'" @click="showReplyBox(item)">回复</span>
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
            回复<span class="at" :to="`/user/${reply.replyUserId}`">@{{ reply.replyUserName }}</span> :
          </span>
          <span class="reply-content" v-for="content in handleMention(reply.content, reply.atUserIds, reply.atUsernames)">
            <span v-if="!content.key">{{ content.value }}</span>
            <span v-else class="at" :to="`/user/${content.key}`">{{ content.value }}</span>
          </span>
        </span>
        <div class="reply-info">
          <span class="reply-time">{{ formatRelativeTime(reply.createdAt) }}</span>
          <span class="reply-btn" :class="isLoggedIn ? '' : 'btn-disabled'" @click="showReplyBox(item, reply)">回复</span>
        </div>
      </div>
    </div>
    <!-- 加载更多 -->
    <div v-if="item.replyCount > 3" class="view-more">
      <span v-if="!item.page || item.page < 1">共{{ item.replyCount }}条回复，</span>
      <span v-if="!item.noMore" class="view-more-btn" @click="replyPageChange(item, 1)">
        {{ (!item.page || item.page < 1) ? '点击查看' : '加载更多' }} </span>
          <span v-if="item.page" class="view-more-btn retract" @click="retract(item)">收起</span>
    </div>
    <!-- 回复输入框 -->
    <div v-show="item.showReplyBox" :id="`reply-box-${item.id}`" class="comment-box reply-box">
      <n-input class="comment-input" v-model:value="commentForm.content" resize="none" :rows="3" type="textarea"
        :placeholder="replyTip" />
      <button class="comment-submit" @click="submitReply(item)">回复</button>
    </div>
    <div class="bottom-line"></div>
  </div>
</template>

<script setup lang="ts">
import { nextTick, onBeforeMount, onBeforeUnmount, reactive, ref } from "vue";
import { handleMention } from '@/utils/mention';
import { statusCode } from '@/utils/status-code';
import { formatRelativeTime } from "@/utils/format";
import CommonAvatar from "@/components/common-avatar/index.vue";
import { addVideoCommentAPI, getVideoCommentAPI, getVideoReplyAPI, deleteVideoCommentAPI } from "@/api/comment";
import { scrollToViewCenter } from "@/utils/scroll";
import { NInput, useMessage } from "naive-ui";
import { getUserInfoAPI } from "@/api/user";

const emit = defineEmits(["updateCount"])
const props = defineProps<{
  vid: number
}>();

const message = useMessage();

const isLoggedIn = ref(false);
const showCommentBox = ref(false);
const userInfo = ref<UserInfoType>();
const getUserBaseInfo = async () => {
  const res = await getUserInfoAPI();
  if (res.data.code === statusCode.OK) {
    isLoggedIn.value = true;
    userInfo.value = res.data.data.userInfo;
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
  const res = await getVideoCommentAPI(props.vid, pagination.page, pagination.pageSize);
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
  const res = await getVideoReplyAPI(comment.id, comment.page || 1, replyPageSize);
  if (res.data.code === statusCode.OK) {
    if (res.data.data.replies) {
      comment.reply = comment.reply.concat(res.data.data.replies);
      if (res.data.data.replies < replyPageSize) {
        comment.noMore = true;
      }
    } else {
      comment.noMore = true;
    }
  }
}

// 加载更多评论
const replyPageChange = (comment: CommentType, page: number) => {
  if (!comment.page) {
    comment.page = 1;
    comment.reply = []
  } else {
    comment.page++;
  }

  getReplyList(comment)
}

const retract = (comment: CommentType) => {
  comment.page = 0;
  comment.noMore = false;
  comment.reply = comment.reply.slice(0, 3);
}

const commentContent = ref("");
const commentForm = reactive<AddCommentType>({
  cid: props.vid,
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
  showCommentBox.value = false;
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

    if (comment.reply.length < 3 || comment.reply.length <= comment.replyCount) {
      comment.reply.push(reply);
    }
  }
}

const replyTip = ref('善语结善缘，恶言伤人心');
const showReplyBox = async (comment: CommentType, reply?: ReplyType) => {
  // 关闭所有回复框
  showCommentBox.value = false;
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

  const res = await addVideoCommentAPI(commentForm);
  if (res.data.code === statusCode.OK) {
    const comment = res.data.data.comment;
    comment.author = userInfo.value;
    commentForm.content = "";
    commentContent.value = "";
    return comment;
  } else {
    message.error('发送失败');
    return 0;
  }
}

const deleteComment = async (index: number, comment: CommentType, reply?: ReplyType) => {
  const deleteId = reply ? reply.id : comment.id;
  const res = await deleteVideoCommentAPI(deleteId);
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
  getUserBaseInfo();
  getCommentList();
  window.addEventListener("scroll", lazyLoading, true);
})

onBeforeUnmount(() => {
  window.removeEventListener("scroll", lazyLoading, true);
})
</script>

<style lang="scss" scoped>
.comment-navigation {
  margin-bottom: 3.2vmin;
  padding: 0 3.2vmin;

  .nav-bar {
    display: flex;
    align-items: center;
    justify-content: space-between;
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

    .comment-btn {
      cursor: pointer;
      color: #9499a0;
    }
  }
}

.comment-box {
  padding: 0 3.2vmin 2vmin;
  display: flex;
  align-items: center;
  justify-content: space-between;

  .comment-input {
    flex: 1;
    margin-right: 6px;
  }

  .comment-submit {
    width: 80px;
    height: 80px;
    padding: 4px 20px;
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
  }
}

.comment-container {
  position: relative;
  // padding: 16px 0 0 80px;
  padding: 0 3.2vmin 0 calc(6.4vmin + 40px);


  .comment-avatar {
    display: flex;
    justify-content: center;
    position: absolute;
    left: 3.2vmin;
    width: 40px;
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

  }

  .retract {
    padding-left: 6px;
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
}

.reply-btn {
  cursor: pointer;
}

.del-btn {
  cursor: pointer;
  margin-left: 16px;
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