package milvus

import (
	"context"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino-ext/components/indexer/milvus"
	"github.com/cloudwego/eino/schema"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

//参考文档 https://www.cloudwego.io/zh/docs/eino/ecosystem_integration/indexer/indexer_milvus/

// floatVectorSchema 定义浮点向量的数据结构（相当于表的结构）
type floatVectorSchema struct {
	ID       string    `json:"id" milvus:"name:id"`
	Content  string    `json:"content" milvus:"name:content"`
	Vector   []float32 `json:"vector" milvus:"name:vector"`
	Metadata []byte    `json:"metadata" milvus:"name:metadata"`
}

// floatVectorDocumentConverter 将 schema.Document 转换为浮点向量格式
func floatVectorDocumentConverter(ctx context.Context, docs []*schema.Document, vectors [][]float64) ([]interface{}, error) {
	if len(docs) != len(vectors) {
		return nil, fmt.Errorf("docs and vectors length mismatch: %d != %d", len(docs), len(vectors))
	}
	rows := make([]interface{}, 0, len(docs))
	for idx, doc := range docs {
		// 序列化 metadata
		metadata, err := sonic.Marshal(doc.MetaData)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal metadata: %w", err)
		}
		vector := make([]float32, len(vectors[idx]))
		for i, v := range vectors[idx] {
			vector[i] = float32(v)
		}
		row := &floatVectorSchema{
			ID:       doc.ID,
			Content:  doc.Content,
			Vector:   vector,
			Metadata: metadata,
		}
		rows = append(rows, row)
	}

	return rows, nil
}

// IndexerService 封装 Milvus 索引器服务
type IndexerService struct {
	indexer *milvus.Indexer
	config  *milvus.IndexerConfig
}

// NewIndexerServiceWithDimension 创建新的索引器服务（指定维度）
func NewIndexerServiceWithDimension(ctx context.Context, config *milvus.IndexerConfig, dimension int) (*IndexerService, error) {
	if config == nil {
		return nil, fmt.Errorf("config is nil")
	}
	if config.Client == nil {
		return nil, fmt.Errorf("milvus client is nil")
	}
	if config.Collection == "" {
		return nil, fmt.Errorf("collection name is required")
	}
	if config.Embedding == nil {
		return nil, fmt.Errorf("embedding is required")
	}
	if dimension <= 0 {
		return nil, fmt.Errorf("dimension must be positive, got %d", dimension)
	}
	// 构建 Milvus Indexer 配置
	// 使用浮点向量字段以支持 HNSW 等高效索引
	indexerConfig := &milvus.IndexerConfig{
		Client:     config.Client,
		Collection: config.Collection,
		Embedding:  config.Embedding,
		// 显式指定字段，使用 FloatVector 而非默认的 BinaryVector
		Fields: []*entity.Field{
			entity.NewField().
				WithName("id").
				WithDescription("the unique id of the document").
				WithIsPrimaryKey(true).
				WithDataType(entity.FieldTypeVarChar).
				WithMaxLength(255),
			entity.NewField().
				WithName("vector").
				WithDescription("the vector of the document").
				WithIsPrimaryKey(false).
				WithDataType(entity.FieldTypeFloatVector). // 使用浮点向量
				WithDim(int64(dimension)),                 // 使用传入的向量维度
			entity.NewField().
				WithName("content").
				WithDescription("the content of the document").
				WithIsPrimaryKey(false).
				WithDataType(entity.FieldTypeVarChar).
				WithMaxLength(4096), // 增加长度以支持更长的文档
			entity.NewField().
				WithName("metadata").
				WithDescription("the metadata of the document").
				WithIsPrimaryKey(false).
				WithDataType(entity.FieldTypeJSON),
		},
		// 使用 L2 或 COSINE 度量类型（适用于浮点向量）
		MetricType:        milvus.COSINE,
		DocumentConverter: floatVectorDocumentConverter,
	}
	// 创建 Indexer
	indexer, err := milvus.NewIndexer(ctx, indexerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create milvus indexer: %w", err)
	}
	return &IndexerService{
		indexer: indexer,
		config:  config,
	}, nil
}

// Store 存储文档到 Milvus
func (s *IndexerService) Store(ctx context.Context, docs []*schema.Document) ([]string, error) {
	if s.indexer == nil {
		return nil, fmt.Errorf("indexer is not initialized")
	}
	ids, err := s.indexer.Store(ctx, docs)
	if err != nil {
		return nil, fmt.Errorf("failed to store documents: %w", err)
	}
	return ids, nil
}

// GetIndexer 获取底层的 Indexer 实例
func (s *IndexerService) GetIndexer() *milvus.Indexer {
	return s.indexer
}

// GetConfig 获取配置
func (s *IndexerService) GetConfig() *milvus.IndexerConfig {
	return s.config
}
