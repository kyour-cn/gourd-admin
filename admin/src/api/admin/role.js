import config from "@/config";
import http from "@/utils/request";

export default {
    add: {
        url: `${config.API_URL}/admin/role/add`,
        name: "新增角色",
        post: async function (data = {}) {
            return await http.post(this.url, data);
        }
    },
    list: {
        url: `${config.API_URL}/admin/role/list`,
        name: "角色列表",
        get: async function (params = {}) {
            return await http.get(this.url, params);
        }
    },
    edit: {
        url: `${config.API_URL}/admin/role/edit`,
        name: "修改角色",
        post: async function (data = {}) {
            return await http.post(this.url, data);
        }
    },
    delete: {
        url: `${config.API_URL}/admin/role/delete`,
        name: "删除角色",
        post: async function (data = {}) {
            return await http.post(this.url, data);
        }
    }
}
