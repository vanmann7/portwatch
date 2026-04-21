// Package labeler assigns human-readable labels to network ports.
//
// Labels are resolved in priority order:
//  1. Custom rules supplied at construction time (first match wins).
//  2. Built-in well-known service names (ssh, http, https, …).
//  3. A generic "port/<n>" fallback for unrecognised ports.
//
// Rules can be parsed from config strings with ParseRules, which accepts
// entries of the form "<port>:<label>" or "<low>-<high>:<label>".
package labeler
