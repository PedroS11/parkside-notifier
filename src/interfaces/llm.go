package interfaces

type Payload struct {
	Model  string   `json:"model"`
	Prompt string   `json:"prompt"`
	Stream bool     `json:"stream"`
	Images []string `json:"images"`
}

type Response struct {
	Model              string `json:"model"`
	CreatedAt          string `json:"created_at"`
	ResponseData       string `json:"response"`
	Done               bool   `json:"done"`
	DoneReason         string `json:"done_reason"`
	Context            []int  `json:"context"`
	TotalDuration      int64  `json:"total_duration"`
	LoadDuration       int64  `json:"load_duration"`
	PromptEvalCount    int    `json:"prompt_eval_count"`
	PromptEvalDuration int64  `json:"prompt_eval_duration"`
	EvalCount          int    `json:"eval_count"`
	EvalDuration       int64  `json:"eval_duration"`
}

type Product struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type OpenAIInput struct {
	Role    string        `json:"role"`
	Content []interface{} `json:"content"`
}

type OpenAITextContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type OpenAIImageContent struct {
	Type     string `json:"type"`
	ImageURL string `json:"image_url"`
}

type OpenAIRequest struct {
	Model        string        `json:"model"`
	Input        []OpenAIInput `json:"input"`
	Instructions string        `json:"instructions"`
}

// RESPONSE

// Response is the top-level object returned by the /v1/responses endpoint.
type OpenAIResponse struct {
	ID                 string        `json:"id"`
	Object             string        `json:"object"`
	CreatedAt          int64         `json:"created_at"`
	Status             string        `json:"status"`
	Error              Error         `json:"error"`              // interface{} for null
	IncompleteDetails  interface{}   `json:"incomplete_details"` // interface{} for null
	Instructions       interface{}   `json:"instructions"`       // interface{} for null
	MaxOutputTokens    interface{}   `json:"max_output_tokens"`  // interface{} for null
	Model              string        `json:"model"`
	Output             []OutputItem  `json:"output"`
	ParallelToolCalls  bool          `json:"parallel_tool_calls"`
	PreviousResponseID interface{}   `json:"previous_response_id"` // interface{} for null
	Reasoning          Reasoning     `json:"reasoning"`
	Store              bool          `json:"store"`
	Temperature        float64       `json:"temperature"`
	Text               Text          `json:"text"`
	ToolChoice         interface{}   `json:"tool_choice"` // string ("auto") or object
	Tools              []interface{} `json:"tools"`       // []interface{} for empty array
	TopP               float64       `json:"top_p"`
	Truncation         string        `json:"truncation"`
	Usage              Usage         `json:"usage"`
	User               interface{}   `json:"user"`     // interface{} for null
	Metadata           interface{}   `json:"metadata"` // map[string]interface{} or interface{}
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// OutputItem represents an item in the 'output' array.
type OutputItem struct {
	Type    string        `json:"type"`
	ID      string        `json:"id"`
	Status  string        `json:"status"`
	Role    string        `json:"role"`
	Content []ContentItem `json:"content"`
}

// ContentItem represents an item in the 'content' array of an OutputItem.
type ContentItem struct {
	Type        string        `json:"type"`
	Text        string        `json:"text"`
	Annotations []interface{} `json:"annotations"` // []interface{} for empty array
}

// Reasoning represents the 'reasoning' object.
type Reasoning struct {
	Effort  interface{} `json:"effort"`  // interface{} for null
	Summary interface{} `json:"summary"` // interface{} for null
}

// Text represents the 'text' object.
type Text struct {
	Format TextFormat `json:"format"`
}

// TextFormat represents the 'format' object within the 'text' object.
type TextFormat struct {
	Type string `json:"type"`
}

// Usage represents the 'usage' object.
type Usage struct {
	InputTokens         int                 `json:"input_tokens"`
	InputTokensDetails  TokensDetails       `json:"input_tokens_details"`
	OutputTokens        int                 `json:"output_tokens"`
	OutputTokensDetails OutputTokensDetails `json:"output_tokens_details"`
	TotalTokens         int                 `json:"total_tokens"`
}

// TokensDetails represents the 'input_tokens_details' object.
type TokensDetails struct {
	CachedTokens int `json:"cached_tokens"`
}

// OutputTokensDetails represents the 'output_tokens_details' object.
type OutputTokensDetails struct {
	ReasoningTokens int `json:"reasoning_tokens"`
}
