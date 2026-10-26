# Gitolite Forge

This is a WoodPeckerCI Forge Addon for Gitolite.

The main goal is to make as simple and generic as possible (would be good if could work with any git repository). Too achieve that we :

* Implement the minimum required methods
* Keeping separated authentication/authorization from repository access
* Allow someway to provide static (i.e. via file) data for the methods

## Requirements

Mainly the forge must provide the following services:

* Authentication : Thru Generic OAuth2/OpenID provider with as least config possible
* Repository Access: Thru SSH, by using a fixed-provided SSH-key

## Methods to implement

* Name,URL (basic)
* Login (authenticates the user via OAuth2) and provides identity to WoodPecker
* Repo (fetch a repository using SSH, would need to map repository names), may be using static data
* Repos (List available repositories, with permissions) may be using static data
* File (Fetch a file at commit-id from repo, used to fetch pipeine configuration)
* Hook (parse webhook and provide pipeline metadata)

# Implementation Plan

## Configuration needed

We need to have provided:

* OAuth client-id, secret and URL
* SSH key
* URL for repository view (http)
* URL for repository clone (ssh)

## Login and identity

* Will use OpenID autodiscovery to ease configuration and also to keep generic as possible
* Will just authenticate the user, although we get a token, that token is not validated by repos (accessed thru SSH)

# Ideas

* Repository listing could be fetched from gitweb project.list
* Either RClone or bind-mount the forge-addon

# References

* [go-gitolite](https://github.com/VonC/gogitolite)
* [go-package](https://pkg.go.dev/github.com/VonC/gogitolite/gitolite)
* [OIDC Client test](https://github.com/leplusorg/openid-connect-provider-debugger)
* [RClone docker plugin](https://rclone.org/docker/)

We could also mock Repos methods (static)

For SSH need:
- Known hosts
- SSH user and key
- SSH server and port

