<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import SparkMD5 from 'spark-md5'

import {
  addPgcEpisodeApi,
  deletePgcEpisodeApi,
  deletePgcApi,
  getPgcDetailWithEpisodesApi,
  updatePgcEpisodeApi,
  updatePgcEpisodeStatusApi,
  updatePgcStatusApi,
  updatePgcApi,
  type AddPgcEpisodeReq,
  type PgcType,
  type UpdatePgcReq,
} from '@/api/pgc'
import {
  uploadImageApi,
  checkVideoUploadApi,
  uploadVideoChunkApi,
  mergeVideoUploadApi,
  createVideoUploadApi,
} from '@/api/upload'
import { useAuthStore } from '@/stores/auth'
import { useThemeStore } from '@/stores/theme'
import { toast } from '@/utils/toast'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const theme = useThemeStore()

const pgcId = computed(() => String(route.params.pgcId ?? ''))

type CategoryValue = 'cn' | 'jp' | 'documentary' | 'movie' | 'tv'
const categories: { value: CategoryValue; label: string; pgcType: PgcType }[] = [
  { value: 'cn', label: '国创(CN)', pgcType: 1 },
  { value: 'jp', label: '日创(JP)', pgcType: 2 },
  { value: 'documentary', label: '纪录片', pgcType: 3 },
  { value: 'movie', label: '电影', pgcType: 4 },
  { value: 'tv', label: '电视剧', pgcType: 5 },
]

function pick<T = any>(obj: any, keys: string[], fallback?: T): T {
  for (const k of keys) {
    if (obj && obj[k] !== undefined && obj[k] !== null) return obj[k] as T
  }
  return fallback as T
}

function categoryFromPgcType(t: number): CategoryValue {
  if (t === 2) return 'jp'
  if (t === 3) return 'documentary'
  if (t === 4) return 'movie'
  if (t === 5) return 'tv'
  return 'cn'
}

const loading = ref(false)
const loadError = ref('')

const saving = ref(false)
const saveError = ref('')
const saveOkMsg = ref('')
const opBusy = ref(false)
const pgcStatus = ref<number>(0)

const form = ref({
  projectTitle: '',
  category: 'cn' as CategoryValue,
  coverUrl: '',
  releaseYear: String(new Date().getFullYear()),
  area: '',
  rating: 0,
  synopsis: '',
  isOngoing: false,
  totalEpisodes: 1,
})
const currentEpisodes = ref(0)

type EpisodeRow = {
  id: number
  episode_number: number
  title: string
  vid: number
  duration: number
  publish_time: string
  status: number
}

const episodes = ref<EpisodeRow[]>([])

const coverUploading = ref(false)
const coverUploadError = ref('')
const coverFileInput = ref<HTMLInputElement | null>(null)

async function onPickCoverFile(e: Event) {
  coverUploadError.value = ''
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  if (!auth.token) {
    coverUploadError.value = '未登录，无法上传'
    return
  }
  coverUploading.value = true
  try {
    const res = await uploadImageApi(auth.token, file)
    if (res.code !== 200) {
      coverUploadError.value = res.msg || '上传失败'
      return
    }
    const url = (res.data as any)?.url as string | undefined
    if (!url) {
      coverUploadError.value = '上传成功但未返回 url'
      return
    }
    form.value.coverUrl = url
  } finally {
    coverUploading.value = false
    if (coverFileInput.value) coverFileInput.value.value = ''
  }
}

async function loadDetail() {
  loading.value = true
  loadError.value = ''
  try {
    const res = await getPgcDetailWithEpisodesApi(pgcId.value)
    if (res.code !== 200) {
      loadError.value = res.msg || '加载失败'
      return
    }
    const pgc = (res.data as any)?.pgc ?? {}
    const eps = ((res.data as any)?.episodes ?? []) as any[]

    form.value.projectTitle = String(pick(pgc, ['title', 'Title'], ''))
    form.value.coverUrl = String(pick(pgc, ['cover', 'Cover'], ''))
    form.value.synopsis = String(pick(pgc, ['desc', 'Desc'], ''))
    form.value.releaseYear = String(pick(pgc, ['year', 'Year'], new Date().getFullYear()))
    form.value.area = String(pick(pgc, ['area', 'Area'], ''))
    form.value.rating = Number(pick(pgc, ['rating', 'Rating'], 0)) || 0
    form.value.isOngoing = Boolean(pick(pgc, ['is_ongoing', 'IsOngoing'], false))
    form.value.totalEpisodes = Math.max(1, Number(pick(pgc, ['total_episodes', 'TotalEpisodes'], 1)) || 1)
    currentEpisodes.value = Number(pick(pgc, ['current_episodes', 'CurrentEpisodes'], 0)) || 0
    form.value.category = categoryFromPgcType(Number(pick(pgc, ['pgc_type', 'PGCType'], 1)))
    pgcStatus.value = Number(pick(pgc, ['status', 'Status'], 0)) || 0

    episodes.value = eps
      .map((r) => {
        return {
          id: Number(pick(r, ['id', 'ID'], 0)) || 0,
          episode_number: Number(pick(r, ['episode_number', 'EpisodeNumber'], 0)) || 0,
          title: String(pick(r, ['title', 'Title'], '')),
          vid: Number(pick(r, ['vid', 'VID'], 0)) || 0,
          duration: Number(pick(r, ['duration', 'Duration'], 0)) || 0,
          publish_time: String(pick(r, ['publish_time', 'PublishTime'], '')),
          status: Number(pick(r, ['status', 'Status'], 0)) || 0,
        } satisfies EpisodeRow
      })
      .sort((a, b) => a.episode_number - b.episode_number)
  } catch (e: any) {
    loadError.value = e?.message || '加载失败'
  } finally {
    loading.value = false
  }
}


