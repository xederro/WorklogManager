package worklogText

import (
	"context"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/xederro/WorklogManager/internal/config"
	"google.golang.org/genai"
)

var KeyMap = huh.KeyMap{
	Text: huh.TextKeyMap{
		Prev:    key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "back")),
		Next:    key.NewBinding(key.WithKeys("tab", "enter"), key.WithHelp("enter", "next")),
		Submit:  key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "submit")),
		NewLine: key.NewBinding(key.WithKeys("alt+enter"), key.WithHelp("alt+enter", "new line")),
		Editor:  key.NewBinding(key.WithKeys("ctrl+e"), key.WithHelp("ctrl+e", "open editor")),
	},
}

type WorklogText struct {
	*huh.Text
	Ai  key.Binding
	Esc key.Binding
	Key string
}

type WorklogAIMsg struct {
	WorklogText string
	Key         string
	Err         error
}

type WorklogExitMsg struct{}

func NewWorklogText(issueKey string) *WorklogText {
	t := huh.NewText()
	return &WorklogText{
		Text: t,
		Ai: key.NewBinding(
			key.WithKeys("ctrl+a"),
			key.WithHelp("ctrl+a", "generate AI response"),
		),
		Esc: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "exit"),
		),
		Key: issueKey,
	}
}

func (t *WorklogText) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, t.Ai):
			if !config.Conf.UseAi {
				break
			}

			// blocking call to AI service
			// because we want to wait for the response
			// not block the UI
			// resulted in race condition
			// it is to be checked if something can be done about it
			// for example:
			// receive the response in a channel
			// doesn't update the UI
			// unless some other thing is done
			// this can result in overriding the text
			// by user input
			t.getAiMessage()
		case key.Matches(msg, t.Esc):
			config.Ch <- ExitResponse()
		default:
			var cmd tea.Cmd
			_, cmd = t.Text.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return t, tea.Batch(cmds...)
}

func (t *WorklogText) KeyBinds() []key.Binding {
	keys := t.Text.KeyBinds()
	if config.Conf.UseAi {
		keys = append(keys, t.Ai)
	}
	keys = append(keys, t.Esc)
	return keys
}
func ExitResponse() func() tea.Msg {
	return func() tea.Msg {
		return WorklogExitMsg{}
	}
}

func AiResponse(wt string, e error, key string) func() tea.Msg {
	return func() tea.Msg {
		return WorklogAIMsg{
			WorklogText: wt,
			Err:         e,
			Key:         key,
		}
	}
}

func (t *WorklogText) getAiMessage() {
	result, err := config.GoogleClient.Models.GenerateContent(
		context.Background(),
		config.Conf.GoogleAi.DefaultModel,
		genai.Text(config.Conf.GoogleAi.DefaultPrompt+" "+t.GetValue().(string)),
		&genai.GenerateContentConfig{
			ThinkingConfig: &genai.ThinkingConfig{
				ThinkingBudget:  nil,
				IncludeThoughts: false,
			},
		},
	)
	if err != nil {
		config.Ch <- AiResponse("", err, t.Key)
	} else {
		config.Ch <- AiResponse(result.Text(), nil, t.Key)
	}
}

func (t *WorklogText) Title(title string) *WorklogText {
	t.Text = t.Text.Title(title)
	return t
}

func (t *WorklogText) Value(value *string) *WorklogText {
	t.Text = t.Text.Value(value)
	return t
}
func (t *WorklogText) WithHeight(height int) huh.Field {
	t.Text = t.Text.WithHeight(height).(*huh.Text)
	return t
}
