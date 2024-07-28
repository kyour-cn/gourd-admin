import config from "@/config"
import http from "@/utils/request"

export default {
	config: {
		list: {
			url: `${config.API_URL}/app/admin/config/list`,
			name: "获取列表",
			get: async function(data={}){
				return await http.get(this.url, data);
			}
		},
		configList: {
			url: `${config.API_URL}/app/admin/config/configList`,
			name: "获取配置json",
			get: async function(data={}){
				return await http.get(this.url, data);
			}
		},
		getConfigValue: {
			url: `${config.API_URL}/app/admin/config/getConfigValue`,
			name: "根据参数获取配置值",
			post: async function(params){
				return await http.post(this.url, params);
			}
		},
		save: {
			url: `${config.API_URL}/app/admin/config/save`,
			name: "新增",
			post: async function(params){
				return await http.post(this.url, params);
			}
		},
		update: {
			url: `${config.API_URL}/app/admin/config/update`,
			name: "更新",
			post: async function(params){
				return await http.post(this.url, params);
			}
		},
		delete: {
			url: `${config.API_URL}/app/admin/config/delete`,
			name: "删除",
			post: async function(params){
				return await http.post(this.url, params);
			}
		},
		saveSystemConfig: {
			url: `${config.API_URL}/app/admin/config/saveSystemConfig`,
			name: "系统设置保存",
			post: async function(params){
				return await http.post(this.url, params);
			}
		},
		clearCache: {
			url: `${config.API_URL}/app/admin/config/clearCache`,
			name: "清除缓存",
			post: async function(params){
				return await http.post(this.url, params);
			}
		},
		changeStatus: {
			url: `${config.API_URL}/app/admin/config/changeStatus`,
			name: "更改状态",
			put: async function(params){
				return await http.put(this.url, params);
			}
		}
	},
}
