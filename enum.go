package githubql

// ProjectState represents state of the project; either 'open' or 'closed'.
type ProjectState string

// State of the project; either 'open' or 'closed'.
const (
	Open   ProjectState = "OPEN"   // The project is open.
	Closed ProjectState = "CLOSED" // The project is closed.
)

// ProjectOrderField represents properties by which project connections can be ordered.
type ProjectOrderField string

// Properties by which project connections can be ordered.
const (
	CreatedAt ProjectOrderField = "CREATED_AT" // Order projects by creation time.
	UpdatedAt ProjectOrderField = "UPDATED_AT" // Order projects by update time.
	Name      ProjectOrderField = "NAME"       // Order projects by name.
)

// OrderDirection represents possible directions in which to order a list of items when provided an `orderBy` argument.
type OrderDirection string

// Possible directions in which to order a list of items when provided an `orderBy` argument.
const (
	Asc  OrderDirection = "ASC"  // Specifies an ascending order for a given `orderBy` argument.
	Desc OrderDirection = "DESC" // Specifies a descending order for a given `orderBy` argument.
)

// ProjectCardState represents various content states of a ProjectCard.
type ProjectCardState string

// Various content states of a ProjectCard.
const (
	ContentOnly ProjectCardState = "CONTENT_ONLY" // The card has content only.
	NoteOnly    ProjectCardState = "NOTE_ONLY"    // The card has a note only.
	Redacted    ProjectCardState = "REDACTED"     // The card is redacted.
)

// CommentCannotUpdateReason represents the possible errors that will prevent a user from updating a comment.
type CommentCannotUpdateReason string

// The possible errors that will prevent a user from updating a comment.
const (
	InsufficientAccess    CommentCannotUpdateReason = "INSUFFICIENT_ACCESS"     // You must be the author or have write access to this repository to update this comment.
	Locked                CommentCannotUpdateReason = "LOCKED"                  // Unable to create comment because issue is locked.
	LoginRequired         CommentCannotUpdateReason = "LOGIN_REQUIRED"          // You must be logged in to update this comment.
	Maintenance           CommentCannotUpdateReason = "MAINTENANCE"             // Repository is under maintenance.
	VerifiedEmailRequired CommentCannotUpdateReason = "VERIFIED_EMAIL_REQUIRED" // At least one email address must be verified to update this comment.
)

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

// ReactionOrderField represents a list of fields that reactions can be ordered by.
type ReactionOrderField string

// A list of fields that reactions can be ordered by.
const (
	CreatedAt ReactionOrderField = "CREATED_AT" // Allows ordering a list of reactions by when they were created.
)

// GitSignatureState represents the state of a Git signature.
type GitSignatureState string

// The state of a Git signature.
const (
	Valid                GitSignatureState = "VALID"                 // Valid signature and verified by GitHub.
	Invalid              GitSignatureState = "INVALID"               // Invalid signature.
	MalformedSig         GitSignatureState = "MALFORMED_SIG"         // Malformed signature.
	UnknownKey           GitSignatureState = "UNKNOWN_KEY"           // Key used for signing not known to GitHub.
	BadEmail             GitSignatureState = "BAD_EMAIL"             // Invalid email used for signing.
	UnverifiedEmail      GitSignatureState = "UNVERIFIED_EMAIL"      // Email used for signing unverified on GitHub.
	NoUser               GitSignatureState = "NO_USER"               // Email used for signing not known to GitHub.
	UnknownSigType       GitSignatureState = "UNKNOWN_SIG_TYPE"      // Unknown signature type.
	Unsigned             GitSignatureState = "UNSIGNED"              // Unsigned.
	GpgverifyUnavailable GitSignatureState = "GPGVERIFY_UNAVAILABLE" // Internal error - the GPG verification service is unavailable at the moment.
	GpgverifyError       GitSignatureState = "GPGVERIFY_ERROR"       // Internal error - the GPG verification service misbehaved.
	NotSigningKey        GitSignatureState = "NOT_SIGNING_KEY"       // The usage flags for the key that signed this don't allow signing.
	ExpiredKey           GitSignatureState = "EXPIRED_KEY"           // Signing key expired.
)

