# Tennessee / Davidson County (Nashville) residential lease requirements

**Purpose:** Reference notes for a future Brickwise **lease builder** feature focused on residential rentals in **Davidson County / Metro Nashville, Tennessee**.  
**Status:** Research snapshot for product design — **do not treat as legal advice**. Laws change; have counsel review before shipping lease templates or auto-generated disclosures.

**Scope of this note:** What is distinctive about Tennessee state law (especially URLTA) and Metro Nashville / Davidson County practice that a lease builder should account for.

---

## 1. Governing framework — URLTA applies in Davidson County

Tennessee’s **Uniform Residential Landlord and Tenant Act (URLTA)** is **Tenn. Code Ann. Title 66, Chapter 28**. It does **not** apply statewide.

- URLTA applies only in counties with population **> 75,000** under the **2010** federal census (Tenn. Code Ann. § 66-28-102(a)).
- **Davidson County is a URLTA county.** Commentary commonly lists: Anderson, Blount, Bradley, Davidson, Hamilton, Knox, Madison, Maury, Montgomery, Rutherford, Sevier, Shelby, Sullivan, Sumner, Washington, Williamson, and Wilson.
- In covered counties, URLTA **occupies and preempts** the local field of landlord–tenant regulation; local governments generally may not enact rules that **conflict with or add to** URLTA (§ 66-28-102(e)). Metro Nashville still has specific ordinances (see § 9 below) that a lease builder should surface — verify current enforceability with counsel given the preemption language.

**Excluded arrangements** (not governed by URLTA unless structured to evade it) include institutional residences, occupancy under a purchase contract, transient hotel/motel occupancy, condo/co-op ownership occupancy, and primarily agricultural rentals (§ 66-28-102(c)). HUD / public housing units may be partly displaced by federal regulation (§ 66-28-102(d)).

**Lease-builder implication:** Default templates for Nashville properties should assume **URLTA**, not generic “Tennessee non-URLTA” rules (§ 66-7-109), which matter mainly outside large counties.

**Sources:**  
- https://law.justia.com/codes/tennessee/title-66/chapter-28/part-1/section-66-28-102/  
- https://assets-002.noviams.com/novi-file-uploads/caatn/URLTA_2025.pdf  
- https://www.help4tn.org/senior-services/legal-information/long-term-care-housing/urlta-renters  

---

## 2. Required / expected written disclosures (lease packet)

### State (URLTA) — at or before tenancy begins

| Disclosure | Requirement | Cite |
|---|---|---|
| **Manager / owner for notices & service** | In writing: (1) name & address of agent authorized to manage the premises; (2) name & address of an owner or agent authorized to accept service of process and notices/demands; (3) a maintenance phone **or** email **or** online portal | § 66-28-302 |
| **Keep current** | Info must stay current; successors are bound | § 66-28-302(b) |
| **Personal property insurance** | Landlord/agent must advise in writing that the landlord is **not** responsible for and **will not provide** fire/casualty insurance for the tenant’s personal property | § 66-28-201(a) |
| **Security deposit account location** | If a deposit is taken, notify tenant at lease signing of the **location** of the separate deposit account (account **number** not required) | § 66-28-301(h) |

### Federal — housing built before 1978

- **Lead-based paint disclosure** (federal, not a separate Tennessee lease statute): disclose known lead-based paint / hazards, provide EPA pamphlet *Protect Your Family From Lead in Your Home*, and include the signed lead warning statement before the tenant becomes obligated.
- Regulations: **24 CFR Part 35** (HUD) and **40 CFR Part 745** (EPA).

**Lease-builder implication:** Conditionally attach lead packet when `year_built < 1978` (or unknown → prompt user). Always include insurance disclaimer + owner/agent/maintenance contacts + deposit bank location when a deposit is charged.

**Sources:**  
- https://law.justia.com/codes/tennessee/title-66/chapter-28/part-3/section-66-28-302/  
- https://law.justia.com/codes/tennessee/title-66/chapter-28/part-2/section-66-28-201/  
- https://law.justia.com/codes/tennessee/title-66/chapter-28/part-3/section-66-28-301/  
- https://www.epa.gov/lead/real-estate-disclosure  

---

## 3. Security deposits

