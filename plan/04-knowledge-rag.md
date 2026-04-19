# 知识库检索与 RAG 实现计划

> 任务编号: 4. 知识库检索与 RAG 实现 [P0]
> 预估工作量: 1-2 周
> 阻塞影响: Agent 无法理解项目上下文

---

## 1. 概述

实现知识库检索与 RAG (Retrieval-Augmented Generation) 机制，让 Agent 能够基于项目知识库提供上下文相关的回答和建议，提高开发效率和准确性。

---

## 2. 技术选型

### 2.1 核心技术
- **向量数据库**: Chroma / Qdrant / Weaviate (可选)
- **嵌入模型**: sentence-transformers / OpenAI Embeddings
- **分词器**: tiktoken / tokenizers
- **文本检索**: BM25 / TF-IDF
- **相似度计算**: 余弦相似度 / 欧氏距离

### 2.2 选型理由
- 向量检索提供语义搜索能力
- 嵌入模型捕捉文本语义
- BM25 提供关键词检索补充
- 混合检索提高准确率

---

## 3. 架构设计

### 3.1 RAG 架构
```
┌─────────────────────────────────────────────────────────┐
│                   Query Layer                           │
│  User Query -> Query Preprocessor                       │
└─────────────────────────────────────────────────────────┘
                           │
┌─────────────────────────────────────────────────────────┐
│               Retrieval Layer                          │
├─────────────────────────────────────────────────────────┤
│  ├── Vector Search (Embedding + Similarity)            │
│  ├── Keyword Search (BM25)                             │
│  ├── Hybrid Search (Re-ranking)                        │
│  └── Context Selection                                 │
└─────────────────────────────────────────────────────────┘
                           │
┌─────────────────────────────────────────────────────────┐
│               RAG Layer                                │
├─────────────────────────────────────────────────────────┤
│  ├── Context Injection                                 │
│  ├── Prompt Template Management                        │
│  ├── Context Length Control                           │
│  └── Multi-turn Context Management                     │
└─────────────────────────────────────────────────────────┘
                           │
┌─────────────────────────────────────────────────────────┐
│              Generation Layer                          │
│  Context + Query -> AI Runtime -> LLM -> Answer        │
└─────────────────────────────────────────────────────────┘
```

---

## 4. 实施步骤

### 4.1 Phase 1: 文档分片算法 (2-3 天)

#### 1.1 智能文本分片
**文件**: `backend/domain/knowledge/chunker.go`

```go
type Chunker interface {
    Chunk(text string, maxTokens int) ([]*Chunk, error)
}

type TextChunker struct {
    tokenizer *tokenizer.Tokenizer
    overlap   int
}

type Chunk struct {
    ID        string
    Content   string
    TokenCount int
    Metadata  map[string]interface{}
}

func NewTextChunker(overlap int) *TextChunker
func (c *TextChunker) Chunk(text string, maxTokens int) ([]*Chunk, error)
func (c *TextChunker) splitByParagraph(text string) []string
func (c *TextChunker) splitBySentence(text string) []string
func (c *TextChunker) mergeShortChunks(chunks []*Chunk, minTokens int) []*Chunk
func (c *TextChunker) addOverlap(chunks []*Chunk, overlapTokens int) []*Chunk
```

#### 1.2 代码文件分片
**文件**: `backend/domain/knowledge/code_chunker.go`

```go
type CodeChunker struct {
    parser    *code.Parser
    tokenizer *tokenizer.Tokenizer
}

type CodeChunk struct {
    ID          string
    FilePath    string
    Language    string
    Type        string  // function, class, module
    Name        string
    Content     string
    StartLine   int
    EndLine     int
    TokenCount  int
    Metadata    map[string]interface{}
}

func NewCodeChunker() *CodeChunker
func (c *CodeChunker) Chunk(code, language, filePath string) ([]*CodeChunk, error)
func (c *CodeChunker) extractFunctions(code, language string) []*CodeChunk
func (c *CodeChunker) extractClasses(code, language string) []*CodeChunk
func (c *CodeChunker) extractImports(code, language string) []*CodeChunk
func (c *CodeChunker) estimateTokenCount(content string) int
```

#### 1.3 分片质量评估
**文件**: `backend/domain/knowledge/chunk_evaluator.go`

```go
type ChunkEvaluator struct{}

type ChunkQuality struct {
    Score        float64
    Completeness float64
    Coherence    float64
    Redundancy   float64
    Issues       []string
}

func NewChunkEvaluator() *ChunkEvaluator
func (e *ChunkEvaluator) Evaluate(chunk *Chunk) *ChunkQuality
func (e *ChunkEvaluator) checkCompleteness(chunk *Chunk) float64
func (e *ChunkEvaluator) checkCoherence(chunk *Chunk) float64
func (e *ChunkEvaluator) checkRedundancy(chunks []*Chunk) []float64
```

---

### 4.2 Phase 2: 索引构建 (2-3 天)

#### 2.1 文本索引构建
**文件**: `backend/infrastructure/knowledge/text_index.go`

