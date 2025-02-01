import { onBeforeUnmount, onMounted, ref } from "vue";
import { throttle } from "lodash-es";
import { NLayoutContent, NLayoutHeader } from "naive-ui";
import useSystemStore from "@/stores/modules/system-store";

export default ({
    contentPadding,
}: {
    /** 内容padding尺寸，即layout-content-main内的padding尺寸 */
    contentPadding: number;
}) => {
    const style_content_padding = contentPadding;
    const style_content_height = ref("100%");

    const systemStore = useSystemStore();

    const headerRef = ref<InstanceType<typeof NLayoutHeader> | null>(null);
    const contentRef = ref<InstanceType<typeof NLayoutContent> | null>(null);

    const resizeHandler = throttle(() => {
        const headerEl: HTMLElement = headerRef.value?.$el;
        const contentEl: HTMLElement = contentRef.value?.$el;

        const headerHeight = headerEl.offsetHeight;
        const contentHeight = window.innerHeight - headerHeight;
        const contentBorderBoxHeight = contentHeight - style_content_padding * 2;
        const contentWidth = contentEl.offsetWidth - style_content_padding * 2;

        style_content_height.value = `${contentHeight}px`;

        systemStore.$patch({ headerHeight, contentHeight, contentBorderBoxHeight, contentWidth });
    }, 200);

    const changeContentWidth = () => {
        const contentEl = contentRef.value?.$el as HTMLElement;
        systemStore.contentWidth = contentEl.offsetWidth - style_content_padding * 2;
    };

    onMounted(() => {
        window.addEventListener("resize", resizeHandler);
        resizeHandler();
    });

    onBeforeUnmount(() => {
        window.removeEventListener("resize", resizeHandler);
    });

    return {
        /** layout-content-main内的padding尺寸 */
        style_content_padding,
        /** layout-content-main的内容高度，初始值为100% */
        style_content_height,
        /** n-layout-header  ref */
        headerRef,
        /** n-layout-content  ref */
        contentRef,
        /** 改变contentWidth方法 */
        changeContentWidth,
    };
};
