package system

import (
	"mime/multipart"
	"net/http"
	"strconv"

	"app/internal/http/admin/dto"
	"app/internal/http/admin/services"
	"app/internal/http/common/controller"
)

// File 应用控制器
type File struct {
	controller.Base //继承基础控制器
}

func (c *File) MenuList(w http.ResponseWriter, r *http.Request) {

	service := services.NewFileService(r.Context())
	res, err := service.GetMenuList()
	if err != nil {
		_ = c.Fail(w, 500, "获取列表失败", err.Error())
	}
	_ = c.Success(w, "", res)
}

func (c *File) List(w http.ResponseWriter, r *http.Request) {
	req := &dto.FileListReq{}
	if err := c.QueryReqUnmarshal(r, req); err != nil {
		_ = c.Fail(w, 101, "请求参数异常："+err.Error(), "")
		return
	}

	service := services.NewFileService(r.Context())
	res, err := service.GetList(req)
	if err != nil {
		_ = c.Fail(w, 500, "获取列表失败", err.Error())
		return
	}
	_ = c.Success(w, "", res)
}

func (c *File) Upload(w http.ResponseWriter, r *http.Request) {
	claims, err := c.GetJwt(r)
	if err != nil {
		_ = c.Fail(w, 101, err.Error(), nil)
		return
	}

	// 解析 multipart form
	err = r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		_ = c.Fail(w, 501, "文件太大或解析失败", err.Error())
		return
	}

	// 获取上传文件
	file, header, err := r.FormFile("file")
	if err != nil {
		_ = c.Fail(w, 501, "上传文件读取失败", err.Error())
		return
	}
	defer func(file multipart.File) {
		_ = file.Close()
	}(file)

	req := &dto.FileUploadReq{
		Claims:     claims,
		File:       file,
		FileHeader: header,
	}
	req.MenuId, _ = strconv.ParseInt(r.FormValue(""), 10, 64)

	service := services.NewFileService(r.Context())
	res, err := service.Upload(req)
	if err != nil {
		_ = c.Fail(w, 1, "创建失败", err.Error())
		return
	}
	_ = c.Success(w, "success", map[string]any{
		"id":       res.ID,
		"src":      res.URL,
		"fileName": res.FileName,
	})
}

func (c *File) Delete(w http.ResponseWriter, r *http.Request) {
	req := struct {
		Ids []int64 `json:"ids"`
	}{}
	if err := c.JsonReqUnmarshal(r, &req); err != nil {
		_ = c.Fail(w, 101, "请求参数异常", err.Error())
		return
	}

	service := services.NewFileService(r.Context())
	res, err := service.Delete(req.Ids)
	if err != nil {
		_ = c.Fail(w, 1, "删除失败", err.Error())
		return
	}
	_ = c.Success(w, "success", res)
}
