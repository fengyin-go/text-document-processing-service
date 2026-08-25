package service

import "texttool/internal/resource"

func (s *Service) ReadTextBatch(factory *resource.Factory, names []string) ([]string, error) {
	result := make([]string, 0, len(names))
	for _, name := range names {
		handle, err := factory.Open(name)
		if err != nil {
			return nil, err
		}
		defer handle.Close()
		result = append(result, handle.Read())
	}
	return result, nil
}
