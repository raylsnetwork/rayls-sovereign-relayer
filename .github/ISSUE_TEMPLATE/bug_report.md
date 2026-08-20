---
name: Bug report
about: Something is broken or behaving unexpectedly
title: "[Bug] <short description>"
labels: bug
---

<!--
KEEP IT SHORT. A good bug report fits on one screen.
Facts over prose. If a section doesn't apply, write "n/a" — don't pad it.
If you used an AI to draft this: trim it. No restating the same thing three ways.
Replace the _italic example text_ in each section with your own.
-->

## What's wrong
<!-- One or two sentences. What did you expect, what happened instead. -->
_e.g. Expected the executor to retry the batch after an RPC timeout, but it dropped the message silently._

## Where
<!-- Tick what's affected. Delete the rest. -->
- Service: `private-relayer` / `public-relayer` / `cts`
- Component: <!-- e.g. listener, executor, batcher, keyqueue, enygma/DvP, migrations, other -->
- Environment: <!-- local dev / remote dev / qa  -->
- Version / commit: <!-- git SHA or tag -->

## Steps to reproduce
<!-- Numbered, minimal, in order. Someone else should be able to follow these. -->
1. _Start the private-relayer against ..._
2. _Send a message with ..._
3. _Observe ..._

## Logs / errors
<!-- Paste the RELEVANT lines only, not the whole log. Trim stack traces to the meaningful frames. -->
```
paste log lines here
```

## Impact
<!-- Who/what is blocked? Is there a workaround? One line. -->
_e.g. Blocks DvP settlement on qa; no known workaround._