// StatusState represents the possible commit status states.
type StatusState string

// The possible commit status states.
const (
	Expected StatusState = "EXPECTED" // Status is expected.
	Error    StatusState = "ERROR"    // Status is errored.
	Failure  StatusState = "FAILURE"  // Status is failing.
	Pending  StatusState = "PENDING"  // Status is pending.
	Success  StatusState = "SUCCESS"  // Status is successful.
)

// IssueState represents the possible states of an issue.
type IssueState string

// The possible states of an issue.
const (
	Open   IssueState = "OPEN"   // An issue that is still open.
	Closed IssueState = "CLOSED" // An issue that has been closed.
)

// IssueOrderField represents properties by which issue connections can be ordered.
type IssueOrderField string

// Properties by which issue connections can be ordered.
const (
	CreatedAt IssueOrderField = "CREATED_AT" // Order issues by creation time.
	UpdatedAt IssueOrderField = "UPDATED_AT" // Order issues by update time.
	Comments  IssueOrderField = "COMMENTS"   // Order issues by comment count.
)

// SubscriptionState represents the possible states of a subscription.
type SubscriptionState string

// The possible states of a subscription.
const (
	Unsubscribed SubscriptionState = "UNSUBSCRIBED" // The User is only notified when particpating or @mentioned.
	Subscribed   SubscriptionState = "SUBSCRIBED"   // The User is notified of all conversations.
	Ignored      SubscriptionState = "IGNORED"      // The User is never notified.
)

// RepositoryPrivacy represents the privacy of a repository.
type RepositoryPrivacy string

// The privacy of a repository.
const (
	Public  RepositoryPrivacy = "PUBLIC"  // Public.
	Private RepositoryPrivacy = "PRIVATE" // Private.
)

// RepositoryOrderField represents properties by which repository connections can be ordered.
type RepositoryOrderField string

// Properties by which repository connections can be ordered.
const (
	CreatedAt RepositoryOrderField = "CREATED_AT" // Order repositories by creation time.
	UpdatedAt RepositoryOrderField = "UPDATED_AT" // Order repositories by update time.
	PushedAt  RepositoryOrderField = "PUSHED_AT"  // Order repositories by push time.
	Name      RepositoryOrderField = "NAME"       // Order repositories by name.
)

// RepositoryAffiliation represents the affiliation of a user to a repository.
type RepositoryAffiliation string

// The affiliation of a user to a repository.
const (
	Owner              RepositoryAffiliation = "OWNER"               // Repositories that are owned by the authenticated user.
	Collaborator       RepositoryAffiliation = "COLLABORATOR"        // Repositories that the user has been added to as a collaborator.
	OrganizationMember RepositoryAffiliation = "ORGANIZATION_MEMBER" // Repositories that the user has access to through being a member of an organization. This includes every repository on every team that the user is on.
)

// PullRequestState represents the possible states of a pull request.
type PullRequestState string

// The possible states of a pull request.
const (
	Open   PullRequestState = "OPEN"   // A pull request that is still open.
	Closed PullRequestState = "CLOSED" // A pull request that has been closed without being merged.
	Merged PullRequestState = "MERGED" // A pull request that has been closed by being merged.
)

// MergeableState represents whether or not a PullRequest can be merged.
type MergeableState string

// Whether or not a PullRequest can be merged.
const (
	Mergeable   MergeableState = "MERGEABLE"   // The pull request can be merged.
	Conflicting MergeableState = "CONFLICTING" // The pull request cannot be merged due to merge conflicts.
	Unknown     MergeableState = "UNKNOWN"     // The mergeability of the pull request is still being calculated.
)

// IssuePubSubTopic represents the possible PubSub channels for an issue.
type IssuePubSubTopic string

// The possible PubSub channels for an issue.
const (
	Updated    IssuePubSubTopic = "UPDATED"    // The channel ID for observing issue updates.
	Markasread IssuePubSubTopic = "MARKASREAD" // The channel ID for marking an issue as read.
)

