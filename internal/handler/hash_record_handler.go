package handler

import (
	"net/http"

	"texttool/internal/model"
	"texttool/pkg/httpx"
)

func (s *Server) registerHashRecordRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/hashes", s.createHashRecord)
	mux.HandleFunc("GET /api/hashes", s.listHashRecords)
	mux.HandleFunc("GET /api/hashes/{id}", s.getHashRecord)
	mux.HandleFunc("POST /api/hashes/compute", s.computeHash)
	mux.HandleFunc("POST /api/hashes/verify", s.verifyHash)
	mux.HandleFunc("DELETE /api/hashes/{id}", s.deleteHashRecord)
}

type createHashRecordRequest struct {
	DocumentID string `json:"document_id"`
	Algorithm  string `json:"algorithm"`
	Hash       string `json:"hash"`
}

func (s *Server) createHashRecord(w http.ResponseWriter, r *http.Request) {
	var req createHashRecordRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	h, err := s.svc.CreateHashRecord(model.HashRecord{DocumentID: req.DocumentID, Algorithm: req.Algorithm, Hash: req.Hash})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, h)
}

func (s *Server) listHashRecords(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.HashRecordFilter{
		DocumentID: r.URL.Query().Get("document_id"),
		Algorithm:  r.URL.Query().Get("algorithm"),
	}
	items, total, err := s.svc.ListHashRecords(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getHashRecord(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	h, err := s.svc.GetHashRecord(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, h)
}

type computeHashRequest struct {
	DocumentID string `json:"document_id"`
	Algorithm  string `json:"algorithm"`
}

func (s *Server) computeHash(w http.ResponseWriter, r *http.Request) {
	var req computeHashRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	h, err := s.svc.ComputeHash(req.DocumentID, req.Algorithm)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, h)
}

type verifyHashRequest struct {
	DocumentID string `json:"document_id"`
	Algorithm  string `json:"algorithm"`
	Expected   string `json:"expected"`
}

func (s *Server) verifyHash(w http.ResponseWriter, r *http.Request) {
	var req verifyHashRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	ok, err := s.svc.VerifyHash(req.DocumentID, req.Algorithm, req.Expected)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]bool{"matched": ok})
}

func (s *Server) deleteHashRecord(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteHashRecord(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