```go
type TextIndex struct {
    vectorStore vector.Store
    bm25Index   *BM25Index
}

type BM25Index struct {
    documents  map[string]*Document
    idf        map[string]float64
    avgDocLen  float64
    k1         float64
    b          float64
}

func NewTextIndex(vectorStore vector.Store) *TextIndex
func (i *TextIndex) BuildIndex(chunks []*Chunk) error
func (i *TextIndex) AddChunks(chunks []*Chunk) error
func (i *TextIndex) UpdateChunks(chunks []*Chunk) error
func (i *TextIndex) DeleteChunks(ids []string) error
func (i *TextIndex) GetIndexStats() *IndexStats
```

#### 2.2 代码结构索引
**文件**: `backend/infrastructure/knowledge/code_index.go`

```go
type CodeIndex struct {
    astIndex    *ASTIndex
    callGraph   *CallGraph
    typeGraph   *TypeGraph
}

type ASTIndex struct {
    functions   map[string]*FunctionInfo
    classes     map[string]*ClassInfo
    modules     map[string]*ModuleInfo
}

type FunctionInfo struct {
    ID          string
    Name        string
    File        string
    Lines       []int
    Params      []ParamInfo
    Returns     []TypeInfo
    Calls       []string
    CalledBy    []string
}

func NewCodeIndex() *CodeIndex
func (i *CodeIndex) BuildIndex(chunks []*CodeChunk) error
func (i *CodeIndex) FindFunctions(query string) []*FunctionInfo
func (i *CodeIndex) FindClasses(query string) []*ClassInfo
func (i *CodeIndex) GetCallGraph(functionName string) *CallGraph
```

#### 2.3 元数据索引
**文件**: `backend/infrastructure/knowledge/metadata_index.go`

```go
type MetadataIndex struct {
    index map[string]map[interface{}][]string  // field -> value -> chunk IDs
}

func NewMetadataIndex() *MetadataIndex
func (i *MetadataIndex) AddChunk(id string, metadata map[string]interface{}) error
func (i *MetadataIndex) Query(field string, value interface{}) []string
func (i *MetadataIndex) QueryRange(field string, min, max interface{}) []string
func (i *MetadataIndex) DeleteChunk(id string) error
```

#### 2.4 索引更新机制
**文件**: `backend/infrastructure/knowledge/index_updater.go`

```go
type IndexUpdater struct {
    textIndex   *TextIndex
    codeIndex   *CodeIndex
    metaIndex   *MetadataIndex
    updateQueue chan *UpdateTask
}

type UpdateTask struct {
    Type      string  // add, update, delete
    Chunk     interface{}
}

func NewIndexUpdater(textIndex *TextIndex, codeIndex *CodeIndex, metaIndex *MetadataIndex) *IndexUpdater
func (u *IndexUpdater) Start(ctx context.Context)
func (u *IndexUpdater) AddChunk(chunk interface{}) error
func (u *IndexUpdater) UpdateChunk(chunk interface{}) error
func (u *IndexUpdater) DeleteChunk(id string) error
func (u *IndexUpdater) RebuildIndex() error
```

---

### 4.3 Phase 3: 检索功能 (2-3 天)

#### 3.1 向量检索
**文件**: `backend/application/knowledge/vector_search.go`

```go
type VectorSearch struct {
    embedder   embedding.Embedder
    store      vector.Store
}

type SearchResult struct {
    ChunkID    string
    Content    string
    Score      float64
    Metadata   map[string]interface{}
}

func NewVectorSearch(embedder embedding.Embedder, store vector.Store) *VectorSearch
func (s *VectorSearch) Search(ctx context.Context, query string, topK int) ([]*SearchResult, error)
func (s *VectorSearch) SearchWithFilter(ctx context.Context, query string, filter map[string]interface{}, topK int) ([]*SearchResult, error)
func (s *VectorSearch) GetEmbedding(ctx context.Context, text string) ([]float64, error)
```

#### 3.2 BM25 检索
**文件**: `backend/application/knowledge/bm25_search.go`

```go
type BM25Search struct {
    index *BM25Index
}

func NewBM25Search(index *BM25Index) *BM25Search
func (s *BM25Search) Search(query string, topK int) ([]*SearchResult, error)
func (s *BM25Search) tokenize(text string) []string
func (s *BM25Search) calculateScore(docID string, queryTerms []string) float64
```

#### 3.3 混合检索
**文件**: `backend/application/knowledge/hybrid_search.go`

```go
type HybridSearch struct {
    vectorSearch *VectorSearch
    bm25Search   *BM25Search
    alpha        float64  // vector search weight
}

func NewHybridSearch(vectorSearch *VectorSearch, bm25Search *BM25Search, alpha float64) *HybridSearch
func (s *HybridSearch) Search(ctx context.Context, query string, topK int) ([]*SearchResult, error)
func (s *HybridSearch) reRank(results []*SearchResult, query string) []*SearchResult
```

#### 3.4 相关性评分
**文件**: `backend/application/knowledge/scorer.go`

