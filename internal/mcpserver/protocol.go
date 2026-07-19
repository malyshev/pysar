// Package mcpserver is a minimal MCP-over-stdio server, hand-rolled against
// the same JSON-RPC wire shape haft's own MCP server uses (initialize,
// tools/list, tools/call, newline-delimited JSON) -- no third-party MCP SDK,
// matching this project's plain-Go-toolchain preference (dec-20260718-817bfb74)
// and haft's own proven precedent of implementing the protocol directly.
package mcpserver

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
)

const protocolVersion = "2024-11-05"

type jsonrpcRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
	ID      interface{}     `json:"id"`
}

type jsonrpcResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *rpcError   `json:"error,omitempty"`
	ID      interface{} `json:"id"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// tool describes one callable tool for tools/list.
type tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
}

// callToolResult is what a tool handler returns for tools/call.
type callToolResult struct {
	Content []contentItem `json:"content"`
	IsError bool          `json:"isError,omitempty"`
}

type contentItem struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func errorResult(format string, args ...interface{}) callToolResult {
	return callToolResult{
		Content: []contentItem{{Type: "text", Text: fmt.Sprintf(format, args...)}},
		IsError: true,
	}
}

func textResult(format string, args ...interface{}) callToolResult {
	return callToolResult{Content: []contentItem{{Type: "text", Text: fmt.Sprintf(format, args...)}}}
}

type toolHandler func(args json.RawMessage) callToolResult

// Server is a stateful MCP server over an arbitrary reader/writer pair
// (production: stdin/stdout; tests: in-memory buffers, never real stdio).
type Server struct {
	name    string
	version string
	baseDir string
	homeDir string
	in      io.Reader
	out     io.Writer

	tools    []tool
	handlers map[string]toolHandler
}

// New creates a Server with every Pysar tool registered. baseDir is the
// project root every project-scoped tool's file operations are relative to
// (e.g. .pysar/voice.md). homeDir is the operator's home directory, used by
// tools that manage cross-project state (e.g. ~/.pysar/templates/,
// dec-20260719-3e36577e) rather than this one project.
func New(name, version, baseDir, homeDir string, in io.Reader, out io.Writer) *Server {
	s := &Server{
		name:     name,
		version:  version,
		baseDir:  baseDir,
		homeDir:  homeDir,
		in:       in,
		out:      out,
		handlers: make(map[string]toolHandler),
	}
	s.registerSaveVoiceProfile()
	s.registerSaveVoiceTemplate()
	s.registerListVoiceTemplates()
	s.registerSaveStyleProfile()
	s.registerSaveStyleTemplate()
	s.registerListStyleTemplates()
	return s
}

func (s *Server) register(t tool, h toolHandler) {
	s.tools = append(s.tools, t)
	s.handlers[t.Name] = h
}

// Run reads newline-delimited JSON-RPC requests until EOF or a read error.
func (s *Server) Run() error {
	scanner := bufio.NewScanner(s.in)
	scanner.Buffer(make([]byte, 1<<20), 1<<20)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var req jsonrpcRequest
		if err := json.Unmarshal(line, &req); err != nil {
			s.sendError(nil, -32700, "Parse error")
			continue
		}
		s.handle(req)
	}
	return scanner.Err()
}

func (s *Server) handle(req jsonrpcRequest) {
	defer func() {
		if r := recover(); r != nil && req.ID != nil {
			s.sendError(req.ID, -32603, fmt.Sprintf("internal error: %v", r))
		}
	}()

	switch req.Method {
	case "initialize":
		s.sendResult(req.ID, map[string]interface{}{
			"protocolVersion": protocolVersion,
			"capabilities":    map[string]interface{}{"tools": map[string]interface{}{}},
			"serverInfo":      map[string]string{"name": s.name, "version": s.version},
		})
	case "tools/list":
		s.sendResult(req.ID, map[string]interface{}{"tools": s.tools})
	case "tools/call":
		s.handleToolsCall(req)
	case "notifications/initialized":
		// no-op
	default:
		if req.ID != nil {
			s.sendError(req.ID, -32601, "Method not found")
		}
	}
}

func (s *Server) handleToolsCall(req jsonrpcRequest) {
	var params struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}
	if err := json.Unmarshal(req.Params, &params); err != nil {
		s.sendError(req.ID, -32700, "Invalid params")
		return
	}

	h, ok := s.handlers[params.Name]
	if !ok {
		s.sendResult(req.ID, errorResult("unknown tool: %s", params.Name))
		return
	}
	s.sendResult(req.ID, h(params.Arguments))
}

func (s *Server) sendResult(id interface{}, result interface{}) {
	s.send(jsonrpcResponse{JSONRPC: "2.0", ID: id, Result: result})
}

func (s *Server) sendError(id interface{}, code int, message string) {
	s.send(jsonrpcResponse{JSONRPC: "2.0", ID: id, Error: &rpcError{Code: code, Message: message}})
}

func (s *Server) send(resp jsonrpcResponse) {
	data, err := json.Marshal(resp)
	if err != nil {
		return
	}
	fmt.Fprintf(s.out, "%s\n", data)
}
