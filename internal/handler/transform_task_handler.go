package handler

import (
	"net/http"

	"texttool/internal/model"
	"texttool/pkg/httpx"
)

func (s *Server) registerTransformTaskRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/tasks", s.createTransformTask)
	mux.HandleFunc("GET /api/tasks", s.listTransformTasks)
	mux.HandleFunc("GET /api/tasks/{id}", s.getTransformTask)
	mux.HandleFunc("POST /api/tasks/{id}/execute", s.executeTransformTask)
	mux.HandleFunc("DELETE /api/tasks/{id}", s.deleteTransformTask)
}

type createTransformTaskRequest struct {
	DocumentID string `json:"document_id"`
	Operation  string `json:"operation"`
	Input      string `json:"input"`
}

func (s *Server) createTransformTask(w http.ResponseWriter, r *http.Request) {
	var req createTransformTaskRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.CreateTransformTask(model.TransformTask{DocumentID: req.DocumentID, Operation: req.Operation, Input: req.Input})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, t)
}

func (s *Server) listTransformTasks(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.TransformTaskFilter{
		DocumentID: r.URL.Query().Get("document_id"),
		Operation:  r.URL.Query().Get("operation"),
		Status:     r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListTransformTasks(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getTransformTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	t, err := s.svc.GetTransformTask(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

func (s *Server) executeTransformTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	t, err := s.svc.ExecuteTransformTask(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

func (s *Server) deleteTransformTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteTransformTask(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
