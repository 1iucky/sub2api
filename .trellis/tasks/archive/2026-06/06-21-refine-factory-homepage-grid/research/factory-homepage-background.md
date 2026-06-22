# Factory.ai Homepage Background Notes

Checked: 2026-06-21

Source: https://factory.ai/

## Observations

* Factory.ai currently serves a dark-first homepage (`data-theme="dark"`).
* Theme variables in the downloaded CSS include `--dark-base-primary: #020202`, `--light-base-primary: #eee`, and orange accents around `#ef6f2e`.
* Browser inspection shows the visible hero background is mostly the dark body surface (`#020202`) plus extremely subtle overlays; it should not be interpreted as an orange color wash.
* The header uses a blurred, masked top surface instead of a hard solid bar.
* The visual read is quiet and industrial: low-contrast surfaces, tiny gray-white points, sparse linework, and a black/white-gray metallic texture. It is not a bright or uniformly dense graph-paper pattern.

## Application To This Repo

* Keep the implementation pure CSS to avoid copying or depending on Factory.ai assets.
* Replace the homepage's dense 18px uniform grid with a lower-contrast, larger dot-and-line texture layered over the base surface.
* Avoid orange in the background treatment; keep orange only where existing UI accents already use the product primary color.
* Preserve the existing warm palette tokens already present in `tailwind.config.js`.
* Keep the effect confined to the first screen and faded before the marquee.