// PullRequestReviewState represents the possible states of a pull request review.
type PullRequestReviewState string

// The possible states of a pull request review.
const (
	Pending          PullRequestReviewState = "PENDING"           // A review that has not yet been submitted.
	Commented        PullRequestReviewState = "COMMENTED"         // An informational review.
	Approved         PullRequestReviewState = "APPROVED"          // A review allowing the pull request to merge.
	ChangesRequested PullRequestReviewState = "CHANGES_REQUESTED" // A review blocking the pull request from merging.
	Dismissed        PullRequestReviewState = "DISMISSED"         // A review that has been dismissed.
)

// PullRequestPubSubTopic represents the possible PubSub channels for a pull request.
type PullRequestPubSubTopic string

// The possible PubSub channels for a pull request.
const (
	Updated    PullRequestPubSubTopic = "UPDATED"    // The channel ID for observing pull request updates.
	Markasread PullRequestPubSubTopic = "MARKASREAD" // The channel ID for marking an pull request as read.
)

// DeploymentStatusState represents the possible states for a deployment status.
type DeploymentStatusState string

// The possible states for a deployment status.
const (
	Pending  DeploymentStatusState = "PENDING"  // The deployment is pending.
	Success  DeploymentStatusState = "SUCCESS"  // The deployment was successful.
	Failure  DeploymentStatusState = "FAILURE"  // The deployment has failed.
	Inactive DeploymentStatusState = "INACTIVE" // The deployment is inactive.
	Error    DeploymentStatusState = "ERROR"    // The deployment experienced an error.
)

// DeploymentState represents the possible states in which a deployment can be.
type DeploymentState string

// The possible states in which a deployment can be.
const (
	Abandoned DeploymentState = "ABANDONED" // The pending deployment was not updated after 30 minutes.
	Active    DeploymentState = "ACTIVE"    // The deployment is currently active.
	Destroyed DeploymentState = "DESTROYED" // An inactive transient deployment.
	Error     DeploymentState = "ERROR"     // The deployment experienced an error.
	Failure   DeploymentState = "FAILURE"   // The deployment has failed.
	Inactive  DeploymentState = "INACTIVE"  // The deployment is inactive.
	Pending   DeploymentState = "PENDING"   // The deployment is pending.
)

// OrganizationInvitationRole represents the possible organization invitation roles.
type OrganizationInvitationRole string

// The possible organization invitation roles.
const (
	DirectMember   OrganizationInvitationRole = "DIRECT_MEMBER"   // The user is invited to be a direct member of the organization.
	Admin          OrganizationInvitationRole = "ADMIN"           // The user is invited to be an admin of the organization.
	BillingManager OrganizationInvitationRole = "BILLING_MANAGER" // The user is invited to be a billing manager of the organization.
	Reinstate      OrganizationInvitationRole = "REINSTATE"       // The user's previous role will be reinstated.
)

// DefaultRepositoryPermissionField represents the possible default permissions for organization-owned repositories.
type DefaultRepositoryPermissionField string

// The possible default permissions for organization-owned repositories.
const (
	Read  DefaultRepositoryPermissionField = "READ"  // Members have read access to org repos by default.
	Write DefaultRepositoryPermissionField = "WRITE" // Members have read and write access to org repos by default.
	Admin DefaultRepositoryPermissionField = "ADMIN" // Members have read, write, and admin access to org repos by default.
)

// TeamPrivacy represents the possible team privacy values.
type TeamPrivacy string

// The possible team privacy values.
const (
	Secret  TeamPrivacy = "SECRET"  // A secret team can only be seen by its members.
	Visible TeamPrivacy = "VISIBLE" // A visible team can be seen and @mentioned by every member of the organization.
)

// UserOrderField represents properties by which user connections can be ordered.
type UserOrderField string

// Properties by which user connections can be ordered.
const (
	Login  UserOrderField = "LOGIN"  // Allows ordering a list of users by their login.
	Action UserOrderField = "ACTION" // Allows ordering a list of users by their ability action.
)

// TeamOrderField represents properties by which team connections can be ordered.
type TeamOrderField string