function validate() {
  if (!form.value.projectTitle.trim()) return '请填写项目标题'
  if (!form.value.coverUrl.trim()) return '请上传封面或填写封面 URL'
  const year = Number(form.value.releaseYear)
  if (!Number.isFinite(year) || year < 1900 || year > 2100) return '上映/首播年份不合法'
  if (form.value.rating < 0 || form.value.rating > 10) return '评分必须在 0-10 之间'
  if (!Number.isFinite(form.value.totalEpisodes) || form.value.totalEpisodes <= 0) return '总集数必须大于 0'
  if (form.value.totalEpisodes < currentEpisodes.value) return `总集数不能小于当前已更新集数（${currentEpisodes.value}）`
  return ''
}

async function saveMeta() {
  saveError.value = ''
  saveOkMsg.value = ''
  const msg = validate()
  if (msg) {
    saveError.value = msg
    return
  }
  if (!auth.token) {
    saveError.value = '未登录'
    return
  }

  const body: UpdatePgcReq = {
    pgc_id: pgcId.value,
    title: form.value.projectTitle.trim(),
    cover: form.value.coverUrl.trim(),
    desc: form.value.synopsis.trim(),
    year: Number(form.value.releaseYear),
    area: form.value.area.trim(),
    rating: Number(form.value.rating) || 0,
    is_ongoing: Boolean(form.value.isOngoing),
    total_episodes: Number(form.value.totalEpisodes) || 1,
  }

  saving.value = true
  try {
    const res = await updatePgcApi(auth.token, body)
    if (res.code !== 200) {
      saveError.value = res.msg || '保存失败'
      return
    }
    saveOkMsg.value = '保存成功'
  } catch (e: any) {
    saveError.value = e?.message || '保存失败'
  } finally {
    saving.value = false
  }
}

async function togglePgcOffline() {
  if (!auth.token) {
    toast('error', '未登录')
    return
  }
  // 下架: -1，上架: 300(PGCAuditApproved)
  const toStatus = pgcStatus.value === -1 ? 300 : -1
  if (!confirm(toStatus === -1 ? '确定下架该内容？' : '确定上架该内容？')) return
  opBusy.value = true
  try {
    const res = await updatePgcStatusApi(auth.token, pgcId.value, toStatus)
    if (res.code !== 200) {
      toast('error', res.msg || '操作失败')
      return
    }
    toast('success', toStatus === -1 ? '已下架' : '已上架')
    await loadDetail()
  } finally {
    opBusy.value = false
  }
}

async function deletePgc() {
  if (!auth.token) {
    toast('error', '未登录')
    return
  }
  if (!confirm('确定删除该内容？此操作不可恢复。')) return
  opBusy.value = true
  try {
    const res = await deletePgcApi(auth.token, pgcId.value)
    if (res.code !== 200) {
      toast('error', res.msg || '删除失败')
      return
    }
    toast('success', '已删除')
    await router.replace('/content')
  } finally {
    opBusy.value = false
  }
}

const addingEp = ref(false)
const addEpError = ref('')
const addEpMode = ref<'upload' | 'bind'>('upload')
const addEpVidInput = ref('')
const addEp = ref<AddPgcEpisodeReq>({
  episode_number: 1,
  title: '',
  vid: 0,
  duration: 0,
  publish_time: '',
})
const addEpFile = ref<File | null>(null)
const addEpUploading = ref(false)
const addEpUploadPercent = ref(0)
const addEpUploadStage = ref<'idle' | 'checking' | 'uploading' | 'merging' | 'creating'>('idle')