- **No statutory dollar cap** on residential security deposits under URLTA (§ 66-28-301 sets process rules, not a multiple-of-rent limit). Market practice is often 1–2 months’ rent (secondary guides).
- **Separate account:** Deposits must be held in an account used **only** for security deposits, at a bank/institution regulated by Tennessee or a U.S. agency (§ 66-28-301(a)).
- **No statutory interest** requirement found for residential deposits.
- **Move-out inspection / listing of damages:** Tenant has rights to inspect and to a listing of ascertainable damage and estimated repair cost; detailed timing rules (landlord notice, inspection within ~4 days of vacating when scheduled, etc.) are in § 66-28-301(b).
- **Forfeiture of deductions:** Landlord who fails to deposit in a required account **and** fails to provide the required damage listing **may not retain any portion** of the deposit (§ 66-28-301(c)).
- **Later-discovered damage window:** Physical damage costs generally only recoverable if discovered before the earlier of **30 days** after tenant vacates/abandons **or** **7 days** after a new tenant takes possession (§ 66-28-301(g)).
- Secondary guides commonly state refund/itemization practice around **30 days** after tenancy ends and unclaimed balance rules around **60 days** after refund notice — confirm against the full text of § 66-28-301 when implementing workflows.

**Lease-builder implication:** Collect “deposit amount,” “bank/institution name & city,” inspection/checklist clauses aligned with § 66-28-301, and warn that noncompliance can forfeit the right to keep the deposit.

**Source:** https://law.justia.com/codes/tennessee/title-66/chapter-28/part-3/section-66-28-301/

---

## 4. Rent, late fees, and lease terms

- Parties may agree on rent, term, and other terms **not prohibited** by URLTA or other law (§ 66-28-201(a)).
- **Tenant may not waive** URLTA rights/remedies in the rental agreement (§ 66-28-201(a); also reinforced in prohibited-provisions rules).
- **Late fees:**
  - **5-day grace period** counted from (and including) the rent due date before any late fee may be charged.
  - If day 5 is a **Sunday or legal holiday**, no late fee if rent is paid the **next business day**.
  - Cap: **≤ 10% of the amount of rent past due** (§ 66-28-201(d)).
- Unless otherwise agreed, periodic rent for terms of one month or less is payable at the beginning of the term; longer terms in equal monthly installments at the beginning of each month (§ 66-28-201(c)).
- **Notice of nonpayment** may be waived in the written lease **only** if the waiver is in **12-point bold** (or larger); even then, the 5-day grace period still applies (§ 66-28-505(b); § 66-28-201).

**Lease-builder implication:** Late-fee UI should enforce grace ≥ 5 days and fee ≤ 10% of past-due rent; offer an optional bold-type nonpayment-notice waiver clause with clear labeling.

**Source:** https://law.justia.com/codes/tennessee/title-66/chapter-28/part-2/section-66-28-201/

---

## 5. Notice periods — ending / changing tenancies

### Periodic tenancies (URLTA)

| Tenancy | Minimum written notice | Cite |
|---|---|---|
| Week-to-week | **10 days** before termination date in the notice | § 66-28-512(a) |
| Month-to-month | **30 days** before the periodic rental date specified in the notice | § 66-28-512(b) |

### Fixed-term leases

- Fixed term ends on the agreed end date unless renewed; holdover remedies in § 66-28-512(c).
- Early termination for breach: see eviction overview (§ 7).

### Metro Nashville / Davidson County ordinances (Chapter 11.22)

Metro Code **Ch. 11.22 – Residential Landlord and Tenant Regulations** (confirm against Municode / counsel; note URLTA preemption tension in § 1):

| Topic | Local rule (as published) | Notes |
|---|---|---|
| **Rent increase notice** | If the written lease does **not** already require **≥ 60 days** written notice of a rent increase, landlord must give **≥ 90 days** advance written notice | Applies to residential leases entered into on/after **July 1, 2020**; increase cannot take effect before original lease term ends; excludes certain government-regulated housing |
| **Nonrenewal notice** | If the written lease does **not** already require **≥ 60 days** written notice of nonrenewal, landlord must give **≥ 60 days** written notice before lease end of intent not to renew | Applies to leases with term of **one year or greater** |
| **Lease assignment notice** | Assignee generally must notify tenant within **14 days** with contact and payment address (unless lease already addresses assignment notice) | Leases on/after **April 1, 2021** |

**Lease-builder implication:** For Davidson County 12-month+ leases, default nonrenewal notice to **≥ 60 days** (or longer if user prefers), and rent-increase notice to **≥ 90 days** unless the template itself embeds a ≥ 60-day increase clause (still verify desired product default with counsel).

**Sources:**  
- https://law.justia.com/codes/tennessee/title-66/chapter-28/part-5/section-66-28-512/  
- https://library.municode.com/tn/metro_government_of_nashville_and_davidson_county/codes/code_of_ordinances?nodeId=CD_TIT11PUPEMOWE_DIVIGERE_CH11.22RELATERE  

---

## 6. Habitability / landlord maintenance duties

Landlord shall (§ 66-28-304):

1. Comply with building/housing codes **materially affecting health and safety**;
2. Make repairs and keep premises in a **fit and habitable** condition;
3. Keep **common areas** clean and safe;
4. Maintain facilities/appliances supplied under the rental agreement in good/safe working order (see full statute for electrical, plumbing, sanitary, HVAC, etc.);
5. In multi-unit buildings of **four or more units**, provide/maintain appropriate trash receptacles for common collection points (per statutory wording).

