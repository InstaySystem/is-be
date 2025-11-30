package implement

import (
	"context"
	"errors"

	"github.com/InstaySystem/is-be/internal/model"
	"github.com/InstaySystem/is-be/internal/repository"
	"github.com/InstaySystem/is-be/internal/types"
	"gorm.io/gorm"
)

type chatRepoImpl struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) repository.ChatRepository {
	return &chatRepoImpl{db}
}

func (r *chatRepoImpl) CreateChatTx(tx *gorm.DB, chat *model.Chat) error {
	return tx.Create(chat).Error
}

func (r *chatRepoImpl) CreateMessageTx(tx *gorm.DB, message *model.Message) error {
	return tx.Create(message).Error
}

func (r *chatRepoImpl) FindChatByIDTx(tx *gorm.DB, chatID int64) (*model.Chat, error) {
	var chat model.Chat
	if err := tx.Where("id = ?", chatID).First(&chat).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &chat, nil
}

func (r *chatRepoImpl) FindAllUnreadMessageIDsByChatIDAndSenderTypeTx(tx *gorm.DB, chatID, staffID int64, senderType string) ([]int64, error) {
	var ids []int64
	if err := tx.Where("chat_id = ? AND sender_type", chatID, senderType).Where(
		"id NOT IN (?)", tx.Model(&model.MessageStaff{}).
			Select("message_id").
			Where("staff_id = ?", staffID),
	).Model(&model.Message{}).Pluck("id", &ids).Error; err != nil {
		return nil, err
	}

	return ids, nil
}

func (r *chatRepoImpl) CreateMessageStaffsTx(tx *gorm.DB, messageStaffs []*model.MessageStaff) error {
	return tx.Create(messageStaffs).Error
}

func (r *chatRepoImpl) UpdateChatTx(tx *gorm.DB, chatID int64, updateData map[string]any) error {
	return tx.Model(&model.Chat{}).Where("id = ?", chatID).Updates(updateData).Error
}

func (r *chatRepoImpl) UpdateMessagesByChatIDAndSenderTypeTx(tx *gorm.DB, chatID int64, senderType string, updateData map[string]any) error {
	return tx.Model(&model.Message{}).Where("id = ? AND sender_type = ? AND is_read = false", chatID, senderType).Updates(updateData).Error
}

func (r *chatRepoImpl) FindAllChatsByDepartmentIDWithDetailsPaginated(ctx context.Context, query types.ChatPaginationQuery, staffID, departmentID int64) ([]*model.Chat, int64, error) {
	var chats []*model.Chat
	var total int64

	db := r.db.WithContext(ctx).Where("department_id = ?", departmentID).Model(&model.Chat{})
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (query.Page - 1) * query.Limit
	if err := db.Order("last_message_at DESC").
		Limit(int(query.Limit)).
		Offset(int(offset)).
		Preload("OrderRoom.Room").
		Preload("OrderRoom.Booking").
		Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Select("messages.*").
				Joins("JOIN chats ON chats.id = messages.chat_id AND chats.last_message_at = messages.created_at")
		}).
		Preload("Messages.Sender").
		Preload("Messages.StaffsRead", "staff_id = ?", staffID).
		Find(&chats).Error; err != nil {
		return nil, 0, err
	}

	return chats, total, nil
}

func (r *chatRepoImpl) FindChatByIDWithDetails(ctx context.Context, chatID, staffID int64) (*model.Chat, error) {
	var chat model.Chat

	if err := r.db.WithContext(ctx).
		Preload("OrderRoom.Room").
		Preload("OrderRoom.Booking").
		Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Preload("Messages.Sender").
		Preload("Messages.StaffsRead", "staff_id = ?", staffID).
		Where("id = ?", chatID).First(&chat).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &chat, nil
}

func (r *chatRepoImpl) FindChatByCodeWithDetails(ctx context.Context, chatCode string) (*model.Chat, error) {
	var chat model.Chat
	if err := r.db.WithContext(ctx).Preload("Department").Preload("Messages", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at ASC")
	}).Where("code = ?", chatCode).First(&chat).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &chat, nil
}

func (r *chatRepoImpl) FindAllChatsByOrderRoomIDWithDetails(ctx context.Context, orderRoomID int64) ([]*model.Chat, error) {
	var chats []*model.Chat
	if err := r.db.WithContext(ctx).
		Where("order_room_id = ?", orderRoomID).
		Order("last_message_at DESC").
		Preload("Department").
		Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Select("messages.*").
				Joins("JOIN chats ON chats.id = messages.chat_id AND chats.last_message_at = messages.created_at")
		}).
		Find(&chats).Error; err != nil {
		return nil, err
	}

	return chats, nil
}
