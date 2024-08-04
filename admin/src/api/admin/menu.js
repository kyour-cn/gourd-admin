import config from "@/config";
import http from "@/utils/request";

export default {
    add: {
        url: `${config.API_URL}/admin/menu/add`,
        name: "新增菜单",
        post: async function (data = {}) {
            return await http.post(this.url, data);
        }
    },
    list: {
        url: `${config.API_URL}/admin/menu/list`,
        name: "菜单列表",
        get: async function (params = {}) {
            return await http.get(this.url, params);
        }
    },
    edit: {
        url: `${config.API_URL}/admin/menu/edit`,
        name: "修改菜单",
        post: async function (data = {}) {
            return await http.post(this.url, data);
        }
    },
    delete: {
        url: `${config.API_URL}/admin/menu/delete`,
        name: "删除菜单",
        post: async function (data = {}) {
            return await http.post(this.url, data);
        }
    }
}
