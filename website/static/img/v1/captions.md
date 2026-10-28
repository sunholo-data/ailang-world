# Release screenshot captions

Each section below is one promoted image, in the order the release capture takes them.
The alt text is the short description a screen reader reads. The caption states what the
image shows and what it does not. Every image is a static capture of recorded history from
a demo store; none is a live view.

## workbench-light.png

Alt: Light workbench for the selected demo world, with its header and first timeline entries

Desktop light theme at 1440x1100. The header shows the selected world, its revision, log head and state root, and the timeline shows its first entries: the part of the page the one-column layout places in the first viewport. This is a static capture of recorded history. The headless capture shows a vertical scrollbar.

## workbench-dark.png

Alt: Dark workbench for the same demo world, with its header and first timeline entries

The same regions as the light workbench, measured with the dark theme tokens, at 1440x1100. The scheme check confirms a dark background; it is not a review of the content. This is a static capture of recorded history. The headless capture shows a vertical scrollbar.

## packet.png

Alt: Decision packet for a validated demo proposal, with its pending approval chain and four-grade evidence roster

The decision packet at 1440x3700, a tall framing that puts the validated sentence and the proof-envelope roster inside the image. It shows the proposal and request links, the status `approvals chain: pending`, and the four-grade evidence roster: proven, tested, attested and claimed. The CLAIMED row is demo-authored commentary and is never review. Nothing in the image is an approval decision. This is a static capture of recorded history. The headless capture shows a vertical scrollbar.

## graph.png

Alt: World graph section for the selected demo entry, with its checked graph targets

The world graph section at 1440x2100, a tall framing that puts the section and its checked graph targets inside the image. It shows the selected entry and its graph targets from the demo store, and no evidence graph that the store does not contain. Node labels are truncated where the SVG text runs past its node box. That is a known cosmetic issue, not a content failure. This is a static capture of recorded history. The headless capture shows a vertical scrollbar.

## provenance.png

Alt: Stored invocation record opened from the approval entry, with its leading fields legible

The provenance view at 1440x3000, a tall framing. It shows the stored invocation record, opened from the approval entry, with its leading fields legible; the complete reference set is the accompanying `why` JSON. At 1440 px the payload runs past the right edge, so the record has a horizontal scrollbar and its later fields are clipped. This is a static capture of recorded history. The headless capture also shows a vertical scrollbar.

## phone.png

Alt: Phone-width workbench for the selected demo world, with its header and first timeline entries

The phone layout at 390x844: the selected world and a readable timeline. At this width the page overflows horizontally, and the title is clipped to "AILANG World Workbenc". That is a disclosed cosmetic defect in the renderer, not a content failure. This is a static capture of recorded history. The headless capture shows a vertical scrollbar.

## Disclosures

- Every image is a static capture of recorded history from a demo store. It is not a live view and not a review of the store's content.
- Only the CLAIMED row is demo-authored commentary, and it is never review. The TESTED row comes from a real `ailang test` run in the seed, the ATTESTED row from a recorded effect, and the PROVEN row from a subject-bound proof the daemon's validator accepted. A validated packet is the state of one demo store at capture time.
- Tall framing: packet.png is 1440x3700, graph.png is 1440x2100 and provenance.png is 1440x3000. The other three are 1440x1100 and 390x844.
- The provenance claim is narrowed to: the stored invocation record, opened from the approval entry, with its leading fields legible; the complete reference set is the accompanying `why` JSON.
- Separate pins: the capsule and proof interpreter is v0.41.0, and the workspace tools use the separate v0.54.15 tool pin. The two are not interchangeable.
- Cosmetic risks, all disclosed: vertical scrollbars in the headless shots; a horizontal scrollbar and clipped later fields in provenance; phone overflow with the title clipped; truncated graph labels.
- Byte equality is claimed only between the two passes over one frozen store. Fresh seeds differ, so a new seed is a new capture, not a repeat of these images.
