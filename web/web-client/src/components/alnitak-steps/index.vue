<template>
  <div class="alnitak-steps">
    <div class="alnitak-step" v-for="(item, index) in props.data">
      <div class="step-indicator" :class="index + 1 <= props.current ? 'step-indicator-active' : ''">
        <div class="step-indicator-slot">
          <div class="slot-icon" v-if="index + 1 === props.current && props.status === 'error'">
            <close-icon size="14"></close-icon>
          </div>
          <div class="slot-icon" v-else-if="index + 1 === props.current && props.status === 'finish'">
            <check-icon size="14"></check-icon>
          </div>
          <div v-else class="slot-index">{{ index + 1 }}</div>
        </div>
      </div>
      <div class="step-content">
        <div class="step-content-header">
          <div class="header-title">{{ item }}</div>
          <div v-if="index < props.data.length - 1" class="step-splitor"
            :class="index + 1 <= props.current - 1 ? 'step-splitor-active' : ''"></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Close as CloseIcon, Check as CheckIcon } from "@icon-park/vue-next";

const props = withDefaults(defineProps<{
  current: number,
  status?: string,
  data: string[]
}>(), {
  status: "process",
})

</script>

<style lang="scss" scoped>
.alnitak-steps {
  width: 100%;
  display: flex;

  .alnitak-step {
    position: relative;
    display: flex;
    flex: 1;

    .step-indicator {
      height: 28px;
      width: 28px;
      color: #c2c2c2;
      background-color: #fff;
      box-shadow: 0 0 0 1px #c2c2c2;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      transition: background-color .3s cubic-bezier(.4, 0, .2, 1), box-shadow .3s cubic-bezier(.4, 0, .2, 1);

      .step-indicator-slot {
        position: relative;
        width: 18px;
        height: 18px;
        font-size: 18px;
        line-height: 18px;

        .slot-index {
          display: inline-block;
          text-align: center;
          position: absolute;
          left: 0;
          top: 0;
          white-space: nowrap;
          font-size: 16px;
          width: 18px;
          height: 18px;
          line-height: 18px;
          transition: color .3s cubic-bezier(.4, 0, .2, 1);
        }

        .slot-icon {
          display: inline-block;
          text-align: center;
          position: absolute;
          left: 2px;
          top: 0;
          white-space: nowrap;
        }
      }
    }

    .step-indicator-active {
      color: #fff;
      background-color: var(--primary-color);
      box-shadow: 0 0 0 1px var(--primary-color);
    }

    .step-content {
      flex: 1;

      .step-content-header {
        position: relative;
        display: flex;
        color: var(--font-primary-1);
        margin: 6px 0 0 9px;
        line-height: 16px;
        font-size: 16px;
        font-weight: 500;
        transition: color .3s cubic-bezier(.4, 0, .2, 1), background-color .3s cubic-bezier(.4, 0, .2, 1);

        .header-title {
          white-space: nowrap;
          flex: 0;
        }

        .step-splitor {
          background-color: #c2c2c2;
          margin: 8px 12px 0;
          height: 1px;
          flex: 1;
          align-self: flex-start;
          transition: color .3s cubic-bezier(.4, 0, .2, 1), background-color .3s cubic-bezier(.4, 0, .2, 1);
        }

        .step-splitor-active {
          height: 2px;
          background-color: var(--primary-color);
        }
      }
    }
  }
}
</style>