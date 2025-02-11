package cosyvoice

import "context"

type SyncSynthesizer struct {
	asyncSynthesizer *AsyncSynthesizer
}

func (s *SyncSynthesizer) Close() error {
	return s.asyncSynthesizer.Close()
}

func (s *SyncSynthesizer) Call(ctx context.Context, text string) (<-chan Result, error) {
	output, err := s.asyncSynthesizer.RunTask(ctx)
	if err != nil {
		return nil, err
	}

	err = s.asyncSynthesizer.SendText(ctx, text)
	if err != nil {
		return nil, err
	}

	err = s.asyncSynthesizer.FinishTask(ctx)
	if err != nil {
		return nil, err
	}

	return output, nil
}
