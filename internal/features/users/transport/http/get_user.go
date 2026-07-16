package users_transport_http

import (
	"net/http"
)

type GetUserResponse UserDTOResponse

func (h *UsersHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	// log := logger.FromContext(ctx)
	// sender := http_response.NewHTTPSender(log, w)

	// idStr := chi.URLParam(r, "id")

	// userID, err := uuid.Parse(idStr)
	// if err != nil {
	// 	sender.Error(domain.ValidationErr("user_id", nil))
	// 	return
	// }
	// user, err := h.usersService.GetUser(ctx, userID)
	// if err != nil {
	// 	sender.Error(err)
	// 	return
	// }

	// response := GetUserResponse(userDTOFromDomain(user))
	// sender.OK(http.StatusOK, response)
}
