<template>
    <span :class="'icon saas-icon sa-' + name"></span>
</template>

<script setup>
// 接收 props
const props = defineProps({
    name: String
})
</script>