// Properties by which team connections can be ordered.
const (
	Name TeamOrderField = "NAME" // Allows ordering a list of teams by name.
)

// TeamRole represents the role of a user on a team.
type TeamRole string

// The role of a user on a team.
const (
	Admin  TeamRole = "ADMIN"  // User has admin rights on the team.
	Member TeamRole = "MEMBER" // User is a member of the team.
)

// StarOrderField represents properties by which star connections can be ordered.
type StarOrderField string

// Properties by which star connections can be ordered.
const (
	StarredAt StarOrderField = "STARRED_AT" // Allows ordering a list of stars by when they were created.
)

// GistPrivacy represents the privacy of a Gist.
type GistPrivacy string

// The privacy of a Gist.
const (
	Public GistPrivacy = "PUBLIC" // Public.
	Secret GistPrivacy = "SECRET" // Secret.
	All    GistPrivacy = "ALL"    // Gists that are public and secret.
)

// MilestoneState represents the possible states of a milestone.
type MilestoneState string

// The possible states of a milestone.
const (
	Open   MilestoneState = "OPEN"   // A milestone that is still open.
	Closed MilestoneState = "CLOSED" // A milestone that has been closed.
)

// RepositoryLockReason represents the possible reasons a given repsitory could be in a locked state.
type RepositoryLockReason string

// The possible reasons a given repsitory could be in a locked state.
const (
	Moving    RepositoryLockReason = "MOVING"    // The repository is locked due to a move.
	Billing   RepositoryLockReason = "BILLING"   // The repository is locked due to a billing related reason.
	Rename    RepositoryLockReason = "RENAME"    // The repository is locked due to a rename.
	Migrating RepositoryLockReason = "MIGRATING" // The repository is locked due to a migration.
)

// RepositoryCollaboratorAffiliation represents the affiliation type between collaborator and repository.
type RepositoryCollaboratorAffiliation string

// The affiliation type between collaborator and repository.
const (
	All     RepositoryCollaboratorAffiliation = "ALL"     // All collaborators of the repository.
	Outside RepositoryCollaboratorAffiliation = "OUTSIDE" // All outside collaborators of an organization-owned repository.
)

// LanguageOrderField represents properties by which language connections can be ordered.
type LanguageOrderField string

// Properties by which language connections can be ordered.
const (
	Size LanguageOrderField = "SIZE" // Order languages by the size of all files containing the language.
)

// SearchType represents represents the individual results of a search.
type SearchType string

// Represents the individual results of a search.
const (
	Issue      SearchType = "ISSUE"      // Returns results matching issues in repositories.
	Repository SearchType = "REPOSITORY" // Returns results matching repositories.
	User       SearchType = "USER"       // Returns results matching users on GitHub.
)

// PullRequestReviewEvent represents the possible events to perform on a pull request review.
type PullRequestReviewEvent string

// The possible events to perform on a pull request review.
const (
	Comment        PullRequestReviewEvent = "COMMENT"         // Submit general feedback without explicit approval.
	Approve        PullRequestReviewEvent = "APPROVE"         // Submit feedback and approve merging these changes.
	RequestChanges PullRequestReviewEvent = "REQUEST_CHANGES" // Submit feedback that must be addressed before merging.
	Dismiss        PullRequestReviewEvent = "DISMISS"         // Dismiss review so it now longer effects merging.
)

// TopicSuggestionDeclineReason represents reason that the suggested topic is declined.
type TopicSuggestionDeclineReason string

// Reason that the suggested topic is declined.
const (
	NotRelevant        TopicSuggestionDeclineReason = "NOT_RELEVANT"        // The suggested topic is not relevant to the repository.
	TooSpecific        TopicSuggestionDeclineReason = "TOO_SPECIFIC"        // The suggested topic is too specific for the repository (e.g. #ruby-on-rails-version-4-2-1).
	PersonalPreference TopicSuggestionDeclineReason = "PERSONAL_PREFERENCE" // The viewer does not like the suggested topic.
	TooGeneral         TopicSuggestionDeclineReason = "TOO_GENERAL"         // The suggested topic is too general for the repository.
)
