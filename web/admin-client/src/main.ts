import { createApp } from "vue";
import App from "@/App.vue";

// pinia
import stores from "@/stores";

// router
import router from "@/router";
import "@/router/guards";

createApp(App).use(router).use(stores).mount("#app");
