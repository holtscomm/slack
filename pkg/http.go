package slack

import (
	"net/http"
)

type Payload struct {
	Token          string
	TeamID         string
	TeamDomain     string
	EnterpriseID   string
	EnterpriseName string
	ChannelID      string
	ChannelName    string
	UserID         string
	UserName       string
	Command        string
	Text           string
	ResponseURL    string
	TriggerID      string
}

// PayloadFromRequest returns a struct with all of the variables Slack might send with a request to your application.
// N.B. This will only work for Slash Commands and other commands that send a url encoded post body. This won't work
// with JSON bodies.
func PayloadFromRequest(r *http.Request) (*Payload, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}
	return &Payload{
		Token:          r.FormValue("token"),
		TeamID:         r.FormValue("team_id"),
		TeamDomain:     r.FormValue("team_domain"),
		EnterpriseID:   r.FormValue("enterprise_id"),
		EnterpriseName: r.FormValue("enterprise_name"),
		ChannelID:      r.FormValue("channel_id"),
		ChannelName:    r.FormValue("channel_name"),
		UserID:         r.FormValue("user_id"),
		UserName:       r.FormValue("user_name"),
		Command:        r.FormValue("command"),
		Text:           r.FormValue("text"),
		ResponseURL:    r.FormValue("response_url"),
		TriggerID:      r.FormValue("trigger_id"),
	}, nil
}
