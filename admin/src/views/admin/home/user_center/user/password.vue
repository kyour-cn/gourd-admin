<template>
	<el-card shadow="never" header="修改密码">
		<el-alert title="密码更新成功后，您将被重定向到登录页面，您可以使用新密码重新登录。" type="info" show-icon style="margin-bottom: 15px;"/>
		<el-form ref="formRef" :model="form" :rules="rules" label-width="120px" style="margin-top:20px;">
			<el-form-item label="当前密码" prop="userPassword">
				<el-input v-model="form.userPassword" type="password" show-password placeholder="请输入当前密码"></el-input>
				<div class="el-form-item-msg">必须提供当前登录用户密码才能进行更改</div>
			</el-form-item>
			<el-form-item label="新密码" prop="newPassword">
				<el-input v-model="form.newPassword" type="password" show-password placeholder="请输入新密码"></el-input>
				<sc-password-strength v-model="form.newPassword"></sc-password-strength>
				<div class="el-form-item-msg">请输入包含英文、数字的8位以上密码</div>
			</el-form-item>
			<el-form-item label="确认新密码" prop="confirmNewPassword">
				<el-input v-model="form.confirmNewPassword" type="password" show-password placeholder="请再次输入新密码"></el-input>
			</el-form-item>
			<el-form-item>
				<el-button type="primary" @click="save">保存密码</el-button>
			</el-form-item>
		</el-form>
	</el-card>
</template>

<script setup>
import scPasswordStrength from '@/components/scPasswordStrength'
import {ref} from "vue";
import {ElMessage, ElMessageBox} from "element-plus";
import userApi from "@/api/common/user.js";
import router from "@/router/index.js";

const formRef = ref(null);
const form = ref({
    userPassword: "",
    newPassword: "",
    confirmNewPassword: ""
})

const rules = {
    userPassword: [
        {required: true, message: '请输入当前密码'}
    ],
    newPassword: [
        {required: true, message: '请输入新密码'}
    ],
    confirmNewPassword: [
        {required: true, message: '请再次输入新密码'},
        {
            validator: (rule, value, callback) => {
                if (value !== form.value.newPassword) {
                    callback(new Error('两次输入密码不一致'));
                } else {
                    callback();
                }
            }
        }
    ]
}

const save = () => {
    formRef.value.validate(async valid => {
        if (valid) {
            const res = await userApi.password.post({
                user_password: form.value.userPassword,
                new_password: form.value.newPassword,
                confirm_new_password: form.value.confirmNewPassword
            })
            if (res.code !== 0) {
                ElMessage({
                    type: 'error',
                    message: res.message || '密码修改失败，请稍后再试'
                })
                return false
            }

            ElMessageBox.alert("密码修改成功，即将跳转至登录页使用新密码登录", "修改成功", {
                type: 'success',
                center: true
            }).then(() => {
                router.replace({
                    path: '/login'
                })
            }).catch(() => {
            })
        } else {
            return false
        }
    })
}

</script>
