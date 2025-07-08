import config from "@/config";
import http from "@/utils/request";

export default {
  app: {
    add: {
      url: `${config.API_URL}/admin/system/app/add`,
      name: "新增应用",
      post: async function (data = {}) {
        return await http.post(this.url, data);
      }
    },
    list: {
      url: `${config.API_URL}/admin/system/app/list`,
      name: "应用列表",
      get: async function (params = {}) {
        return await http.get(this.url, params);
      }
    },
    edit: {
      url: `${config.API_URL}/admin/system/app/edit`,
      name: "修改应用",
      post: async function (data = {}) {
        return await http.post(this.url, data);
      }
    },
    delete: {
      url: `${config.API_URL}/admin/system/app/delete`,
      name: "删除应用",
      post: async function (data = {}) {
        return await http.post(this.url, data);
      }
    }
  },
  log: {
    typeList: {
      url: `${config.API_URL}/admin/system/log/typeList`,
      name: "日志类型列表",
      get: async function (params) {
        return await http.get(this.url, params);
      }
    },
    list: {
      url: `${config.API_URL}/admin/system/log/list`,
      name: "日志列表",
      get: async function (params) {
        return await http.get(this.url, params);
      }
    },
    logStat: {
      url: `${config.API_URL}/admin/system/log/logStat`,
      name: "日志页详情",
      get: async function (params) {
        return await http.get(this.url, params);
      }
    }
  },
  role: {
    add: {
      url: `${config.API_URL}/admin/system/role/add`,
      name: "新增角色",
      post: async function (data = {}) {
        return await http.post(this.url, data);
      }
    },
    list: {
      url: `${config.API_URL}/admin/system/role/list`,
      name: "角色列表",
      get: async function (params = {}) {
        return await http.get(this.url, params);
      }
    },
    edit: {
      url: `${config.API_URL}/admin/system/role/edit`,
      name: "修改角色",
      post: async function (data = {}) {
        return await http.post(this.url, data);
      }
    },
    editPermission: {
      url: `${config.API_URL}/admin/system/role/edit?type=permission`,
      name: "修改角色权限",
      post: async function (data = {}) {
        return await http.post(this.url, data);
      }
    },
    delete: {
      url: `${config.API_URL}/admin/system/role/delete`,
      name: "删除角色",
      post: async function (data = {}) {
        return await http.post(this.url, data);
      }
    }
  },
  user: {
    add: {
      url: `${config.API_URL}/admin/system/user/add`,
      name: "新增用户",
      post: async function (data = {}) {
        return await http.post(this.url, data);
      }
    },
    list: {
      url: `${config.API_URL}/admin/system/user/list`,
      name: "用户列表",
      get: async function (params = {}) {
        return await http.get(this.url, params);
      }
    },
    edit: {
      url: `${config.API_URL}/admin/system/user/edit`,
      name: "修改用户",
      post: async function (data = {}) {
        return await http.post(this.url, data);
      }
    },
    delete: {
      url: `${config.API_URL}/admin/system/user/delete`,
      name: "删除用户",
      post: async function (data = {}) {
        return await http.post(this.url, data);
      }
    }
  },
  menu: {
    add: {
      url: `${config.API_URL}/admin/system/menu/add`,
      name: "新增菜单",
      post: async function (data = {}) {
        return await http.post(this.url, data);
      }
    },
    list: {
      url: `${config.API_URL}/admin/system/menu/list`,
      name: "菜单列表",
      get: async function (params = {}) {
        return await http.get(this.url, params);
      }
    },
    edit: {
      url: `${config.API_URL}/admin/system/menu/edit`,
      name: "修改菜单",
      post: async function (data = {}) {
        return await http.post(this.url, data);
      }
    },
    delete: {
      url: `${config.API_URL}/admin/system/menu/delete`,
      name: "删除菜单",
      post: async function (data = {}) {
        return await http.post(this.url, data);
      }
    }
  }
}
