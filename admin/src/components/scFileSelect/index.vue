<template>
  <div class="sc-file-select">
    <div class="sc-file-select__side" v-loading="state.menuLoading">
      <div class="sc-file-select__side-menu">
        <el-tree ref="group" class="menu" :data="state.menu" :node-key="state.treeProps.key" :props="state.treeProps"
                 :current-node-key="state.menu.length>0?state.menu[0][state.treeProps.key]:''" highlight-current
                 @node-click="groupClick">
          <template #default="{ node }">
            <span class="el-tree-node__label">
              <el-icon class="icon"><el-icon-folder/></el-icon>{{ node.label }}
            </span>
          </template>
        </el-tree>
      </div>
      <div class="sc-file-select__side-msg" v-if="props.multiple">
        已选择 <b>{{state.value.length}}</b> / <b>{{props.max}}</b> 项
      </div>
    </div>
    <div class="sc-file-select__files" v-loading="state.listLoading">
      <div class="sc-file-select__top">
        <div class="upload" v-if="!hideUpload">
          <el-upload class="sc-file-select__upload" action="" multiple :show-file-list="false" :accept="props.accept"
                     :on-change="uploadChange" :before-upload="uploadBefore" :on-progress="uploadProcess"
                     :on-success="uploadSuccess" :on-error="uploadError" :http-request="uploadRequest">
            <el-button type="primary" icon="el-icon-upload">本地上传</el-button>
          </el-upload>
          <span class="tips"><el-icon><el-icon-warning /></el-icon>大小不超过{{props.maxSize}}MB</span>
        </div>
        <div class="keyword">
          <el-input v-model="state.keyword" prefix-icon="el-icon-search" placeholder="文件名搜索" clearable
                    @keyup.enter="search" @clear="search"></el-input>
        </div>
        <div v-if="props.mode === 'manage'" class="multiple-toggle">
          <el-switch
            v-model="multiple"
            active-text="多选"
            inactive-text="单选"
          />
        </div>
      </div>
      <div class="sc-file-select__list">
        <el-scrollbar ref="scrollbar">
          <el-empty v-if="state.fileList.length===0 && state.data.length===0" description="无数据" :image-size="80"/>
          <div v-for="(file, index) in state.fileList" :key="index" class="sc-file-select__item">
            <div class="sc-file-select__item__file">
              <div class="sc-file-select__item__upload">
                <el-progress type="circle" :percentage="file.progress" :width="70"></el-progress>
              </div>
              <el-image :src="file.tempImg" fit="contain"></el-image>
            </div>
            <p>{{file.name}}</p>
          </div>
          <div v-for="item in state.data" :key="item[state.fileProps.key]" class="sc-file-select__item"
               :class="{active: state.value.includes(item) }" @click="select(item)">
            <div class="sc-file-select__item__file">
              <div class="sc-file-select__item__checkbox" v-if="multiple">
                <el-icon><el-icon-check/></el-icon>
              </div>
              <div class="sc-file-select__item__select" v-else>
                <el-icon><el-icon-check/></el-icon>
              </div>
              <div class="sc-file-select__item__box"></div>
              <el-image v-if="_isImg(item[state.fileProps.url])" :src="tool.resUrl(item[state.fileProps.url])"
                        fit="contain" lazy/>
              <div v-else class="item-file item-file-doc">
                <i v-if="state.files[_getExt(item[state.fileProps.url])]"
                   :class="state.files[_getExt(item[state.fileProps.url])].icon"
                   :style="{color:state.files[_getExt(item[state.fileProps.url])].color}"/>
                <i v-else class="sc-icon-file-list-fill" style="color: #999;"/>
              </div>
            </div>
            <p :title="item[state.fileProps.fileName]">{{item[state.fileProps.fileName]}}</p>
          </div>
        </el-scrollbar>
      </div>
      <div class="sc-file-select__pagination">
        <el-pagination size="small" background layout="prev, pager, next" :total="state.total"
                       :page-size="state.pageSize" v-model:currentPage="state.currentPage" @current-change="reload"/>
      </div>
      <div class="sc-file-select__do">
        <slot name="do"></slot>
<!--        <el-button v-if="props.mode === 'manage' && state.value.length === 1" type="primary" @click="rename">-->
<!--          重命名-->
<!--        </el-button>-->
        <el-button v-if="props.mode === 'manage'" type="danger" :disabled="state.value.length<=0" @click="deleteFile">
          删 除
        </el-button>
        <el-button v-if="props.mode === 'select'" type="primary" :disabled="state.value.length<=0" @click="submit">
          确 定
        </el-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import config from "@/config/fileSelect"
import tool from "@/utils/tool.js";
import {onMounted, reactive, ref, watch, computed} from "vue";
import {ElMessageBox, ElNotification} from "element-plus";

const scrollbar = ref(null)

