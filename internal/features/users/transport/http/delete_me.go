package users_transport_http

import (
	core_context "messenger/internal/core/context"
	"messenger/internal/core/logger"
	http_response "messenger/internal/core/transport/http/response"
	"net/http"
)

func (h *UsersHandler) DeleteMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	claims := core_context.ClaimsRequired(ctx)
	sender := http_response.NewHTTPSender(log, w, errorMapper)

	if err := h.usersService.DeleteAccount(ctx, claims.UserID); err != nil {
		sender.Error(err)
		return
	}
	h.cookieManager.ClearRefreshToken(w)
	sender.OK(http.StatusNoContent, nil)
}
