<template>
  <div class="upload-video">
    <p class="title">视频管理</p>
    <div class="tab-bar">
      <span
        v-for="tab in tabs"
        :key="tab.key"
        class="tab-item"
        :class="{ active: activeTab === tab.key }"
        @click="changeTab(tab.key)"
      >
        {{ tab.label }}
      </span>
    </div>
    <div class="video-box">
      <el-scrollbar ref="scrollbarRef" @scroll="onScroll">
        <ul v-if="videoList.length" class="video-list">
          <li class="video-item" v-for="(item, index) in videoList" :key="index">
            <div class="item-left">
              <div class="cover">
                <oss-image v-if="item.cover" :src="item.cover" alt="封面" />
              </div>
            </div>
            <div class="item-center">
              <template v-if="item.status !== reviewCode.AUDIT_APPROVED">
                <span class="item-title unlinked">{{ item.title }}</span>
              </template>
              <template v-else>
                <nuxt-link class="item-title" :to="`/watch?v=${item.shortId}`">{{ item.title }}</nuxt-link>
              </template>
              <span class="desc">简介：{{ item.desc }}</span>
              <div class="desc">
                <span>创建于：{{ formatTime(item.createdAt) }}</span>
                <span class="status" v-if="getStatusText(item.status)"
                  :style="`color: ${getStatusTextColor(item.status)}`">{{ getStatusText(item.status) }}</span>
                <span class="status status-btn" v-if="item.status === reviewCode.REVIEW_FAILED || item.status === reviewCode.PROCESSING_FAIL"
                  @click="showReason(item.vid)">查看原因</span>
              </div>
              <div class="progress-box" v-if="item.status === reviewCode.CREATED_VIDEO || item.status === reviewCode.VIDEO_PROCESSING || item.status === reviewCode.SUBMIT_REVIEW">
                <div class="progress-detail" v-if="(item.transcodingDetails || []).length">
                  <template v-for="group in groupDetails(item.transcodingDetails || [])" :key="`${item.vid}-${group.resourceId}`">
                    <div class="resource-group">
                      <div class="group-header" @click="toggleResourceGroup(item.vid, group.resourceId)">
                        <span class="group-arrow" :class="{ expanded: isResourceGroupExpanded(item.vid, group.resourceId) }">▶</span>
                        <span class="group-title">{{ group.title }}</span>
                        <span class="group-count">{{ group.items.length }} 个画质</span>
                        <span class="group-status" v-if="group.allWaiting">排队中</span>
                        <span class="group-status success" v-else-if="group.allDone">完成</span>
                        <span class="group-status warn" v-else-if="group.failCount > 0 && group.successCount > 0">部分完成 ({{ group.successCount }}/{{ group.items.length }})</span>
                        <span class="group-status error" v-else-if="group.failCount > 0 && group.successCount === 0">失败</span>
                        <span class="group-progress">
                          <el-progress
                            :percentage="group.allWaiting ? 0 : group.overallPct"
                            :stroke-width="14"
                            :show-text="true"
                            :status="group.allDone ? 'success' : (group.failCount > 0 && group.successCount === 0 ? 'exception' : undefined)" />
                        </span>
                      </div>
                      <div class="group-detail" v-if="isResourceGroupExpanded(item.vid, group.resourceId)">
                        <div class="detail-item" v-for="detail in group.items" :key="`${detail.resourceId}-${detail.quality}`">
                          <div class="detail-title">
                            <span class="detail-quality">{{ detail.quality }}</span>
                            <el-tag v-if="detail.status === 'waiting'" size="small" type="info">排队中</el-tag>
                            <el-tag v-else-if="detail.status === 'fail'" size="small" type="danger">失败</el-tag>
                            <span v-if="detail.status === 'fail'" class="retry-btn" @click.stop="retryQuality(group.resourceId, detail.quality)">重试</span>
                          </div>
                          <el-progress :percentage="detail.status === 'waiting' ? 0 : Number((detail.progress || 0).toFixed(1))" :stroke-width="8"
                            :show-text="true"
                            :status="detail.status === 'fail' ? 'exception' : (detail.status === 'success' ? 'success' : undefined)" />
                        </div>
                        <!-- 上传进度 -->
                        <div class="upload-section" v-if="group.upload">
                          <div class="detail-title"><span class="detail-quality">OSS 上传</span></div>
                          <el-progress
                            :percentage="Math.round((group.upload.progress || 0))"
                            :stroke-width="8"
                            :show-text="true"
                            :status="group.upload.status === 'fail' ? 'exception' : (group.upload.status === 'success' ? 'success' : undefined)" />
                        </div>
                      </div>
                    </div>
                  </template>
                </div>
              </div>
            </div>
            <div class="item-right">
              <el-dropdown>
                <el-button class="more-btn" plain>
                  <el-icon size="16">
                    <more-icon></more-icon>
                  </el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
