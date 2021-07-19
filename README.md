# github-sync

This tool is heavily inspired on [tempelis](https://github.com/kubernetes-sigs/slack-infra/tree/master/tempelis)

It syncronizes the configuration described in a YAML file against your GitHub Organization.
Combined with a CI system, it can be used to implement GitOps for GitHub.

At this stage, it can:

### Org Members
  - Add users to an organization
  - Remove users from an organization

### Collaborators
  - Add collaborators to a repository
  - Remove collaborators from a repository


## Config

### Authentication

It expects a config file in the location given by `--auth` that looks like this:

```
{
    "authToken": "ic3hu6ydebbsib1yd7x5wn1nro",
    "org": "YOUR_GITHUB_ORG"
}
```

`authToken` is the [GitHub Token](https://docs.github.com/en/github/authenticating-to-github/keeping-your-account-and-data-secure/creating-a-personal-access-token)


#### Users

It expects a complete list of Org Members to be provided. If a org member exists on
GitHub that is not in the yaml users list, it will error out.


```yaml
users:
- username: cpanato
  role: admin # possibles are `admin` / `direct_member` / `billing_manager`
  email: dont@honk.at.me
- username: honk_user
  role: admin  # possibles are `admin` / `direct_member` / `billing_manager`
  email: test4@example.com
```

#### Collaborators

It expects a complete list of Repositories and Collaborators to be provided. If a collaborator exists on
GitHub that is not in the yaml repositories list, it will error out.

```yaml
repositories:
- name: playground
  collaborators:
  - username: username_1
    email: test@example.com
    permission: push # possibles are `pull` / `push` / `admin` / `maintain` / `triage`
  - username: username_2
    email: test2@example.com
    permission: push # possibles are `pull` / `push` / `admin` / `maintain` / `triage`
- name: another_repo
  collaborators:
  - username: username_3
    email: test3@example.com
    permission: push # possibles are `pull` / `push` / `admin` / `maintain` / `triage`

```


## Future Work

Add support:

- Repositories configuration
- Teams
