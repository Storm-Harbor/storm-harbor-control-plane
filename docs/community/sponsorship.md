# Sponsorship and sustainability

Storm Harbor welcomes support for its open-source recovery platform. The project owner manages receiving-account verification and onboarding privately. An existing Stripe account is intended for payment processing.

## Current status

Payment collection is **not activated** by this documentation. No working Storm Harbor payment link or approved GitHub Sponsors organization profile has been verified yet. Until a live destination is configured, sponsorship links lead here rather than to an unconfigured checkout.

For sponsorship interest, open a [sponsorship inquiry](https://github.com/Storm-Harbor/.github/issues/new?title=Sponsorship%20inquiry). Do not include payment, banking, tax, infrastructure or private project information in a public issue. A confidential contact channel must be agreed before exchanging commercial details.

## Funding priorities

- Go runtime and CLI development;
- provider adapter contracts and isolated multicloud recovery tests;
- resource/controller extension compatibility and security;
- public documentation, releases and community maintenance.

Support does not grant access to other organizations' repositories, knowledge, infrastructure or recovery evidence. Sponsorship does not bypass security policy, imply a production SLA or purchase recovery authorization.

## Activation checklist

1. Verify the existing Stripe account belongs to the intended receiving entity and is enabled for live payments.
2. Agree sponsorship amounts, currency, one-time/recurring terms and any public recognition with the owner.
3. Create and verify hosted Stripe Payment Links in that account; keep payment data on Stripe.
4. Publish only verified live URLs on the website and in the organization-level `FUNDING.yml`.
5. Optionally complete [GitHub Sponsors organization onboarding](https://docs.github.com/en/sponsors/receiving-sponsorships-through-github-sponsors/setting-up-github-sponsors-for-your-organization), including payout and tax information directly in GitHub/its payment provider and GitHub approval.
6. Use a signed, idempotent Go webhook receiver if future product entitlements depend on payment events. A browser success redirect is not proof of payment. Simple voluntary sponsorship links do not require a product backend.

Potential later revenue includes support, managed operations and implementation services, with separately agreed terms and a clear distinction from voluntary open-source sponsorship. Those services and pricing are not currently offered by this page.
