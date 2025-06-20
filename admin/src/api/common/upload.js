import config from '@/config';
import http from '@/utils/request';

export default {
    uploadImage: {
        url: `${config.API_URL}/admin/Common/uploadImage`,
        name: '图片上传',
        post: async function (params) {
            return await http.post(this.url, params);
        }
    }
};
