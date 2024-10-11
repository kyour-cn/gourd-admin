<template>
    <el-container>
        <el-aside width="220px" style="border-radius: 10px 0 0 0;">
            <el-tree
                ref="category"
                class="menu"
                node-key="label"
                :data="category"
                :default-expanded-keys="['app', 'system']"
                :highlight-current="true"
                :expand-on-click-node="false"
                @current-change="handleCurrentChange"
            >
                <template #default="{node, data}">
					<span class="custom-tree-node">
                        <sc-status-indicator
                            v-if="data.color"
                            type="success"
                            :style="{background: data.color, marginRight: '5px'}"
                        />
						<span class="label">
							{{ data.name }}
						</span>
					</span>
                </template>
            </el-tree>
        </el-aside>
        <el-container>
            <el-main class="nopadding">
                <el-container>
                    <el-header>
                        <div class="left-panel">
                            <el-date-picker
                                v-model="date"
                                format="YYYY-MM-DD HH:mm:ss"
                                type="datetimerange"
                                range-separator="至"
                                start-placeholder="开始日期"
                                end-placeholder="结束日期"
                                @change="timeChange"
                            />
                        </div>
                        <div class="right-panel">

                        </div>
                    </el-header>
                    <el-header style="height:150px;">
                        <scEcharts height="100%" width="100%" :option="logsChartOption"/>
                    </el-header>
                    <el-main class="nopadding">
                        <scTable
                            ref="table" :apiObj="apiObj"
                            :params="{start_time:date?.[0] || null,end_time:date?.[1] || null}"
                            stripe
                            highlightCurrentRow @row-click="rowClick">
                            <el-table-column label="ID" prop="id" width="100"/>
                            <el-table-column label="名称" prop="type_name" width="100"/>
                            <el-table-column label="标题" show-overflow-tooltip prop="title"/>
                            <el-table-column label="请求来源" show-overflow-tooltip prop="request_source"
                                             width="250"/>
                            <el-table-column label="请求ip" prop="request_ip" width="150"/>
                            <el-table-column label="用户" prop="request_user" width="100"/>
                            <el-table-column label="请求时间" prop="create_time" width="170">
                                <template #default="{row,$index}">
                                    {{ $TOOL.dateFormat(row.create_time * 1000) }}
                                </template>
                            </el-table-column>
                        </scTable>
                    </el-main>
                </el-container>
            </el-main>
        </el-container>
    </el-container>

    <el-drawer v-model="infoDrawer" title="日志详情" :size="600" destroy-on-close>
        <info ref="info"/>
    </el-drawer>
</template>

<script>
import info from './info'
import scEcharts from '@/components/scEcharts'
import ScStatusIndicator from "@/components/scMini/scStatusIndicator.vue";
import logApi from '@/api/admin/log'

