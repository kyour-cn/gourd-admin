<template>
	<el-form ref="loginForm" :model="state.form" :rules="state.rules" label-width="0" size="large" @keyup.enter="onVerify">
		<el-form-item prop="user">
			<el-input v-model="state.form.user" prefix-icon="el-icon-user" clearable :placeholder="$t('login.userPlaceholder')">
				<!--				<template #append>-->
				<!--					<el-select v-model="userType" style="width: 130px;">-->
				<!--						<el-option :label="$t('login.admin')" value="admin"></el-option>-->
				<!--						<el-option :label="$t('login.user')" value="user"></el-option>-->
				<!--					</el-select>-->
				<!--				</template>-->
			</el-input>
		</el-form-item>

		<el-form-item prop="password">
			<el-input v-model="state.form.password" prefix-icon="el-icon-lock" clearable show-password
								:placeholder="$t('login.PWPlaceholder')"></el-input>
		</el-form-item>

		<el-form-item style="margin-bottom: 10px;">
			<el-col :span="12">
				<el-checkbox :label="$t('login.rememberMe')" v-model="state.form.autologin"></el-checkbox>
			</el-col>
			<el-col :span="12" class="login-forgot">
				<router-link to="/reset_password">{{ $t('login.forgetPassword') }}？</router-link>
			</el-col>
		</el-form-item>
		<el-form-item>
            <el-popover :visible="state.captchaShow" placement="top-start" width="350">
                <CaptchaSlide
                    v-if="state.captchaShow"
                    :data="state.captchaData"
                    :events="{
                        close: closeCaptcha,
                        refresh: refreshCaptcha,
                        confirm: confirmEvent,
                    }"
                />
                <template #reference>
                    <el-button type="primary" style="width: 100%;" :loading="islogin" round @click="onVerify">
                        {{ $t('login.signIn') }}
                    </el-button>
                </template>
            </el-popover>
		</el-form-item>
		<div class="login-reg">
			{{ $t('login.noAccount') }}
			<router-link to="/user_register">{{ $t('login.createAccount') }}</router-link>
		</div>
	</el-form>
</template>

<script setup>
import config from "@/config"
import {Slide as CaptchaSlide} from 'go-captcha-vue'
import {getCurrentInstance, reactive, ref} from "vue";
import authApi from "@/api/common/auth.js"

const proxy = getCurrentInstance().proxy

const loginForm = ref(null)

const state = reactive({
    form: {
        user: "",
        password: "",
        autologin: false,
    },
    rules: {
        user: [
            {required: true, message: proxy.$t('login.userError'), trigger: 'blur'}
        ],
        password: [
            {required: true, message: proxy.$t('login.PWError'), trigger: 'blur'}
        ]
    },
    islogin: false,
    captchaShow: false,
    captchaConfig: null,
    captchaData: null
})

const onVerify = () => {
    state.captchaShow = false
    authApi.captcha.get().then(res => {
        if(res.code === 0) {

            state.captchaData = {
                image: res.data.image_base64,
                thumb: res.data.tile_base64,
                captKey: res.data.captcha_key,
                thumbX: res.data.tile_x,
                thumbY: res.data.tile_y,
                thumbWidth: res.data.tile_width,
                thumbHeight: res.data.tile_height,
            }

            state.captchaShow = true
        }else{
            proxy.$message.error(res.message)
        }
    })
}
const refreshCaptcha = () => {
    onVerify()
}
const closeCaptcha = () => {
    state.captchaShow = false
}
const confirmEvent = async(point, clear) => {
    closeCaptcha()

    const validate = await loginForm.value.validate().catch()
    if (!validate) {
        return false
    }

    state.islogin = true
    const data = {
        username: state.form.user,
        password: proxy.$TOOL.crypto.MD5(state.form.password),
        md5: true,
        point: point,
        captcha_key: state.captchaData.captKey
    };
    //获取token
    const user = await authApi.login.post(data);
    if (user.code === 0) {
        proxy.$TOOL.cookie.set("TOKEN", user.data.token, {
            expires: state.form.autologin ? user.data.expire : 0
        })
        proxy.$TOOL.data.set("USER_INFO", user.data.userInfo)
    } else {
        if(user.code === 102) {
            refreshCaptcha()
        }
        state.islogin = false
        proxy.$message.warning(user.message)
        return false
    }
    //获取菜单
    const res = await authApi.menu.get();
    if (res.code === 0) {
        if (res.data.menu.length === 0) {
            state.islogin = false
            await proxy.$alert("当前用户无任何菜单权限，请联系系统管理员", "无权限访问", {
                type: 'error',
                center: true
            })
            return false
        }
        proxy.$TOOL.data.set("MENU", res.data.menu)
        proxy.$TOOL.data.set("PERMISSIONS", res.data.permissions)
    } else {
        state.islogin = false
        proxy.$message.warning(res.message)
        return false
    }
    //默认路由地址
    const defaultRoute = config.DASHBOARD_URL;

    //递归菜单及子菜单，判断是否存在默认路由
    const findDefaultRoute = (menu) => {
        for (let i = 0; i < menu.length; i++) {
            if (menu[i].path === defaultRoute) {
                return true
            }
            if (menu[i].children && menu[i].children.length) {
                return findDefaultRoute(menu[i].children)
            }
        }
        return false
    }
    //不存在默认路由，跳转到第一个菜单
    if (!findDefaultRoute(res.data.menu)) {
        //取第一个菜单
        let menuItem = res.data.menu[0];
        if (menuItem.children?.length) {
            menuItem = menuItem.children[0]
        }
        proxy.$router.replace({
            path: menuItem.path
        })
    } else {
        proxy.$router.replace({
            path: defaultRoute
        })
    }

    proxy.$message.success("登录成功")
    state.islogin = false
}

</script>

<style>
</style>
