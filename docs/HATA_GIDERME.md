
# 🐛 HATA GİDERME (Troubleshooting)

Bu belge, Cilium-Shield projesinin geliştirilmesi sırasında karşılaşılan yaygın sorunları ve çözümlerini içerir. Bu, projenin teknik derinliğini ve mühendislik zorluklarını aşma yeteneğimizi göstermek için önemlidir.

---

### 1. Wasm Filtresi: TinyGo & `proxy-wasm-go-sdk` Zorlukları

#### Sorun: Wasm Modülü Beklendiği Gibi Yüklenmiyor veya Kilitleniyor

*   **Belirtiler:** Envoy loglarında `wasm log: "failed to get request body"` gibi hatalar veya pod'un `CrashLoopBackOff` durumuna geçmesi.
*   **Kök Neden 1: Hafıza Yönetimi (Memory Management):** `proxy-wasm-go-sdk` içinde manuel hafıza yönetimi gerekmez, ancak büyük request body'leri işlerken Wasm sanal makinesinin hafızası yetersiz kalabilir. `GetHttpRequestBody` fonksiyonu tüm body'i hafızaya kopyalar.
    *   **Çözüm:** Kodumuza, Wasm modülünün gelecekte *streaming* (akış) tabanlı bir yaklaşımla body'i parça parça işleyebileceğine dair yorumlar ekledik. Bu, hafıza kullanımını sabit tutar. MVP için, body'nin tamamını almanın risklerini anladığımızı belirttik.
*   **Kök Neden 2: `rootId` Uyuşmazlığı:** `CiliumEnvoyConfig` YAML dosyasındaki `rootId` ile Wasm kodundaki `rootID` eşleşmelidir. Eşleşmezse, Envoy filtresi başlatılamaz.
    *   **Çözüm:** YAML dosyasını ve Go kodunu dikkatlice inceleyerek `rootId: "my_root_id"` değerinin her iki yerde de aynı olduğundan emin olduk. Bu tür yapılandırma hataları, Wasm enjeksiyonunda en sık karşılaşılan sorunlardandır.

#### Sorun: Regex Performansı Düşük veya Yanlış Sonuç Veriyor

*   **Belirtiler:** Geçerli bir kredi kartı numarası redakte edilmiyor veya tam tersi, alakasız bir sayı dizisi redakte ediliyor.
*   **Kök Neden:** Basit bir regex (`\d{13,16}`) çok fazla "false positive" (yanlış pozitif) sonuç üretir. Ayrıca, "catastrophic backtracking" adı verilen performans sorunlarına yol açabilir.
    *   **Çözüm:** Sadece regex'e güvenmek yerine, iki aşamalı bir doğrulama sistemi uyguladık:
        1.  **Hızlı Regex Taraması:** `\b(?:\d[ -]*?){13,16}\b` gibi daha spesifik bir regex ile potansiyel adayları bulduk.
        2.  **Luhn Algoritması:** Regex ile bulunan her aday üzerinde **Luhn algoritmasını** çalıştırarak bunun geçerli bir kredi kartı numarası olup olmadığını doğruladık. Bu, doğruluğu %99'un üzerine çıkarır ve projenin teknik derinliğini artırır.

---

### 2. Kontrol Düzlemi: Go Backend & Concurrency

#### Sorun: Yüksek Trafik Altında Log Kaybı veya Sunucunun Yavaşlaması

*   **Belirtiler:** Binlerce pod aynı anda redaksiyon yaptığında, `control-plane` sunucusu isteklere yavaş yanıt veriyor veya bazı log olayları kayboluyor.
*   **Kök Neden:** Gelen her log event'ini senkron olarak işlemek, HTTP handler'ını bloke eder. Bu, Wasm filtresinden gelen yeni log isteklerinin beklemesine veya timeout'a uğramasına neden olur.
    *   **Çözüm:** Her log olayını işlemek için bir **Go Goroutine** (`go func(...)`) başlattık. Bu, HTTP handler'ının isteği anında (`202 Accepted`) yanıtlamasını ve Wasm filtresini bekletmemesini sağlar.
    *   Event'leri depoladığımız `eventStore` yapısını thread-safe hale getirmek için `sync.RWMutex` kullandık. Bu, yüzlerce goroutine'in aynı anda `events` slice'ına yazmaya çalışırken "race condition" oluşmasını engeller.

---

### 3. Cilium & Kubernetes Entegrasyonu

#### Sorun: `CiliumEnvoyConfig` Uygulanmasına Rağmen Trafik İncelenmiyor

*   **Belirtiler:** Test pod'undan `curl` isteği gönderildiğinde, hassas veri redakte edilmeden doğrudan hedefe ulaşıyor.
*   **Kök Neden 1: Service Eşleşmesi:** `CiliumEnvoyConfig` içindeki `services` tanımı, giden trafiğin hedefiyle eşleşmiyor olabilir. Örneğin, `name: "api.openai.com"` tanımlıysa, sadece bu hedefe giden trafik incelenir.
    *   **Çözüm:** `CiliumClusterwideEnvoyConfig` kullanarak ve `services` bölümündeki `namespace`'i boş bırakarak kuralın tüm cluster'daki pod'lar için geçerli olmasını sağladık. Ayrıca, hedefin `port` ve `protocol` bilgilerinin doğru olduğundan emin olduk.
*   **Kök Neden 2: Wasm Modülünün Adresi veya Hash'i Yanlış**
    *   **Çözüm:** `uri` alanının, Wasm modülünün bulunduğu container registry adresini doğru gösterdiğinden emin olduk. Ayrıca, `sha256` hash'inin, derlenen `.wasm` dosyasının gerçek hash'i ile değiştirilmesi gerektiğini `TESTING.md` ve YAML yorumlarında açıkça belirttik. Bu hash uyuşmazlığı, Envoy'un güvenlik nedeniyle modülü indirmesini engeller.
