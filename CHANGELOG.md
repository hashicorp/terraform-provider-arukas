## 1.2.0 (Unreleased)
## 1.1.0 (April 03, 2019)

IMPROVEMENTS:

* Update to be terraform 0.12 compliant [[#10](https://github.com/terraform-providers/terraform-provider-arukas/issues/10)]  
The provider is now compatible with Terraform v0.12, while retaining compatibility with prior versions.

BUG FIXES:

* Upgrade go-arukas to v0.1.0 ([#9](https://github.com/terraform-providers/terraform-provider-arukas/issues/9))
* resource/arukas_container: fixed wrong value setting of endpoint_full_xxx attributes ([#4](https://github.com/terraform-providers/terraform-provider-arukas/issues/4))

## 1.0.0 (September 04, 2018)

BACKWARDS INCOMPATIBILITIES / NOTES:

* provider: This version supports Arukas' GA version. If you were using Arukas' beta version, you will need to register your account again. For details, please refer to [Arukas updates: Relaunch of the Arukas beta service](https://arukas.io/en/updates-en/201802_arukas_beta_relaunch-en/) [[#2](https://github.com/terraform-providers/terraform-provider-arukas/issues/2)].
* provider: Use go v1.10 for testing and building [[#2](https://github.com/terraform-providers/terraform-provider-arukas/issues/2)].
* provider: Use `hashicorp/terraform/helper/validation` package for validation [[#2](https://github.com/terraform-providers/terraform-provider-arukas/issues/2)].

* resource/arukas_container: `memory` attribute was removed. Please use `plan` attribute instead [[#2](https://github.com/terraform-providers/terraform-provider-arukas/issues/2)].
* resource/arukas_container: `app_id` attribute was removed [[#2](https://github.com/terraform-providers/terraform-provider-arukas/issues/2)].

ENHANCEMENTS:

* resource/arukas_container: Added `plan` attribute [[#2](https://github.com/terraform-providers/terraform-provider-arukas/issues/2)].
* resource/arukas_container: Added `service_id` attribute [[#2](https://github.com/terraform-providers/terraform-provider-arukas/issues/2)].
* resource/arukas_container: Added `region` attribute [[#2](https://github.com/terraform-providers/terraform-provider-arukas/issues/2)].

## 0.1.0 (June 20, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
