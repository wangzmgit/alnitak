<template>
  <div id="mask" class="leaf-captcha" :class="props.show ? 'captcha-show' : ''" @click="maskClick">
    <div class="captcha-container">
      <div class="captcha-img">
        <refresh class="refresh" @click="reset()"></refresh>
        <img class="bg-img" :src="state.bg" alt="背景图">
        <img class="sub-img" :style="`top: ${state.top}px;left: ${state.left}px`" :src="state.subImg" alt="拼图">
        <div class="msg" :style="`background-color: ${state.msgBgColor}`">{{ state.msgText }}</div>
      </div>
      <div ref="sliderRef" class="slider-bar">
        <span class="range-text">拖动滑块完成拼图</span>
        <div class="slider-played" :style="`width: ${state.activeWidth}px`">
          <div ref="blockRef" class="slider-block"></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref, watch } from "vue";
import { statusCode } from "@/utils/status-code";
import { Refresh } from "@vicons/ionicons5";
import { getSliderAPI, validateSliderAPI } from "@/api/captcha";

const emits = defineEmits(["update:show", "success"]);
const props = defineProps<{
  captchaId: string,
  show: boolean,
}>();

const initialOffset = 3;//left初始偏移量
const initialSliderSize = 50; //初始滑块大小
const state = reactive({
  top: 0,
  bg: "",
  subImg: "",
  msgBgColor: "",
  msgText: "",
  canMove: true,
  left: initialOffset,
  activeWidth: initialSliderSize
});

//控制显示/关闭
const show = ref(props.show);
watch(() => props.show, (value) => {
  if (value) getCaptcha();
  show.value = value;
});

//点击遮罩关闭
const maskClick = (e: Event) => {
  if ((e.target as HTMLElement).id === "mask") {
    close();
  }
}

const close = () => {
  reset();
  show.value = false;
  emits("update:show", show.value);
}

//获取验证码图片
const getCaptcha = () => {
  getSliderAPI(props.captchaId).then((res) => {
    if (res.data.code === statusCode.OK) {
      state.top = res.data.data.slider_captcha.y;
      state.bg = `data:image/jpg;base64,${res.data.data.slider_captcha.bg_img}`;
      state.subImg = `data:image/jpg;base64,${res.data.data.slider_captcha.slider_img}`;
    }
  })
}

//提示文本
const setMsg = (status: string) => {
  switch (status) {
    case "wait":
      state.msgText = "等待服务端响应";
      state.msgBgColor = "rgba(106, 160, 255, .8)";
      break;
    case "success":
      state.msgText = "验证通过";
      state.msgBgColor = "rgba(131, 206, 63, .8)";
      emits("success");
      break;
    case "fail":
      state.msgText = "验证不通过";
      state.msgBgColor = "rgba(254, 108, 108, .8)";
      break;
    default:
      state.msgText = "";
      state.msgBgColor = "";
      break;
  }
}

// 处理滑动
const blockRef = ref<HTMLElement | null>(null);
const sliderRef = ref<HTMLElement | null>(null);
// 滑动滑动条
const slidingSlider = () => {
  //鼠标事件
  blockRef.value!.onmousedown = function () {
    document.onmousemove = function (e) {
      sliderValueChange(e);
    };
    document.onmouseup = function () {
      moveEnd();
      document.onmousemove = document.onmouseup = null;
    };
  };
  // 触摸事件
  blockRef.value!.ontouchstart = function () {
    document.ontouchmove = function (e) {
      sliderValueChange(e);
    };
    document.ontouchend = function () {
      moveEnd();
      document.ontouchmove = document.ontouchend = null;
    };
  };
}

//滑动条值改变
const sliderValueChange = (e: MouseEvent | TouchEvent) => {
  if (!state.canMove) return;
  const sliderWidth = sliderRef.value!.clientWidth;
  const clientX = e instanceof MouseEvent ? e.clientX : e.changedTouches[0].clientX;
  let activeSize = (clientX - sliderRef.value!.getBoundingClientRect().left);

  activeSize = Math.max(50, activeSize);
  activeSize = Math.min(sliderWidth, activeSize);

  state.activeWidth = activeSize;
  state.left = Math.min(activeSize - 50 - initialOffset, sliderWidth);
}

//滑动结束
const moveEnd = () => {
  state.canMove = false;
  setMsg("wait");
  validateSliderAPI(props.captchaId, Math.round(state.left)).then((res) => {
    if (res.data.code === 200) {
      setMsg("success");
      setTimeout(() => {
        close();
      }, 500)
    } else {
      setMsg("fail")
      setTimeout(() => {
        reset();
      }, 800)
    }
  });
}

//重置
const reset = () => {
  setMsg("")
  getCaptcha();
  state.canMove = true;
  state.left = initialOffset;
  state.activeWidth = initialSliderSize;
}

onMounted(() => { slidingSlider(); });

</script>

<style lang="scss" scoped>
.leaf-captcha {
  position: fixed;
  top: 0;
  left: 0;
  bottom: 0;
  right: 0;
  background-color: rgba(0, 0, 0, .3);
  z-index: 999;
  opacity: 0;
  pointer-events: none;
  transition: opacity .2s;
}

.captcha-show {
  opacity: 1;
  pointer-events: auto;
}

.captcha-container {
  position: absolute;
  top: 40%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 350px;
  height: 270px;
  box-sizing: border-box;
  // 禁止选中
  user-select: none;
  -moz-user-select: -moz-none;
  -moz-user-select: none;
  -o-user-select: none;
  -khtml-user-select: none;
  -webkit-user-select: none;
  -ms-user-select: none;

  padding: 20px 20px 0 20px;
  background-color: #fff;
  border-radius: 3px;
  box-shadow: 0 1px 3px rgb(0 0 0 / 30%);
}

.captcha-img {
  width: 310px;
  height: 160px;
  position: relative;

  .refresh {
    position: absolute;
    right: 6px;
    top: 6px;
    width: 26px;
    height: 26px;
    color: #43cf96;
    cursor: pointer;
  }

  .bg-img {
    width: 100%;
    height: 100%;
  }

  .sub-img {
    position: absolute;
    width: 50px;
    height: 50px;
    box-shadow: 0 1px 3px rgb(0 0 0 / 30%);
  }

  .msg {
    position: absolute;
    bottom: 0;
    width: 100%;
    text-align: center;
    transition: all .2s;
    color: #fff;
  }
}

.slider-bar {
  position: relative;
  width: 100%;
  height: 50px;
  margin-top: 20px;
  background-color: #eef1f8;

  .range-text {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    font-size: 14px;
    color: #b7bcd1;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    text-align: center;
    width: 100%;
  }

  .slider-played {
    position: absolute;
    height: inherit;
    background-color: rgba(106, 160, 255, .8);


    .slider-block {
      position: absolute;
      width: 50px;
      height: 50px;
      background-color: #fff;
      box-shadow: 0 0 4px #ccc;
      cursor: pointer;
      border-radius: 3px;
      right: 0;
      display: flex;
      align-items: center;
      justify-content: center;
    }
  }
}
</style>