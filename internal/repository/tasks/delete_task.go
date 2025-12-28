package tasks

import "context"

func (r *Repo) DeleteTask(ctx context.Context, taskID uint) error {
	r.mu.Lock()
	delete(r.storage, taskID)
	r.mu.Unlock()

	return nil
}
