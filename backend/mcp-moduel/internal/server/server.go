package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	"github.com/google/uuid"
	"github.com/mark3labs/mcp-go/mcp"
)

// ServerConfig 服务器配置
type ServerConfig struct {
	// Address 服务器监听地址，格式: "host:port"
	Address string

	// ServerName 服务器名称
	ServerName string

	// ServerVersion 服务器版本
	ServerVersion string

	// EnableCORS 是否启用 CORS
	EnableCORS bool

	// AllowedOrigins CORS 允许的来源
	AllowedOrigins []string
}

// DefaultServerConfig 返回默认服务器配置
func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Address:       ":8080",
		ServerName:    "mcp-tool-server",
		ServerVersion: "1.0.0",
		EnableCORS:    true,
		AllowedOrigins: []string{
			"*",
		},
	}
}

// Validate 验证配置
func (c *ServerConfig) Validate() error {
	if c.Address == "" {
		c.Address = ":8080"
	}
	if c.ServerName == "" {
		c.ServerName = "mcp-tool-server"
	}
	if c.ServerVersion == "" {
		c.ServerVersion = "1.0.0"
	}
	return nil
}

// SSESession SSE 会话
type SSESession struct {
	ID          string
	EventChan   chan string
	RequestChan chan *JSONRPCRequest
	Done        chan struct{}
	ResponseMap sync.Map // map[string]chan interface{}
	mu          sync.RWMutex
}

// Server MCP 服务器实现
type Server struct {
	// 工具注册表
	toolRegistry *ToolRegistry

	// 服务器配置
	config *ServerConfig

	// 请求处理锁
	mu sync.RWMutex

	// HTTP 服务器
	httpServer *http.Server

	// SSE 会话管理
	sseSessions sync.Map // map[string]*SSESession
}

// NewServer 创建新的 MCP 服务器实例
func NewServer(config *ServerConfig) *Server {
	if config == nil {
		config = DefaultServerConfig()
	}

	return &Server{
		toolRegistry: NewToolRegistry(),
		config:       config,
	}
}

// RegisterTool 注册工具到服务器
func (s *Server) RegisterTool(tool Tool) error {
	return s.toolRegistry.Register(tool)
}

// RegisterTools 批量注册工具
func (s *Server) RegisterTools(tools []Tool) error {
	for _, tool := range tools {
		if err := s.toolRegistry.Register(tool); err != nil {
			return fmt.Errorf("failed to register tool %s: %w", tool.Name(), err)
		}
	}
	return nil
}

// Start 启动 HTTP 服务器
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// 注册 MCP 协议端点
	mux.HandleFunc("/mcp", s.handleMCPRequest)
	mux.HandleFunc("/mcp/sse", s.handleSSEConnection)
	mux.HandleFunc("/mcp/sse/", s.handleSSERequest)
	mux.HandleFunc("/health", s.handleHealthCheck)

	// 创建 HTTP 服务器
	s.httpServer = &http.Server{
		Addr:    s.config.Address,
		Handler: mux,
	}

	log.Printf("MCP Server starting on %s", s.config.Address)
	return s.httpServer.ListenAndServe()
}

// Stop 停止服务器
func (s *Server) Stop(ctx context.Context) error {
	if s.httpServer == nil {
		return nil
	}
	return s.httpServer.Shutdown(ctx)
}

// JSONRPCRequest JSON-RPC 2.0 请求结构
type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// JSONRPCError JSON-RPC 2.0 错误结构
type JSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

// JSONRPCErrorResponse JSON-RPC 2.0 错误响应
type JSONRPCErrorResponse struct {
	JSONRPC string       `json:"jsonrpc"`
	ID      interface{}  `json:"id"`
	Error   JSONRPCError `json:"error"`
}

