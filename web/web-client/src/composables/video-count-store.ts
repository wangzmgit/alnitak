import { defineStore } from 'pinia';

export const useVideoCountStore = defineStore({
  id: "videoCount",
  state: () => {
    return {
      videoCount: 0,
    }
  },
  actions: {
    setVideoCountState(val: number) {
      this.videoCount = val;
    }
  },
})
