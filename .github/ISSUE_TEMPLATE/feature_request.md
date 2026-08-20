---
name: Feature / change request
about: Propose new functionality or a change to existing behavior
title: "[Feature] <short description>"
labels: enhancement
---

<!--
KEEP IT SHORT. Describe the problem before the solution.
If you used an AI to draft this: cut the filler. We want the idea, not an essay.
Replace the _italic example text_ in each section with your own.
-->

## Problem
<!-- What can't be done today, or what's painful? One short paragraph. -->
_e.g. There's no way to replay a failed batch without restarting the whole service._

## Proposed change
<!-- What you'd like to see. Bullets are fine. -->
_e.g. Add a CLI command to requeue transactions in a failed state by batch ID._

## Scope
<!-- Tick what's touched. Delete the rest. -->
- Service: `private-relayer` / `public-relayer` / `cts`
- Component: <!-- e.g. listener, executor, repository, config, contractclient -->

## Alternatives considered
<!-- Optional. One line each. Write "none" if you didn't. -->
_e.g. Manual DB update — too error-prone._

## Notes
<!-- Anything else worth knowing: dependencies, risks, related issues. Keep it brief. -->
_e.g. Related to #123._
