// Package eventclassifier assigns classification labels (normal, suspicious,
// critical) to alert events based on configurable port-range and protocol rules.
// Rules are evaluated in insertion order; the first matching rule determines
// the class. Events that match no rule receive the ClassUnknown label.
package eventclassifier
