<script setup lang="ts">
// 重定向到用户ID所在的位置
definePageMeta({
  middleware: [(to) => {
    if (!to.params.id) {
      const userId = useCookie('user_id');
      if (userId.value) { 
        return navigateTo({ name: "user-id-manuscript", params: { id: userId.value } });
      } else {
        return navigateTo('/login');
      }
    }
  }]
})
</script>