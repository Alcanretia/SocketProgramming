# Pull request descriptions

- `.github/pull_request_template.md` is the default template GitHub fills in for any new PR.
- The files here are ready-to-paste, pre-filled versions for the five branches we planned. Fill in the `<!-- ... -->` parts in your own words.

| File | Branch | Into | Reviewer |
|---|---|---|---|
| `01-feature-app-p4-initial.md` | `feature/app-p4-initial` | `main` | everyone |
| `02-review-p1-connection-core.md` | `review/p1-connection-core` | `feature/app-p4-initial` (or `main` once merged) | P4 |
| `03-review-p2-users-dm.md` | `review/p2-users-dm` | same | P1 |
| `04-review-p3-groups-qr.md` | `review/p3-groups-qr` | same | P2 |
| `05-learn-p4-explain-back.md` | `learn/p4-explain-back` | same | P3 |

A review branch has no diff against `feature/app-p4-initial` until its owner commits something, so open that PR only after pushing your own commits.

## Opening a PR with the GitHub CLI (PowerShell)

```powershell
gh pr create --base main --head feature/app-p4-initial `
  --title "Chat app: server, client, tests, docs" `
  --body-file docs/pr/01-feature-app-p4-initial.md

gh pr create --base feature/app-p4-initial --head review/p1-connection-core `
  --title "P1: review and learning for the connection core" `
  --body-file docs/pr/02-review-p1-connection-core.md
```

Without the CLI: open the branch on GitHub, press "Compare & pull request", and paste the file's text into the description box.
