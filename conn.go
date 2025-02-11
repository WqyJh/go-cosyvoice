package cosyvoice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	openairt "github.com/WqyJh/go-openai-realtime"
)

type wsConn struct {
	ctx          context.Context
	cancel       context.CancelFunc
	conn         openairt.WebSocketConn
	logger       openairt.Logger
	pingInterval time.Duration

	sendBarrier chan struct{}

	wg sync.WaitGroup
}

type Task struct {
	data string
}

func (c *wsConn) sendRunTaskCmd(ctx context.Context, taskID string, voiceConfig SynthesizerConfig) error {
	runTaskCmd, err := generateRunTaskCmd(taskID, voiceConfig)
	if err != nil {
		return err
	}

	err = c.conn.WriteMessage(ctx, openairt.MessageText, []byte(runTaskCmd))
	if err != nil {
		return err
	}

	return nil
}

func (c *wsConn) handleReceiveMessage(resultCh chan<- Result) {
	defer func() {
		close(resultCh)
		c.wg.Done()
	}()

	for {
		select {
		case <-c.ctx.Done():
			c.logger.Infof("handleReceiveMessage stopped: %v", c.ctx.Err())
			resultCh <- Result{Err: c.ctx.Err()}
			return
		default:
		}

		msgType, message, err := c.conn.ReadMessage(c.ctx)
		if err != nil {
			resultCh <- Result{Err: err}
			return
		}

		if msgType == openairt.MessageBinary {
			c.logger.Debugf("received binary message: %d bytes", len(message))
			resultCh <- Result{Data: message}
		} else {
			// process event message
			var event Event
			err = json.Unmarshal(message, &event)
			if err != nil {
				resultCh <- Result{Err: err}
				return
			}
			ok, err := c.handleEvent(event)
			if err != nil {
				resultCh <- Result{Err: err}
				return
			}
			if ok {
				return
			}
		}
	}
}

func (c *wsConn) handleEvent(event Event) (bool, error) {
	c.logger.Debugf("received event: %s", event.Header.Event)
	switch event.Header.Event {
	case "task-started":
		// only received the task-start event event,we can send other message.
		close(c.sendBarrier)
		return false, nil

	case "result-generated":
		// pass; ignore result-generated event

	case "task-finished":
		return true, nil

	case "task-failed":
		return true, fmt.Errorf("received task-failed event: %w", event.Header.Error)
	}
	return false, nil
}

func (c *wsConn) Start(inputCh <-chan Task, outputCh chan<- Result) {
	c.sendBarrier = make(chan struct{})

	c.wg.Add(2)
	go c.sendMessage(inputCh)
	go c.handleReceiveMessage(outputCh)
}

func (c *wsConn) sendMessage(sendCh <-chan Task) {
	defer c.wg.Done()

	ticker := time.NewTicker(c.pingInterval)
	defer ticker.Stop()

	select {
	case <-c.ctx.Done():
		return
	case <-c.sendBarrier:
		// wait for the barrier
	}

	for {
		select {
		case <-c.ctx.Done():
			c.logger.Infof("sendMessage stopped: %v", c.ctx.Err())
			return

		case <-ticker.C:
			err := c.conn.Ping(c.ctx)
			if err != nil {
				c.logger.Errorf("ping err: %v", err)
			} else {
				c.logger.Debugf("ping success")
			}

			ticker.Reset(c.pingInterval)

		case task, ok := <-sendCh:
			if !ok {
				return
			}
			err := c.conn.WriteMessage(c.ctx, openairt.MessageText, []byte(task.data))
			if err != nil {
				var permanent *openairt.PermanentError
				if errors.As(err, &permanent) {
					c.logger.Errorf("send message permanent error: %w", permanent.Err)
					return
				}
				c.logger.Errorf("send message error: %+v", err)
				continue
			}
			c.logger.Debugf("send message success: %s", task.data)

			ticker.Reset(c.pingInterval)
		}
	}
}

func (c *wsConn) Close() error {
	c.cancel()
	return c.conn.Close()
}

func (c *wsConn) Wait() {
	c.wg.Wait()
}
