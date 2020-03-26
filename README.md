# Repo Manager

I'd like to have an easy way to sync which repos I have on my local machine, so I can easily track:

- which remotes I have configured
- which branches I have locally
- whether I have unstaged/staged/unpushed code
- whether I have worktrees configured

Ideally this would also allow me to enable/disable repos on different computers and support $GOPATH, but that's complicated.

Inspired by tiny-care-terminals config options of:

- `TTC_REPOS` - a comma separated list of repos to look at for git commits.
- `TTC_REPOS_DEPTH` - the max directory-depth to look for git repositories in the directories defined with TTC_REPOS (by default 1). Note that the deeper the directory depth, the slower the results will be fetched. seeing your commits in tiny-terminal-care, set this to gitlog

I know there are a few similar projects, that help manage multiple repos / submodules:

- TODO
