package handler

import (
	"net/http"

	"texttool/internal/model"
	"texttool/pkg/httpx"
)

func (s *Server) registerTextDocumentRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/documents", s.createTextDocument)
	mux.HandleFunc("GET /api/documents", s.listTextDocuments)
	mux.HandleFunc("GET /api/documents/{id}", s.getTextDocument)
	mux.HandleFunc("PUT /api/documents/{id}", s.updateTextDocument)
	mux.HandleFunc("DELETE /api/documents/{id}", s.deleteTextDocument)
	mux.HandleFunc("GET /api/documents/{id}/stats", s.getTextStats)
	mux.HandleFunc("GET /api/documents/{id}/dedup", s.dedupTextDocument)
}

type createTextDocumentRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Encoding string `json:"encoding"`
}

func (s *Server) createTextDocument(w http.ResponseWriter, r *http.Request) {
	var req createTextDocumentRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	d, err := s.svc.CreateTextDocument(model.TextDocument{Title: req.Title, Content: req.Content, Encoding: req.Encoding})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, d)
}

func (s *Server) listTextDocuments(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.TextDocumentFilter{
		Encoding: r.URL.Query().Get("encoding"),
		Keyword:  r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListTextDocuments(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getTextDocument(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	d, err := s.svc.GetTextDocument(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, d)
}

type updateTextDocumentRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Encoding string `json:"encoding"`
}

func (s *Server) updateTextDocument(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateTextDocumentRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	d, err := s.svc.UpdateTextDocument(id, model.TextDocument{Title: req.Title, Content: req.Content, Encoding: req.Encoding})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, d)
}

func (s *Server) deleteTextDocument(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteTextDocument(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) getTextStats(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	stats, err := s.svc.AnalyzeText(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, stats)
}

func (s *Server) dedupTextDocument(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	result, err := s.svc.DedupLines(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]string{"result": result})
}
