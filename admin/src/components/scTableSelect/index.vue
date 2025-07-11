<!--
 * @Descripttion: 表格选择器组件
 * @version: 1.3
 * @Author: sakuya
 * @Date: 2021年6月10日10:04:07
 * @LastEditors: sakuya
 * @LastEditTime: 2022年6月6日21:50:36
-->

<template>
	<el-select
        ref="select"
        v-model="defaultValue"
        :size="size"
        :clearable="clearable"
        :multiple="multiple"
        :collapse-tags="collapseTags"
        :max-collapse-tags="3"
        :collapse-tags-tooltip="collapseTagsTooltip"
        :filterable="filterable"
        :placeholder="placeholder"
        :disabled="disabled"
        label="name"
        :filter-method="filterMethod"
        @remove-tag="removeTag"
        @visible-change="visibleChange"
        @clear="clear"
    >
		<template #empty>
			<div class="sc-table-select__table" :style="{width: tableWidth+'px'}" v-loading="loading">
				<div class="sc-table-select__header">
					<slot name="header" :form="formData" :submit="formSubmit"></slot>
				</div>
				<el-table ref="table" :data="tableData" :height="245" :highlight-current-row="!multiple" @row-click="click" @select="select" @select-all="selectAll">
					<el-table-column v-if="multiple" type="selection" width="45"></el-table-column>
					<el-table-column v-else type="index" width="45">
						<template #default="scope"><span>{{scope.$index+(currentPage - 1) * pageSize + 1}}</span></template>
					</el-table-column>
					<slot></slot>
				</el-table>
				<div class="sc-table-select__page">
					<el-pagination size="small" background layout="prev, pager, next" :total="total" :page-size="pageSize" v-model:currentPage="currentPage" @current-change="reload"></el-pagination>
				</div>
			</div>
		</template>
        <template v-if="!multiple" #label="{ value }">
            {{value[defaultProps.label]}}
        </template>
	</el-select>
</template>