function getFileBaseName(name: string) {
  const idx = name.lastIndexOf('.')
  return idx > 0 ? name.slice(0, idx) : name
}

async function readVideoDuration(file: File): Promise<number> {
  return new Promise((resolve) => {
    const url = URL.createObjectURL(file)
    const video = document.createElement('video')
    video.preload = 'metadata'
    video.onloadedmetadata = () => {
      const d = Number(video.duration)
      URL.revokeObjectURL(url)
      resolve(Number.isFinite(d) && d > 0 ? Math.round(d) : 0)
    }
    video.onerror = () => {
      URL.revokeObjectURL(url)
      resolve(0)
    }
    video.src = url
  })
}

async function md5Hex(file: File): Promise<string> {
  const buf = await file.arrayBuffer()
  return SparkMD5.ArrayBuffer.hash(buf)
}

async function onPickEpisodeFile(e: Event) {
  addEpError.value = ''
  const input = e.target as HTMLInputElement
  const file = input.files?.[0] ?? null
  addEpFile.value = file
  if (!file) return
  addEp.value.title = getFileBaseName(file.name)
  addEp.value.duration = await readVideoDuration(file)
  addEp.value.vid = 0
  addEp.value.publish_time = ''
}

async function uploadEpisodeFileAndFill() {
  if (!auth.token) {
    toast('error', '未登录')
    return
  }
  if (!addEpFile.value) {
    addEpError.value = '请先选择视频文件'
    return
  }
  addEpUploading.value = true
  addEpUploadPercent.value = 0
  addEpUploadStage.value = 'checking'
  addEpError.value = ''
  try {
    const file = addEpFile.value
    const hash = await md5Hex(file)
    const size = file.size

    const check = await checkVideoUploadApi(auth.token, { hash, size })
    if (check.code !== 200) {
      addEpError.value = check.msg || '检测上传状态失败'
      return
    }
    const chunks = ((check.data as any)?.chunks ?? []) as number[]
    const fileID = Number((check.data as any)?.fileID ?? 0) || undefined
    const isInstant = chunks.includes(-1)

    if (!isInstant) {
      addEpUploadStage.value = 'uploading'
      const up = await uploadVideoChunkApi(auth.token, {
        file,
        hash,
        chunkIndex: 0,
        totalChunks: 1,
        onProgress: (percent) => {
          addEpUploadPercent.value = percent
        },
      })
      if (up.code !== 200) {
        addEpError.value = up.msg || '上传视频失败'
        return
      }
      addEpUploadPercent.value = 100
      addEpUploadStage.value = 'merging'
      const merge = await mergeVideoUploadApi(auth.token, { hash, size, fileID })
      if (merge.code !== 200) {
        addEpError.value = merge.msg || '合并视频失败'
        return
      }
    } else {
      addEpUploadPercent.value = 100
    }

    addEpUploadStage.value = 'creating'
    const create = await createVideoUploadApi(auth.token, { hash, size, fileID })
    if (create.code !== 200) {
      addEpError.value = create.msg || '创建视频资源失败'
      return
    }
    const resource = (create.data as any)?.resource
    addEp.value.vid = Number(resource?.vid ?? 0)
    if (resource?.title) addEp.value.title = String(resource.title)
    if (resource?.duration) addEp.value.duration = Number(resource.duration) || addEp.value.duration
    addEp.value.publish_time = ''
    addEpError.value = ''
    toast('success', '视频上传成功，已自动填入剧集信息')
  } catch (e: any) {
    addEpError.value = e?.message || '上传失败'
  } finally {
    addEpUploadStage.value = 'idle'
    addEpUploading.value = false
  }
}

function initNextEpisodeNumber() {
  const max = episodes.value.reduce((m, r) => Math.max(m, r.episode_number), 0)
  addEp.value.episode_number = max + 1
}

async function addEpisode() {
  addEpError.value = ''
  if (!auth.token) {
    addEpError.value = '未登录'
    toast('error', '未登录')
    return
  }
  if (!addEp.value.episode_number || addEp.value.episode_number <= 0) addEp.value.episode_number = 1
  // 绑定已有视频模式：从手动输入的 vid 取值
  if (addEpMode.value === 'bind') {
    const parsed = Number(addEpVidInput.value)
    if (!parsed || parsed <= 0 || !Number.isInteger(parsed)) {
      addEpError.value = '请输入有效的视频 ID'
      return
    }
    addEp.value.vid = parsed
  }
  if (!addEp.value.vid || addEp.value.vid <= 0) {
    addEpError.value = '请先上传视频或输入已有视频 ID'
    return
  }
  addingEp.value = true
  try {
    const res = await addPgcEpisodeApi(auth.token, pgcId.value, {
      episode_number: Number(addEp.value.episode_number),
      title: addEp.value.title.trim(),
      vid: Number(addEp.value.vid),
      duration: Number(addEp.value.duration) || 0,
      publish_time: '',
    })
    if (res.code !== 200) {
      addEpError.value = res.msg || '添加失败'
      toast('error', addEpError.value)
      return
    }
    toast('success', '已添加剧集')
    await loadDetail()
    addEp.value.title = ''
    addEp.value.vid = 0
    addEp.value.duration = 0
    addEp.value.publish_time = ''
    addEpFile.value = null
    initNextEpisodeNumber()
  } catch (e: any) {
    addEpError.value = e?.message || '添加失败'
    toast('error', addEpError.value)
  } finally {
    addingEp.value = false
  }
}

