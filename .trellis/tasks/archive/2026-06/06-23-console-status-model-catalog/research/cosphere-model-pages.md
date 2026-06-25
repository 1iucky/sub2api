# CoSphere Model Pages Research

## Source Files Inspected

* `/Users/liuliang/work/openS/CoSphere/web/src/pages/Model/index.jsx`
* `/Users/liuliang/work/openS/CoSphere/web/src/components/table/models/index.jsx`
* `/Users/liuliang/work/openS/CoSphere/web/src/components/table/models/ModelsColumnDefs.jsx`
* `/Users/liuliang/work/openS/CoSphere/web/src/components/table/models/ModelsTabs.jsx`
* `/Users/liuliang/work/openS/CoSphere/web/src/components/table/models/ModelsFilters.jsx`
* `/Users/liuliang/work/openS/CoSphere/web/src/components/table/models/ModelsActions.jsx`
* `/Users/liuliang/work/openS/CoSphere/web/src/pages/Setting/Model/SettingGlobalModel.jsx`
* `/Users/liuliang/work/openS/CoSphere/web/src/pages/Setting/Ratio/ModelRatioSettings.jsx`

## CoSphere Patterns

* Model management is a dedicated page routed at `/console/models`.
* The page is table-centric: tabs at the top, actions and filters above the table, pagination below.
* Vendor/category is a primary organizational axis. Vendor tabs include counts and vendor management actions.
* Model rows include icon, model name, matching rule, sync participation, description, vendor, tags, endpoints, bound channels, enabled groups, quota types, created time, updated time, and operations.
* Filters include model-name search and vendor search.
* Actions include add model, missing model discovery, sync, prefill group management, compact mode, batch delete, copy names, and add selected models to prefill group.
* Model marketplace visibility can be controlled from operation settings, including whether unauthenticated users can access the marketplace.
* Ratio/pricing settings remain separate from model metadata settings. This separation is useful for this project because channel pricing already owns billing values.

## Mapping to sub2api

* sub2api currently has channel pricing entries with `platform` and model name arrays. Those names can be used as the association key into a model catalog.
* sub2api should keep billing behavior in channel pricing and use model catalog data for metadata, navigation, marketplace display, and admin organization.
* Vendor tabs map well to providers/platforms such as OpenAI, Anthropic, Gemini, Antigravity, and future platforms.
* Missing-model discovery can be useful by scanning channel pricing entries with no associated catalog record, but it can be deferred if the initial task needs a smaller MVP.
* Bound channels can be computed by finding channel pricing entries that include or match a catalog model.

## Recommended MVP Shape

* Persistent model catalog records:
  * name/model ID
  * display name
  * platform/provider
  * vendor
  * category/type
  * tags
  * endpoints/capabilities
  * description
  * icon key or URL
  * status/visibility
* Admin model management page:
  * table with search, platform/vendor filters, create/edit/delete or enable/disable
  * per-row related pricing/channel count
* User-facing model marketplace page:
  * browse by platform/vendor/category
  * show pricing associations where available
* Channel pricing association:
  * exact normalized model-name match first
  * unmatched pricing names still display as raw names with an "unmatched catalog" state
  * avoid changing pricing calculation or model routing behavior

## Deferred Enhancements

* Official upstream sync wizard.
* Vendor CRUD if initial vendor can be represented as plain text.
* Prefill groups and bulk import/export.
* Complex match rules beyond exact normalized model name.