// handleMCPRequest 处理 MCP 协议请求
func (s *Server) handleMCPRequest(w http.ResponseWriter, r *http.Request) {
	// 只接受 POST 请求
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")

	// 解析 JSON-RPC 请求
	var jsonrpcReq JSONRPCRequest
	if err := json.NewDecoder(r.Body).Decode(&jsonrpcReq); err != nil {
		s.sendJSONRPCError(w, nil, -32700, "Parse error", err.Error())
		return
	}

	// 处理不同的请求方法
	switch jsonrpcReq.Method {
	case "initialize":
		s.handleInitialize(w, &jsonrpcReq)
	case "tools/list":
		s.handleListTools(w, &jsonrpcReq)
	case "tools/call":
		s.handleCallTool(w, &jsonrpcReq)
	default:
		s.sendJSONRPCError(w, jsonrpcReq.ID, -32601, "Method not found",
			fmt.Sprintf("Unknown method: %s", jsonrpcReq.Method))
	}
}

// handleInitialize 处理初始化请求
func (s *Server) handleInitialize(w http.ResponseWriter, request *JSONRPCRequest) {
	var params mcp.InitializeParams
	if len(request.Params) > 0 {
		if err := sonic.Unmarshal(request.Params, &params); err != nil {
			s.sendJSONRPCError(w, request.ID, -32602, "Invalid params", err.Error())
			return
		}
	}

	result := mcp.InitializeResult{
		ProtocolVersion: mcp.LATEST_PROTOCOL_VERSION,
		ServerInfo: mcp.Implementation{
			Name:    s.config.ServerName,
			Version: s.config.ServerVersion,
		},
		Capabilities: mcp.ServerCapabilities{
			Tools: &struct {
				ListChanged bool `json:"listChanged,omitempty"`
			}{},
		},
	}

	requestID := mcp.NewRequestId(request.ID)
	response := mcp.NewJSONRPCResultResponse(requestID, result)
	s.sendResponse(w, response)
}

// handleListTools 处理工具列表请求
func (s *Server) handleListTools(w http.ResponseWriter, request *JSONRPCRequest) {
	tools := s.toolRegistry.ListTools()

	// 转换为 MCP 格式的工具列表
	mcpTools := make([]mcp.Tool, 0, len(tools))
	for _, tool := range tools {
		// 将 map[string]interface{} 转换为 ToolInputSchema
		schema := tool.InputSchema()
		schemaBytes, _ := sonic.Marshal(schema)

		mcpTool := mcp.NewToolWithRawSchema(
			tool.Name(),
			tool.Description(),
			schemaBytes,
		)
		mcpTools = append(mcpTools, mcpTool)
	}

	result := mcp.ListToolsResult{
		Tools: mcpTools,
	}

	requestID := mcp.NewRequestId(request.ID)
	response := mcp.NewJSONRPCResultResponse(requestID, result)
	s.sendResponse(w, response)
}

// handleCallTool 处理工具调用请求
func (s *Server) handleCallTool(w http.ResponseWriter, request *JSONRPCRequest) {
	var params mcp.CallToolParams
	if err := sonic.Unmarshal(request.Params, &params); err != nil {
		s.sendJSONRPCError(w, request.ID, -32602, "Invalid params", err.Error())
		return
	}

	// 获取工具
	tool, err := s.toolRegistry.GetTool(params.Name)
	if err != nil {
		s.sendJSONRPCError(w, request.ID, -32601, "Tool not found", err.Error())
		return
	}

	// 执行工具
	ctx := context.Background()

	// 类型断言参数
	arguments, ok := params.Arguments.(map[string]interface{})
	if !ok {
		// 如果参数为空，使用空 map
		if params.Arguments == nil {
			arguments = make(map[string]interface{})
		} else {
			s.sendJSONRPCError(w, request.ID, -32602, "Invalid params", "arguments must be a map")
			return
		}
	}

	result, err := tool.Execute(ctx, arguments)
	if err != nil {
		// 返回错误结果
		errorResult := mcp.NewToolResultError(fmt.Sprintf("Tool execution error: %v", err))
		requestID := mcp.NewRequestId(request.ID)
		response := mcp.NewJSONRPCResultResponse(requestID, errorResult)
		s.sendResponse(w, response)
		return
	}

	// 返回成功结果
	if result == nil {
		errorResult := mcp.NewToolResultError("Tool returned nil result")
		requestID := mcp.NewRequestId(request.ID)
		response := mcp.NewJSONRPCResultResponse(requestID, errorResult)
		s.sendResponse(w, response)
		return
	}

	requestID := mcp.NewRequestId(request.ID)
	response := mcp.NewJSONRPCResultResponse(requestID, result)
	s.sendResponse(w, response)
}

