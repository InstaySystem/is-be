package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/InstaySystem/is-be/internal/common"
	"github.com/InstaySystem/is-be/internal/service"
	"github.com/InstaySystem/is-be/internal/types"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userSvc service.UserService
}

func NewUserHandler(userSvc service.UserService) *UserHandler {
	return &UserHandler{userSvc}
}

// CreateUser godoc
// @Summary      Create User
// @Description  Tạo một người dùng mới (admin hoặc staff)
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        payload  body      types.CreateUserRequest  true  "Thông tin người dùng mới"
// @Success      201      {object}  types.APIResponse{data=object{id=int64}}  "Tạo user thành công"
// @Failure      400      {object}  types.APIResponse  "Bad Request (validation error hoặc staff thiếu department)"
// @Failure      401      {object}  types.APIResponse  "Unauthorized"
// @Failure      404      {object}  types.APIResponse  "Department không tìm thấy"
// @Failure      409      {object}  types.APIResponse  "Conflict (email/username/SĐT đã tồn tại)"
// @Failure      500      {object}  types.APIResponse  "Internal Server Error"
// @Router       /users   [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req types.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mess := common.HandleValidationError(err)
		common.ToAPIResponse(c, http.StatusBadRequest, mess, nil)
		return
	}

	if req.DepartmentID != nil && req.Role == "admin" {
		req.DepartmentID = nil
	}
	if req.DepartmentID == nil && req.Role == "staff" {
		common.ToAPIResponse(c, http.StatusBadRequest, common.ErrDepartmentRequired.Error(), nil)
		return
	}

	id, err := h.userSvc.CreateUser(ctx, req)
	if err != nil {
		switch err {
		case common.ErrEmailAlreadyExists, common.ErrUsernameAlreadyExists, common.ErrPhoneAlreadyExists:
			common.ToAPIResponse(c, http.StatusConflict, err.Error(), nil)
		case common.ErrDepartmentNotFound:
			common.ToAPIResponse(c, http.StatusNotFound, err.Error(), nil)
		default:
			common.ToAPIResponse(c, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	common.ToAPIResponse(c, http.StatusCreated, "User created successfully", gin.H{
		"id": id,
	})
}

// GetUserByID godoc
// @Summary      Get User By ID
// @Description  Lấy thông tin chi tiết của một người dùng bằng ID
// @Tags         Users
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id          path      int  true  "User ID"
// @Success      200         {object}  types.APIResponse{data=object{user=types.UserResponse}}  "Lấy thông tin user thành công"
// @Failure      400         {object}  types.APIResponse  "Bad Request (ID không hợp lệ)"
// @Failure      401         {object}  types.APIResponse  "Unauthorized"
// @Failure      404         {object}  types.APIResponse  "User Not Found"
// @Failure      409         {object}  types.APIResponse  "Invalid Information"
// @Failure      500         {object}  types.APIResponse  "Internal Server Error"
// @Router       /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		common.ToAPIResponse(c, http.StatusBadRequest, common.ErrInvalidID.Error(), nil)
		return
	}

	user, err := h.userSvc.GetUserByID(ctx, userID)
	if err != nil {
		switch err {
		case common.ErrUserNotFound:
			common.ToAPIResponse(c, http.StatusNotFound, err.Error(), nil)
		default:
			common.ToAPIResponse(c, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Get user information successfully", gin.H{
		"user": common.ToUserResponse(user),
	})
}

// GetUsers godoc
// @Summary      Get User List
// @Description  Lấy danh sách người dùng có phân trang và lọc
// @Tags         Users
// @Produce      json
// @Security     ApiKeyAuth
// @Param        query  query     types.UserPaginationQuery  false  "Query phân trang và lọc"
// @Success      200    {object}  types.APIResponse{data=types.UserListResponse}  "Lấy danh sách user thành công"
// @Failure      400    {object}  types.APIResponse  "Bad Request (query không hợp lệ)"
// @Failure      401    {object}  types.APIResponse  "Unauthorized"
// @Failure      409    {object}  types.APIResponse  "Invalid Information"
// @Failure      500    {object}  types.APIResponse  "Internal Server Error"
// @Router       /users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var query types.UserPaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		mess := common.HandleValidationError(err)
		common.ToAPIResponse(c, http.StatusBadRequest, mess, nil)
		return
	}

	users, meta, err := h.userSvc.GetUsers(ctx, query)
	if err != nil {
		common.ToAPIResponse(c, http.StatusInternalServerError, "internal server error", nil)
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "Get user list successfully", common.ToUserListResponse(users, meta))
}

