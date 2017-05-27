package githubql

// TODO: Generate the rest from schema according to this formula.

// ReactionContent represents emojis that can be attached to Issues, Pull Requests and Comments.
type ReactionContent string

// Emojis that can be attached to Issues, Pull Requests and Comments.
const (
	ThumbsUp   ReactionContent = "THUMBS_UP"   // Represents the 👍 emoji.
	ThumbsDown ReactionContent = "THUMBS_DOWN" // Represents the 👎 emoji.
	Laugh      ReactionContent = "LAUGH"       // Represents the 😄 emoji.
	Hooray     ReactionContent = "HOORAY"      // Represents the 🎉 emoji.
	Confused   ReactionContent = "CONFUSED"    // Represents the 😕 emoji.
	Heart      ReactionContent = "HEART"       // Represents the ❤️ emoji.
)
