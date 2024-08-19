import config from "@/config";
import http from "@/utils/request";

export default {
    levelList: {
        url: `${config.API_URL}/admin/log/levelList`,
        name: "日志级别列表",
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
    logPageInfo: {
        url: `${config.API_URL}/admin/log/logStat`,
        name: "日志页详情",
        get: async function (params) {
            return await http.get(this.url, params);
        }
    }
}
