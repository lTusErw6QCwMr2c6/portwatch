// Package eventlabel provides a rule-based labeling engine that attaches
// key/value string labels to port events. Rules are registered by name
// and evaluated against each event; matching rules contribute their
// label to the result map returned by Apply.
package eventlabel
