## 1.0.0 (Unreleased)

BACKWARDS INCOMPATIBILITIES / NOTES:

* provider: This version supports Arukas' GA version. If you were using Arukas' beta version, you will need to register your account again. For details, please refer to [Arukas updates: Relaunch of the Arukas beta service](https://arukas.io/en/updates-en/201802_arukas_beta_relaunch-en/) [GH-2].
* provider: Use go v1.10 for testing and building [GH-2].
* provider: Use `hashicorp/terraform/helper/validation` package for validation [GH-2].

* resource/arukas_container: `memory` attribute was removed. Please use `plan` attribute instead [GH-2].
* resource/arukas_container: `app_id` attribute was removed [GH-2].

ENHANCEMENTS:

* resource/arukas_container: Added `plan` attribute [GH-2].
* resource/arukas_container: Added `service_id` attribute [GH-2].
* resource/arukas_container: Added `region` attribute [GH-2].

## 0.1.0 (June 20, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
