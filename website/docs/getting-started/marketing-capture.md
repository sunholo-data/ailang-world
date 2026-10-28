---
title: Release screenshots
sidebar_position: 4
description: The six release screenshots of the AILANG World workbench, what each one shows and does not show, and how the capture is reproduced.
---

import useBaseUrl from '@docusaurus/useBaseUrl';

# Release screenshots

:::note
These images are static captures of recorded history, published by the attended v1.0.0 release capture. Until that capture is promoted, the images are not yet in the site.
:::

The six images below are in the order the release capture takes them. Each is a static capture of recorded history from a demo store, not a live view. The captions say what each image shows and what it does not.

### Light workbench

<img src={useBaseUrl('/img/v1/workbench-light.png')} alt="Light workbench for the selected demo world, with its header and first timeline entries" />

Desktop light theme at 1440x1100: the selected world, its revision, log head and state root, and its first timeline entries. The headless capture shows a vertical scrollbar.

### Dark workbench

<img src={useBaseUrl('/img/v1/workbench-dark.png')} alt="Dark workbench for the same demo world, with its header and first timeline entries" />

The same regions in the dark theme at 1440x1100. The scheme check confirms a dark background; it is not a review of the content.

### Decision packet

<img src={useBaseUrl('/img/v1/packet.png')} alt="Decision packet for a validated demo proposal, with its pending approval chain and four-grade evidence roster" />

The decision packet at 1440x3700, a tall framing that puts the validated sentence and the four-grade evidence roster inside the image. The status is `approvals chain: pending`. The CLAIMED row is demo-authored commentary and is never review.

### World graph

<img src={useBaseUrl('/img/v1/graph.png')} alt="World graph section for the selected demo entry, with its checked graph targets" />

The world graph section at 1440x2100 for the selected entry and its checked graph targets. It shows no evidence graph that the store does not contain. Node labels are truncated where they run past their boxes.

### Provenance

<img src={useBaseUrl('/img/v1/provenance.png')} alt="Stored invocation record opened from the approval entry, with its leading fields legible" />

The provenance view at 1440x3000. It shows the stored invocation record, opened from the approval entry, with its leading fields legible; the complete reference set is the accompanying `why` JSON. At this width the payload runs past the right edge, so later fields are clipped and a horizontal scrollbar appears.

### Phone

<img src={useBaseUrl('/img/v1/phone.png')} alt="Phone-width workbench for the selected demo world, with its header and first timeline entries" />

The phone layout at 390x844, showing the selected world and a readable timeline. The page overflows horizontally and the title is clipped to "AILANG World Workbenc". That is a known cosmetic defect, not a content failure.

## What the images do and do not show

- They are static captures of recorded history from one demo store. They are not a live view, and they do not show anything that happened after the capture.
- Only the CLAIMED row is demo-authored: it is commentary and is never review. The TESTED row comes from a real `ailang test` run in the seed, the ATTESTED row from a recorded effect, and the PROVEN row from a subject-bound proof the daemon's validator accepted. A validated packet is the state of that store at capture time, not an approval.
- The images show the packet awaiting a decision. No approval decision appears in any of them.
- The packet, graph and provenance images are tall (1440x3700, 1440x2100 and 1440x3000) so that the required content is in the image rather than below the fold.
- The provenance claim is narrowed to the stored invocation record with its leading fields legible. The complete reference set is in the accompanying `why` JSON.
- Cosmetic risks are disclosed and not fixed here: vertical scrollbars in the headless shots, a horizontal scrollbar and clipped fields in provenance, phone overflow with a clipped title, and truncated graph labels.
- Two capture pins are involved and they are not interchangeable: the capsule and proof interpreter is v0.41.0, and the workspace tools use the separate v0.54.15 tool pin.
- Byte equality is claimed only between the two passes over one frozen store. Fresh seeds differ, so a new seed produces a new capture rather than a repeat of these images.

## Reproduce them

The reproduction runbook is not duplicated here. Follow it in the repository:

- [Prepare the release demo seed (section 12)](https://github.com/sunholo-data/ailang-world/blob/dev/docs/QUICKSTART.md#12-prepare-the-release-demo-seed-incremental-tooling)
- [Capture the release screenshots (section 13)](https://github.com/sunholo-data/ailang-world/blob/dev/docs/QUICKSTART.md#13-capture-the-release-screenshots-attended)
