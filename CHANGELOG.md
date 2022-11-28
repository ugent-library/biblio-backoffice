# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- CHANGELOG.md document.

### Fixed

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
- #842: ESCI-ID missing for V classified publications.
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

[unreleased]:  https://github.com/ugent-library/biblio-backoffice/compare/v1.0.11...HEAD
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