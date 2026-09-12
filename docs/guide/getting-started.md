# Primeros Pasos e Instalación

Bienvenido al manual oficial de **Vula**. Vula transforma Ubuntu 24.04 LTS en un entorno de desarrollo de vanguardia, altamente estético, ergonómico y dotado de inteligencia artificial nativa.

---

## 💻 Requisitos del Sistema

- **Sistema Operativo:** Ubuntu 24.04 LTS (Noble Numbat) o superior.
- **Entorno de Escritorio:** GNOME Shell 46+.
- **Arquitectura:** Linux x86_64.
- **Dependencias requeridas:** Go 1.24+ (para compilar desde código fuente), `curl`, `git`, `ffmpeg`, `dconf-cli`.

---

## ⚡ Instalación Automática

Ejecuta el script de instalación en tu terminal:

```bash
wget -qO- https://raw.githubusercontent.com/hjaguen/vula/main/scripts/install.sh | bash
```

El instalador automático se encargará de:
1. Compilar el binario optimizado `vula` y colocarlo en `~/.local/bin/vula`.
2. Configurar las fuentes Nerd Fonts y extensiones de GNOME Shell necesarias (`Tiling Assistant`, `Tactile`, `Blur my Shell`).
3. Registrar los accesos directos globales en GNOME (`Super + Space` para el HUD).

---

## 🧪 Verificación del Sistema (`vula doctor`)

Para verificar que todos los subsistemas (GNOME, Audio, Ollama, Proveedores de IA y Fuentes) estén correctamente configurados, ejecuta:

```bash
vula doctor
```

Para ver la tarjeta informativa de rendimiento del sistema:

```bash
vula fetch
```

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