Tenant duties (care of unit, waste, fixtures, etc.) are in § 66-28-401 et seq.

**Lease-builder implication:** Include habitability acknowledgment and repair-request procedure pointing to the disclosed maintenance phone/email/portal; avoid clauses that waive landlord health/safety duties.

**Source:** URLTA PDF / Nolo overview — https://www.nolo.com/landlord-tenant/landlord-tenant-laws-tennessee.html  

---

## 7. Eviction procedures overview (not a how-to for unlawful eviction)

Self-help eviction (changing locks, shutting off essentials) is risky/unlawful under URLTA remedies for unlawful exclusion / diminution of services (§ 66-28-504 area). Typical path: **notice → court detainer → writ**.

### Tenant noncompliance / nonpayment (URLTA § 66-28-505)

- **Curable by paying money** (rent, repairs, damages, amounts due): written notice that agreement terminates if not remedied within **14 days** after receipt.
- **Not remediable by payment:** termination date **≥ 14 days** after receipt of notice.
- **Repeat** substantially same breach within **6 months:** termination on **≥ 7 days’** written notice.
- Optional **lease waiver** of the § 66-28-505 notice (12-point bold) allows filing a detainer warrant after the grace period without that notice — still does **not** shorten the 5-day late-fee grace (§ 66-28-505(b)).

### Violence / serious danger

Separate shorter-notice tools exist under broader Title 66 provisions commonly discussed alongside URLTA (e.g., aspects of § 66-7-109 for violent acts / present danger — confirm which statute applies to a given case). Secondary guides sometimes cite **3-day** notices for violent/drug-related danger scenarios; **verify before encoding**.

### After judgment

Possession is obtained through the court process (detainer warrant / writ of possession). Abandoned property storage rules appear in URLTA (e.g., storage periods around **30 days** before disposal/sale in abandonment contexts — see full chapter text).

**Lease-builder implication:** Templates can describe notice types and optional bold waiver language; do **not** automate court filings or encourage lockouts.

**Sources:**  
- https://law.justia.com/codes/tennessee/title-66/chapter-28/part-5/section-66-28-505/  
- https://law.justia.com/codes/tennessee/title-66/chapter-7/section-66-7-109/  

---

## 8. Landlord entry

- Landlord may enter to inspect, make repairs, supply services, or exhibit as provided by statute/agreement (§ 66-28-403).
- **Showing to prospective tenants** in the **final 30 days** of the tenancy: allowed if the **rental agreement grants** that right **and** tenant gets **≥ 24 hours’** notice before entry (§ 66-28-403).
- Unlawful/harassing entry gives tenant remedies (injunction, termination, damages, attorney’s fees).

**Lease-builder implication:** Include an express showing-access clause if the landlord wants end-of-term showings; default entry-notice language even where statute is silent for ordinary repairs (contractual clarity).

---

## 9. Davidson County / Metro Nashville–specific items

### Landlord registration (state statute scoped to Metro + large population)

- **Tenn. Code Ann. § 66-28-107** requires landlords of one or more dwelling units in a county with a **metropolitan government** and population **> 500,000** (2000 or later census) to register with the local codes agency — this is **Davidson County**.
- Register: landlord/agent **name, phone, physical address (no P.O. box)** and each unit’s street address/unit number.
- Fee authorized **≤ $10/year**; Metro collects **$10** covering **all** of an owner’s units.
- Update ownership/info within **30 days** of change.
- Failure: up to **$50/week** per unit after hearing opportunity (15 days’ mailed notice).

Metro Codes (Property Standards) administers registration via **ePermits** or mail. Metro materials state registration is required for **single-family, duplex, triplex, quadplex** (and similar); **apartment complexes are not required to register** per Property Standards Division manual — confirm current practice with Codes.

**Official pages:**  
- https://www.nashville.gov/departments/codes/property-standards/landlord-registration  
- https://law.justia.com/codes/tennessee/title-66/chapter-28/part-1/section-66-28-107/  
- Metro Ch. 16.33 rental / registration context: Municode rental inspections chapter  

### Rent increase / nonrenewal / assignment notices

See § 5 (Metro Ch. 11.22).

### Codes & property standards

Habitability enforcement often runs through **Metro Codes / Property Standards** and health departments — lease builder can link to HUB Nashville / Codes complaint channels as informational resources, not as legal process.

---

## 10. Rent control status

- **No rent control** in Nashville / Davidson County for private residential rentals.
- Tennessee **prohibits** local governments from enacting/maintaining/enforcing rent control on private residential or commercial property (**Tenn. Code Ann. § 66-35-102**), including certain below-market mandate / permit-contingency schemes (voluntary incentive programs remain a separate policy path).

