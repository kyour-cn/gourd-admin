import config from "@/config";
import http from "@/utils/request";

export default {
    typeList: {
        url: `${config.API_URL}/admin/log/typeList`,
        name: "日志类型列表",
        get: async function (params) {
            return await http.get(this.url, params);
        }
    },
    list: {
        url: `${config.API_URL}/admin/log/list`,
        name: "日志列表",
        get: async function (params) {
            return await http.get(this.url, params);
        }
    },
    logStat: {
        url: `${config.API_URL}/admin/log/logStat`,
        name: "日志页详情",
        get: async function (params) {
            return await http.get(this.url, params);
        }
    }
}
