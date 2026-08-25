package service

import (
	"context"
	"texttool/internal/scanner"
)

func (s *Service) ScanText(ctx context.Context, runner *scanner.Runner, content string) (string, error) {
	var result string
	var err error
	for attempt := 0; attempt < 2; attempt++ {
		result, err = runner.Scan(ctx, content)
		if err == nil {
			return result, nil
		}
	}
	return "", err
}