async function saveEpisodeTitle(ep: EpisodeRow) {
  if (!auth.token) return
  const res = await updatePgcEpisodeApi(auth.token, pgcId.value, ep.id, { title: ep.title.trim() })
  if (res.code !== 200) {
    toast('error', res.msg || '保存标题失败')
    return
  }
  toast('success', `第 ${ep.episode_number} 集标题已保存`)
}

const deletingId = ref<number | null>(null)
const episodeActionOpenId = ref<number | null>(null)
const episodeActionAnchor = ref<HTMLElement | null>(null)
const episodeActionPos = ref<{ top: number; left: number }>({ top: 0, left: 0 })
const selectedEpisode = computed(() => episodes.value.find((x) => x.id === episodeActionOpenId.value) ?? null)

function syncEpisodeActionPos() {
  const el = episodeActionAnchor.value
  if (!el) return
  const r = el.getBoundingClientRect()
  episodeActionPos.value = {
    top: r.bottom + 8,
    left: Math.max(8, r.right - 144), // 144 = w-36
  }
}

function toggleEpisodeActionMenu(id: number, anchor?: HTMLElement | null) {
  if (episodeActionOpenId.value === id) {
    closeEpisodeActionMenu()
    return
  }
  episodeActionOpenId.value = id
  episodeActionAnchor.value = anchor ?? null
  syncEpisodeActionPos()
}

function closeEpisodeActionMenu() {
  episodeActionOpenId.value = null
  episodeActionAnchor.value = null
}

function onEpisodeDocClick(e: MouseEvent) {
  const target = e.target as HTMLElement | null
  if (!target) return
  if (target.closest('[data-episode-action-menu-root]')) return
  closeEpisodeActionMenu()
}

async function deleteEpisode(ep: EpisodeRow) {
  if (!confirm(`确定删除第 ${ep.episode_number} 集？`)) return
  if (!auth.token) return
  deletingId.value = ep.id
  try {
    const res = await deletePgcEpisodeApi(auth.token, pgcId.value, ep.id)
    if (res.code !== 200) {
      toast('error', res.msg || '删除失败')
      return
    }
    toast('success', `已删除第 ${ep.episode_number} 集`)
    closeEpisodeActionMenu()
    await loadDetail()
    initNextEpisodeNumber()
  } finally {
    deletingId.value = null
  }
}

async function toggleEpisodeOffline(ep: EpisodeRow) {
  if (!auth.token) return
  const next = ep.status === -1 ? 0 : -1
  if (!confirm(next === -1 ? `确定下架第 ${ep.episode_number} 集？` : `确定上架第 ${ep.episode_number} 集？`)) return
  deletingId.value = ep.id
  try {
    const res = await updatePgcEpisodeStatusApi(auth.token, pgcId.value, ep.id, next)
    if (res.code !== 200) {
      toast('error', res.msg || '操作失败')
      return
    }
    toast('success', next === -1 ? `已下架第 ${ep.episode_number} 集` : `已上架第 ${ep.episode_number} 集`)
    closeEpisodeActionMenu()
    await loadDetail()
    initNextEpisodeNumber()
  } finally {
    deletingId.value = null
  }
}

onMounted(async () => {
  document.addEventListener('click', onEpisodeDocClick)
  window.addEventListener('scroll', syncEpisodeActionPos, true)
  window.addEventListener('resize', syncEpisodeActionPos)
  if (!pgcId.value || pgcId.value === 'undefined' || pgcId.value === 'null') {
    await router.replace('/content')
    return
  }
  await loadDetail()
  initNextEpisodeNumber()
})

onBeforeUnmount(() => {
  document.removeEventListener('click', onEpisodeDocClick)
  window.removeEventListener('scroll', syncEpisodeActionPos, true)
  window.removeEventListener('resize', syncEpisodeActionPos)
})
</script>

