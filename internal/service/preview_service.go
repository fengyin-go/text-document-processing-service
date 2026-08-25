package service

func (s *Service) PreparePreview(id, source string) error {
	_, err := s.preview.Prepare(id, source)
	return err
}

func (s *Service) RenderPreview(id string) (string, error) {
	return s.preview.Render(id)
}
