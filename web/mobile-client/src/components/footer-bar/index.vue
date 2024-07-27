<template>
  <footer class="footer">
    <div class="links link-list">
      <span>
        <a :href="webUrl">桌面端</a>
      </span>
      <span>|</span>
      <span>
        Powered by <a href="https://github.com/wangzmgit/alnitak" target="_blank">Alnitak</a>
      </span>
    </div>
    <div class="copyright link-list">
      <span v-show="globalConfig.icp">
        <a href="https://beian.miit.gov.cn/#/Integrated/index" target="_blank">{{ globalConfig.icp }}</a>
      </span>
      <span v-show="globalConfig.icp">|</span>
      <span v-show="globalConfig.security" class="security">
        <img class="security" src="@/assets/images/common/filing.png" alt="公网安备" />
        <a href="https://www.beian.gov.cn/portal/registerSystemInfo" target="_blank">
          {{ globalConfig.security }}
        </a>
      </span>
    </div>
  </footer>
</template>

<script setup lang="ts">
import { globalConfig } from '@/utils/global-config';

const webUrl = `http${globalConfig.https ? 's' : ''}://${globalConfig.domain}/`;
</script>

<style lang="scss" scoped>
.footer {
  color: #666;
  width: 100%;
  height: 80px;
  text-align: center;
  background-color: #f6f7f8;

  .links {
    padding-top: 16px;
  }

  .copyright {
    padding-top: 8px;

    .security {
      img {
        width: 12px;
        padding: 0 3px;
      }
    }
  }

  .link-list {
    span {
      margin: 0 3px;

      a {
        color: #666;
        text-decoration: none;

        &:hover {
          text-decoration: underline;
        }
      }
    }
  }
}
</style>