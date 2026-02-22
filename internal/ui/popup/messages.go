package popup

import (
	"crypto/rand"
	"math/big"
)

// defaultMessages is a list of motivational lines shown in the initial popup.
// Leave it empty to use the configured popup message.
var defaultMessages = []string{
	"Stand up. Your spine is not a banana.",
	"Unclench your jaw. Seriously.",
	"Stretch for a sec. I’ll wait.",
	"Posture check. Gravity is winning.",
	"Remember that joke about the shrimp fried rice, and its actually a shrimp frying the rice the kitchen. yea. You look like that shrimp.",
	"Neck stretch time. Pretend you’re a Brontosaurus reaching a tall leaf.",
	"Your back called. It’s tired of being hunched.",
	"Sit like a human, not a shrimp.",
	"Posture check: Are you a question mark?",
	"Move your body. Any movement.",
	"Stretch like a confused cat.",
	"Helo, Your bones want attention.",
	"You have been sitting in one shape for too long.",
	"Blink. Again. One more time.",
	"Look at something 20 feet away for 20 seconds.",
	"Your eyes called. They want a break.",
	"Blink. Do not make this weird.",
	"Your eyes are dry. Emotionally and physically.",
	"Stare at something that is not a screen.",
	"The wall is free. Go look at it.",
	"Your eyeballs are begging. Politely.",
	"Hydration check, You are not a cactus. Probably.",
	"When did you last drink water? Be honest.",
	"Go pee. This is your permission.",
	"Slow breath. Drop your shoulders.",
	"You’re doing fine. Take a breather.",
	"git commit -m 'took a break'",
	"This break is mandatory. The Overlords said so.",
	"Task failed successfully. Take a break.",
	"Imagine a frog. Imagine it jumping away, go chase it.",
	"Stretch like a cat that owns the place.",
	"Go stare out a window dramatically.",
	// "It’s late. Consider stopping. Gently.",
	"You have been productive. That is suspicious.",
	"Take a break. Even CPUs throttle.",
	"Step away before you start arguing with your code.",
	"Pause. Think of nothing. Failed? Good. Take a Breather.",
	"Your brain needs a snack. Or a break. Or both.",
	"You are not a machine. Probably. Hopefully...",
	"Life is short. Stretch.",
	"Take a break before you become a chair.",
	"Congratulations. You have unlocked: Standing.",
	"Get up. I am not asking.",
	"Stand. Stretch. Obey.",
	"You have ignored this long enough.",
	"Last warning. Your posture is watching.",
	"Break. Now.",
	"MOVE",
	"BLINK",
	"STRETCH",
	"PLEASEE",
	"Pretend you're the giraffe's ancestor and the tree is trynna grow taller than you",
	"Stand up before your spine files a complaint.",
	"You have been still, for too long...",
	"This is your sign. Yes, this one.",
	"Take a break. The vibes are off. I can feel them.",
	"Pause. Rehydrate. Continue existing.",
	"NPC behavior detected.",
	"CAPTCHA: Confirm you are not a chair.",
	"Select all images where you are still sitting.",
	"Prove you are human. Stand up.",
	"Verification failed. Please stretch.",
	"CAPTCHA required: Blink twice.",
	"Human check: Move your neck.",
	"Error: Human activity not detected.",
	"Please rotate your shoulders to continue.",
	"Select all squares containing good posture.",
	"System suspects you are furniture.",
	"You are doing side quests anyway. Take a break.",
	"Stand up challenge (100% win rate).",
	"You are doing great. Now stop for a second.",
	"Friendly reminder to not become a statue.",
	"Break time. Argue with the wall.",
	"Move before you merge with the chair.",
	"Your posture just lost subscriber count.",
	"Time to wiggle.",
	"Stand up and exist violently.",
	"Do one weird movement. I dare you.",
	"Become uncompressed.",
	"This is your last gentle reminder.",
	"We could have been normal. Stand up.",
	"You have disappointed the spine council.",
	"Achievement unlocked: Chair Prisoner.",
}

var exitAttemptMessages = []string{
	"Nice Try...",
	"Attempt Denied...",
	"No. (yeah the one with Bugs Bunny)",
	"Breaks are non-negotiable.",
	"sudo: permission denied (you are not rooted, go stretch).",
}

func pickRandomMessage(fallback string) string {
	if len(defaultMessages) == 0 {
		return fallback
	}
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(defaultMessages))))
	if err != nil {
		return defaultMessages[0]
	}
	return defaultMessages[int(n.Int64())]
}

func pickRandomExitAttemptMessage() string {
	if len(exitAttemptMessages) == 0 {
		return "nice try"
	}
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(exitAttemptMessages))))
	if err != nil {
		return exitAttemptMessages[0]
	}
	return exitAttemptMessages[int(n.Int64())]
}