// handleHealthCheck 处理健康检查请求
func (s *Server) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "healthy",
		"tools":   s.toolRegistry.Count(),
		"version": s.config.ServerVersion,
	})
}

// sendResponse 发送成功响应
func (s *Server) sendResponse(w http.ResponseWriter, response interface{}) {
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

// sendJSONRPCError 发送 JSON-RPC 错误响应
func (s *Server) sendJSONRPCError(w http.ResponseWriter, id interface{}, code int, message, data string) {
	response := JSONRPCErrorResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: JSONRPCError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}

	w.WriteHeader(http.StatusOK) // MCP 协议错误也返回 200
	json.NewEncoder(w).Encode(response)
}

// handleSSEConnection 处理 SSE 连接
func (s *Server) handleSSEConnection(w http.ResponseWriter, r *http.Request) {
	// 只接受 GET 请求
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 创建 SSE 会话
	sessionID := uuid.New().String()
	session := &SSESession{
		ID:          sessionID,
		EventChan:   make(chan string, 100),
		RequestChan: make(chan *JSONRPCRequest, 100),
		Done:        make(chan struct{}),
	}

	s.sseSessions.Store(sessionID, session)

	// 设置 SSE 响应头
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Cache-Control")

	// 发送会话 ID
	fmt.Fprintf(w, "data: %s\n\n", fmt.Sprintf(`{"type":"session","sessionId":"%s"}`, sessionID))
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}

	// 启动响应处理协程
	go s.handleSSEResponses(session, w)

	// 等待连接关闭
	<-r.Context().Done()
	close(session.Done)
	s.sseSessions.Delete(sessionID)
	log.Printf("SSE session %s closed", sessionID)
}

// handleSSERequest 处理通过 SSE 发送的请求
func (s *Server) handleSSERequest(w http.ResponseWriter, r *http.Request) {
	// 只接受 POST 请求
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 从 URL 路径获取会话 ID
	sessionID := r.URL.Path[len("/mcp/sse/"):]
	if sessionID == "" {
		http.Error(w, "Session ID required", http.StatusBadRequest)
		return
	}

	// 获取会话
	sessionInterface, ok := s.sseSessions.Load(sessionID)
	if !ok {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}
	session := sessionInterface.(*SSESession)

	// 解析 JSON-RPC 请求
	var jsonrpcReq JSONRPCRequest
	if err := json.NewDecoder(r.Body).Decode(&jsonrpcReq); err != nil {
		http.Error(w, "Invalid JSON-RPC request", http.StatusBadRequest)
		return
	}

	// 将请求发送到会话的请求通道
	select {
	case session.RequestChan <- &jsonrpcReq:
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("Request queued"))
	default:
		http.Error(w, "Request queue full", http.StatusServiceUnavailable)
	}
}

// handleSSEResponses 处理 SSE 响应
func (s *Server) handleSSEResponses(session *SSESession, w http.ResponseWriter) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		return
	}

	// 处理请求
	go func() {
		for {
			select {
			case <-session.Done:
				return
			case req := <-session.RequestChan:
				s.processSSERequest(session, req)
			}
		}
	}()

	// 发送事件
	for {
		select {
		case <-session.Done:
			return
		case event := <-session.EventChan:
			fmt.Fprintf(w, "data: %s\n\n", event)
			flusher.Flush()
		}
	}
}

// processSSERequest 处理 SSE 请求
func (s *Server) processSSERequest(session *SSESession, request *JSONRPCRequest) {
	var response interface{}

	// 处理不同的请求方法
	switch request.Method {
	case "initialize":
		response = s.handleInitializeSSE(request)
	case "tools/list":
		response = s.handleListToolsSSE(request)
	case "tools/call":
		response = s.handleCallToolSSE(request)
	default:
		response = JSONRPCErrorResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: JSONRPCError{
				Code:    -32601,
				Message: "Method not found",
				Data:    fmt.Sprintf("Unknown method: %s", request.Method),
			},
		}
	}

	// 序列化响应
	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Printf("Failed to marshal SSE response: %v", err)
		return
	}

	// 发送响应
	select {
	case session.EventChan <- string(responseJSON):
	case <-time.After(5 * time.Second):
		log.Printf("Timeout sending SSE response for request %v", request.ID)
	}
}