```go
type Scorer struct{}

func (s *Scorer) CalculateSimilarity(embedding1, embedding2 []float64) float64
func (s *Scorer) CosineSimilarity(a, b []float64) float64
func (s *Scorer) EuclideanDistance(a, b []float64) float64
func (s *HybridSearch) NormalizeScores(scores []float64) []float64
```

---

### 4.4 Phase 4: RAG 机制 (2-3 天)

#### 4.1 上下文注入策略
**文件**: `backend/application/knowledge/context_injector.go`

```go
type ContextInjector struct {
    searchService SearchService
    templateMgr   *template.Manager
}

type ContextConfig struct {
    MaxChunks       int
    MaxTokens       int
    MinSimilarity   float64
    IncludeMetadata bool
}

type InjectedContext struct {
    Query      string
    Chunks     []*Chunk
    Context    string
    TokenCount int
}

func NewContextInjector(searchService SearchService, templateMgr *template.Manager) *ContextInjector
func (i *ContextInjector) Inject(ctx context.Context, query string, config *ContextConfig) (*InjectedContext, error)
func (i *ContextInjector) SelectChunks(chunks []*Chunk, maxTokens int) []*Chunk
func (i *ContextInjector) FormatContext(chunks []*Chunk) string
```

#### 4.2 提示词模板管理
**文件**: `backend/application/knowledge/template_manager.go`

```go
type TemplateManager struct {
    templates map[string]*PromptTemplate
}

type PromptTemplate struct {
    Name        string
    Description string
    Template    string
    Variables   []string
    MaxTokens   int
}

func NewTemplateManager() *TemplateManager
func (m *TemplateManager) LoadTemplates(dir string) error
func (m *TemplateManager) GetTemplate(name string) (*PromptTemplate, error)
func (m *TemplateManager) Render(name string, vars map[string]interface{}) (string, error)
func (m *TemplateManager) AddTemplate(template *PromptTemplate) error
```

#### 4.3 上下文长度控制
**文件**: `backend/application/knowledge/context_controller.go`

```go
type ContextController struct {
    tokenizer *tokenizer.Tokenizer
}

func NewContextController(tokenizer *tokenizer.Tokenizer) *ContextController
func (c *ContextController) TruncateContext(context string, maxTokens int) string
func (c *ContextController) OptimizeChunks(chunks []*Chunk, maxTokens int) []*Chunk
func (c *ContextController) EstimateTokenCount(text string) int
```

#### 4.4 多轮对话上下文管理
**文件**: `backend/application/knowledge/conversation_manager.go`

```go
type ConversationManager struct {
    conversations map[string]*Conversation
    maxHistory    int
    maxTokens     int
}

type Conversation struct {
    ID      string
    Messages []*Message
    Context map[string]interface{}
}

type Message struct {
    Role      string  // user, assistant, system
    Content   string
    Timestamp time.Time
    TokenCount int
}

func NewConversationManager(maxHistory, maxTokens int) *ConversationManager
func (m *ConversationManager) CreateConversation(id string) *Conversation
func (m *ConversationManager) GetConversation(id string) *Conversation
func (m *ConversationManager) AddMessage(id string, message *Message) error
func (m *ConversationManager) GetContext(id string) (string, error)
func (m *ConversationManager) TrimConversation(id string) error
```

---

## 5. 时间安排

| 阶段 | 工作内容 | 预计时间 |
|------|----------|----------|
| Phase 1 | 文档分片算法 | 2-3 天 |
| Phase 2 | 索引构建 | 2-3 天 |
| Phase 3 | 检索功能 | 2-3 天 |
| Phase 4 | RAG 机制 | 2-3 天 |

**总计**: 1-2 周

---

## 6. 验收标准

### 6.1 功能验收
- [ ] 文档分片算法正常工作
- [ ] 代码文件分片正确识别结构
- [ ] 索引构建和更新正常
- [ ] 向量检索准确
- [ ] BM25 检索正常
- [ ] 混合检索结果准确
- [ ] 上下文注入正确
- [ ] 提示词模板管理正常
- [ ] 多轮对话上下文管理正常

### 6.2 性能验收
- [ ] 文档分片速度 ≥ 1000 文档/分钟
- [ ] 索引构建速度 ≥ 500 文档/分钟
- [ ] 检索响应时间 < 500ms
- [ ] 上下文注入时间 < 200ms
- [ ] 支持并发检索

### 6.3 质量验收
- [ ] 单元测试覆盖率 ≥ 80%
- [ ] 检索准确率 ≥ 85%
- [ ] 代码审查通过
- [ ] 文档完整

---

## 7. 风险与缓解

### 7.1 技术风险
**风险**: 向量检索准确率不稳定
**缓解**:
- 使用混合检索提高准确率
- 实现重排序机制
- 优化嵌入模型

### 7.2 性能风险
**风险**: 大规模数据检索性能差
**缓解**:
- 实现索引分片
- 使用缓存机制
- 优化查询性能

### 7.3 质量风险
**风险**: 上下文注入不准确
**缓解**:
- 实现上下文质量评估
- 提供人工审核机制
- 优化上下文选择策略
