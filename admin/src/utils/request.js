import axios from 'axios';
import {ElMessageBox, ElNotification} from 'element-plus';
import sysConfig from "@/config";
import tool from '@/utils/tool';
import router from '@/router';

axios.defaults.baseURL = ''

axios.defaults.timeout = sysConfig.TIMEOUT

// HTTP request 拦截器
axios.interceptors.request.use(
  (config) => {
    let token = tool.cookie.get("TOKEN");
    if(token){
      config.headers[sysConfig.TOKEN_NAME] = sysConfig.TOKEN_PREFIX + token
    }
    if(!sysConfig.REQUEST_CACHE && config.method === 'get'){
      config.params = config.params || {};
      config.params['_'] = new Date().getTime();
    }
    Object.assign(config.headers, sysConfig.HEADERS)
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

//FIX 多个API同时401时疯狂弹窗BUG
let MessageBox_401_show = false

// HTTP response 拦截器
axios.interceptors.response.use(
  (response) => {
    // 对响应数据做处理，并返回处理后的数据
    if(sysConfig.DATA_ENCRYPTION && response.data.data){
      const decrypt = tool.crypto.AES.decrypt(response.data.data, sysConfig.DATA_ENCRYPTION_KEY)
      response.data.data = JSON.parse(decrypt)
    }
    return response;
  },
  (error) => {
    if (error.response) {
      if (error.response.status === 403) {
        ElNotification.error({
          title: '无权限访问',
          message: "Status:403，您没有权限访问该资源！"
        });
      } else if (error.response.status === 404) {
        ElNotification.error({
          title: '请求错误',
          message: "Status:404，正在请求不存在的服务器记录！"
        });
      } else if (error.response.status === 500) {
        ElNotification.error({
          title: '请求错误',
          message: error.response.data.message || "Status:500，服务器发生错误！"
        });
      } else if (error.response.status === 401) {
        if(!MessageBox_401_show){
          MessageBox_401_show = true
          ElMessageBox.confirm('当前用户已被登出或无权限访问当前资源，请尝试重新登录后再操作。', '无权限访问', {
            type: 'error',
            closeOnClickModal: false,
            center: true,
            confirmButtonText: '重新登录',
            beforeClose: (action, instance, done) => {
              MessageBox_401_show = false
              done()
            }
          }).then(() => {
            router.replace({path: '/login'}).then();
          }).catch(() => {})
        }
      } else {
        ElNotification.error({
          title: '请求错误',
          message: error.message || `Status:${error.response.status}，未知错误！`
        });
      }
    } else {
      ElNotification.error({
        title: '请求错误',
        message: "请求服务器无响应！"
      });
    }

    return Promise.reject(error.response);
  }
);

const http = {

  /** get 请求
   * @param  {string} url 接口地址
   * @param  {object} params 请求参数
   * @param  {object} config 参数
   */
  get: function (url, params = {}, config = {}) {
    return new Promise((resolve, reject) => {
      axios({
        method: 'get',
        url: url,
        params: params,
        ...config
      }).then((response) => {
        resolve(response.data);
      }).catch((error) => {
        reject(error);
      })
    })
  },

  /** post 请求
   * @param  {string} url 接口地址
   * @param  {object} data 请求参数
   * @param  {object} config 参数
   */
  post: function (url, data = {}, config = {}) {
    return new Promise((resolve, reject) => {
      axios({
        method: 'post',
        url: url,
        data: data,
        ...config
      }).then((response) => {
        resolve(response.data);
      }).catch((error) => {
        reject(error);
      })
    })
  },

  /** put 请求
   * @param  {string} url 接口地址
   * @param  {object} data 请求参数
   * @param  {object} config 参数
   */
  put: function (url, data = {}, config = {}) {
    return new Promise((resolve, reject) => {
      axios({
        method: 'put',
        url: url,
        data: data,
        ...config
      }).then((response) => {
        resolve(response.data);
      }).catch((error) => {
        reject(error);
      })
    })
  },

  /** patch 请求
   * @param  {string} url 接口地址
   * @param  {object} data 请求参数
   * @param  {object} config 参数
   */
  patch: function (url, data = {}, config = {}) {
    return new Promise((resolve, reject) => {
      axios({
        method: 'patch',
        url: url,
        data: data,
        ...config
      }).then((response) => {
        resolve(response.data);
      }).catch((error) => {
        reject(error);
      })
    })
  },

  /** delete 请求
   * @param  {string} url 接口地址
   * @param  {object} data 请求参数
   * @param  {object} config 参数
   */
  delete: function (url, data = {}, config = {}) {
    return new Promise((resolve, reject) => {
      axios({
        method: 'delete',
        url: url,
        data: data,
        ...config
      }).then((response) => {
        resolve(response.data);
      }).catch((error) => {
        reject(error);
      })
    })
  },

  /** jsonp 请求
   * @param  {string} url 接口地址
   * @param  {string} name JSONP回调函数名称
   */
  jsonp: function (url, name = 'jsonp') {
    return new Promise((resolve) => {
      const script = document.createElement('script');
      script.id = `jsonp${Math.ceil(Math.random() * 1000000)}`
      script.type = 'text/javascript'
      script.src = url
      window[name] = (response) => {
        resolve(response)
        document.getElementsByTagName('head')[0].removeChild(script)
        try {
          delete window[name];
        } catch (e) {
          window[name] = undefined;
        }
      }
      document.getElementsByTagName('head')[0].appendChild(script)
    })
  }
};

export default http;
