# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

### Fixed

## [v1.0.22] - 2023-02-17

### Added

- #1017: Add locked / unlocked filter to datasets
- #983: File upload progress
- Show error dialog if upload is too large
- Make maximum file size configurable
- #930: Reindex command
- #955: Better WOS Type facet (reindex needed)
- #1015: Search on organization ID (reindex needed)
- #956: Show status of facet filter in badge
- #1004: Add Reset filters button
- #957: Show top 3 selected facet values in badge
- Cypress tests

### Fixed

- Fix conflict error after file upload cancel
- #1006: Expose abstract language to frontoffice

## [v1.0.21] - 2023-01-25

### Added

### Fixed

- Fix bug in some repeated text fields

## [v1.0.20] - 2023-01-25

### Added

- Switch from deprecated biblio-backend prefix to biblio-backoffice

### Fixed

## [v1.0.19] - 2023-01-25

### Added

- Switch to more secure github.com/ugent-library/oidc for authentication
- #936: Differentiate between sytem and user updates in ui
- Add request id to log statements; improve logging in general

### Fixed

- #966: Add missing external contributor names to dataset xlsx export
- #932: Add missing has_message curation only facet to datasets
- #921: Add status and vabb_year fields to xslx export
- #991: Be more forgiving when decoding boolean values from MongoDB
- Various css fixes, fix typos

## [v1.0.18] - 2023-01-20

### Added

- Simple batch interface for curators (can currently only add projects) 
- ulid wrapper package is no longer needed and has been removed

### Fixed

## [v1.0.17] - 2023-01-20

### Added

### Fixed

- #984: Fix decoding of projects from Elasticsearch
- #986: Reviewer tags facet was missing a 'Select all' button
- #975: Fix exposing of licenses to frontoffice

## [v1.0.16] - 2023-01-19

### Added

### Fixed

- #984: Fix decoding of projects from MongoDB

## [v1.0.15] - 2023-01-18

### Added

- #848: Show legacy flag to curators and display prettier boolean flags.
- #950: Show the chosen license in the Full text & Files list.
- #926: Search on WoS ID.

### Fixed

- #925: Use frontoffice Elasticsearch and MongoDB directlto relieve pressure on
  frontoffice app.
- #937: Fix field extern display.
- Fix typos.
- #943: Fix timestamp format in frontoffice handler.

## [v1.0.14] - 2022-12-20

### Added

### Fixed

- Previous fix for #910 had a bug where publication version field didn't appear,
  this is now resolved.

## [v1.0.13] - 2022-12-20

### Added

- #928: Allow transferring a single publication between users
- #881: Add a publication transfer command that rewrites history and assigns
  publications to another user id
- #875: Improve error reporting by including an error id.
- #850: Add "Deselect all" action to facet filter dialog.
- Make facets configurable.

### Fixed

- #900: Empty list item spotted in authors and editors – probably for UGent ID.
- #901: Not all departments of people are showing in the overview.
- #887: Fix handle creation for datasets.
- #902: Import language from WoS.
- #910: File document type defaults to full text.
- #924: Order year facet new to old.
- #918: Set most open license as copyright statement in frontoffice.
- Various ux fixes.

## [v1.0.12] - 2022-11-30

### Added

- CHANGELOG.md document.

### Fixed

- #600: Improve search by removing punctuation and icu folding (requires index switch).
- #865: Remove "Publication short title" from dissertation details display.
- #866: Add missing "Journal title" and "Short journal title" labels for issue_editor.
- #863: Only show "Lay summaries" and "Conference details" links in sidebar menu.
  when the type uses these fields.
