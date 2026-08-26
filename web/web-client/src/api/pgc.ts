import request from '@/utils/request';
import { baseURL } from '@/utils/request';
import { useAsyncData } from 'nuxt/app';

export interface PGCRecommendQuery {
  page: number;
  pageSize: number;
  pgcType?: number;
  seedPgcId?: string;
  scene?: string;
}

export interface PGCSearchQuery {
  page: number;
  pageSize: number;
  keyword: string;
  pgcType?: number;
}

function buildRecommendQuery(query: PGCRecommendQuery) {
  const params = new URLSearchParams();
  params.set('page', String(query.page));
  params.set('page_size', String(query.pageSize));
  if (query.pgcType) params.set('pgc_type', String(query.pgcType));
  if (query.seedPgcId) params.set('seed_pgc_id', query.seedPgcId);
  if (query.scene) params.set('scene', query.scene);
  return params.toString();
}

// 首页/详情页 PGC 推荐（新接口）
export const getPGCRecommendAPI = (query: PGCRecommendQuery) => {
  return request.get(`v1/pgc/recommend?${buildRecommendQuery(query)}`);
};

// SSR 首屏获取
export const asyncGetPGCRecommendAPI = async (query: PGCRecommendQuery) => {
  return await useAsyncData(`pgc-recommend-${query.page}-${query.pageSize}-${query.seedPgcId || ''}-${query.scene || ''}`, () =>
    $fetch(`${baseURL}/api/v1/pgc/recommend?${buildRecommendQuery(query)}`),
  );
};

// 搜索 PGC（后端要求 page / page_size / pgc_type 都是数字参数）
export const searchPGCAPI = (query: PGCSearchQuery) => {
  const params = new URLSearchParams();
  params.set('page', String(query.page));
  params.set('page_size', String(query.pageSize));
  params.set('keyword', query.keyword);
  params.set('pgc_type', String(query.pgcType ?? 0));
  return request.get(`v1/pgc/search?${params.toString()}`);
};

export const asyncGetPGCRecommendByVideoAPI = async (vid: number, page = 1, pageSize = 10) => {
  return await useAsyncData(`pgc-recommend-by-video-${vid}-${page}-${pageSize}`, () =>
    $fetch(`${baseURL}/api/v1/pgc/recommend-by-video?vid=${vid}&page=${page}&page_size=${pageSize}`),
  );
};

export const getPGCRecommendByVideoAPI = (vid: number, page = 1, pageSize = 10) => {
  return request.get(`v1/pgc/recommend-by-video?vid=${vid}&page=${page}&page_size=${pageSize}`);
};

export const getPGCDetailWithEpisodesAPI = (pgcId: string) => {
  return request.get(`v1/pgc/detail-with-episodes?pgc_id=${encodeURIComponent(pgcId)}`);
};

export const getPGCPlayPanelByVideoAPI = (vid: number | string, seasonId?: string) => {
  const params = new URLSearchParams();
  params.set('vid', String(vid));
  if (seasonId) params.set('season_id', seasonId);
  return request.get(`v1/pgc/play-panel-by-video?${params.toString()}`);
};

export const getPGCEpisodeDetailAPI = (epId: number | string) => {
  return request.get(`v1/pgc/episode-detail?ep_id=${encodeURIComponent(String(epId))}`);
};

