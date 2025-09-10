<!--
 * @Descripttion: 处理iframe持久化，涉及store(VUEX)
 * @version: 1.0
 * @Author: sakuya
 * @Date: 2021年6月30日13:20:41
 * @LastEditors:
 * @LastEditTime:
-->

<template>
	<div v-show="$route.meta.type === 'iframe'" class="iframe-pages">
		<iframe v-for="item in iframeList" :key="item.meta.url" v-show="$route.meta.url === item.meta.url" :src="item.meta.url" frameborder='0'></iframe>
	</div>
</template>

<script setup>
import { computed, watch, onMounted } from 'vue'
import { useStore } from 'vuex'
import { useRoute } from 'vue-router'

const store = useStore()
const route = useRoute()

// 计算属性
const iframeList = computed(() => store.state.iframe.iframeList)
const ismobile = computed(() => store.state.global.ismobile)
const layoutTags = computed(() => store.state.global.layoutTags)

// 方法
const push = (routeItem) => {
	if(routeItem.meta.type === 'iframe'){
		if(ismobile.value || !layoutTags.value){
			store.commit("setIframeList", routeItem)
		}else{
			store.commit("pushIframeList", routeItem)
		}
	}
}

// 监听器
watch(route, (newRoute) => {
	push(newRoute)
})

// 生命周期
onMounted(() => {
	push(route)
})
</script>

<style scoped>
.iframe-pages {
	position: relative;
	width: 100%;
	height: 100%;
}

.iframe-pages iframe {
	width: 100%;
	height: 100%;
	border: none;
}
</style>
