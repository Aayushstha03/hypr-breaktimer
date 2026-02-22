package logic

import (
	"crypto/rand"
	"math/big"
)

var nudgeMessages = []string{
	"Posture check. You are not a folded Ethernet cable.",
	"Straighten your back. Gravity is not your enemy.",
	"This is a reminder to stop cosplaying a question mark.",
	"Relax your shoulders. You are not compiling stress.",
	"Jaw status: clenched. Please release before system damage.",
	"You are scrolling again. No, this does not count as research.",
	"Reminder: staring at text is not the same as processing it.",
	"This tab switch achieved nothing. Congratulations.",
	"You have been ‘about to start’ for a concerning amount of time.",
	"Stop multitasking. You are a single-core processor with limited RAM.",
	"Parallel tasks detected. Performance degradation inevitable.",
	"Context switching is not a productivity feature.",
	"Pick one task. This is not a load balancer.",
	"Your brain does not support true concurrency.",
	"No YouTube.",
	"You do not need a video in the background to think.",
	"That video will still exist after the task. Shocking, I know.",
	"Watching others be productive is not productivity.",
	"Close the tab.",
	"This is your sign to recalibrate.",
	"Autopilot disengaged. Please take control.",
	"You are drifting. This is a correction ping.",
}

var nudgeMessengers = []string{
	"Posture Police",
	"Break Bot 3000",
	"The Voice",
	"Your Conscience",
	"Ergonomics Enforcer",
	"Reality Check Service",
	"Attention Scheduler",
	"Cognitive Garbage Collector",
	"Context Switch Police",
	"SingleCore Supervisor",
	"Userland Nudge Engine",
	"The Nudge Gremlin",
	"Unpaid Supervisor",
	"Thought Auditor",
	"Brain.exe",
	"Me, But Meaner",
	"Internal Monologue",
}

func RandomNudgeMessage() (string, string) {
	mIdx := 0
	if n, err := rand.Int(rand.Reader, big.NewInt(int64(len(nudgeMessengers)))); err == nil {
		mIdx = int(n.Int64())
	}

	msgIdx := 0
	if n, err := rand.Int(rand.Reader, big.NewInt(int64(len(nudgeMessages)))); err == nil {
		msgIdx = int(n.Int64())
	}

	return nudgeMessengers[mIdx], nudgeMessages[msgIdx]
}