**Lease-builder implication:** Do not add rent-cap compliance UI for Davidson County; still honor **contractual** and **Metro notice** rules for increases.

**Source:** Secondary summary citing § 66-35-102 — https://skylinecapitalinvestments.com/market-insights/nashville-tn-landlord-tenant-laws/ (confirm against current TCA text when implementing).

---

## 11. Prohibited / high-risk lease clauses (URLTA)

Among other limits, rental agreements generally **may not**:

- Have the tenant **waive or forgo** URLTA rights/remedies (§ 66-28-201(a)).
- Include **confession of judgment** or landlord **exculpation / indemnity** for landlord liability arising under law (§ 66-28-203 — see full section).
- Willfully use provisions known to be prohibited — can trigger **actual damages** (§ 66-28-203(b) area).

Attorney’s-fee clauses favoring the landlord for breach/nonpayment are commonly treated as allowed when provided in the agreement (§ 66-28-505(d) context) — still draft carefully.

**Lease-builder implication:** Maintain a “forbidden clause” blocklist and warn when users paste custom language that waives statutory rights.

---

## 12. Other practical items for a future lease builder

| Topic | Note for product |
|---|---|
| **Written agreement** | Strongly preferred; URLTA contemplates written waivers and disclosures. Oral arrangements create proof problems. |
| **Occupants / guests / pets / smoking** | Mostly contractual; keep consistent with fair housing (federal + Tennessee Human Rights Act — separate compliance track). |
| **Utilities** | Spell out who pays which utilities; interruption of essential services can create landlord liability if unlawful. |
| **Renter’s insurance** | Landlord must disclose they don’t insure tenant belongings; leases often **require** tenant to carry renters insurance (generally treated as allowed — confirm). |
| **Fire/casualty** | Substantial impairment → termination pathways and prepaid rent/deposit handling under § 66-28-503 — consider optional clause + workflow later. |
| **Fair housing** | Federal FHA + state law; templates must avoid discriminatory preferences. |
| **Short-term rentals** | Separate Metro STR regulatory regime — out of scope for standard residential lease builder unless explicitly added later. |
| **Properties outside Davidson** | If Brickwise expands counties, branch templates: URLTA vs non-URLTA (§ 66-7-109 notice scheme). |

---

## 13. Suggested disclosure / clause checklist (Davidson County residential)

Use as a future feature checklist (not a shipped template):

1. Parties, premises address, term (start/end), rent amount & due date  
2. Security deposit amount + **bank/institution location** disclosure  
3. Late fee terms compliant with **5-day grace** and **≤ 10%**  
4. Written notice that landlord **does not insure** tenant personal property  
5. Manager name/address; owner/agent for service; maintenance phone/email/portal  
6. Optional: renters-insurance requirement  
7. Habitability / repair request procedure  
8. Entry / showing clause (24-hour showings in last 30 days if desired)  
9. Nonrenewal notice (**≥ 60 days** for ≥1-year leases if not already in lease)  
10. Rent-increase notice (**≥ 90 days** default unless lease embeds ≥60-day rule — confirm desired approach)  
11. Lead-based paint packet if pre-1978  
12. Optional bold **nonpayment notice waiver** (12-pt) — with education UI  
13. No waiver of URLTA rights; no confession of judgment / unlawful exculpation  
14. Reminder: Metro **landlord registration** for covered unit types (owner ops, not necessarily a lease clause)

---

## 14. Primary sources to re-check before building

| Resource | URL |
|---|---|
| TCA Title 66, Ch. 28 (Justia index) | https://law.justia.com/codes/tennessee/title-66/chapter-28/ |
| URLTA compiled PDF (industry copy, 2025) | https://assets-002.noviams.com/novi-file-uploads/caatn/URLTA_2025.pdf |
| Metro Ch. 11.22 | https://library.municode.com/tn/metro_government_of_nashville_and_davidson_county/codes/code_of_ordinances?nodeId=CD_TIT11PUPEMOWE_DIVIGERE_CH11.22RELATERE |
| Metro landlord registration | https://www.nashville.gov/departments/codes/property-standards/landlord-registration |
| § 66-28-107 registration statute | https://law.justia.com/codes/tennessee/title-66/chapter-28/part-1/section-66-28-107/ |
| EPA lead disclosure | https://www.epa.gov/lead/real-estate-disclosure |
| Help4TN URLTA overview | https://www.help4tn.org/senior-services/legal-information/long-term-care-housing/urlta-renters |

---

## 15. Explicit non-goals for this document

- **Do not** implement the lease builder feature from this file alone.  
- **Do not** treat secondary blogs as controlling over the Tennessee Code / Metro Code.  
- **Do** revisit statutes and Metro ordinances when the feature is scheduled; note the research date below.

*Research compiled: 2026-11-20 (America/New_York). Re-verify before product use.*