<script>
	import config from "@/config/tableSelect";

	export default {
		props: {
			modelValue: null,
			apiObj: { type: Object, default: () => {} },
			params: { type: Object, default: () => {} },
			placeholder: { type: String, default: "请选择" },
			size: { type: String, default: "default" },
			clearable: { type: Boolean, default: false },
			multiple: { type: Boolean, default: false },
			filterable: { type: Boolean, default: false },
			collapseTags: { type: Boolean, default: false },
			collapseTagsTooltip: { type: Boolean, default: false },
			disabled: { type: Boolean, default: false },
			tableWidth: {type: Number, default: 400},
			mode: { type: String, default: "popover" },
			props: { type: Object, default: () => {} }
		},
		data() {
			return {
				loading: false,
				keyword: null,
				defaultValue: [],
				tableData: [],
				pageSize: config.pageSize,
				total: 0,
				currentPage: 1,
				defaultProps: {
					label: config.props.label,
					value: config.props.value,
					page: config.request.page,
					pageSize: config.request.pageSize,
					keyword: config.request.keyword
				},
				formData: {}
			}
		},
		watch: {
			modelValue:{
				handler(){
					this.defaultValue = this.modelValue
					this.autoCurrentLabel()
				},
				deep: true
			}
		},
		mounted() {
			this.defaultProps = Object.assign(this.defaultProps, this.props);
			this.defaultValue =  this.modelValue
			this.autoCurrentLabel()
		},
		methods: {
			//表格显示隐藏回调
			visibleChange(visible){
				if(visible){
					this.currentPage = 1
					this.keyword = null
					this.formData = {}
					this.getData()
				}else{
					this.autoCurrentLabel()
				}
			},
			//获取表格数据
			async getData(){
				this.loading = true;
                const reqData = {
                    [this.defaultProps.page]: this.currentPage,
                    [this.defaultProps.pageSize]: this.pageSize,
                    [this.defaultProps.keyword]: this.keyword
                };
                Object.assign(reqData, this.params, this.formData)
                const res = await this.apiObj.get(reqData);
                const parseData = config.parseData(res);
                this.tableData = parseData.rows;
				this.total = parseData.total;
				this.loading = false;
				//表格默认赋值
				await this.$nextTick(() => {
                    if (this.multiple) {
                        this.defaultValue.forEach(row => {
                            const setrow = this.tableData.filter(item => item[this.defaultProps.value] === row[this.defaultProps.value])
                            if (setrow.length > 0) {
                                this.$refs.table.toggleRowSelection(setrow[0], true);
                            }
                        })
                    } else {
                        const setrow = this.tableData.filter(item => item[this.defaultProps.value] === this.defaultValue[this.defaultProps.value]);
                        this.$refs.table.setCurrentRow(setrow[0]);
                    }
                    this.$refs.table.setScrollTop(0)
                })
			},
			//插糟表单提交
			formSubmit(){
				this.currentPage = 1
				this.keyword = null
				this.getData()
			},
			//分页刷新表格
			reload(){
				this.getData()
			},
			//自动模拟options赋值
			autoCurrentLabel(){
				this.$nextTick(() => {
					if(this.multiple) {
                        this.defaultValue.forEach(item => {
                            item.label = item[this.defaultProps.label]
                        })
                    }
				})
			},
			//表格勾选事件
			select(rows, row){
                if(rows.length && rows.indexOf(row) !== -1){
					this.defaultValue.push(row)
				}else{
					this.defaultValue.splice(this.defaultValue.findIndex(item => item[this.defaultProps.value] == row[this.defaultProps.value]), 1)
				}
				this.autoCurrentLabel()
				this.$emit('update:modelValue', this.defaultValue);
				this.$emit('change', this.defaultValue);
			},
			//表格全选事件
			selectAll(rows){
				var isAllSelect = rows.length > 0
				if(isAllSelect){
					rows.forEach(row => {
						var isHas = this.defaultValue.find(item => item[this.defaultProps.value] == row[this.defaultProps.value])
						if(!isHas){
							this.defaultValue.push(row)
						}
					})
				}else{
					this.tableData.forEach(row => {
						var isHas = this.defaultValue.find(item => item[this.defaultProps.value] == row[this.defaultProps.value])
						if(isHas){
							this.defaultValue.splice(this.defaultValue.findIndex(item => item[this.defaultProps.value] == row[this.defaultProps.value]), 1)
						}
					})
				}
				this.autoCurrentLabel()
				this.$emit('update:modelValue', this.defaultValue);
				this.$emit('change', this.defaultValue);
			},
			click(row){
				if(this.multiple){
					//处理多选点击行
				}else{
					this.defaultValue = row
					this.$refs.select.blur()
					this.autoCurrentLabel()
					this.$emit('update:modelValue', this.defaultValue);
					this.$emit('change', this.defaultValue);
				}
			},
			//tags删除后回调
			removeTag(tag){
				var row = this.findRowByKey(tag[this.defaultProps.value])
				this.$refs.table.toggleRowSelection(row, false);
				this.$emit('update:modelValue', this.defaultValue);
			},
			//清空后的回调
			clear(){
				this.$emit('update:modelValue', this.defaultValue);
			},
			// 关键值查询表格数据行
			findRowByKey (value) {
				return this.tableData.find(item => item[this.defaultProps.value] === value)
			},
			filterMethod(keyword){
				if(!keyword){
					this.keyword = null;
					return false;
				}
				this.keyword = keyword;
				this.getData()
			},
			// 触发select隐藏
			blur(){
				this.$refs.select.blur();
			},
			// 触发select显示
			focus(){
				this.$refs.select.focus();
			}
		}
	}
</script>

<style scoped>
	.sc-table-select__table {padding:12px;}
	.sc-table-select__page {padding-top: 12px;}
</style>
