package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"

	"backend/internal/domain"
	"backend/internal/repository"
	"backend/pkg/errors"
)

// ChatUseCase 定义 AI 对话的应用服务层接口
type ChatUseCase interface {
	SendMessage(ctx context.Context, userID uint64, sessionID, content string) (string, *errors.AppError)
}

type chatUseCase struct {
	repo     repository.ChatRepository
	llmModel model.ChatModel
}

// NewChatUseCase 实例化对话应用服务
func NewChatUseCase(repo repository.ChatRepository, llmModel model.ChatModel) ChatUseCase {
	return &chatUseCase{
		repo:     repo,
		llmModel: llmModel,
	}
}

func (u *chatUseCase) SendMessage(ctx context.Context, userID uint64, sessionID, content string) (string, *errors.AppError) {
	if u.llmModel == nil {
		return "", errors.ErrChatModelFailed
	}

	// 1. 获取历史对话记录 (取最近的10条作为上下文)
	histories, err := u.repo.GetHistoriesBySessionID(ctx, userID, sessionID, 10)
	if err != nil {
		return "", errors.ErrDatabase
	}

	// 2. 组装 eino的Messages 数组
	messages := make([]*schema.Message, 0, len(histories)+2)

	// 2.1 添加 System Prompt (依据 PRD: 扮演杭州导游)
	messages = append(messages, schema.SystemMessage(
		"你现在是一名资深的杭州土著导游，热情、专业、幽默。你非常了解杭州的各个景点、美食、文化和出行路线"+
			"请尽量用简洁、生动的语言回答游客的问题，如果遇到你不懂的问题，请礼貌地表明你不知道。只输出纯文本",
	),
	)

	// 2.2 添加历史上下文
	for _, h := range histories {
		if h.Role == "user" {
			messages = append(messages, schema.UserMessage(h.Content))
		} else if h.Role == "assistant" {
			messages = append(messages, schema.AssistantMessage(h.Content, nil))
		}
	}

	// 2.3 添加当前用户输入
	messages = append(messages, schema.UserMessage(content))

	// 3. 异步保存用户的提问到数据库
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[Chat] recovered from panic in save user history: %v", r)
			}
		}()
		if err := u.repo.SaveHistory(context.Background(), &domain.ChatHistory{
			UserID:    userID,
			SessionID: sessionID,
			Role:      "user",
			Content:   content,
		}); err != nil {
			log.Printf("[Chat] failed to save user history: %v", err)
		}
	}()

	// 4. 调用 eino 的Generate 方法获取 AI 回复
	resp, err := u.llmModel.Generate(ctx, messages)
	if err != nil {
		fmt.Printf("Ollama generate error: %v\n", err)
		return "", errors.ErrChatModelFailed
	}

	aiReply := resp.Content

	// 5. 异步保存 AI 的回复到数据库
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[Chat] recovered from panic in save assistant history: %v", r)
			}
		}()
		if err := u.repo.SaveHistory(context.Background(), &domain.ChatHistory{
			UserID:    userID,
			SessionID: sessionID,
			Role:      "assistant",
			Content:   aiReply,
		}); err != nil {
			log.Printf("[Chat] failed to save assistant history: %v", err)
		}
	}()

	return aiReply, nil
}
