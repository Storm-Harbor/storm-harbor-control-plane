# CNCF readiness

Assessment date: **2026-11-12**. Storm Harbor is **not a CNCF project** and has not submitted a Sandbox application. This page prepares a possible future application; it does not claim affiliation, acceptance or a guaranteed date.

## Route and current gaps

The normal project entry route is [CNCF Sandbox](https://contribute.cncf.io/projects/submit-project/). Project acceptance is distinct from a company's [CNCF membership](https://www.cncf.io/about/join/).

The [current Sandbox process](https://github.com/cncf/sandbox) and [application form](https://github.com/cncf/sandbox/blob/main/.github/ISSUE_TEMPLATE/application.yml) require Apache 2.0, a directly linked maintainer file with names/GitHub IDs/organizations, an actively developed repository at least six months old, and a reusable project rather than a company-specific platform or reference architecture. Recheck the live criteria before applying.

The technical repository was created on **2026-08-18**. It reaches six calendar months on **2027-02-18**; age alone does not establish readiness. The current documentation/mock scaffold needs a reusable working implementation and evidence of continued development.

| Area | Evidence or remaining work |
| --- | --- |
| License | publish Apache 2.0 consistently in the participating repositories |
| Maintainers | publish names, GitHub identities and affiliations; broaden maintainership over time |
| Reusable implementation | deliver a Go MVP independent of any adopter’s identity, data and cloud accounts |
| Cloud-native role | explain declarative recovery, controllers and interoperability without claiming to replace Kubernetes |
| Governance | public contribution process, conduct policy, security reporting and decision ownership |
| Documentation | English technical documentation, installation, runnable examples and truthful support matrix |
| Adoption | documented consenting adopters and reproducible non-company-specific scenarios |
| Differentiation | compare recovery scope with Velero, Crossplane, Cluster API, Flux and related projects using current primary sources |
| Supply chain | reproducible releases, dependency/license inventory, signed artifacts and vulnerability handling |
| Legal contribution | owner review of the contribution agreement, trademarks and any required parent-project separation vote |

The TOC must approve an application and the required contribution agreement must be signed before the project represents itself as donated or contributed to CNCF. Submission, legal commitments and trademark changes require a separate explicit owner decision.

## Proposed sequence

Publish the project foundations and website, deliver the Go resource/controller and adapter contracts, prove an isolated recovery journey, invite external contributors/adopters, then reassess eligibility and prepare the application with linked evidence. No CNCF badges or logos should imply an existing affiliation.
