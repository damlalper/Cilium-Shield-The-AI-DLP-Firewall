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


### Yapay Zekanın oluşturacağı mvpden sonra test etmen gereken şeyler
#### 🛠️ I. Entegrasyon ve "Yapıştırıcı Kod" (Glue Code)

Yapay zekanın ürettiği ayrı ayrı kod bloklarını (Wasm, Go, Next.js) bir araya getirip aralarındaki iletişimi kurmalısın.

---

##### 1. ⚙️ Wasm Enjeksiyonu ve Doğrulama

###### ✅ Wasm Modülünü Derleme
- TinyGo veya Rust ile yazılan Wasm filtresini doğru hedefe derlemelisin.
- Derlenen modülü bir container registry'ye (örneğin Docker Hub) yüklemelisin.

###### ✅ CiliumEnvoyConfig İnce Ayarı
- Cilium CRD'sindeki URL'nin registry'deki Wasm modülü adresiyle *birebir eşleştiğinden* emin olmalısın.
- Bu YAML dosyasında bir **boşluk bile tüm filtre zincirini (filter chain) kırabilir**.

###### ✅ İletişim Kanalı
- Wasm filtresinin maskelediği verinin bilgisini:
  - HTTP Header
  - veya gRPC/HTTP isteği
  aracılığıyla Go Backend'e gönderdiğinden emin olmalısın.
- Bu "yapıştırıcı kod"un hatasız çalışması gerekir.

---

##### 2. 🔌 Cilium Politikalarının Yönetimi
- Go Backend yalnızca log almakla kalmamalı,
- gerektiğinde Cilium Network Politikalarını (CNP) okuyup yönetebilmelidir.

> Basit bir senaryo için zorunlu değil, ancak **Teknik Derinlik puanını artırır.**

---

#### 💻 II. Optimizasyon ve Mühendislik Mükemmeliyeti

Yapay zekanın kolayına kaçtığı veya "hackathon yeterli" dediği yerleri **üretim kalitesine** taşımalısın.

---

##### 1. 🔍 RegEx Optimizasyonu
- Yazılan regex kurallarının gerçek:
  - PII (Kişisel Tanımlayıcı Bilgi)
  - API key formatlarını
  doğru yakaladığını manuel olarak test etmelisin.
- Basit regex yerine:
  - kredi kartı için **Luhn Algoritması**
  gibi doğrulama yöntemlerini TinyGo Wasm içinde uygulamalısın.

---

##### 2. ⚡ Performans Testleri
- Envoy proxy'ye saniyede yüzlerce istek göndererek load test yapılmalı.
- Wasm filtresinin latency (gecikme) değerleri ölçülmeli.
- Amaç:
  - Wasm filtresinin trafiği yavaşlatmadığını göstermek.

> Bu, projenin Teknik Derinlik kriterini karşılayan en somut kanıttır.

---

##### 3. 📉 Hata ve Bellek Yönetimi
- Go Backend'de goroutine'lerin gerektiği gibi temizlendiğinden emin olmalısın.
- Wasm tarafında:
  - hafıza sızıntısı (memory leak) olmadığından emin olmalısın.
  - TinyGo'daki malloc/free benzeri yapıyı optimize etmelisin.

---

#### 📝 III. Doğrulama ve Sunum

Bu aşama, projenin:
- Clarity (Netlik)
- Yaratıcılık
kriterlerini karşılaması için hayati önem taşır.

---

##### 1. 📹 Video Senaryosu (Storytelling)

###### 🎬 Önce Hata
- Demo videosunda sistem devre dışıyken
- bir kredi kartı numarasının LLM API'sine gönderildiğini göstermelisin.

###### ✅ Sonra Çözüm
- Sistemi deploy edip
- aynı isteğin gönderildiğinde verinin:
  - anında `[REDACTED]` olarak maskelendiğini
  - Next.js Dashboard'da loglandığını
  göstermelisin.

> Bu, projenin değerini anında kanıtlar.

---

##### 2. 🖼️ Mimari Görseller
- ARCHITECTURE.md dosyasındaki diyagramlar:
  - net
  - profesyonel
  - projenin değerini vurgulayan
  olmalıdır.
- Mermaid veya harici araçlarla çizilebilir.

> Bir görsel, binlerce kelimeden daha etkilidir.

---

##### 3. 🐛 Hata Belgeleme
- AI'ın kaçırdığı veya senin çözdüğün entegrasyon hatalarını
  README veya **HATA_GİDERME.md** dosyasına eklemelisin.
- Örnek:
  - “TinyGo Wasm'da JSON parse hatasını nasıl çözdüm”

> Bu, jüriye gerçek mühendislik sorunlarıyla karşılaştığını ve onları aştığını gösterir.



## 5. Future Roadmap (Post-Hackathon)
* **Q1 2026:** Integration with enterprise Vault for custom regex management.
* **Q2 2026:** "Block Mode" vs "Redact Mode" toggle via UI.
* **Q3 2026:** Support for gRPC traffic and proprietary LLM endpoints.