const props = defineProps({
  modelValue: null,
  mode: {type: String, default: 'select'}, // 模式 select 选择器, manage 管理模式
  hideUpload: {type: Boolean, default: false},
  multiple: {type: Boolean, default: false},
  max: {type: Number, default: config.max},
  onlyImage: {type: Boolean, default: false},
  maxSize: {type: Number, default: config.maxSize},
})

const state = reactive({
  keyword: null,
  pageSize: 20,
  total: 0,
  currentPage: 1,
  data: [],
  menu: [{
    id: 0,
    name: '全部'
  }],
  menuId: 0,
  value: [],
  fileList: [],
  accept: props.onlyImage ? "image/gif, image/jpeg, image/png" : "",
  listLoading: false,
  menuLoading: false,
  treeProps: config.menuProps,
  fileProps: config.fileProps,
  files: config.files
})

const emit = defineEmits(['update:modelValue', 'update:multiple'])

// 计算属性 - multiple的双向绑定
const multiple = computed({
  get: () => props.multiple,
  set: (val) => {
    emit('update:multiple', val)
  }
})

// 监听器
watch(() => props.multiple, (newVal) => {
  state.value = newVal ? [] : ''
  emit('update:modelValue', JSON.parse(JSON.stringify(state.value)));
})

// 生命周期
onMounted(() => {
  getMenu()
  getData()
})

const getMenu = async () => {
  state.menuLoading = true
  const res = await config.menuApiObj.get();

  // 保留第一个，追加接口返回数据
  state.menu.splice(1, state.menu.length - 1)
  res.data.forEach(item => {
    state.menu.push(item)
  })
  state.menuLoading = false
}

const getData = async () => {
  state.listLoading = true
  const reqData = {
    [config.request.menuKey]: state.menuId,
    [config.request.page]: state.currentPage,
    [config.request.pageSize]: state.pageSize,
    [config.request.keyword]: state.keyword
  };
  if (props.onlyImage) {
    reqData.type = 'image'
  }
  const res = await config.listApiObj.get(reqData);
  const parseData = config.listParseData(res);
  state.data = parseData.rows
  state.total = parseData.total
  state.listLoading = false
  scrollbar.value?.setScrollTop(0)
}

const groupClick = (data) => {
  state.menuId = data.id
  state.currentPage = 1
  state.keyword = null
  getData()
}

const reload = () => {
  getData()
}

const search = () => {
  state.currentPage = 1
  getData()
}
const select = (item) => {
  if (multiple.value) {
    // 判断id是否存在
    if (state.value.includes(item)) {
      state.value.splice(state.value.findIndex(f => f === item), 1)
    } else {
      state.value.push(item)
    }
  } else {
    state.value = [item]
  }
  emit('update:modelValue', JSON.parse(JSON.stringify(state.value)));
}

const deleteFile = async () =>  {
  try {
    await ElMessageBox.confirm(`确定删除选中的 ${state.value.length} 项吗？`, '提示', {
      type: 'warning',
      confirmButtonText: '删除',
      confirmButtonClass: 'el-button--danger'
    })
  }catch (e) {
    return
  }

  // 调用删除接口
  const res = await config.deleteApiObj.post({
    ids: state.value.map(f => f[state.fileProps.key])
  })
  if (res.code !== config.successCode) {
    ElNotification.error(res.message || '删除失败')
    return
  }
  ElNotification.success(res.message || '删除成功')

  // 删除选中的文件
  for (const item of state.value) {
    // 从文件列表中删除
    const index = state.data.findIndex(f => f[state.fileProps.key] === item[state.fileProps.key])
    if (index !== -1) {
      state.data.splice(index, 1)
    }

    // 从上传列表中删除
    const uploadIndex = state.fileList.findIndex(f => f.uid === item.uid)
    if (uploadIndex !== -1) {
      state.fileList.splice(uploadIndex, 1)
    }
  }

  // 清空选择状态
  state.value = []
  emit('update:modelValue', JSON.parse(JSON.stringify(state.value)))

  // 更新总数
  state.total = state.total - state.value.length
}

const rename = () => {
  // TODO 重命名
}

const submit = () => {
  const value = JSON.parse(JSON.stringify(state.value))
  emit('update:modelValue', value);
  emit('submit', value);
}

const uploadChange = (file, fileList) => {
  file.tempImg = URL.createObjectURL(file.raw);
  state.fileList = fileList
}

const uploadBefore = (file) => {
  const maxSize = file.size / 1024 / 1024 < state.maxSize;
  if (!maxSize) {
    ElNotification.warning(`上传文件大小不能超过 ${state.maxSize}MB!`);
    return false;
  }
}

