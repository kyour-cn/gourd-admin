<template>
  <el-container>
    <el-header>
      <div class="left-panel">
        <el-button type="primary" icon="el-icon-plus" @click="add"/>
        <el-button type="danger" plain icon="el-icon-delete" :disabled="!state.selection.length"
                   @click="batchDel"/>
      </div>
      <div class="right-panel">
        <div class="right-panel-search">
          <el-input v-model="state.search.keyword" placeholder="登录账号 / 昵称 / 手机号" clearable
                    @clear="clearSearch"/>
          <el-button type="primary" icon="el-icon-search" @click="upSearch"/>
        </div>
      </div>
    </el-header>
    <el-main class="nopadding">
      <sc-table
        ref="table"
        :apiObj="apiObj"
        :params="state.tableParams"
        row-key="id"
        @selection-change="selectionChange"
        stripe
      >
        <el-table-column type="selection" width="50"/>
        <el-table-column label="ID" prop="id" width="80" sortable='custom'/>
        <el-table-column label="头像" width="80" column-key="filterAvatar">
          <template #default="scope">
            <el-avatar :src="tool.resUrl(scope.row.avatar)" size="small"></el-avatar>
          </template>
        </el-table-column>
        <el-table-column label="登录账号" prop="username" width="150" column-key="filterUserName"/>
        <el-table-column label="昵称" prop="nickname" width="150"/>
        <el-table-column label="所属角色" prop="role_id" width="200"/>
        <el-table-column label="注册时间" prop="create_time" width="170">
          <template #default="{row}">
            {{ $TOOL.dateFormat(row.create_time * 1000) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" fixed="right" align="right" width="300">
          <template #default="scope">
            <el-button-group>
              <el-button text plain type="success" size="small" @click="tableShow(scope.row)">查看</el-button>
              <el-button text plain type="primary" size="small" @click="tableEdit(scope.row)">编辑</el-button>
              <el-popconfirm title="确定删除吗？" @confirm="tableDel(scope.row)">
                <template #reference>
                  <el-button text plain type="danger" size="small">删除</el-button>
                </template>
              </el-popconfirm>
            </el-button-group>
          </template>
        </el-table-column>
      </sc-table>
    </el-main>
  </el-container>

  <save-dialog
    v-if="dialog.save"
    ref="saveDialogRef"
    @success="handleSaveSuccess"
    @closed="dialog.save=false"
  />

</template>

<script setup>

import {nextTick, reactive, ref} from "vue"
import SaveDialog from './save'
import ScTable from "@/components/scTable/index.vue"
import systemApi from "@/api/admin/system.js";
import {ElMessage, ElMessageBox} from "element-plus";
import tool from "@/utils/tool.js";

defineOptions({
  name: 'product_list',
})

const saveDialogRef = ref(null)
const table = ref(null)

const state = reactive({
  selection: [],
  search: {
    keyword: null
  },
  tableParams: {
    keyword: null
  }
})

const apiObj = systemApi.user.list

const dialog = reactive({
  save: false
})

const selectionChange = (val) => {
  state.selection = val
}

const tableShow = (row) => {
  dialog.save = true
  nextTick(() => {
    saveDialogRef.value.open('show')
    saveDialogRef.value.setData(row)
  })
}

//添加
const add = () => {
  dialog.save = true
  nextTick(() => {
    saveDialogRef.value.open()
  })
}

const tableEdit = (row) => {
  dialog.save = true
  nextTick(() => {
    saveDialogRef.value.open('edit')
    saveDialogRef.value.setData(row)
  })
}

//删除
const tableDel = async (row) => {
  const res = await systemApi.user.delete.post({
    ids: [row.id]
  })
  if (res.code === 0) {
    ElMessage.success("操作成功");
    table.value.upData()
  }
}

//批量删除
const batchDel = async () => {
  const confirmRes = await ElMessageBox.confirm(`确定删除选中的 ${state.selection.length} 项吗？`, '提示', {
    type: 'warning',
    confirmButtonText: '删除',
    confirmButtonClass: 'el-button--danger'
  })
  if (!confirmRes) return false

  const ids = state.selection.map(v => v.id)
  const res = await systemApi.user.delete.post({ids})
  if (res.code === 0) {
    table.value.removeKeys(ids)
    ElMessage.success("操作成功");
  } else {
    await ElMessageBox.alert(res.message, "提示", {type: 'error'});
  }
}

//搜索
const upSearch = () => {
  table.value.upData({
    keyword: state.search.keyword
  }, 1)
}

// 删除搜索
const clearSearch = () => {
  table.value.reload({
    keyword: ''
  }, 1)
}

//本地更新数据
const handleSaveSuccess = () => {
  table.value.refresh()
}

</script>