export default {
    name: 'log',
    components: {
        ScStatusIndicator,
        info,
        scEcharts
    },
    data() {
        return {
            map: [],
            dateMaps: [],
            infoDrawer: false,
            echartsData: [],
            // 动态图表数据
            seriesData: [],
            logsChartOption: {
                color: ['#409eff', '#e6a23c', '#f56c6c'],
                grid: {
                    top: '0px',
                    left: '10px',
                    // right: '10px',
                    bottom: '0px'
                },
                tooltip: {
                    show: true,
                    trigger: 'axis',
                    formatter: (params) => {
                        let result = ''
                        const {seriesData, map} = this
                        if (map.length !== 0) {
                            const xAxis_val = `${params[0].axisValue}</br>`
                            for (let i = 0; i < seriesData.length; i++) {
                                if (!params[i].data) continue;
                                result += `<span style="display:inline-block;margin-right:5px;border-radius:10px;width:10px;height:10px;background-color:${seriesData[i].color}"></span>`
                                result += `${map[i].name}:${params[i].data}`
                                result += `</br>`
                            }
                            result = xAxis_val + result
                        } else {
                            result = '暂无数据'
                        }
                        return result
                    }
                },
                xAxis: {
                    type: 'category',
                    boundaryGap: false,
                    data: ['']
                },
                yAxis: {
                    show: false,
                    type: 'value'
                },
                series: []
            },
            category: [],
            types: [],
            date: [
                this.$TOOL.dateFormat(this.getCurrentMonthFirst()),
                this.$TOOL.dateFormat(new Date()),
            ],
            apiObj: this.$API.admin.log.list,
            search: {
                keyword: ""
            }
        }
    },
    mounted() {
        logApi.typeList.get({page_size: 1000}).then((res) => {
            this.types = res.data.rows
            this.category = this.renderTreeMenu(res.data.rows)
            this.echartsRender();
        })
    },
    methods: {
        async echartsRender() {
            const start_time = this.date[0]
            const end_time = this.date[1]

            let res = await this.$API.admin.log.logStat.get({start_time, end_time});

            const dateMaps = {}
            for (const i in res.data.days) {
                dateMaps[res.data.days[i]] = i
            }
            // 填充x轴的数据
            this.logsChartOption.xAxis.data = res.data.days;

            // 填充图表数据
            if (res.data.rows.length !== 0) {

                let seriesData = {}
                const typeMap = {}
                for (const i in this.types) {
                    typeMap[String(this.types[i].id)] = this.types[i]
                }

                for (const i in res.data.rows) {
                    const item = res.data.rows[i]
                    if (seriesData[item.type_id]) {
                        continue
                    }
                    seriesData[item.type_id] = {
                        id: item.type_id,
                        name: item.type_name,
                        color: typeMap[String(item.type_id)].color,
                        data: this.arrayPad([],res.data.days.length, 0)
                    }
                }

                // 填充堆叠图表数据
                for (const key in res.data.rows) {
                    const item = res.data.rows[key]
                    seriesData[item.type_id].data[dateMaps[item.date]] = item.count
                }

                this.map = []
                for (const key in seriesData) {
                    this.map.push(seriesData[key])
                }

                for (const key in this.map) {
                    this.seriesData.push({
                        data: this.map[key].data,
                        type: 'bar',
                        stack: 'log',
                        barWidth: '15px',
                        color: this.map[key].color
                    })
                }
                setTimeout(() => {
                    this.logsChartOption.series = this.seriesData
                }, 500)
            } else {
                this.seriesData = [];
                this.logsChartOption.series = [];
                this.logsChartOption.xAxis.data = [0];
                this.$message("暂无更多数据");
            }
        },
        arrayPad(arr, len, val) {
            if (arr.length >= len) {
                return arr
            }
            return arr.concat(Array(len - arr.length).fill(val))
        },

        // 左侧树形菜单
        renderTreeMenu(data) {
            const sysList = []
            const appList = []
            for (const i in data) {
                const item = data[i]
                if (item.app_id > 0) {
                    appList.push(item)
                } else {
                    sysList.push(item)
                }
            }

            return [
                {
                    id: 1,
                    name: "系统日志",
                    label: "system",
                    children: sysList
                },
                {
                    id: 2,
                    name: "应用日志",
                    label: "app",
                    children: appList
                },
            ]
        },
        timeChange() {
            // 点击时间清除判断是否为null
            if (this.date !== null) {
                this.$refs.table.upData({
                    start_time: this.date[0],
                    end_time: this.date[1]
                })
                this.echartsRender();
            } else {
                this.$refs.table.upData({
                    start_time: this.date?.[0] || '',
                    end_time: this.date?.[1] || ''
                })
            }
        },
        //获取当前月的第一天
        getCurrentMonthFirst() {
            let date = new Date();
            date.setDate(1);
            date.setHours(0, 0, 0, 0)
            return date;
        },
        rowClick(row) {
            this.infoDrawer = true
            this.$nextTick(() => {
                this.$refs.info.setData(row)
            })
        },
        handleCurrentChange(data) {
            this.$refs.table.upData({
                type_id: data.id
            })
        },
    }
}

</script>
