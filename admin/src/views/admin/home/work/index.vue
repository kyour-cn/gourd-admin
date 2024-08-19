<template>
	<el-main>
		<el-alert title="根据角色配置,可让不同角色访问不同的控制台视图,参数值在登录成功后返回 dashboard:{type}" type="success" style="margin-bottom:20px;"></el-alert>
		<el-row :gutter="15">
			<el-col :lg="24">
				<el-card shadow="never" header="我的常用">
					<myapp/>
				</el-card>
			</el-col>
		</el-row>
	</el-main>
</template>

<script setup>
	import myapp from './components/myapp.vue';
    import {onMounted} from "vue";

    const emits = defineEmits(['on-mounted'])

    onMounted(() => {
        emits('on-mounted')
    })
</script>
