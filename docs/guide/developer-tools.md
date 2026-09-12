# Herramientas Dev (Fuentes, WebApps & Medios)

Vula integra herramientas dedicadas para optimizar el flujo de trabajo de desarrolladores.

---

## 🔤 Gestor de Fuentes Dev (`vula font`)

Sincroniza fuentes **Nerd Fonts** de alta calidad (con glifos e íconos) entre terminales (`Ghostty`, `Alacritty`), editores (`VS Code`) e interfaz de GNOME.

```bash
# Listar fuentes disponibles
vula font list

# Cambiar fuente dev (descarga automática de Nerd Font)
vula font set jetbrains      # Opciomnes: cascadia, firacode, jetbrains, meslo

# Cambiar tamaño de fuente en pt
vula font size 12
```

---

## 🌐 Creador de Aplicaciones Web (`vula webapp`)

Convierte cualquier servicio o sitio web en una aplicación de escritorio nativa con ícono HD automático y acceso rápido desde el HUD (`Super + Space`).

```bash
# Crear una Web App (ejemplo Notion o ChatGPT)
vula webapp add Notion https://notion.so
vula webapp add ChatGPT https://chatgpt.com

# Listar Web Apps creadas por Vula
vula webapp list

# Eliminar una Web App
vula webapp remove Notion
```

---

## 🎬 Conversor de Grabaciones de Pantalla (`vula media`)

Convierte los videos de grabación de pantalla `.webm` generados por GNOME a formato `.mp4` (H.264 + AAC) compatible con la web y redes sociales.

```bash
# Convertir automáticamente la última grabación realizada
vula media webm2mp4

# Convertir un archivo específico
vula media webm2mp4 ~/Videos/mi-grabacion.webm
```

---

## 💻 Herramientas CLI & Dotfiles (`vula apps` / `vula dotfiles`)

```bash
# Instalar el stack de herramientas CLI modernas (fzf, zoxide, btop, lazygit, ripgrep)
vula apps install-cli

# Abrir el App Store TUI interactivo
vula apps ui

# Instalar dotfiles curados (Fish, Starship, Neovim, Tmux)
vula dotfiles install
```

---

## 🤖 Utilidades de IA para Git & Código (`vula ai`)

```bash
# Generar commit convencional a partir del git diff activo
vula ai commit

# Diagnosticar un error de terminal y obtener solución
vula ai fix "error message text"

# Explicar código o texto seleccionado en el portapapeles
vula ai explain
```
