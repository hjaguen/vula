# Primeros Pasos e Instalación

Bienvenido al manual oficial de **Vula**. Vula transforma Ubuntu 24.04 LTS en un entorno de desarrollo de vanguardia, altamente estético, ergonómico y dotado de inteligencia artificial nativa.

---

## 💻 Requisitos del Sistema

- **Sistema Operativo:** Ubuntu 24.04 LTS (Noble Numbat) o superior.
- **Entorno de Escritorio:** GNOME Shell 46+.
- **Arquitectura:** Linux x86_64.
- **Dependencias requeridas:** Go 1.24+ (para compilar desde código fuente), `curl`, `git`, `ffmpeg`, `dconf-cli`.

---

## ⚡ Instalación Universal One-Liner (Cero Fricción)

Ejecuta el script oficial en cualquier terminal de Ubuntu 24.04+ LTS:

```bash
curl -fsSL https://raw.githubusercontent.com/hjaguen/vula/main/scripts/install.sh | bash
```

El instalador automático se encarga de:
1. **Detección de Arquitectura:** Descarga el binario precompilado más reciente (`amd64` o `arm64`) desde GitHub Releases en solo 3 segundos (o compila automáticamente con Go si no hay binario de release).
2. **Snapshot de Seguridad:** Ejecuta de forma preventiva `vula backup create pre-vula-install` para congelar tus dotfiles y preferencias de GNOME previas.
3. **Entorno Universal de Shell:** Configura `~/.config/vula/shell.sh` e inyecta los atajos clave en `~/.bashrc`, `~/.zshrc` y `config.fish`:
   - `do`: Acciones en lenguaje natural (`vula do`).
   - `c` / `cx`: Lanzadores de agentes autónomos OpenCode y Claude Code.
   - `hud`: Lanzador HUD flotante (`Super + Space`).
   - `fleet`: Monitor de la flota de agentes Herdr (`vula agents status`).
4. **Habilidad para Agentes:** Genera `~/.config/vula/SKILL.md` (`vula skill`).
5. **Diagnóstico Automatizado:** Ejecuta `vula doctor` para certificar la salud del sistema.

---

## 🔄 Actualización en Caliente (`vula update`)

Mantener tu entorno al día no requiere volver a clonar repositorios ni compilar a mano. Simplemente ejecuta:

```bash
vula update
```

Vula comprobará si existe una nueva versión oficial en GitHub Releases, descargará la actualización y reemplazará el binario de forma atómica y segura. También puedes usar el alias `vula upgrade`.

---

## 🧪 Verificación del Sistema (`vula doctor`)

Para verificar que todos los subsistemas (GNOME, Audio, Ollama, Proveedores de IA y Fuentes) estén correctamente configurados, ejecuta:

```bash
vula doctor
```

Para ver la tarjeta informativa de rendimiento y especificaciones del sistema:

```bash
vula fetch
```

<p align="center">
  <img src="/vula_fetch_showcase.png" alt="Vula Fetch System Output" style="border-radius: 10px; box-shadow: 0 4px 16px rgba(0,0,0,0.3); margin: 1.5rem 0;" />
</p>

---

## 🧠 Arquitectura de IA Híbrida

Vula incluye un motor de IA **Híbrido** inteligente:
- **Local (Privacidad y Rapidez):** Usa modelos locales mediante **Ollama** (ej. `qwen2.5-coder:1.5b` o `llama3.2`).
- **Nube (Razonamiento Complejo):** Se conecta automáticamente a proveedores en la nube como **Google Gemini 2.0 Flash**, **Groq** u **OpenAI** cuando las tareas requieren alto cómputo.

Puedes verificar o cambiar la configuración de IA en cualquier momento:

```bash
# Ver estado de los proveedores
vula ai status

# Configurar el modo híbrido y tu API Key de Gemini
vula ai config --mode=hybrid --provider=gemini --key=gemini:TU_API_KEY
```
