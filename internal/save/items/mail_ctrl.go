package items

import (
	"fmt"
	"strconv"
	"strings"
)

// MailCtrl mirrors the C# Mail_Ctrl class from CustomClasses/Save/Items/Mail_Ctrl.cs.
type MailCtrl struct {
	Inbox        []string
	LastID       int
	UnreadCount  int
	PendingMails int
	PmailTimers  int
}

// FromProperties populates the MailCtrl from a map of SII properties.
func (m *MailCtrl) FromProperties(props map[string][]string) error {
	for key, vals := range props {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]

		switch {
		case key == "inbox":
			// capacity hint; ignored in Go, slices grow dynamically
		case strings.HasPrefix(key, "inbox["):
			m.Inbox = append(m.Inbox, val)
		case key == "last_id":
			m.LastID = parseInt(val)
		case key == "unread_count":
			m.UnreadCount = parseInt(val)
		case key == "pending_mails":
			m.PendingMails = parseInt(val)
		case key == "pmail_timers":
			m.PmailTimers = parseInt(val)
		}
	}
	return nil
}

// ToProperties converts the MailCtrl struct to a map of properties.
func (m *MailCtrl) ToProperties() map[string][]string {
	props := make(map[string][]string)

	props["inbox"] = []string{strconv.Itoa(len(m.Inbox))}
	for i, v := range m.Inbox {
		props[fmt.Sprintf("inbox[%d]", i)] = []string{v}
	}

	props["last_id"] = []string{strconv.Itoa(m.LastID)}
	props["unread_count"] = []string{strconv.Itoa(m.UnreadCount)}
	props["pending_mails"] = []string{strconv.Itoa(m.PendingMails)}
	props["pmail_timers"] = []string{strconv.Itoa(m.PmailTimers)}

	return props
}

