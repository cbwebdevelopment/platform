## HEAD

## v1.9.0 (2017-08-10)

- Bump `hash_deactivate_old` data deduplicator to version 1.1.0
- Update `hash_deactivate_old` data deduplicator to use archived dataset id and time fields to accurately:
  - Deactivate deduplicate data on dataset addition
  - Activate undeduplicated data on dataset deletion
  - Record entire deduplication history
- Update mongo queries related to `hash_deactivate_old` data deduplicator
- Remove backwards-compatible legacy deduplicator name test in `DeduplicatorDescriptor.IsRegisteredWithNamedDeduplicator` (after `v1.8.0` required migration)
- Add archived dataset id and time fields to base data type
- Add MD5 hash of authentication token to request logger
- Add service middleware to extract select request headers and add as request logger fields
- Defer access to context store sessions and log until actually needed

## v1.8.0 (2017-08-09)

- Add CHANGELOG.md
- **REQUIRED MIGRATION**: `migrate_data_deduplicator_descriptor` - data deduplicator descriptor name and version
- Force `precise` Ubuntu distribution for Travis (update to `trusty` later)
- Add deduplicator version
- Update deduplicator name scheme
- Add github.com/blang/semver package dependency
- Fix dependency import capitalization
- Update dependencies
- Remove unused data store functionality
- Remove unused data deduplicators

## v1.7.0 (2017-06-22)

- See commit history for details on this and all previous releases
