<template>
  <el-row :gutter="40">
    <el-col v-if="!form.id">
      <el-empty description="请选择左侧菜单后操作" :image-size="100"></el-empty>
    </el-col>
    <template v-else>
      <el-col :lg="12">
        <h2>{{ form.meta.title || "新增菜单" }}</h2>
        <el-form :model="form" :rules="rules" ref="dialogForm" label-width="80px" label-position="left">
          <el-form-item label="显示名称" prop="meta.title">
            <el-input v-model="form.meta.title" clearable placeholder="菜单显示名字"></el-input>
          </el-form-item>
          <el-form-item label="上级菜单" prop="parentId">
            <el-cascader v-model="form.parentId" :options="menuOptions" :props="menuProps"
                         :show-all-levels="false"
                         placeholder="顶级菜单" clearable></el-cascader>
          </el-form-item>
          <el-form-item label="类型" prop="meta.type">
            <el-radio-group v-model="form.meta.type">
              <el-radio-button value="menu">菜单</el-radio-button>
              <el-radio-button value="iframe">Iframe</el-radio-button>
              <el-radio-button value="link">外链</el-radio-button>
              <el-radio-button value="rule">权限</el-radio-button>
            </el-radio-group>
          </el-form-item>
          <el-form-item label="别名" prop="name">
            <el-input v-model="form.name" clearable placeholder="别名"></el-input>
            <div class="el-form-item-msg">
              系统唯一且与内置组件名一致，否则导致缓存失效。如类型为Iframe的菜单，别名将代替源地址显示在地址栏
            </div>
          </el-form-item>
          <el-form-item label="排序" prop="sort">
            <el-input v-model="form.sort" clearable placeholder="菜单排序"></el-input>
            <div class="el-form-item-msg">正序排列，数字越小越靠前</div>
          </el-form-item>
          <el-form-item label="菜单图标" prop="meta.icon">
            <sc-icon-select v-model="form.meta.icon" clearable></sc-icon-select>
          </el-form-item>
          <el-form-item label="路由地址" prop="path">
            <el-input v-model="form.path" clearable placeholder=""></el-input>
          </el-form-item>
          <!--					<el-form-item label="重定向" prop="redirect">-->
          <!--						<el-input v-model="form.redirect" clearable placeholder=""></el-input>-->
          <!--					</el-form-item>-->
          <!--					<el-form-item label="菜单高亮" prop="active">-->
          <!--						<el-input v-model="form.active" clearable placeholder=""></el-input>-->
          <!--						<div class="el-form-item-msg">子节点或详情页需要高亮的上级菜单路由地址</div>-->
          <!--					</el-form-item>-->
          <el-form-item label="视图" prop="component">
            <el-input v-model="form.component" clearable placeholder="">
              <template #prepend>views/</template>
            </el-input>
            <div class="el-form-item-msg">如父节点、链接或Iframe等没有视图的菜单不需要填写</div>
          </el-form-item>
          <!--					<el-form-item label="颜色" prop="color">-->
          <!--						<el-color-picker v-model="form.meta.color" :predefine="predefineColors"></el-color-picker>-->
          <!--					</el-form-item>-->
          <el-form-item label="是否隐藏" prop="meta.hidden">
            <el-checkbox v-model="form.meta.hidden">隐藏菜单</el-checkbox>
            <el-checkbox v-model="form.meta.hiddenBreadcrumb">隐藏面包屑</el-checkbox>
            <div class="el-form-item-msg">菜单不显示在导航中，但用户依然可以访问，例如详情页</div>
          </el-form-item>
          <el-form-item label="整页路由" prop="fullPage">
            <el-switch v-model="form.meta.fullPage"/>
          </el-form-item>
          <!--                    <el-form-item label="标签" prop="tag">-->
          <!--                        <el-input v-model="form.meta.tag" clearable placeholder=""></el-input>-->
          <!--                    </el-form-item>-->
          <el-form-item>
            <el-button type="primary" @click="save" :loading="loading">保 存</el-button>
          </el-form-item>
        </el-form>
      </el-col>
      <el-col :lg="12" class="apilist">
        <h2>接口权限</h2>
        <sc-form-table v-model="form.apiList" :addTemplate="apiListAddTemplate" placeholder="暂无匹配接口权限">
          <el-table-column prop="tag" label="标识" width="150">
            <template #default="scope">
              <el-input v-model="scope.row.tag" placeholder="请输入内容"></el-input>
            </template>
          </el-table-column>
          <el-table-column prop="path" label="Api url">
            <template #default="scope">
              <el-input v-model="scope.row.path" placeholder="请输入内容"></el-input>
            </template>
          </el-table-column>
        </sc-form-table>
      </el-col>
    </template>
  </el-row>

</template>

<script>
import ScIconSelect from '@/components/scIconSelect'
import ScFormTable from '@/components/scFormTable'

export default {
  components: {
    ScIconSelect,
    ScFormTable,
  },
  props: {
    menu: {
      type: Object, default: () => {
      }
    },
  },
  data() {
    return {
      form: {
        id: "",
        app_id: 0,
        parentId: "",
        name: "",
        path: "",
        component: "",
        redirect: "",
        sort: "0",
        meta: {
          title: "",
          icon: "",
          active: "",
          color: "",
          type: "menu",
          fullPage: false,
          tag: "",
          hidden: false,
          hiddenBreadcrumb: false
        },
        apiList: []
      },
      checkPid: 0, //用于比对是否修改上级
      menuOptions: [],
      menuProps: {
        value: 'id',
        label: 'title',
        checkStrictly: true
      },
      rules: {
        required: false //避免表单验证错误
      },
      apiListAddTemplate: {
        tag: "",
        path: ""
      },
      loading: false,
      appId: 0
    }
  },
  watch: {
    menu: {
      handler() {
        this.menuOptions = this.treeToMap(this.menu)
      },
      deep: true
    }
  },
  methods: {
    //简单化菜单
    treeToMap(tree) {
      let map = [];
      tree.forEach(item => {
        const obj = {
          id: item.id,
          parentId: item.parentId,
          title: item.meta.title,
          children: item.children && item.children.length > 0 ? this.treeToMap(item.children) : null
        }
        map.push(obj)
      })
      return map
    },
    //保存
    async save() {
      this.loading = true
      const res = await this.$API.admin.system.menu.edit.post(this.form)
      this.loading = false
      if (res.code === 0) {
        this.$message.success("保存成功")
        if (this.checkPid !== this.form.parentId) {
          this.$emit('refreshMenu')
          this.checkPid = this.form.parentId;
        }
      } else {
        this.$message.warning(res.message)
      }
    },
    //表单注入数据
    setData(data, pid) {
      this.form = data
      this.form.apiList = data.apiList || []
      this.form.parentId = pid
      this.checkPid = pid
    },
    unsetData() {
      this.form.id = 0
    }
  }
}
</script>

<style scoped>
h2 {
  font-size: 17px;
  color: #3c4a54;
  padding: 0 0 30px 0;
}

.apilist {
  border-left: 1px solid #eee;
}

[data-theme="dark"] h2 {
  color: #fff;
}

[data-theme="dark"] .apilist {
  border-color: #434343;
}
</style>