<template>
  <div class="mx-auto max-w-[1280px] space-y-6">
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h2 class="text-xl font-semibold text-studio-fg">编辑内容</h2>
        <p class="mt-1 text-sm text-studio-muted">PGC ID：{{ pgcId }}</p>
      </div>
      <div class="flex items-center gap-3">
        <button
          type="button"
          class="rounded-xl border border-studio-border bg-studio-card px-4 py-2 text-sm text-studio-fg-muted transition hover:bg-studio-elevated"
          @click="router.push('/content')"
        >
          返回列表
        </button>
        <button
          type="button"
          class="rounded-xl border border-studio-border bg-studio-card px-4 py-2 text-sm text-studio-fg-muted transition hover:bg-studio-elevated disabled:opacity-60"
          :disabled="opBusy"
          @click="togglePgcOffline"
        >
          下架/上架
        </button>
        <button
          type="button"
          class="rounded-xl border border-rose-500/30 bg-rose-500/10 px-4 py-2 text-sm text-rose-400 transition hover:bg-rose-500/15 disabled:opacity-60"
          :disabled="opBusy"
          @click="deletePgc"
        >
          删除
        </button>
        <button
          type="button"
          class="rounded-xl bg-gradient-to-r from-cyan-500 to-cyan-600 px-4 py-2 text-sm font-semibold text-slate-950 shadow-lg shadow-cyan-500/20 transition hover:from-cyan-400 hover:to-cyan-500 disabled:opacity-60"
          :disabled="saving || loading"
          @click="saveMeta"
        >
          {{ saving ? '保存中…' : '保存元信息' }}
        </button>
      </div>
    </div>

    <div v-if="loadError" class="rounded-2xl border border-rose-500/30 bg-rose-500/10 p-4 text-sm text-rose-400">
      {{ loadError }}
    </div>
    <div v-else-if="loading" class="rounded-2xl border border-studio-border bg-studio-card p-6 text-sm text-studio-muted">
      加载中…
    </div>

    <section v-if="!loading && !loadError" class="rounded-2xl border border-studio-border bg-studio-card p-6">
      <h3 class="text-sm font-semibold text-studio-fg">核心元数据</h3>

      <div v-if="saveError" class="mt-4 rounded-xl border border-rose-500/30 bg-rose-500/10 p-3 text-sm text-rose-400">
        {{ saveError }}
      </div>
      <div v-if="saveOkMsg" class="mt-4 rounded-xl border border-emerald-500/30 bg-emerald-500/10 p-3 text-sm text-emerald-400">
        {{ saveOkMsg }}
      </div>

      <div class="mt-4 grid gap-4 sm:grid-cols-2">
        <div class="sm:col-span-2">
          <label class="mb-1.5 block text-xs font-medium text-studio-fg-muted">项目标题</label>
          <input
            v-model="form.projectTitle"
            type="text"
            class="w-full rounded-xl border border-studio-border bg-studio-input px-4 py-2.5 text-sm text-studio-fg placeholder:text-studio-muted focus:border-cyan-500/50 focus:outline-none focus:ring-1 focus:ring-cyan-500/30"
          />
        </div>

        <div>
          <label class="mb-1.5 block text-xs font-medium text-studio-fg-muted">类型</label>
          <select
            v-model="form.category"
            class="w-full rounded-xl border border-studio-border bg-studio-input px-4 py-2.5 text-sm text-studio-fg focus:border-cyan-500/50 focus:outline-none focus:ring-1 focus:ring-cyan-500/30"
          >
            <option v-for="c in categories" :key="c.value" :value="c.value">
              {{ c.label }}
            </option>
          </select>
          <p class="mt-1 text-xs text-amber-200/90">提示：后端当前 `update` 接口不支持改 `pgc_type`，这里只做展示/预留。</p>
        </div>

        <div>
          <label class="mb-1.5 block text-xs font-medium text-studio-fg-muted">上映 / 首播年份</label>
          <input
            v-model="form.releaseYear"
            type="text"
            class="w-full rounded-xl border border-studio-border bg-studio-input px-4 py-2.5 text-sm text-studio-fg focus:border-cyan-500/50 focus:outline-none focus:ring-1 focus:ring-cyan-500/30"
          />
        </div>

        <div class="sm:col-span-2">
          <label class="mb-1.5 block text-xs font-medium text-studio-fg-muted">封面（上传或填写 URL）</label>
          <div class="flex flex-col gap-3 sm:flex-row sm:items-center">
            <input
              ref="coverFileInput"
              type="file"
              accept="image/*"
              class="block w-full text-sm text-studio-muted file:mr-3 file:rounded-lg file:border-0 file:bg-studio-elevated file:px-3 file:py-2 file:text-sm file:font-medium file:text-studio-fg hover:file:text-cyan-300"
              :disabled="coverUploading"
              @change="onPickCoverFile"
            />
            <span class="text-xs text-studio-muted whitespace-nowrap">
              {{ coverUploading ? '上传中…' : '或' }}
            </span>
            <input
              v-model="form.coverUrl"
              type="text"
              class="w-full rounded-xl border border-studio-border bg-studio-input px-4 py-2.5 text-sm text-studio-fg placeholder:text-studio-muted focus:border-cyan-500/50 focus:outline-none focus:ring-1 focus:ring-cyan-500/30"
            />
          </div>
          <p v-if="coverUploadError" class="mt-2 text-xs text-rose-400">{{ coverUploadError }}</p>
        </div>

        <div>
          <label class="mb-1.5 block text-xs font-medium text-studio-fg-muted">地区（选填）</label>
          <input
            v-model="form.area"
            type="text"
            class="w-full rounded-xl border border-studio-border bg-studio-input px-4 py-2.5 text-sm text-studio-fg placeholder:text-studio-muted focus:border-cyan-500/50 focus:outline-none focus:ring-1 focus:ring-cyan-500/30"
          />
        </div>

        <div>
          <label class="mb-1.5 block text-xs font-medium text-studio-fg-muted">评分（0-10）</label>
          <input
            v-model.number="form.rating"
            type="number"
            min="0"
            max="10"
            step="0.1"
            class="w-full rounded-xl border border-studio-border bg-studio-input px-4 py-2.5 text-sm text-studio-fg focus:border-cyan-500/50 focus:outline-none focus:ring-1 focus:ring-cyan-500/30"
          />
        </div>

        <div>
          <label class="mb-1.5 block text-xs font-medium text-studio-fg-muted">总集数</label>
          <input
            v-model.number="form.totalEpisodes"
            type="number"
            min="1"
            step="1"
            class="w-full rounded-xl border border-studio-border bg-studio-input px-4 py-2.5 text-sm text-studio-fg focus:border-cyan-500/50 focus:outline-none focus:ring-1 focus:ring-cyan-500/30"
          />
          <p class="mt-1 text-xs text-studio-muted">当前已更新：{{ currentEpisodes }} 集</p>
        </div>

        <div class="sm:col-span-2">
          <label class="mb-1.5 block text-xs font-medium text-studio-fg-muted">简介</label>
          <textarea
            v-model="form.synopsis"
            rows="4"
            class="w-full resize-y rounded-xl border border-studio-border bg-studio-input px-4 py-3 text-sm text-studio-fg placeholder:text-studio-muted focus:border-cyan-500/50 focus:outline-none focus:ring-1 focus:ring-cyan-500/30"
          />
        </div>

        <div class="sm:col-span-2">
          <label class="inline-flex items-center gap-2 text-sm text-studio-fg-muted">
            <input
              v-model="form.isOngoing"
              type="checkbox"
              class="sr-only"
            />
            <span
              class="relative inline-flex h-4 w-4 items-center justify-center rounded-[4px] border bg-studio-card transition"
              :class="form.isOngoing ? 'border-cyan-500' : 'border-studio-border'"
            >
              <span
                class="text-[11px] leading-none text-cyan-400 transition-opacity"
                :class="form.isOngoing ? 'opacity-100' : 'opacity-0'"
              >
                ✓
              </span>
            </span>
            是否连载中
          </label>
        </div>
      </div>
    </section>

    <section v-if="!loading && !loadError" class="rounded-2xl border border-studio-border bg-studio-card p-6">
      <div class="flex flex-col justify-between gap-2 sm:flex-row sm:items-center">
        <div>
          <h3 class="text-sm font-semibold text-studio-fg">剧集管理</h3>
          <p class="mt-1 text-xs text-studio-muted">支持编辑标题、上下架、删除单集；新增单集请从已上传资源中选择。</p>
        </div>
        <span class="text-xs text-studio-muted">共 {{ episodes.length }} 集</span>
      </div>

      <div class="mt-4 overflow-auto rounded-xl border border-studio-border">
        <table class="w-full min-w-[860px] border-collapse text-left text-sm">
          <thead class="sticky top-0 bg-studio-elevated text-xs text-studio-muted">
            <tr>
              <th class="w-24 px-3 py-2 font-medium">集号</th>
              <th class="px-3 py-2 font-medium">标题</th>
              <th class="w-40 px-3 py-2 font-medium">vid</th>
              <th class="w-40 px-3 py-2 font-medium">时长(秒)</th>
              <th class="w-56 px-3 py-2 font-medium">发布时间</th>
              <th class="w-24 px-3 py-2 font-medium text-right">操作</th>
            </tr>
          </thead>
          <tbody class="text-studio-fg-muted">
            <tr v-for="ep in episodes" :key="ep.id">
              <td class="whitespace-nowrap px-3 py-2 font-mono text-studio-muted">第 {{ ep.episode_number }} 集</td>
              <td class="px-3 py-2">
                <input
                  v-model="ep.title"
                  type="text"
                  class="w-full rounded-lg border border-transparent bg-studio-elevated/30 px-2 py-1 text-sm text-studio-fg outline-none focus:border-cyan-500/40"
                  placeholder="输入剧集标题"
                  @blur="saveEpisodeTitle(ep)"
                />
              </td>
              <td class="px-3 py-2">{{ ep.vid }}</td>
              <td class="px-3 py-2">{{ ep.duration }}</td>
              <td class="px-3 py-2">{{ ep.publish_time || '—' }}</td>
              <td class="px-3 py-2 text-right">
                <div class="relative inline-flex" data-episode-action-menu-root>
                  <button
                    type="button"
                    class="inline-flex h-8 w-8 items-center justify-center rounded-lg text-studio-muted transition hover:bg-studio-elevated/70 hover:text-studio-fg-muted disabled:opacity-50"
                    :disabled="deletingId === ep.id"
                    @click.stop="toggleEpisodeActionMenu(ep.id, $event.currentTarget as HTMLElement)"
                  >
                    <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-width="2" d="M12 5v.01M12 12v.01M12 19v.01" stroke-linecap="round" />
                    </svg>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="mt-6 rounded-xl border border-studio-border bg-studio-elevated/30 p-4">
        <h4 class="text-sm font-semibold text-studio-fg">新增一集</h4>

        <!-- 模式切换 -->
        <div class="mt-2 flex gap-2">
          <button
            type="button"
            class="rounded-lg px-3 py-1 text-xs font-medium transition"
            :class="addEpMode === 'upload'
              ? 'bg-cyan-500/20 text-cyan-300'
              : 'bg-studio-elevated/40 text-studio-muted hover:text-studio-fg-subtle'"
            @click="addEpMode = 'upload'"
          >上传新视频</button>
          <button
            type="button"
            class="rounded-lg px-3 py-1 text-xs font-medium transition"
            :class="addEpMode === 'bind'
              ? 'bg-cyan-500/20 text-cyan-300'
              : 'bg-studio-elevated/40 text-studio-muted hover:text-studio-fg-subtle'"
            @click="addEpMode = 'bind'"
          >绑定已有视频</button>
        </div>

        <p v-if="addEpError" class="mt-2 text-xs text-rose-400">{{ addEpError }}</p>

        <div class="mt-3 grid grid-cols-1 gap-3 sm:grid-cols-[minmax(0,5rem)_minmax(0,1fr)_minmax(0,1.2fr)]">
          <input
            v-model.number="addEp.episode_number"
            type="number"
            min="1"
            class="rounded-xl border border-studio-border bg-studio-elevated/40 px-3 py-2 text-sm text-studio-fg-muted"
            placeholder="集号"
            readonly
          />

          <!-- 上传模式：文件选择 -->
          <div v-if="addEpMode === 'upload'" class="min-w-0">
            <input
              type="file"
              accept="video/*"
              class="block w-full text-sm text-studio-muted file:mr-3 file:rounded-lg file:border-0 file:bg-studio-elevated file:px-3 file:py-2 file:text-sm file:font-medium file:text-studio-fg hover:file:text-cyan-300"
              :disabled="addEpUploading"
              @change="onPickEpisodeFile"
            />
          </div>

          <!-- 绑定模式：vid 输入 -->
          <div v-else class="min-w-0">
            <input
              v-model="addEpVidInput"
              type="number"
              min="1"
              step="1"
              class="w-full rounded-xl border border-studio-border bg-studio-input px-3 py-2 text-sm text-studio-fg"
              placeholder="输入已有视频 ID"
            />
          </div>

          <!-- 上传模式标题 auto-fill / 绑定模式可手动填 -->
          <input
            v-model="addEp.title"
            type="text"
            :readonly="addEpMode === 'upload'"
            class="min-w-0 rounded-xl border border-studio-border px-3 py-2 text-sm"
            :class="addEpMode === 'upload'
              ? 'bg-studio-elevated/40 text-studio-fg-muted'
              : 'bg-studio-input text-studio-fg'"
            :placeholder="addEpMode === 'upload' ? '自动带出标题' : '剧集标题（选填）'"
          />
        </div>

        <!-- 上传模式：上传按钮 + 进度 -->
        <template v-if="addEpMode === 'upload'">
          <div class="mt-3 flex items-center justify-between gap-3">
            <p class="text-xs text-studio-muted">
              {{ addEpFile ? (addEpUploading ? '正在上传视频，请稍候…' : '已选择视频文件，点击右侧按钮执行上传') : '请先选择视频文件' }}
            </p>
            <button
              type="button"
              class="inline-flex items-center gap-2 rounded-xl bg-gradient-to-r from-cyan-500 to-cyan-600 px-4 py-2 text-sm font-semibold text-slate-950 shadow-lg shadow-cyan-500/20 transition hover:from-cyan-400 hover:to-cyan-500 disabled:cursor-not-allowed disabled:opacity-60"
              :disabled="addEpUploading || !addEpFile"
              @click="uploadEpisodeFileAndFill"
            >
              <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-width="2" d="M12 16V7m0 0l-4 4m4-4l4 4M4 17v1a2 2 0 002 2h12a2 2 0 002-2v-1" />
              </svg>
              {{ addEpUploading ? '上传中…' : '上传视频并自动填入' }}
            </button>
          </div>
          <div v-if="addEpUploading" class="mt-3">
            <div class="flex items-center justify-between text-xs text-studio-muted">
              <span>
                {{
                  addEpUploadStage === 'checking'
                    ? '校验文件中…'
                    : addEpUploadStage === 'uploading'
                    ? '上传中…'
                    : addEpUploadStage === 'merging'
                    ? '服务端合并中…'
                    : '创建视频资源中…'
                }}
              </span>
              <span>{{ addEpUploadPercent }}%</span>
            </div>
            <div class="mt-1 h-2 w-full overflow-hidden rounded bg-studio-input">
              <div
                class="h-full rounded bg-gradient-to-r from-cyan-500 to-cyan-400 transition-all duration-200"
                :style="{ width: `${addEpUploadPercent}%` }"
              ></div>
            </div>
          </div>
        </template>

        <!-- 绑定模式：提示 -->
        <p v-else class="mt-3 text-xs text-studio-muted">
          输入已存在于系统的视频 ID，添加后该视频将被标记为 PGC 资产并从 UGC 分区移除。
        </p>

        <div class="mt-3 flex justify-end">
          <button
            type="button"
            class="rounded-xl bg-gradient-to-r from-cyan-500 to-cyan-600 px-4 py-2 text-sm font-semibold text-slate-950 shadow-lg shadow-cyan-500/20 transition hover:from-cyan-400 hover:to-cyan-500 disabled:opacity-60"
            :disabled="addingEp || (addEpMode === 'upload' ? false : !addEpVidInput)"
            @click="addEpisode"
          >
            {{ addingEp ? '添加中…' : '添加' }}
          </button>
        </div>
      </div>
    </section>
  </div>

  <teleport to="body">
    <div
      v-if="episodeActionOpenId != null"
      data-episode-action-menu-root
      class="fixed z-[9999] w-36 overflow-hidden rounded-xl border shadow-[0_8px_24px_rgba(0,0,0,0.25)] backdrop-blur-sm"
      :class="
        theme.mode === 'dark'
          ? 'border-studio-border/80 bg-[#0a1426]/95'
          : 'border-slate-200/80 bg-white/95'
      "
      :style="{ top: `${episodeActionPos.top}px`, left: `${episodeActionPos.left}px` }"
    >
      <button
        type="button"
        class="flex w-full items-center gap-2 px-3 py-2 text-left text-sm transition"
        :class="
          theme.mode === 'dark'
            ? 'text-studio-fg-subtle hover:bg-studio-elevated/60 hover:text-amber-200'
            : 'text-slate-700 hover:bg-slate-100 hover:text-amber-600'
        "
        :disabled="deletingId === episodeActionOpenId || !selectedEpisode"
        @click="selectedEpisode && toggleEpisodeOffline(selectedEpisode)"
      >
        {{ (selectedEpisode?.status ?? 0) === -1 ? '上架' : '下架' }}
      </button>
      <button
        type="button"
        class="flex w-full items-center gap-2 px-3 py-2 text-left text-sm transition disabled:opacity-50"
        :class="
          theme.mode === 'dark'
            ? 'text-rose-300 hover:bg-studio-elevated/60 hover:text-rose-200'
            : 'text-rose-600 hover:bg-rose-50 hover:text-rose-700'
        "
        :disabled="deletingId === episodeActionOpenId || !selectedEpisode"
        @click="selectedEpisode && deleteEpisode(selectedEpisode)"
      >
        {{ deletingId === episodeActionOpenId ? '删除中…' : '删除' }}
      </button>
    </div>
  </teleport>
</template>

