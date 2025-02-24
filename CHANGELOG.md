# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v1.7.0] - 2025-02-21

- New sync to synchronize a golden (expected) JSON with a new JSON string while preserving pattern matching expressions from the golden JSON.

## [v1.6.0] - 2025-02-15

- Errors are now public
- Reporting problems is transparent

## [v1.5.1] - 2025-02-11

- report path in correct order
- change reporting of the path start '.' for root. Reusable by yq.

## [v1.5.0] - 2025-02-09

- refactor tests to test against real errors. Use t.Run to see each subtest as separate subtest.
- library reports errors as a list of errors.

## [v1.4.0] - 2025-02-09

- create empty matcher, which matches undefined, null, missing field.

## [v1.3.1] - 2025-02-08

- BUG Fix the date matcher creation

## [v1.3.0] - 2025-02-07

### Added

- Date pattern: `@date@`

## [1.1.0] - 2019-07-07

### Added

- Email pattern: `@email@`

## [1.0.0] - 2019-01-27

### Added

- Initial release with support for patterns:
  - `@string@`
  - `@number@`
  - `@bool@`
  - `@array@`
  - `@uuid@`
  - `@wildcard@`
  - `@...@`