// handleInitializeSSE 处理 SSE 初始化请求
func (s *Server) handleInitializeSSE(request *JSONRPCRequest) interface{} {
	var params mcp.InitializeParams
	if len(request.Params) > 0 {
		if err := sonic.Unmarshal(request.Params, &params); err != nil {
			return JSONRPCErrorResponse{
				JSONRPC: "2.0",
				ID:      request.ID,
				Error: JSONRPCError{
					Code:    -32602,
					Message: "Invalid params",
					Data:    err.Error(),
				},
			}
		}
	}

	result := mcp.InitializeResult{
		ProtocolVersion: mcp.LATEST_PROTOCOL_VERSION,
		ServerInfo: mcp.Implementation{
			Name:    s.config.ServerName,
			Version: s.config.ServerVersion,
		},
		Capabilities: mcp.ServerCapabilities{
			Tools: &struct {
				ListChanged bool `json:"listChanged,omitempty"`
			}{},
		},
	}

	requestID := mcp.NewRequestId(request.ID)
	return mcp.NewJSONRPCResultResponse(requestID, result)
}

// handleListToolsSSE 处理 SSE 工具列表请求
func (s *Server) handleListToolsSSE(request *JSONRPCRequest) interface{} {
	tools := s.toolRegistry.ListTools()

	// 转换为 MCP 格式的工具列表
	mcpTools := make([]mcp.Tool, 0, len(tools))
	for _, tool := range tools {
		schema := tool.InputSchema()
		schemaBytes, _ := sonic.Marshal(schema)

		mcpTool := mcp.NewToolWithRawSchema(
			tool.Name(),
			tool.Description(),
			schemaBytes,
		)
		mcpTools = append(mcpTools, mcpTool)
	}

	result := mcp.ListToolsResult{
		Tools: mcpTools,
	}

	requestID := mcp.NewRequestId(request.ID)
	return mcp.NewJSONRPCResultResponse(requestID, result)
}

// handleCallToolSSE 处理 SSE 工具调用请求
func (s *Server) handleCallToolSSE(request *JSONRPCRequest) interface{} {
	var params mcp.CallToolParams
	if err := sonic.Unmarshal(request.Params, &params); err != nil {
		return JSONRPCErrorResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: JSONRPCError{
				Code:    -32602,
				Message: "Invalid params",
				Data:    err.Error(),
			},
		}
	}

	// 获取工具
	tool, err := s.toolRegistry.GetTool(params.Name)
	if err != nil {
		return JSONRPCErrorResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: JSONRPCError{
				Code:    -32601,
				Message: "Tool not found",
				Data:    err.Error(),
			},
		}
	}

	// 执行工具
	ctx := context.Background()

	// 类型断言参数
	arguments, ok := params.Arguments.(map[string]interface{})
	if !ok {
		if params.Arguments == nil {
			arguments = make(map[string]interface{})
		} else {
			return JSONRPCErrorResponse{
				JSONRPC: "2.0",
				ID:      request.ID,
				Error: JSONRPCError{
					Code:    -32602,
					Message: "Invalid params",
					Data:    "arguments must be a map",
				},
			}
		}
	}

	result, err := tool.Execute(ctx, arguments)
	if err != nil {
		errorResult := mcp.NewToolResultError(fmt.Sprintf("Tool execution error: %v", err))
		requestID := mcp.NewRequestId(request.ID)
		return mcp.NewJSONRPCResultResponse(requestID, errorResult)
	}

	if result == nil {
		errorResult := mcp.NewToolResultError("Tool returned nil result")
		requestID := mcp.NewRequestId(request.ID)
		return mcp.NewJSONRPCResultResponse(requestID, errorResult)
	}

	requestID := mcp.NewRequestId(request.ID)
	return mcp.NewJSONRPCResultResponse(requestID, result)
}
