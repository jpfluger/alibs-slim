# FIPS 140-3 Compliance Guide

## **What is FIPS 140-3?**
[FIPS 140-3](https://csrc.nist.gov/publications/detail/fips/140/3/final) (Federal Information Processing Standard 140-3) is the latest U.S. government standard for cryptographic module security. It defines requirements for encryption, key management, authentication, and other security controls.

FIPS 140-3 **ensures cryptographic modules are tested and certified** for security, making it critical for:
- **Government agencies & contractors** (FedRAMP, DoD, etc.).
- **Healthcare & finance sectors** (HIPAA, PCI-DSS).
- **Software handling sensitive user data** (passwords, JWT tokens, VPNs).

## **Why Does This Matter?**
If your application is used in **FIPS-mandated environments**, you must use **only FIPS-validated cryptographic libraries**. This means:
✅ Using **FIPS-certified cryptographic modules**.  
✅ Avoiding **non-compliant algorithms** (e.g., MD5, non-FIPS RNGs).  
✅ Running software in **FIPS mode** (OS-dependent).  

Failure to comply may result in **audit failures, security vulnerabilities, or software incompatibilities**.

---

## **FIPS 140-3 in Programming: Supported Libraries**
The table below lists cryptographic libraries in **Python** and **Golang** that support FIPS 140-3 compliance.

| **Language** | **Library**          | **FIPS Mode Support** | **Notes** |
|-------------|----------------------|----------------------|----------|
| **Python**  | `pyca/cryptography`   | ✅ **FIPS-certified OpenSSL** backend | Requires OpenSSL in FIPS mode |
| **Python**  | `hashlib` (stdlib)    | ❌ No native FIPS support | Relies on OpenSSL, may work in FIPS mode |
| **Python**  | `pyOpenSSL`           | ✅ Uses FIPS-certified OpenSSL | Ensure OS-level FIPS mode is enabled |
| **Python**  | `bcrypt`              | ✅ Uses OpenSSL | Compatible with FIPS-certified OpenSSL |
| **Golang**  | Standard `crypto/*`   | ✅ **Go Cryptographic Module (experimental)** | FIPS mode available in Go 1.24+ |
| **Golang**  | `golang.org/x/crypto` | ⚠️ Partially supports FIPS | Some packages use non-FIPS algorithms |
| **Golang**  | `cloud.google.com/go/kms` | ✅ Google Cloud KMS is FIPS-certified | For cloud-based key management |
| **Golang**  | `aws/aws-sdk-go` | ✅ AWS KMS supports FIPS mode | Use GovCloud or FIPS endpoints |

### **Go 1.24+ FIPS Mode**
With Go **1.24**, you can enable FIPS mode:
```sh
GODEBUG=fips140=on go run main.go
```
For strict enforcement:
```sh
GODEBUG=fips140=only go run main.go
```

---

## **How to Enable FIPS Mode in Your Environment**
### **Linux (RHEL, Ubuntu, Amazon Linux)**
1. Enable FIPS mode:
   ```sh
   sudo fips-mode-setup --enable
   ```
2. Reboot the system:
   ```sh
   sudo reboot
   ```
3. Verify:
   ```sh
   fips-mode-setup --check
   ```

### **Windows**
- Windows **CryptoAPI** can enforce FIPS compliance via **Group Policy**.
- Use FIPS-certified libraries like **BCrypt** and **CNG API**.

### **MacOS**
- No official FIPS mode, but OpenSSL libraries can be compiled with FIPS support.

---

## **FAQ**
### **Q: Can I use non-FIPS cryptography if my software doesn't handle sensitive data?**
Yes, but if you work in **regulated industries**, it's safer to stick with **FIPS-validated modules**.

### **Q: How do I check if my system is running in FIPS mode?**
On Linux:
```sh
cat /proc/sys/crypto/fips_enabled
```
If the output is `1`, FIPS mode is enabled.

### **Q: What happens if I run non-FIPS cryptography on a FIPS-enabled system?**
- The system **may block non-FIPS modules**.
- Your application **may fail to start** if it relies on non-compliant crypto.

---

## **Next Steps**
✅ Use only **FIPS-certified** cryptographic libraries.  
✅ Enable **FIPS mode** if required.  
✅ Test your software with **FIPS-compliant dependencies**.  
✅ Stay updated on **Go Cryptographic Module certification** (Go 1.24+).

For more information, refer to:
- [NIST FIPS 140-3](https://csrc.nist.gov/publications/detail/fips/140/3/final)
- [Go FIPS 140 Guide](https://go.dev/doc/security/fips140)
- [Python Cryptography Docs](https://cryptography.io/en/latest/)
