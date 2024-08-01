import api from "../api";

//上传配置
export default {
    apiObj: 'api.admin.common.upload',			//TODO:上传请求API对象
    filename: "file",					//form请求时文件的key
    successCode: 0,					    //请求完成代码
    maxSize: 10,						//最大文件大小 默认10MB
    parseData: function (res) {
        return {
            code: res.code,				//分析状态字段结构
            fileName: res.data.fileName,//分析文件名称
            src: res.data.url,			//分析图片远程地址结构
            msg: res.message			//分析描述字段结构
        }
    },
    apiObjFile: api.common.upload.uploadFile,	//附件上传请求API对象
    maxSizeFile: 10						//最大文件大小 默认10MB
}
