# Security Policy & Architecture

## Security Guarantees in Vula

Vula is designed with strict security defaults for modern open-source software:

### 1. Principle of Least Privilege
* **No Root AI Daemon:** The AI engine, Voice subsystem, and HUD run strictly under the unprivileged user space.
* **Explicit Sudo Boundaries:** `vula` never runs uncontained `sudo` loops. Privileged actions (like `apt install`) require explicit user consent and are scoped strictly to package managers.

### 2. Local-First Privacy
* **Local Ingestion:** By default, all AI completions and embeddings are processed locally via Ollama (`http://localhost:11434`).
* **No Telemetry / Phoning Home:** No telemetry, tracking, or user code snippets are sent to external servers.
* **API Key Isolation:** Cloud API keys (Gemini, Claude, OpenAI) are encrypted or stored locally in `~/.config/vula/config.yaml` with `0600` file permissions.

### 3. Human-in-the-Loop Command Execution
* The AI assistant never executes system-modifying or destructive commands (`rm`, `dd`, `mkfs`, `chmod`, `sudo`) autonomously.
* Suggested commands must be reviewed and triggered interactively by the user.

### 4. Supply Chain & Dependency Verification
* All binary downloads (e.g., STT/TTS models, third-party utilities) are checked against cryptographic SHA-256 hashes.
* Automated CI workflows run `govulncheck`, `trivy`, and `gitleaks` on every Pull Request.

---

## Reporting a Vulnerability

If you discover a security vulnerability within Vula, please follow responsible disclosure:

1. **Do NOT open a public GitHub issue.**
2. Send an email with details and reproduction steps to `security@vula-os.org` (or create a private GitHub Security Advisory).
3. We will respond within 48 hours to acknowledge receipt and coordinate a patch release.
