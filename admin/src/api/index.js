/**
 * @description 自动import导入所有 api 模块
 */

// 后台模块
let files = import.meta.glob('./admin/*.js', { eager: true });
const admin = {}
Object.keys(files).forEach(key => {
	admin[key.replace(/^\.\/admin\/(.*)\.js$/g, '$1')] = files[key].default
})

// 公共模块
files = import.meta.glob('./common/*.js', { eager: true });
const common = {}
Object.keys(files).forEach(key => {
	common[key.replace(/^\.\/common\/(.*)\.js$/g, '$1')] = files[key].default
})

export default {
	admin,
	common,
}
