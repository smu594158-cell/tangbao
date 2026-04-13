package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"

	"backend/internal/domain"
	"backend/pkg/errors"
)

type tourUsecase struct {
	repo     domain.TourRepository
	llmModel model.ChatModel
}

// NewTourUsecase 实例化旅游内容服务
func NewTourUsecase(repo domain.TourRepository, ollamaBaseURL string) domain.TourUsecase {
	chatModel, err := ollama.NewChatModel(context.Background(), &ollama.ChatModelConfig{
		BaseURL: ollamaBaseURL,
		Model:   "qwen:4b",
	})
	if err != nil {
		log.Fatalf("failed to init ollama chat model for tour: %v", err)
	}

	return &tourUsecase{
		repo:     repo,
		llmModel: chatModel,
	}
}

func (u *tourUsecase) GetAttractionInfo(ctx context.Context, id uint64) (*domain.Attraction, error) {
	attraction, err := u.repo.GetAttractionByID(ctx, id)
	if err != nil {
		return nil, errors.ErrDatabase
	}
	if attraction == nil {
		return nil, errors.ErrNotFound
	}
	return attraction, nil
}

func (u *tourUsecase) ListAttractions(ctx context.Context, page, size int) ([]*domain.Attraction, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}

	list, total, err := u.repo.ListAttractions(ctx, page, size)
	if err != nil {
		return nil, 0, errors.ErrDatabase
	}
	return list, total, nil
}

func (u *tourUsecase) CreateAttraction(ctx context.Context, name, description, address string, lng, lat float64) (*domain.Attraction, error) {
	attraction := &domain.Attraction{
		Name:        name,
		Description: description,
		Address:     address,
		LocationLng: lng,
		LocationLat: lat,
		HeatLevel:   0,
	}

	if err := u.repo.CreateAttraction(ctx, attraction); err != nil {
		return nil, errors.ErrDatabase
	}
	return attraction, nil
}

func (u *tourUsecase) UpdateAttraction(ctx context.Context, id uint64, name, description, address string, lng, lat float64) (*domain.Attraction, error) {
	attraction, err := u.repo.GetAttractionByID(ctx, id)
	if err != nil {
		return nil, errors.ErrDatabase
	}
	if attraction == nil {
		return nil, errors.ErrNotFound
	}

	attraction.Name = name
	attraction.Description = description
	attraction.Address = address
	attraction.LocationLng = lng
	attraction.LocationLat = lat

	if err := u.repo.UpdateAttraction(ctx, attraction); err != nil {
		return nil, errors.ErrDatabase
	}
	return attraction, nil
}

func (u *tourUsecase) DeleteAttraction(ctx context.Context, id uint64) error {
	if err := u.repo.DeleteAttraction(ctx, id); err != nil {
		return errors.ErrDatabase
	}
	return nil
}

func (u *tourUsecase) GenerateAttractionText(ctx context.Context, req *domain.GenerateTextRequest) (*domain.GeneratedText, error) {
	var targetName string
	var targetDesc string

	// 1. 获取景点基础信息 (支持ID或直接传名称)
	if req.AttractionID > 0 {
		attraction, err := u.repo.GetAttractionByID(ctx, req.AttractionID)
		if err != nil {
			return nil, errors.ErrDatabase
		}
		if attraction == nil {
			return nil, errors.ErrNotFound
		}
		targetName = attraction.Name
		targetDesc = attraction.Description
	} else if req.LocationName != "" {
		targetName = req.LocationName
		targetDesc = "热门旅游目的地，各具特色"
	} else {
		return nil, errors.ErrInvalidParam
	}

	// 2. 模拟爬虫抓取小红书内容 (实际项目中应调用 Python 爬虫服务或使用代理IP池)
	originalContent := fmt.Sprintf("模拟抓取到的关于 %s 的小红书游记：风景优美，非常值得打卡！", targetName)

	// 3. 构造大模型生成提示词
	wordCount := req.WordCount
	if wordCount == 0 {
		wordCount = 300 // 默认字数
	}

	systemPrompt := "你是一个专业的旅游内容编辑和社交媒体运营专家。请根据提供的基础信息，创作一篇高质量的原创旅游推文。"
	userPrompt := fmt.Sprintf("目的地: %s\n基础描述: %s\n\n要求:\n1. 必须包含以下要素：【景点特色】、【最佳游览时间】、【交通提示】。\n2. 语言风格生动有吸引力，适合年轻人在社交媒体（如小红书、微博）阅读，可适当使用 Emoji 表情。\n3. 字数控制在 %d 字左右。\n4. 直接输出正文，不要有任何多余的开头或结尾寒暄语。",
		targetName, targetDesc, wordCount)

	messages := []*schema.Message{
		schema.SystemMessage(systemPrompt),
		schema.UserMessage(userPrompt),
	}

	// 4. 调用大模型生成文本
	resp, err := u.llmModel.Generate(ctx, messages)
	if err != nil {
		fmt.Printf("Ollama generate text error: %v\n", err)
		return nil, errors.ErrInternalServer
	}

	generatedContent := resp.Content
	actualWordCount := len([]rune(generatedContent)) // 简单统计字符数

	// 5. 保存生成的文本记录
	genText := &domain.GeneratedText{
		AttractionID:     req.AttractionID, // 如果是任意地点，此处可能为0
		SourceURL:        req.SourceURL,
		OriginalContent:  originalContent,
		GeneratedContent: generatedContent,
		WordCount:        actualWordCount,
		PlagiarismScore:  0.0, // 初始值
	}

	if req.AttractionID > 0 {
		err = u.repo.SaveGeneratedText(ctx, genText)
		if err != nil {
			fmt.Printf("Save generated text error: %v\n", err)
			// 不阻断流程，仅记录日志
		}
	}

	return genText, nil
}
