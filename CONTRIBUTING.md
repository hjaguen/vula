# Contributing to Vula

We welcome contributions to Vula! Whether you are writing Go code, designing UI in Charm, optimizing GNOME extensions, or enhancing AI/voice models, here is how you can help.

---

## Development Setup

1. **Prerequisites:**
   * Ubuntu 24.04 LTS (Noble Numbat)
   * Go 1.24+
   * `make`, `git`, `curl`

2. **Clone & Build:**
   ```bash
   git clone https://github.com/vula-os/vula.git
   cd vula
   make build
   ```

3. **Run System Diagnostics:**
   ```bash
   make doctor
   ```

4. **Launch the Floating HUD:**
   ```bash
   make hud
   ```

---

## Code Quality Standards

* **Go Code Formatting:** Run `make fmt` before committing.
* **Testing:** Ensure `make test` passes without race conditions.
* **Idempotence:** Any new installer or system modification step must verify state before modifying user files.
* **Charm Aesthetics:** Follow standard Lipgloss color guidelines defined in `internal/ui/styles.go`.

---

## Pull Request Guidelines

1. Fork the repo and create a descriptive branch: `git checkout -b feature/voice-streaming`.
2. Write tests covering new packages or logic.
3. Submit the PR with a clear summary of changes and test steps.
4. Ensure all GitHub Actions CI checks pass.
