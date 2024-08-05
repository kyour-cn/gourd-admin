import config from "@/config";
import http from "@/utils/request";

export default {
    add: {
        url: `${config.API_URL}/admin/user/add`,
        name: "新增用户",
        post: async function (data = {}) {
            return await http.post(this.url, data);
        }
    },
    list: {
        url: `${config.API_URL}/admin/user/list`,
        name: "用户列表",
        get: async function (params = {}) {
            return await http.get(this.url, params);
        }
    },
    edit: {
        url: `${config.API_URL}/admin/user/edit`,
        name: "修改用户",
        post: async function (data = {}) {
            return await http.post(this.url, data);
        }
    },
    delete: {
        url: `${config.API_URL}/admin/user/delete`,
        name: "删除用户",
        post: async function (data = {}) {
            return await http.post(this.url, data);
        }
    }
}
