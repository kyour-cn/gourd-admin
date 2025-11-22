<template>
  <el-main>
    <el-card shadow="never">
<!--          <el-alert title="资源库选择器即将弃用,将不会维护更新,与1.7版本之后将移除此组件" type="error" style="margin-bottom:20px;"/>-->
      <sc-file-select v-model:multiple="multiple" :hideUpload="false" mode="manage"/>
    </el-card>
  </el-main>
</template>

<script setup>
import scFileSelect from '@/components/scFileSelect'
import {ref} from "vue";

const multiple = ref(false)

</script>
