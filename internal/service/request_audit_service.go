package service

import "texttool/internal/requestmeta"

func (s *Service) ScheduleIdentityAudit(pool *requestmeta.Pool, identity string, gate <-chan struct{}) <-chan string {
	meta := pool.Acquire(identity)
	result := make(chan string, 1)
	go func() {
		<-gate
		result <- meta.Identity
	}()
	pool.Release(meta)
	return result
}
