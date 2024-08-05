import config from "@/config";
import http from "@/utils/request";

export default {
    list: {
        url: `${config.API_URL}/app/admin/system/logList`,
        name: "日志列表",
        get: async function (params) {
            return await http.get(this.url, params);
        }
    },
    logPageInfo: {
        url: `${config.API_URL}/app/admin/system/logPageInfo`,
        name: "日志页详情",
        get: async function (params) {
            return await http.get(this.url, params);
        }
    }
}
