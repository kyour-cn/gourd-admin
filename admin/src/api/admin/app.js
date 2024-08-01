import config from "@/config";
import http from "@/utils/request";

export default {
    add: {
        url: `${config.API_URL}/admin/app/add`,
        name: "新增应用",
        post: async function (data = {}) {
            return await http.post(this.url, data);
        }
    },
    list: {
        url: `${config.API_URL}/admin/app/list`,
        name: "应用列表",
        get: async function (params = {}) {
            return await http.get(this.url, params);
        }
    },
    edit: {
        url: `${config.API_URL}/admin/app/edit`,
        name: "修改应用",
        post: async function (data = {}) {
            return await http.post(this.url, data);
        }
    },
    delete: {
        url: `${config.API_URL}/admin/app/delete`,
        name: "删除应用",
        post: async function (data = {}) {
            return await http.post(this.url, data);
        }
    }
}