import config from "@/config"
import http from "@/utils/request"

export default {
    password: {
        url: `${config.API_URL}/admin/user/password`,
        name: "修改密码",
        post: async function (data = {}) {
            return await http.post(this.url, data);
        }
    }
}