<el-dropdown-item @click="modifyVideo(item.shortId || item.vid)">编辑</el-dropdown-item>
<el-dropdown-item @click="openSubtitleManage(item.shortId || item.vid)">字幕管理</el-dropdown-item>
                    <el-dropdown-item @click="deleteVideo(item, index)">删除稿件</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </li>
        </ul>
        <el-empty v-else-if="!initialLoading" description="暂无视频" />
      </el-scrollbar>
    </div>
    <client-only>
      <el-dialog v-model="deleteDialogVisible" class="delete-dialog" width="500" :before-close="beforeClose">
        <div class="delete-dialog-title">请输入 <strong>{{ deleteVideoInfo?.title }}</strong> 删除此视频</div>
        <div class="delete-dialog-desc">视频删除后将无法恢复，请谨慎操作</div>
        <el-input class="input" v-model="deleteVideoTitle" placeholder="请输入视频标题"></el-input>
        <el-button type="danger" class="delete-btn" @click="submitDelete">确认删除</el-button>
      </el-dialog>
    </client-only>
    <subtitle-manage-dialog :vid="subtitleManageVid" @close="subtitleManageVid = null" />
  </div>
</template>

<script setup lang="ts">
import { onBeforeMount, onBeforeUnmount, ref } from 'vue';
import { getUploadVideoAPI, deleteVideoAPI, reTranscodeResourceAPI } from '@/api/video';
import { MoreOne as MoreIcon } from '@icon-park/vue-next';
import { getVideoReviewRecordAPI } from '@/api/revies';
import { reviewCode } from '@/utils/review-code';
import { statusCode } from '@/utils/status-code';
import { formatTime } from '@/utils/format';
import { getResourceUrl } from '@/utils/resource';
import SubtitleManageDialog from '@/components/subtitle/SubtitleManageDialog.vue';

const subtitleManageVid = ref<string | number | null>(null);

const page = ref(1);
const total = ref(0);
const pageSize = 8;
const noMore = ref(false);
const loading = ref(false);
const initialLoading = ref(true);
const videoList = ref<Array<ManuscriptVideoType>>([]);
let silentRefreshTimer: number | null = null;
const activeTab = ref<'published' | 'pending' | 'rejected' | 'transcoding' | 'transcode_failed'>('published');
const expandedDetail = ref<Record<number, boolean>>({});
const expandedResourceGroups = ref<Set<string>>(new Set());

function toggleResourceGroup(vid: number, resourceId: number) {
  const key = `${vid}-${resourceId}`;
  const s = new Set(expandedResourceGroups.value);
  if (s.has(key)) s.delete(key); else s.add(key);
  expandedResourceGroups.value = s;
}

function isResourceGroupExpanded(vid: number, resourceId: number): boolean {
  return expandedResourceGroups.value.has(`${vid}-${resourceId}`);
}

interface ResourceGroup {
  resourceId: number;
  title: string;
  items: TranscodingProgressDetail[];
  overallPct: number;
  donePct: number;
  successCount: number;
  failCount: number;
  allWaiting: boolean;
  allDone: boolean;
  upload?: UploadProgressInfo;
}

function groupDetails(details: TranscodingProgressDetail[]): ResourceGroup[] {
  const map = new Map<number, TranscodingProgressDetail[]>();
  for (const d of details) {
    if (!map.has(d.resourceId)) map.set(d.resourceId, []);
    map.get(d.resourceId)!.push(d);
  }
  const groups: ResourceGroup[] = [];
  for (const [rid, items] of map) {
    const n = items.length;
    const successCount = items.filter(i => i.status === 'success').length;
    const failCount = items.filter(i => i.status === 'fail').length;
    const allWaiting = items.every(i => i.status === 'waiting');
    const allDone = !allWaiting && items.every(i => i.status === 'success' || i.status === 'fail');
    const donePct = n > 0 ? Math.round((successCount / n) * 100) : 0;
    // 综合进度：已完成的算 100%，正在编码的用实时 progress，排队的算 0%
    const totalProgPct = items.reduce((s, i) => {
      if (i.status === 'success') return s + 100;
      if (i.status === 'fail') return s + 0;
      return s + (i.progress || 0);
    }, 0);
    const overallPct = n > 0 ? Math.round(totalProgPct / n) : 0;
    const upload = items.find(i => (i as any).upload)?.upload;
    groups.push({
      resourceId: rid,
      title: items[0]?.resourceTitle || `分P${rid}`,
      items,
      overallPct, donePct, successCount, failCount,
      allWaiting, allDone,
      upload,
    });
  }
  return groups;
}
const tabs = [
  { key: 'published', label: '已发布' },
  { key: 'pending', label: '待审核' },
  { key: 'rejected', label: '审核失败' },
  { key: 'transcoding', label: '转码中' },
  { key: 'transcode_failed', label: '转码失败' },
] as const;

