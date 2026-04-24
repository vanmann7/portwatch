// Package portwindow implements a sliding-time-window tracker for observed
// open ports. It is useful for detecting ports that appear and disappear
// within a short period, enabling flap detection and short-lived service
// monitoring without persisting state to disk.
package portwindow
