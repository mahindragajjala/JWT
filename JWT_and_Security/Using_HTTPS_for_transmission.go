//Using_HTTPS_for_transmission
/*
- Even if your JWT has a valid signature and strong algorithm like 
  HS256 or RS256, without HTTPS, your token can be stolen during 
  transmission. 
- This leads to token hijacking — a serious security issue.
*/
/*
🔐 Why Is HTTPS Essential?
- HTTPS = HTTP + SSL/TLS encryption
- It encrypts all data between client and server.
- Without HTTPS, JWT sent in headers can be:
    🕵️‍♂️ Sniffed by attackers on public Wi-Fi
    📦 Intercepted via man-in-the-middle (MITM) attacks
    🧠 Read in plain-text (base64 ≠ encryption)
*/