const getCurrentCategory = () => activeTab.value;
const getUploadVideo = async () => {
  if (loading.value || noMore.value) return;
  loading.value = true;
  const res = await getUploadVideoAPI(page.value, pageSize, getCurrentCategory());
  if (res.data.code === statusCode.OK) {
    total.value = res.data.data.total || 0;
    if (res.data.data.videos) {
      videoList.value = videoList.value.concat(res.data.data.videos);
      if (videoList.value.length >= total.value || res.data.data.videos.length < pageSize) {
        noMore.value = true;
      }
    } else {
      noMore.value = true;
    }
  }
  initialLoading.value = false;
  loading.value = false;
}

const scrollLoad = () => {
  if (!loading.value) {
    page.value++;
    getUploadVideo();
  }
}

const scrollbarRef = ref()
const onScroll = () => {
  const wrap = scrollbarRef.value?.wrapRef
  if (!wrap) return
  if (wrap.scrollHeight - wrap.scrollTop - wrap.clientHeight < 50) scrollLoad()
}

const silentRefreshUploadVideo = async () => {
  if (loading.value || activeTab.value !== 'transcoding') return;

  const loadedPages = Math.max(1, page.value);
  const mergedVideos: ManuscriptVideoType[] = [];
  let latestTotal = total.value;
  let lastPageSize = 0;

  try {
    for (let i = 1; i <= loadedPages; i++) {
      const res = await getUploadVideoAPI(i, pageSize, getCurrentCategory());
      if (res.data.code !== statusCode.OK) {
        return;
      }
      const currentVideos: ManuscriptVideoType[] = res.data.data.videos || [];
      if (i === 1) {
        latestTotal = res.data.data.total || 0;
      }
      lastPageSize = currentVideos.length;
      mergedVideos.push(...currentVideos);
      if (currentVideos.length < pageSize) {
        break;
      }
    }

    total.value = latestTotal;
    videoList.value = mergedVideos;
    noMore.value = mergedVideos.length >= latestTotal || lastPageSize < pageSize;

    const maxPage = Math.max(1, Math.ceil((mergedVideos.length || 1) / pageSize));
    if (page.value > maxPage) {
      page.value = maxPage;
    }
  } catch {
    // 静默刷新失败不打断当前页面
  }
}

const startSilentRefresh = () => {
  stopSilentRefresh();
  if (activeTab.value !== 'transcoding') return;
  silentRefreshTimer = window.setInterval(() => {
    silentRefreshUploadVideo();
  }, 3000);
}

const stopSilentRefresh = () => {
  if (silentRefreshTimer !== null) {
    window.clearInterval(silentRefreshTimer);
    silentRefreshTimer = null;
  }
}

const changeTab = (tab: typeof tabs[number]['key']) => {
  if (activeTab.value === tab) return;
  activeTab.value = tab;
  page.value = 1;
  total.value = 0;
  noMore.value = false;
  videoList.value = [];
  initialLoading.value = true;
  expandedDetail.value = {};
  expandedResourceGroups.value = new Set();
  if (tab === 'transcoding') {
    startSilentRefresh();
  } else {
    stopSilentRefresh();
  }
  getUploadVideo();
}

const toggleProgressDetail = (vid: number) => {
  expandedDetail.value[vid] = !expandedDetail.value[vid];
}

const retryQuality = async (resourceId: number, quality: string) => {
  const res = await reTranscodeResourceAPI(resourceId, quality);
  if (res.data.code === statusCode.OK) {
    ElMessage.success(`已触发重试: ${quality}`);
    silentRefreshUploadVideo();
  }
}

const deleteVideoIndex = ref(-1);
const deleteVideoTitle = ref("");
const deleteDialogVisible = ref(false);
const deleteVideoInfo = ref<ManuscriptVideoType>();
const deleteVideo = async (video: ManuscriptVideoType, index: number) => {
  deleteVideoInfo.value = video;
  deleteVideoIndex.value = index;
  deleteDialogVisible.value = true;
}

