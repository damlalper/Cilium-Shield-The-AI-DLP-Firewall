### 3. Ürün Dokümanı: `PRD.md` (Product Requirements Document)
*Bu dosya projeyi bir "Hackathon Projesi"nden "Satılabilir Ürün"e dönüştürür. B2B vizyonunu gösterir.*

```markdown
# 📋 Product Requirements Document (PRD)
**Product Name:** Cilium-Shield
**Version:** 0.1.0 (MVP)
**Status:** Prototype

## 1. Problem Statement
As organizations adopt Large Language Models (LLMs), "Shadow AI" has become a critical risk. Developers might paste database dumps or sensitive customer info into LLM prompts. Traditional DLP (Data Loss Prevention) tools are network-perimeter based and blind to Kubernetes internal traffic flows, while API Gateways are too centralized and create bottlenecks.

## 2. Target Audience (B2B)
* **CISOs (Chief Information Security Officers):** Need compliance (GDPR/CCPA) assurance for AI initiatives.
* **Platform Engineers:** Want security without slowing down developer velocity.
* **FinTech & HealthTech Companies:** Sectors with strict PII regulations.

## 3. User Stories
* **As a Security Engineer**, I want to automatically redact credit card numbers sent to OpenAI so that we remain PCI-DSS compliant.
* **As a DevOps Lead**, I want to see a dashboard of which microservices are trying to leak data, so I can train the developers.
* **As a Developer**, I want the security layer to be transparent (no code changes) so I can focus on building features.

## 4. MVP Capabilities (Hackathon Scope)
| Feature | Priority | Description |
| :--- | :--- | :--- |
| **L7 Traffic Interception** | P0 | Intercept HTTP traffic to `api.openai.com`. |
| **Regex Redaction** | P0 | Detect and replace Credit Cards and Emails with `[REDACTED]`. |
| **Audit Logging** | P1 | Log the source Pod IP and timestamp of the leak. |
| **UI Dashboard** | P1 | Visualize total blocked threats. |

## 5. Future Roadmap (Post-Hackathon)
* **Q1 2026:** Integration with enterprise Vault for custom regex management.
* **Q2 2026:** "Block Mode" vs "Redact Mode" toggle via UI.
* **Q3 2026:** Support for gRPC traffic and proprietary LLM endpoints.