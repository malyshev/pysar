package mcpserver

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

// newTestServer builds a Server with one dummy tool registered directly,
// bypassing New()'s real tool registration -- these tests exercise the
// JSON-RPC dispatch mechanics, not save_voice_profile specifically.
func newTestServer(out *bytes.Buffer) *Server {
	s := &Server{name: "test", version: "0.0.0", handlers: make(map[string]toolHandler)}
	s.out = out
	s.register(
		tool{Name: "echo", Description: "echoes input", InputSchema: map[string]string{"type": "object"}},
		func(args json.RawMessage) callToolResult {
			return textResult("got: %s", string(args))
		},
	)
	return s
}

func runLines(t *testing.T, s *Server, lines ...string) []map[string]interface{} {
	t.Helper()
	s.in = strings.NewReader(strings.Join(lines, "\n") + "\n")
	if err := s.Run(); err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	out := s.out.(*bytes.Buffer)
	var responses []map[string]interface{}
	for _, line := range strings.Split(strings.TrimSpace(out.String()), "\n") {
		if line == "" {
			continue
		}
		var resp map[string]interface{}
		if err := json.Unmarshal([]byte(line), &resp); err != nil {
			t.Fatalf("response line did not parse as JSON: %v\nline: %s", err, line)
		}
		responses = append(responses, resp)
	}
	return responses
}

func TestInitializeReturnsServerInfo(t *testing.T) {
	s := newTestServer(&bytes.Buffer{})
	resp := runLines(t, s, `{"jsonrpc":"2.0","method":"initialize","id":1}`)

	if len(resp) != 1 {
		t.Fatalf("expected 1 response, got %d", len(resp))
	}
	result, ok := resp[0]["result"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected result object, got %v", resp[0])
	}
	if result["protocolVersion"] != protocolVersion {
		t.Fatalf("unexpected protocolVersion: %v", result["protocolVersion"])
	}
	serverInfo, ok := result["serverInfo"].(map[string]interface{})
	if !ok || serverInfo["name"] != "test" {
		t.Fatalf("unexpected serverInfo: %v", result["serverInfo"])
	}
}

func TestToolsListReturnsRegisteredTool(t *testing.T) {
	s := newTestServer(&bytes.Buffer{})
	resp := runLines(t, s, `{"jsonrpc":"2.0","method":"tools/list","id":1}`)

	result := resp[0]["result"].(map[string]interface{})
	tools, ok := result["tools"].([]interface{})
	if !ok || len(tools) != 1 {
		t.Fatalf("expected 1 tool, got %v", result["tools"])
	}
	first := tools[0].(map[string]interface{})
	if first["name"] != "echo" {
		t.Fatalf("expected tool named echo, got %v", first["name"])
	}
}

func TestToolsCallDispatchesToHandler(t *testing.T) {
	s := newTestServer(&bytes.Buffer{})
	resp := runLines(t, s, `{"jsonrpc":"2.0","method":"tools/call","id":1,"params":{"name":"echo","arguments":{"x":1}}}`)

	result := resp[0]["result"].(map[string]interface{})
	content := result["content"].([]interface{})[0].(map[string]interface{})
	if !strings.Contains(content["text"].(string), `"x":1`) {
		t.Fatalf("expected echoed arguments in text, got %v", content["text"])
	}
	if result["isError"] == true {
		t.Fatalf("expected isError false/absent, got true")
	}
}

func TestToolsCallUnknownToolReturnsIsError(t *testing.T) {
	s := newTestServer(&bytes.Buffer{})
	resp := runLines(t, s, `{"jsonrpc":"2.0","method":"tools/call","id":1,"params":{"name":"nope","arguments":{}}}`)

	result := resp[0]["result"].(map[string]interface{})
	if result["isError"] != true {
		t.Fatalf("expected isError true for unknown tool, got %v", result)
	}
}

func TestUnknownMethodReturnsJSONRPCError(t *testing.T) {
	s := newTestServer(&bytes.Buffer{})
	resp := runLines(t, s, `{"jsonrpc":"2.0","method":"bogus/method","id":1}`)

	if resp[0]["error"] == nil {
		t.Fatalf("expected a JSON-RPC error for unknown method, got %v", resp[0])
	}
}

func TestMalformedLineDoesNotCrashServer(t *testing.T) {
	s := newTestServer(&bytes.Buffer{})
	resp := runLines(t, s,
		`not valid json`,
		`{"jsonrpc":"2.0","method":"initialize","id":1}`,
	)

	// One parse-error response (id null) plus one real initialize response.
	if len(resp) != 2 {
		t.Fatalf("expected 2 responses (parse error + initialize), got %d: %v", len(resp), resp)
	}
	if resp[1]["result"] == nil {
		t.Fatalf("expected the second (valid) request to still succeed, got %v", resp[1])
	}
}
