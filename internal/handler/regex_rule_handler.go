package handler

import (
	"net/http"

	"texttool/internal/model"
	"texttool/pkg/httpx"
)

func (s *Server) registerRegexRuleRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/rules", s.createRegexRule)
	mux.HandleFunc("GET /api/rules", s.listRegexRules)
	mux.HandleFunc("GET /api/rules/{id}", s.getRegexRule)
	mux.HandleFunc("PUT /api/rules/{id}", s.updateRegexRule)
	mux.HandleFunc("DELETE /api/rules/{id}", s.deleteRegexRule)
	mux.HandleFunc("POST /api/rules/{id}/match", s.matchWithRule)
	mux.HandleFunc("POST /api/rules/{id}/replace", s.replaceWithRule)
	mux.HandleFunc("POST /api/rules/match", s.matchWithPattern)
	mux.HandleFunc("POST /api/rules/replace", s.replaceWithPattern)
}

type createRegexRuleRequest struct {
	Name        string `json:"name"`
	Pattern     string `json:"pattern"`
	Flags       string `json:"flags"`
	Description string `json:"description"`
}

func (s *Server) createRegexRule(w http.ResponseWriter, r *http.Request) {
	var req createRegexRuleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rule, err := s.svc.CreateRegexRule(model.RegexRule{Name: req.Name, Pattern: req.Pattern, Flags: req.Flags, Description: req.Description})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, rule)
}

func (s *Server) listRegexRules(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.RegexRuleFilter{
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListRegexRules(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getRegexRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	rule, err := s.svc.GetRegexRule(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rule)
}

type updateRegexRuleRequest struct {
	Name        string `json:"name"`
	Pattern     string `json:"pattern"`
	Flags       string `json:"flags"`
	Description string `json:"description"`
}

func (s *Server) updateRegexRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateRegexRuleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rule, err := s.svc.UpdateRegexRule(id, model.RegexRule{Name: req.Name, Pattern: req.Pattern, Flags: req.Flags, Description: req.Description})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rule)
}

func (s *Server) deleteRegexRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteRegexRule(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type matchRequest struct {
	Text string `json:"text"`
}

func (s *Server) matchWithRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req matchRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	result, err := s.svc.MatchWithRule(id, req.Text)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

type replaceRequest struct {
	Text        string `json:"text"`
	Replacement string `json:"replacement"`
}

func (s *Server) replaceWithRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req replaceRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	result, err := s.svc.ReplaceWithRule(id, req.Text, req.Replacement)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]string{"result": result})
}

type matchPatternRequest struct {
	Pattern string `json:"pattern"`
	Flags   string `json:"flags"`
	Text    string `json:"text"`
}

func (s *Server) matchWithPattern(w http.ResponseWriter, r *http.Request) {
	var req matchPatternRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	result, err := s.svc.MatchWithPattern(req.Pattern, req.Flags, req.Text)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}

type replacePatternRequest struct {
	Pattern     string `json:"pattern"`
	Flags       string `json:"flags"`
	Text        string `json:"text"`
	Replacement string `json:"replacement"`
}

func (s *Server) replaceWithPattern(w http.ResponseWriter, r *http.Request) {
	var req replacePatternRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	result, err := s.svc.ReplaceWithPattern(req.Pattern, req.Flags, req.Text, req.Replacement)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]string{"result": result})
}