const beforeClose = () => {
  deleteVideoTitle.value = "";
  deleteVideoIndex.value = -1;
  deleteVideoInfo.value = undefined;
  deleteDialogVisible.value = false;
}

const submitDelete = async () => {
  if (deleteVideoTitle.value === deleteVideoInfo.value?.title) {
    const res = await deleteVideoAPI(deleteVideoInfo.value.vid);
    if (res.data.code === statusCode.OK) {
      videoList.value.splice(deleteVideoIndex.value, 1);
    }

    deleteVideoTitle.value = "";
    deleteVideoIndex.value = -1;
    deleteVideoInfo.value = undefined;
    deleteDialogVisible.value = false;
  } else {
    ElMessage.error("输入标题与原标题不一致");
  }
}

const getStatusText = (status: number) => {
  switch (status) {
    // case reviewCode.CREATED_VIDEO:
    //   return "未提交"
    case reviewCode.VIDEO_PROCESSING: // 200 转码中
    case reviewCode.SUBMIT_REVIEW: // 300 转码中
      return "转码中"
    case reviewCode.WAITING_REVIEW: // 500 转码完成，待审核
      return "待审核"
    case reviewCode.REVIEW_FAILED:
      return "审核不通过"
    case reviewCode.PROCESSING_FAIL:
      return "视频处理失败"
  }
}

const getStatusTextColor = (status: number) => {
  switch (status) {
    case reviewCode.CREATED_VIDEO:
      return "#999"
    case reviewCode.VIDEO_PROCESSING: // 200 转码中
    case reviewCode.SUBMIT_REVIEW: // 300 转码中
      return "#ff9800" // 橙色表示处理中
    case reviewCode.WAITING_REVIEW: // 500 待审核
      return "var(--primary-hover-color)" // 主题色表示待审核
    case reviewCode.REVIEW_FAILED:
      return "#f56c6c"
    case reviewCode.PROCESSING_FAIL:
      return "#f56c6c"
  }
}

const showReason = async (vid: number) => {
  const res = await getVideoReviewRecordAPI(vid);
  if (res.data.code === statusCode.OK) {
    const review = res.data.data.review;
    let message = review.remark || '审核未通过';

    // 如果有冲突信息，添加到消息中
    if (review.conflictResourceId) {
      message += `\n\n【内容冲突】`;
      if (review.conflictReason) {
        message += `\n原因：${review.conflictReason}`;
      }
      message += `\n冲突资源ID：${review.conflictResourceId}`;
    }

    ElMessageBox.alert(message, '审核不通过原因', {
      confirmButtonText: '确认',
      dangerouslyUseHTMLString: false,
    })
  }
}

//前往修改视频
const modifyVideo = (vid: number | string) => {
  navigateTo({ name: "upload-video", query: { vid: vid } });
}

const openSubtitleManage = (vid: number | string) => {
  subtitleManageVid.value = vid;
}

onBeforeMount(() => {
  getUploadVideo();
})

onBeforeUnmount(() => {
  stopSilentRefresh();
})
</script>

