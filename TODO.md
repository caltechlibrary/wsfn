
# Action Items

- [ ] Drop the 32-bit Raspberry Pi build (`Linux-armv7l`, `GOARM=7`) from the Makefile and any release check the next time this module is worked on (workspace DR-0018, proposed 2026-10-26). 64-bit Raspberry Pi OS stays covered by `Linux-aarch64`. Already-published releases are left as they are.

+ [ ] `.mjs` files need to be served as "text/javascript" per https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Modules

## Questions

+ Should LoadWebService overwrite .Access always be overwritten with .AccessFile's contents?
+ Should DumpWebService save .Access if .AccessFile is set?