// GetAllRoles godoc
// @Summary      Get All Roles
// @Description  Lấy danh sách tất cả các vai trò (key và tên hiển thị)
// @Tags         Users
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200    {object}  types.APIResponse{data=object{roles=object}}  "Lấy danh sách vai trò thành công"
// @Failure      401    {object}  types.APIResponse  "Unauthorized"
// @Failure      409    {object}  types.APIResponse  "Invalid Information"
// @Failure      500    {object}  types.APIResponse  "Internal Server Error"
// @Router       /roles [get]
func (h *UserHandler) GetAllRoles(c *gin.Context) {
	rolesMap := map[string]string{
		common.RoleAdminDisplayName: common.RoleAdmin,
		common.RoleStaffDisplayName: common.RoleStaff,
	}

	common.ToAPIResponse(c, http.StatusOK, "Get all roles successfully", gin.H{
		"roles": rolesMap,
	})
}

// UpdateUser godoc
// @Summary      Update User
// @Description  Cập nhật thông tin chi tiết của một người dùng bằng ID
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id          path      int                      true  "User ID"
// @Param        payload     body      types.UpdateUserRequest  true  "Thông tin cần cập nhật"
// @Success      200         {object}  types.APIResponse  "Cập nhật user thành công"
// @Failure      400         {object}  types.APIResponse  "Bad Request (validation, ID, logic error)"
// @Failure      401         {object}  types.APIResponse  "Unauthorized"
// @Failure      404         {object}  types.APIResponse  "Not Found (user hoặc department)"
// @Failure      409         {object}  types.APIResponse  "Conflict (email/username/SĐT đã tồn tại)"
// @Failure      500         {object}  types.APIResponse  "Internal Server Error"
// @Router       /users/{id} [patch]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		common.ToAPIResponse(c, http.StatusBadRequest, common.ErrInvalidID.Error(), nil)
		return
	}

	var req types.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mess := common.HandleValidationError(err)
		common.ToAPIResponse(c, http.StatusBadRequest, mess, nil)
		return
	}

	if err := h.userSvc.UpdateUser(ctx, userID, req); err != nil {
		switch err {
		case common.ErrEmailAlreadyExists, common.ErrUsernameAlreadyExists, common.ErrPhoneAlreadyExists:
			common.ToAPIResponse(c, http.StatusConflict, err.Error(), nil)
		case common.ErrUserNotFound, common.ErrDepartmentNotFound:
			common.ToAPIResponse(c, http.StatusNotFound, err.Error(), nil)
		case common.ErrNeedAdmin, common.ErrDepartmentRequired:
			common.ToAPIResponse(c, http.StatusBadRequest, err.Error(), nil)
		default:
			common.ToAPIResponse(c, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "User updated successfully", nil)
}

// UpdateUserPassword godoc
// @Summary      Update User Password
// @Description  Cập nhật mật khẩu cho một user (thường dùng bởi admin)
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id                   path      int                            true  "User ID"
// @Param        payload              body      types.UpdateUserPasswordRequest  true  "Mật khẩu mới"
// @Success      200                  {object}  types.APIResponse  "Cập nhật mật khẩu thành công"
// @Failure      400                  {object}  types.APIResponse  "Bad Request (validation hoặc ID không hợp lệ)"
// @Failure      401         					{object}  types.APIResponse  "Unauthorized"
// @Failure      404                  {object}  types.APIResponse  "User không tìm thấy"
// @Failure      500                  {object}  types.APIResponse  "Internal Server Error"
// @Router       /users/{id}/password [put]
func (h *UserHandler) UpdateUserPassword(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		common.ToAPIResponse(c, http.StatusBadRequest, common.ErrInvalidID.Error(), nil)
		return
	}

	var req types.UpdateUserPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		mess := common.HandleValidationError(err)
		common.ToAPIResponse(c, http.StatusBadRequest, mess, nil)
		return
	}

	if err = h.userSvc.UpdateUserPassword(ctx, userID, req); err != nil {
		switch err {
		case common.ErrUserNotFound:
			common.ToAPIResponse(c, http.StatusNotFound, err.Error(), nil)
		default:
			common.ToAPIResponse(c, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "User password updated successfully", nil)
}

// DeleteUser godoc
// @Summary      Delete User
// @Description  Xoá một người dùng bằng ID
// @Tags         Users
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  types.APIResponse  "Xoá user thành công"
// @Failure      400  {object}  types.APIResponse  "Bad Request (ID không hợp lệ)"
// @Failure      401  {object}  types.APIResponse  "Unauthorized"
// @Failure      404  {object}  types.APIResponse  "User không tìm thấy"
// @Failure      409  {object}  types.APIResponse  "Conflict (không thể xoá bản ghi được bảo vệ)"
// @Failure      500  {object}  types.APIResponse  "Internal Server Error"
// @Router       /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		common.ToAPIResponse(c, http.StatusBadRequest, common.ErrInvalidID.Error(), nil)
		return
	}

	if err := h.userSvc.DeleteUser(ctx, userID); err != nil {
		switch err {
		case common.ErrUserNotFound:
			common.ToAPIResponse(c, http.StatusNotFound, err.Error(), nil)
		case common.ErrProtectedRecord:
			common.ToAPIResponse(c, http.StatusConflict, err.Error(), nil)
		default:
			common.ToAPIResponse(c, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	common.ToAPIResponse(c, http.StatusOK, "User deleted successfully", nil)
}
