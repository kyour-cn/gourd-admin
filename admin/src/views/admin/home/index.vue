<template>
    <div v-if="pageLoading">
        <el-main>
            <el-card shadow="never">
                <el-skeleton :rows="1"/>
            </el-card>
            <el-card shadow="never" style="margin-top: 15px;">
                <el-skeleton/>
            </el-card>
        </el-main>
    </div>
    <work v-if="dashboard==='1'" @on-mounted="onMounted"></work>
    <widgets v-else @on-mounted="onMounted"></widgets>
</template>

<script setup>
import {defineAsyncComponent, ref} from 'vue';
import tool from '@/utils/tool';

defineOptions({
    name: 'dashboard',
})

const work = defineAsyncComponent(() => import('./work/index.vue'));
const widgets = defineAsyncComponent(() => import('./widgets/index.vue'));

const pageLoading = ref(true);
const dashboard = ref(tool.data.get("USER_INFO").dashboard || '0');

const onMounted = () => {
    pageLoading.value = false
}

</script>
