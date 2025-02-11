package cosyvoice

import (
	"context"
	"errors"

	"github.com/lithammer/shortuuid"
)

type Result struct {
	Data []byte
	Err  error
}

type AsyncSynthesizer struct {
	conn     *wsConn
	config   SynthesizerConfig
	chanSize int

	taskID string
	input  chan Task
	output chan Result
}

func (s *AsyncSynthesizer) RunTask(ctx context.Context) (<-chan Result, error) {
	if s.taskID != "" {
		return nil, errors.New("task already running")
	}

	s.taskID = shortuuid.New()
	err := s.conn.sendRunTaskCmd(ctx, s.taskID, s.config)
	if err != nil {
		return nil, err
	}

	s.input = make(chan Task, s.chanSize)
	s.output = make(chan Result, s.chanSize)

	s.conn.Start(s.input, s.output)

	return s.output, nil
}

func (s *AsyncSynthesizer) SendText(ctx context.Context, text string) error {
	cmd, err := generateContinueTaskCmd(s.taskID, text)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case s.input <- Task{data: cmd}:
	}
	return nil
}

func (s *AsyncSynthesizer) FinishTask(ctx context.Context) error {
	cmd, err := generateFinishTaskCmd(s.taskID)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case s.input <- Task{data: cmd}:
	}
	close(s.input)

	s.conn.logger.Infof("waiting for task %s to finish", s.taskID)
	s.conn.Wait()
	s.conn.logger.Infof("task %s finished", s.taskID)
	s.taskID = ""
	s.input = nil
	s.output = nil
	return nil
}

func (s *AsyncSynthesizer) Close() error {
	return s.conn.Close()
}
