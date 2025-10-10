package ssg

type SSGKeys struct {
	WorkspacePath  string
	DocsPath       string
	MarkdownPath   string
	HTMLPath       string
	LayoutPath     string
	HeaderStyle    string
	AssetsPath     string
	ImagesPath     string
	BlocksMaxItems string
	IndexMaxItems  string

	SearchGoogleEnabled string
	SearchGoogleID      string

	PublishRepoURL         string
	PublishBranch          string
	PublishPagesSubdir     string
	PublishAuthMethod      string
	PublishAuthToken       string
	PublishCommitUserName  string
	PublishCommitUserEmail string
	PublishCommitMessage   string
}

var SSGKey = SSGKeys{
	WorkspacePath:  "ssg.workspace.path",
	DocsPath:       "ssg.docs.path",
	MarkdownPath:   "ssg.markdown.path",
	HTMLPath:       "ssg.html.path",
	LayoutPath:     "ssg.layout.path",
	HeaderStyle:    "ssg.header.style",
	AssetsPath:     "ssg.assets.path",
	ImagesPath:     "ssg.images.path",
	BlocksMaxItems: "ssg.blocks.maxitems",
	IndexMaxItems:  "ssg.index.maxitems",

	SearchGoogleEnabled: "ssg.search.google.enabled",
	SearchGoogleID:      "ssg.search.google.id",

	PublishRepoURL:         "ssg.publish.repo.url",
	PublishBranch:          "ssg.publish.branch",
	PublishPagesSubdir:     "ssg.publish.pages.subdir",
	PublishAuthMethod:      "ssg.publish.auth.method",
	PublishAuthToken:       "ssg.publish.auth.token",
	PublishCommitUserName:  "ssg.publish.commit.user.name",
	PublishCommitUserEmail: "ssg.publish.commit.user.email",
	PublishCommitMessage:   "ssg.publish.commit.message",
}
