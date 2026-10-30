# Gitolite Forge

This (will be) a WoodPeckerCI Forge Addon for Gitolite.

The main goal is to make as simple and generic as possible (would be good if could work with any git repository). Too achieve that we :

* Implement the minimum required methods
* Try to keeping differentiation between authentication/identity/authorization from repository access
* Allow someway to provide static (i.e. via file) data for the methods

# Implementation Plan

## Configuration needed

We need to have provided:

* OAuth client-id, secret and URL
* SSH key
* URL for repository view (http)
* URL for repository clone (ssh)

# References

Below some relevant WoodPecker discussions/issues relevant, please note that some are quite old and not applicable on newer versions, but provide a good point in some hints:

* [Gerrit Forge](https://github.com/woodpecker-ci/woodpecker/issues/4232)
* [Support for generic git servers](https://github.com/woodpecker-ci/woodpecker/discussions/3591)
* [Wish setup](https://github.com/woodpecker-ci/woodpecker/discussions/2651)
* [Explain forge](https://github.com/woodpecker-ci/woodpecker/issues/4444) (although abandoned)
* [How to setup addon](https://github.com/woodpecker-ci/woodpecker/discussions/3266#discussioncomment-8258887)
* [Radicle Addon](https://github.com/cytechmobile/radicle-woodpecker-addon)
* [CodeUp Addon](https://github.com/tsrpc/yunxiao-woodpecker-addon)