const uploadRequest = (param) => {
  const apiObj = config.apiObj;
  const data = new FormData();
  data.append("file", param.file);
  data.append([config.request.menuKey], this.menuId);
  apiObj.post(data, {
    onUploadProgress: e => {
      param.onProgress(e)
    }
  }).then(res => {
    param.onSuccess(res)
  }).catch(err => {
    param.onError(err)
  })
}

const uploadProcess = (event, file) => {
  file.progress = Number((event.loaded / event.total * 100).toFixed(2))
}

const uploadSuccess = (res, file) => {
  state.fileList.splice(state.fileList.findIndex(f => f.uid === file.uid), 1)
  const response = config.uploadParseData(res);
  state.data.unshift({
    [state.fileProps.key]: response.id,
    [state.fileProps.fileName]: response.fileName,
    [state.fileProps.url]: response.url
  })
  if (!multiple.value) {
    state.value = response.url
  }
}

const uploadError = (err) => {
  ElNotification.error({
    title: '上传文件错误',
    message: err
  })
}

const _isImg = (fileUrl) => {
  const imgExt = ['.jpg', '.jpeg', '.png', '.gif', '.bmp']
  const fileExt = fileUrl.substring(fileUrl.lastIndexOf("."))
  return imgExt.indexOf(fileExt) !== -1
}
const _getExt = (fileUrl) => {
  return fileUrl.substring(fileUrl.lastIndexOf(".") + 1)
}

</script>

<style scoped lang="scss">
  .sc-file-select {display: flex;}
  .sc-file-select__files {flex: 1;}

  //.sc-file-select__list {height:500px;}
  .sc-file-select__item {display: inline-block;float: left;margin:0 15px 25px 0;width:110px;cursor: pointer;}
  .sc-file-select__item__file {width:110px;height:110px;position: relative;}
  .sc-file-select__item__file .el-image {width:110px;height:110px;}
  .sc-file-select__item__box {position: absolute;top:0;right:0;bottom:0;left:0;border: 2px solid var(--el-color-success);z-index: 1;display: none;}
  .sc-file-select__item__box::before {content: '';position: absolute;top:0;right:0;bottom:0;left:0;background: var(--el-color-success);opacity: 0.2;display: none;}
  .sc-file-select__item:hover .sc-file-select__item__box {display: block;}
  .sc-file-select__item.active .sc-file-select__item__box {display: block;}
  .sc-file-select__item.active .sc-file-select__item__box::before {display: block;}
  .sc-file-select__item p {margin-top: 10px;white-space:nowrap;text-overflow:ellipsis;overflow:hidden;-webkit-text-overflow:ellipsis;text-align: center;}
  .sc-file-select__item__checkbox {position: absolute;width: 20px;height: 20px;top:7px;right:7px;z-index: 2;background: rgba(0,0,0,0.2);border: 1px solid #fff;display: flex;flex-direction: column;align-items: center;justify-content: center;}
  .sc-file-select__item__checkbox i {font-size: 14px;color: #fff;font-weight: bold;display: none;}
  .sc-file-select__item__select {position: absolute;width: 20px;height: 20px;top:0px;right:0px;z-index: 2;background: var(--el-color-success);display: none;flex-direction: column;align-items: center;justify-content: center;}
  .sc-file-select__item__select i {font-size: 14px;color: #fff;font-weight: bold;}
  .sc-file-select__item.active .sc-file-select__item__checkbox {background: var(--el-color-success);}
  .sc-file-select__item.active .sc-file-select__item__checkbox i {display: block;}
  .sc-file-select__item.active .sc-file-select__item__select {display: flex;}
  .sc-file-select__item__file .item-file {width:110px;height:110px;display: flex;flex-direction: column;align-items: center;justify-content: center;}
  .sc-file-select__item__file .item-file i {font-size: 40px;}
  .sc-file-select__item__file .item-file.item-file-doc {color: #409eff;}

  .sc-file-select__item__upload {position: absolute;top:0;right:0;bottom:0;left:0;z-index: 1;background: rgba(255,255,255,0.7);display: flex;flex-direction: column;align-items: center;justify-content: center;}

  .sc-file-select__side {width: 200px;margin-right: 15px;border-right: 1px solid rgba(128,128,128,0.2);display: flex;flex-flow: column;}
  .sc-file-select__side-menu {flex: 1;}
  .sc-file-select__side-msg {height:32px;line-height: 32px;}

  .sc-file-select__top {margin-bottom: 15px;display: flex;justify-content: space-between;align-items: center;}
  .sc-file-select__upload {display: inline-block;}
  .sc-file-select__top .tips {font-size: 12px;margin-left: 10px;color: #999;}
  .sc-file-select__top .tips i {font-size: 14px;margin-right: 5px;position: relative;bottom: -0.125em;}
  .sc-file-select__top .multiple-toggle {margin-left: 20px;}
  .sc-file-select__pagination {margin:15px 0;}

  .sc-file-select__do {text-align: right;}
</style>
