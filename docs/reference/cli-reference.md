# Referencia Completa de Comandos CLI

Cheat sheet rápido con todos los comandos y subcomandos disponibles en la interfaz CLI de **Vula**.

---

## ⚡ Comandos Principales

```bash
vula                           # Iniciar el HUD flotante interactivo (Super + Space)
vula doctor                    # Diagnóstico completo del sistema (GNOME, Audio, IA, Fuentes)
vula fetch                     # Tarjeta informativa de métricas y hardware
vula install                   # Instalador interactivo TUI
vula hud                       # Lanzar el HUD flotante Raycast TUI
vula do "<instrucción>"        # Ejecutar acción de OS en lenguaje natural mediante IA
vula listen [segundos]         # Iniciar escucha activa conversacional por voz
vula version                   # Mostrar versión de Vula
```

---

## 🤖 Subsistema de IA (`vula ai`)

```bash
vula ai ask "<pregunta>"       # Preguntar a la IA con contexto del sistema
vula ai cmd "<tarea>"          # Traducir tarea a comando shell
vula ai models                 # Listar modelos locales de Ollama
vula ai status                 # Diagnóstico de proveedores (Hybrid / Ollama / Gemini / Groq)
vula ai config                 # Configurar modo, proveedor, modelos y API keys
vula ai commit                 # Generar commit convencional a partir del git diff
vula ai fix "<error>"          # Diagnosticar error de terminal y proponer solución
vula ai explain                # Explicar código o texto activo en portapapeles
```

---

## 🎙️ Subsistema de Voz (`vula voice`)

```bash
vula voice record              # Grabar voz y transcribir directo a ventana activa
vula voice speak "<texto>"     # Sintetizar voz localmente con Piper TTS
vula voice test                # Probar volumen dB de entrada del micrófono
vula voice daemon              # Ejecutar demonio de voz en segundo plano
```

---

## 🎨 Temas & Fondos (`vula theme` / `vula wallpaper`)

```bash
vula theme list                # Listar 19 temas curados
vula theme set <nombre>        # Aplicar tema globalmente
vula theme create              # Abrir Vula Theme Studio TUI
vula theme generate "<prompt>" # Generar tema con IA
vula theme preview <nombre>    # Vista previa de paleta en terminal

vula wallpaper list            # Listar fondos descargados
vula wallpaper set <archivo>   # Establecer fondo específico
vula wallpaper next            # Rotar al siguiente fondo
vula wallpaper fetch           # Descargar colección HD de fondos
```

---

## 🌐 Aplicaciones Web (`vula webapp`)

```bash
vula webapp add <Nom> <URL>    # Crear lanzador de app web con favicon
vula webapp list               # Listar apps web gestionadas por Vula
vula webapp remove <Nom>       # Eliminar app web
```

---

## 🔤 Fuentes Dev (`vula font`)

```bash
vula font list                 # Listar fuentes dev Nerd Fonts
vula font set <nombre>         # Instalar y cambiar fuente (cascadia, firacode, jetbrains, meslo)
vula font size <pt>            # Cambiar tamaño de fuente en Ghostty, Alacritty y VS Code
```

---

## 🎬 Medios (`vula media`)

```bash
vula media webm2mp4 [archivo]  # Convertir grabación .webm a .mp4 con ffmpeg
```

---

## 🖥️ Escritorio & Tiling (`vula desktop` / `vula keys`)

```bash
vula desktop setup             # Aplicar optimizaciones y atajos de GNOME
vula desktop tiling            # Configurar Tiling Assistant & Tactile
vula desktop extensions list   # Listar extensiones GNOME curadas
vula desktop extensions install-curated # Instalar extensiones visuales

vula keys list                 # Listar tabla de atajos globales
vula keys                      # Abrir editor TUI de atajos
vula keys set <id> <shortcut>  # Cambiar atajo en dconf
```

---

## 🧰 Aplicaciones & Dotfiles (`vula apps` / `vula dotfiles`)

```bash
vula apps list                 # Catálogo de aplicaciones dev
vula apps install-cli          # Instalar stack CLI (fzf, zoxide, btop, lazygit)
vula apps ui                   # App Store TUI interactivo
vula dotfiles install          # Instalar dotfiles curados (Fish, Starship, Neovim, Tmux)
```
