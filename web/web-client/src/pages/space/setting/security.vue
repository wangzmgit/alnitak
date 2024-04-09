<template>
  <div class="security">
    <ul class="security-list">
      <li class="security-item">
        <div class="item-title">设置密码</div>
        <div class="item-value">******</div>
        <nuxt-link class="item-btn" to="/setpassword">修改密码</nuxt-link>
      </li>
      <li class="security-item">
        <div class="item-title">绑定邮箱</div>
        <div class="item-value">******</div>
        <div class="item-btn"></div>
      </li>
    </ul>
  </div>
</template>

<script setup lang="ts">
</script>

<style lang="scss" scoped>
.security {
  .security-list {
    padding: 0 20px;
    list-style: none;

    .security-item {
      min-height: 80px;
      display: flex;
      align-items: center;
      justify-content: space-between;
      border-bottom: 1px solid #e5e9ef;

      .item-title {
        box-sizing: border-box;
        font-size: 14px;
        color: #222;
        width: 30%;
        padding: 0 16px;
      }

      .item-value {
        box-sizing: border-box;
        font-size: 12px;
        color: #6d757a;
        flex: 1;
        padding: 0 16px;
      }

      .item-btn {
        box-sizing: border-box;
        width: 160px;
        padding: 0 16px;
        text-align: center;
        cursor: pointer;
        font-size: 12px;
        color: var(--primary-color);

        &:hover {
          color: var(--primary-hover-color);
        }
      }
    }
  }
}
</style>