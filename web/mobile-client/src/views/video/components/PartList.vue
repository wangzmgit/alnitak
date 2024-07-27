<template>
	<div class="part-list">
		<div class="part-item" v-for="(item, index) in resources" :key="index" @click="changePart(index)">
			<span :class="{ 'active': current - 1 === index }">P{{ index + 1 }}</span>
		</div>
	</div>
</template>

<script setup lang="ts">
import { ref } from "vue";

const emits = defineEmits(['change']);
const props = withDefaults(defineProps<{
	active: number,
	resources: Array<BaseResourceType>
}>(), {
	active: 1,
});


const current = ref(props.active);
const changePart = (part: number) => {
	current.value = part + 1;
	emits('change', part + 1)
}
</script>

<style lang="scss" scoped>
.part-list {
	width: 100%;
	overflow: hidden;
	display: flex;
	align-items: center;
	overflow: auto;
	padding-left: 0;
	list-style: none;
	scrollbar-width: none; //隐藏滚动条(火狐)
	-ms-overflow-style: none; //隐藏滚动条(IE)

	&::-webkit-scrollbar {
		//隐藏滚动条(Chrome)
		display: none;
	}

	.part-item {
		text-align: center;
		font-size: 16px;
		padding: 10px 3px;
		// margin: 3px;

		span {
			padding: 0 10px;
			overflow: hidden;
			height: 36px;
			line-height: 36px;
			width: 100px;
			display: block;
			border-radius: 3px;
			border: 1px solid #808080;
		}

		.active {
			color: var(--primary-color);
			border: 1px solid var(--primary-color);
		}
	}

}
</style>