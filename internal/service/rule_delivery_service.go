package service

import "texttool/internal/delivery"

func (s *Service) DeliverRule(tx delivery.Transaction, publisher delivery.Publisher, audit delivery.Audit, payload string) (err error) {
	defer func() { err = delivery.Finalize(tx, err) }()
	for attempt := 0; attempt < 2; attempt++ {
		audit.Record("published:" + payload)
		err = publisher.Publish(payload)
		if err == nil {
			return nil
		}
	}
	return err
}
