package messages_service

import "errors"

var (
	ErrMessageTargetUnavailable = errors.New("message target unavailable")
	ErrMessageConflict          = errors.New("message conflict")
	ErrInternalInconsistency    = errors.New("internal inconsistency")
)
