<template>
	<div class="sc-file-select">
		<div class="sc-file-select__side" v-loading="state.menuLoading">
			<div class="sc-file-select__side-menu">
				<el-tree ref="group" class="menu" :data="state.menu" :node-key="state.treeProps.key" :props="state.treeProps" :current-node-key="state.menu.length>0?state.menu[0][state.treeProps.key]:''" highlight-current @node-click="groupClick">
					<template #default="{ node }">
						<span class="el-tree-node__label">
							<el-icon class="icon"><el-icon-folder /></el-icon>{{node.label}}
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
					<el-upload class="sc-file-select__upload" action="" multiple :show-file-list="false" :accept="props.accept" :on-change="uploadChange" :before-upload="uploadBefore" :on-progress="uploadProcess" :on-success="uploadSuccess" :on-error="uploadError" :http-request="uploadRequest">
						<el-button type="primary" icon="el-icon-upload">本地上传</el-button>
					</el-upload>
					<span class="tips"><el-icon><el-icon-warning /></el-icon>大小不超过{{props.maxSize}}MB</span>
				</div>
				<div class="keyword">
					<el-input v-model="state.keyword" prefix-icon="el-icon-search" placeholder="文件名搜索" clearable @keyup.enter="search" @clear="search"></el-input>
				</div>
			</div>
			<div class="sc-file-select__list">
				<el-scrollbar ref="scrollbar">
					<el-empty v-if="state.fileList.length==0 && state.data.length==0" description="无数据" :image-size="80"></el-empty>
					<div v-for="(file, index) in state.fileList" :key="index" class="sc-file-select__item">
						<div class="sc-file-select__item__file">
							<div class="sc-file-select__item__upload">
								<el-progress type="circle" :percentage="file.progress" :width="70"></el-progress>
							</div>
							<el-image :src="file.tempImg" fit="contain"></el-image>
						</div>
						<p>{{file.name}}</p>
					</div>
					<div v-for="item in state.data" :key="item[state.fileProps.key]" class="sc-file-select__item" :class="{active: state.value.includes(item[state.fileProps.url]) }" @click="select(item)">
						<div class="sc-file-select__item__file">
							<div class="sc-file-select__item__checkbox" v-if="multiple">
								<el-icon><el-icon-check /></el-icon>
							</div>
							<div class="sc-file-select__item__select" v-else>
								<el-icon><el-icon-check /></el-icon>
							</div>
							<div class="sc-file-select__item__box"></div>
							<el-image v-if="_isImg(item[state.fileProps.url])" :src="tool.resUrl(item[state.fileProps.url])" fit="contain" lazy></el-image>
							<div v-else class="item-file item-file-doc">
								<i v-if="state.files[_getExt(item[state.fileProps.url])]" :class="state.files[_getExt(item[state.fileProps.url])].icon" :style="{color:state.files[_getExt(item[state.fileProps.url])].color}"></i>
								<i v-else class="sc-icon-file-list-fill" style="color: #999;"></i>
							</div>
						</div>
						<p :title="item[state.fileProps.fileName]">{{item[state.fileProps.fileName]}}</p>
					</div>
				</el-scrollbar>
			</div>
			<div class="sc-file-select__pagination">
				<el-pagination size="small" background layout="prev, pager, next" :total="state.total" :page-size="state.pageSize" v-model:currentPage="state.currentPage" @current-change="reload"></el-pagination>
			</div>
			<div class="sc-file-select__do">
				<slot name="do"></slot>
				<el-button type="primary" :disabled="state.value.length<=0" @click="submit">确 定</el-button>
			</div>
		</div>
	</div>
</template>

<script setup>
import config from "@/config/fileSelect"
import tool from "@/utils/tool.js";
import {onMounted, reactive, ref, watch} from "vue";
import {ElNotification} from "element-plus";

const scrollbar = ref(null)

const props = defineProps({
  modelValue: null,
  hideUpload: {type: Boolean, default: false},
  multiple: {type: Boolean, default: false},
  max: {type: Number, default: config.max},
  onlyImage: {type: Boolean, default: false},
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
  value: props.multiple ? [] : '',
  fileList: [],
  accept: props.onlyImage ? "image/gif, image/jpeg, image/png" : "",
  listLoading: false,
  menuLoading: false,
  treeProps: config.menuProps,
  fileProps: config.fileProps,
  files: config.files
})

const emit = defineEmits(['update:modelValue'])

// 监听器
watch(() => props.multiple, () => {
  state.value = props.multiple ? [] : ''
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
  // this.menu = res.data
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
  const itemUrl = item[state.fileProps.url]
  if (props.multiple) {
    if (state.value.includes(itemUrl)) {
      state.value.splice(state.value.findIndex(f => f === itemUrl), 1)
    } else {
      state.value.push(itemUrl)
    }
  } else {
    state.value = itemUrl
  }
  emit('update:modelValue', JSON.parse(JSON.stringify(state.value)));
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
  if (!props.multiple) {
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

<style scoped>
	.sc-file-select {display: flex;}
	.sc-file-select__files {flex: 1;}

	.sc-file-select__list {height:400px;}
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

	.sc-file-select__top {margin-bottom: 15px;display: flex;justify-content: space-between;}
	.sc-file-select__upload {display: inline-block;}
	.sc-file-select__top .tips {font-size: 12px;margin-left: 10px;color: #999;}
	.sc-file-select__top .tips i {font-size: 14px;margin-right: 5px;position: relative;bottom: -0.125em;}
	.sc-file-select__pagination {margin:15px 0;}

	.sc-file-select__do {text-align: right;}
</style>