<style lang="scss" scoped>
.upload-video {
  padding: 0 18px 0;
  height: 100%;
  box-sizing: border-box;
  background-color: var(--bg-elev-1);

  .title {
    font-size: 18px;
    margin: 0;
    padding: 16px 0 10px;
  }

  .tab-bar {
    display: flex;
    gap: 8px;
    margin-bottom: 12px;

    .tab-item {
      padding: 6px 12px;
      border-radius: 12px;
      font-size: 13px;
      color: var(--font-primary-3);
      cursor: pointer;
      user-select: none;
    }

    .tab-item.active {
      color: #fff;
      background: var(--primary-hover-color);
    }
  }

  .video-box {
    height: calc(100% - 102px);
  }

  .video-list {
    list-style: none;
    box-sizing: border-box;
    width: 100%;
    margin: 0;
    padding: 0;

    .video-item {
      display: flex;
      padding: 16px 0;
      width: 100%;
      min-height: 80px;
      margin-bottom: 12px;
    border-bottom: 1px solid var(--border-color);
      padding-bottom: 12px;

      .item-left {
        width: 120px;
        height: 80px;
        margin-right: 10px;

        .cover {
          border-radius: 2px;
          width: 100%;
          height: 100%;
          background-color: var(--fill-1, #f1f2f3);

          img {
            width: 100%;
            height: 100%;
            border-radius: 2px;
            object-fit: contain;
          }
        }
      }

      .item-center {
        flex: 1;

        .item-title {
          font-size: 14px;
          color: var(--font-primary-1);
          line-height: 18px;
          margin: 0 0 26px;
          cursor: pointer;
          overflow: hidden;
          text-overflow: ellipsis;
          display: -webkit-box;
          -webkit-line-clamp: 1;
          line-clamp: 1;
          -webkit-box-orient: vertical;

          &:hover {
            color: var(--primary-hover-color);
          }

          &.unlinked {
            cursor: default;
            color: var(--font-primary-3);
            
            &:hover {
              color: var(--font-primary-3);
            }
          }
        }

        .desc {
          font-size: 12px;
          color: var(--font-primary-3);
          overflow: hidden;
          text-overflow: ellipsis;
          display: -webkit-box;
          -webkit-line-clamp: 1;
          line-clamp: 1;
          -webkit-box-orient: vertical;
        }

        .status {
          margin-left: 12px;
          color: var(--primary-hover-color);
        }

        .status-btn {
          cursor: pointer;
        }

            .progress-box {
          margin-top: 8px;

          .progress-detail {
            margin-top: 8px;

            .resource-group {
              margin-bottom: 8px;
              border: 1px solid var(--border-color, #e0e0e0);
              border-radius: 6px;
              overflow: hidden;
              background-color: var(--bg-elev-1);

              .group-header {
                display: flex;
                align-items: center;
                gap: 8px;
                padding: 6px 10px;
                cursor: pointer;
                user-select: none;
                font-size: 12px;
                background-color: var(--fill-1);

                .group-arrow {
                  font-size: 10px;
                  color: var(--font-primary-3);
                  width: 14px;
                  text-align: center;
                  transition: transform .2s;
                  flex-shrink: 0;

                  &.expanded {
                    transform: rotate(90deg);
                  }
                }

                .group-title {
                  font-weight: 500;
                  white-space: nowrap;
                  overflow: hidden;
                  text-overflow: ellipsis;
                  flex-shrink: 1;
                  min-width: 0;
                }

                .group-count {
                  color: var(--font-primary-3);
                  white-space: nowrap;
                  flex-shrink: 0;
                }

                .group-status {
                  font-size: 11px;
                  padding: 1px 6px;
                  border-radius: 3px;
                  background: var(--bg-elev-1);
                  color: var(--font-primary-3);
                  white-space: nowrap;
                  flex-shrink: 0;

                  &.error { color: #f56c6c; background: rgba(245, 108, 108, 0.12); }
                  &.success { color: #67c23a; background: rgba(103, 194, 58, 0.12); }
                  &.warn { color: #e6a23c; background: rgba(230, 162, 60, 0.12); }
                }

                .group-progress {
                  flex: 1;
                  min-width: 120px;
                  max-width: 260px;
                  margin-left: auto;

                  :deep(.el-progress-bar__outer) {
                    background-color: var(--fill-1);
                  }
                }
              }

              .group-detail {
                padding: 8px 12px 4px 26px;
                border-top: 1px solid var(--border-color, #e0e0e0);

                .detail-item {
                  margin-bottom: 8px;
                }

                .detail-title {
                  font-size: 12px;
                  color: var(--font-primary-2);
                  margin-bottom: 4px;
                  display: flex;
                  align-items: center;
                  gap: 6px;

                  .detail-quality {
                    color: var(--font-primary-3);
                  }

                  .retry-btn {
                    color: var(--primary-hover-color);
                    cursor: pointer;
                    font-size: 11px;
                    padding: 0 4px;
                    border-radius: 3px;
                    &:hover {
                      background: var(--primary-color-active);
                    }
                  }
                }

                .upload-section {
                  margin-top: 4px;
                  padding-top: 8px;
                  border-top: 1px dashed var(--border-color, #e0e0e0);
                }
              }
            }
          }
        }
      }

      .item-right {
        width: 90px;
        height: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
      }
    }
  }
}

.delete-dialog {

  .delete-dialog-title {
    font-size: 16px;
  color: var(--font-primary-1);
    text-align: center;
    margin: 20px 0;
  }

  .delete-dialog-desc {
  color: var(--font-primary-3);
    font-size: 13px;
    text-align: center;

  }

  .input {
    margin: 20px 0;

  }

  .delete-btn {
    width: 100%;
    color: #d03050;
    border: none;
    font-family: inherit;
    background-color: rgba(208, 48, 80, 0.16);

    &:hover {
      background-color: rgba(208, 48, 80, 0.22);
    }
  }
}
</style>