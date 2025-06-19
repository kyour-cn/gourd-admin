<template>
	<el-dialog :title="state.titleMap[state.mode]" v-model="state.visible" :width="500" destroy-on-close @closed="$emit('closed')">
		<el-form :model="state.form" :rules="state.rules" :disabled="state.mode==='show'" ref="dialogForm" label-width="100px"
                 label-position="left">
			<!--			<el-form-item label="头像" prop="avatar">-->
			<!--				<sc-upload v-model="form.avatar" title="上传头像"></sc-upload>-->
			<!--			</el-form-item>-->
			<el-form-item label="登录账号" prop="username">
				<el-input v-model="state.form.username" placeholder="用于登录系统" clearable></el-input>
			</el-form-item>
			<el-form-item label="昵称" prop="nickname">
				<el-input v-model="state.form.nickname" placeholder="请输入完整的真实姓名" clearable></el-input>
			</el-form-item>
			<el-form-item label="手机号" prop="mobile">
				<el-input v-model="state.form.mobile" placeholder="请输入手机号" maxlength="11" clearable></el-input>
			</el-form-item>
			<el-form-item label="是否有效" prop="status">
				<el-switch v-model="state.form.status" :active-value="true" :inactive-value="false"></el-switch>
			</el-form-item>
			<template v-if="state.mode==='add'">
				<el-form-item label="登录密码" prop="password">
					<el-input type="password" v-model="state.form.password" clearable show-password></el-input>
				</el-form-item>
			</template>
			<template v-if="state.mode==='edit'">
				<el-form-item label="修改密码">
					<el-input
                        type="password"
                        clearable
                        show-password
                        placeholder="请输入新密码，留空则不修改"
                        v-model="state.form.password"
                        autocomplete="new-password"
                    />
				</el-form-item>
			</template>
			<el-form-item label="用户角色" prop="role_name">
				<roleSelect
                    v-if="state.form.id"
                    :ids="state.form.role_id"
                    :placeholder="state.form.role?.name"
                    @onChange="change"
                />
			</el-form-item>

		</el-form>
		<template #footer>
			<el-button @click="state.visible=false">取 消</el-button>
			<el-button v-if="state.mode!=='show'" type="primary" :loading="state.isSaveing" @click="submit()">保 存</el-button>
		</template>
	</el-dialog>
</template>

<script setup>
import {reactive, ref} from 'vue'
import {ElAlert, ElMessage} from 'element-plus'
import roleSelect from "@/components/admin/roleSelect.vue"
import systemApi from "@/api/admin/system.js";

const emit = defineEmits(['success', 'closed', 'reloadData'])
const dialogForm = ref()

const state = reactive({
	mode: 'add',
	visible: false,
	isSaveing: false,
	titleMap: {
		add: '新增用户',
		edit: '编辑用户',
		show: '查看'
	},
	form: {
		id: "",
		username: "",
		nickname: "",
		status: true,
		password: '',
		mobile: '',
		role_id: 0,
		role: null
	},
	rules: {
		avatar: [
			{ required: true, message: '请上传头像' }
		],
		username: [
			{ required: true, message: '请输入登录账号' }
		],
		nickname: [
			{ required: true, message: '请输入昵称' }
		],
		mobile: [
			{
				validator: (_, value, callback) => {
					if (value === '') return callback()
					const regMobile = /^1\d{10}$/
					if (regMobile.test(value)) return callback()
					callback(new Error('请输入合法的手机号'))
				},
				trigger: 'blur'
			}
		],
		password: [
			{ required: true, message: '请输入登录密码' },
		],
		status: [
			{ required: true, message: '请选择当前状态' }
		]
	},
	groups: [],
	groupsProps: {
		value: "id",
		multiple: true,
		checkStrictly: true
	}
})

const change = (val) => {
    state.form.role_id = val.map(item => item.id).join(",")
}

const open = (mode = 'add') => {
	state.mode = mode
	state.visible = true
	return this
}

const submit = async () => {
	dialogForm.value.validate(async (valid) => {
		if (valid) {
			state.isSaveing = true
			let { username, nickname, status, password, mobile, role_id, id } = state.form

			status = status ? 1 : 0
			let res
			if (state.mode === 'add') {
				res = await systemApi.user.add.post({
					username,
					nickname,
					status,
					password,
					mobile,
					role_id
				})
			} else {
				res = await systemApi.user.edit.post({
					id,
					username,
					nickname,
					status,
					password,
					mobile,
					role_id
				})
			}
			state.isSaveing = false
			if (res.code === 0) {
				emit('success')
				state.visible = false
				ElMessage.success("操作成功")
			} else {
				ElAlert(res.message, "提示", { type: 'error' })
			}
		} else {
			return false
		}
	})
}

const setData = (data) => {
	Object.assign(state.form, {
		id: data.id,
		nickname: data.nickname,
		avatar: data.avatar,
		username: data.username,
		status: data.status === 1,
		mobile: data.mobile,
		role_id: data.role_id,
		role: data.role
	})
}

//暴露给父组件的方法
defineExpose({
	setData,
	open
})
</script>

<style>
</style>