- Various fixes (#867).

## [v1.0.11] - 2022-11-28

### Added

- Add automated keyword cleanup to the cleanup command.

### Fixed

- #859: When uploading files, the file does not always appear, or becomes a separate block.
- #823: Choose which file is the "primary" on a publication.
- #840: A better error page when a user profile can't be retrieved from the Biblio frontend.
- #855: Add the author anyway even if the authors' department couldn't be retrieved from Biblio Frontend.
- #821: Update the bootstrap theme.

## [v1.0.10] - 2022-11-24

### Added

- Cleanup command. Cleanup of author departments and fix missing faculties on publications.
- Expose the file hash to the Biblio frontend.

### Fixed

- #828: Only show "full text missing" label if publication is not extern.
- #846: Remove references to Publication#URL and PublicationFile#URL.
- #792: Various fixes towards data consistency in the gRPC client.

## [v1.0.9] - 2022-11-22

### Fixed

- #837: Dashboard does not display all unclassified records, even when they have a department.
- #842: ESCI ID missing for V classified publications.
- #845: Fix deleting departments with an id containing an asterisk.

## [v1.0.8] - 2022-11-22

### Fixed

- Skip empty string validation on keywords for now.

## [v1.0.7] - 2022-11-22

### Fixed

- Fix WoS import.
- Add missing Department and CreditRole to Publication#Editor and Publication#Supervisor.

## [v1.0.6] - 2022-11-21

### Fixed

- #824: map to 3-letter language code in datacite and crossref importers.
- Various fixes.

## [v1.0.5] - 2022-11-17

### Added

- #822: Added handle create command.

## [v1.0.4] - 2022-11-16

### Added

- Use structured logger in cmd's used in cron (others are todo).

### Changed

- Clear user when system updates record.

## [v1.0.3] - 2022-11-16

### Changed

- Change temporary message on home page.

## [v1.0.2] - 2022-11-16

### Fixed

- #819: Librarian tags are missing, starting from f.
- #818: New home page message.
- #816: Add missing copy_to statements for publication.

## [v1.0.1] - 2022-11-16

### Fixed

- Various fixes

## [v1.0.0] - 2022-11-14

## Added

- Initial release

[unreleased]:  https://github.com/ugent-library/biblio-backoffice/compare/v1.0.22...HEAD
[v1.0.22]:  https://github.com/ugent-library/biblio-backoffice/compare/v1.0.21...v1.0.22
[v1.0.21]:  https://github.com/ugent-library/biblio-backoffice/compare/v1.0.20...v1.0.21
[v1.0.20]:  https://github.com/ugent-library/biblio-backoffice/compare/v1.0.19...v1.0.20
[v1.0.19]:  https://github.com/ugent-library/biblio-backoffice/compare/v1.0.18...v1.0.19
[v1.0.18]:  https://github.com/ugent-library/biblio-backoffice/compare/v1.0.17...v1.0.18
[v1.0.17]:  https://github.com/ugent-library/biblio-backoffice/compare/v1.0.16...v1.0.17
[v1.0.16]:  https://github.com/ugent-library/biblio-backoffice/compare/v1.0.15...v1.0.16
[v1.0.15]:  https://github.com/ugent-library/biblio-backoffice/compare/v1.0.14...v1.0.15
[v1.0.14]:  https://github.com/ugent-library/biblio-backoffice/compare/v1.0.13...v1.0.14
[v1.0.13]:  https://github.com/ugent-library/biblio-backoffice/compare/v1.0.12...v1.0.13
[v1.0.12]:  https://github.com/ugent-library/biblio-backoffice/compare/v1.0.11...v1.0.12
[v1.0.11]:  https://github.com/ugent-library/biblio-backoffice/compare/v1.0.10...v1.0.11
[v1.0.10]:  https://github.com/ugent-library/biblio-backoffice/compare/v1.0.9...v1.0.10
[v1.0.9]:  https://github.com/ugent-library/biblio-backoffice/compare/v1.0.8...v1.0.9
[v1.0.8]:  https://github.com/ugent-library/biblio-backoffice/compare/v1.0.7...v1.0.8
[v1.0.7]:  https://github.com/ugent-library/biblio-backoffice/compare/v1.0.6...v1.0.7
[v1.0.6]:  https://github.com/ugent-library/biblio-backoffice/compare/v1.0.5...v1.0.6
[v1.0.5]:  https://github.com/ugent-library/biblio-backoffice/compare/v1.0.4...v1.0.5
[v1.0.4]:  https://github.com/ugent-library/biblio-backoffice/compare/v1.0.3...v1.0.4
[v1.0.3]:  https://github.com/ugent-library/biblio-backoffice/compare/v1.0.2...v1.0.3
[v1.0.2]:  https://github.com/ugent-library/biblio-backoffice/compare/v1.0.1...v1.0.2
[v1.0.1]:  https://github.com/ugent-library/biblio-backoffice/compare/v1.0.0...v1.0.1
