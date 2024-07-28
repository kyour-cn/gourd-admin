/**
 * @description 自动import导入所有 api 模块
 */

// 系统模块
let files = import.meta.glob('./system/*.js', { eager: true });
const system = {}
Object.keys(files).forEach(key => {
	system[key.replace(/^\.\/system\/(.*)\.js$/g, '$1')] = files[key].default
})

// 系统模块
files = import.meta.glob('./common/*.js', { eager: true });
const common = {}
Object.keys(files).forEach(key => {
	common[key.replace(/^\.\/common\/(.*)\.js$/g, '$1')] = files[key].default
})

export default {
	system,
	common,
}
