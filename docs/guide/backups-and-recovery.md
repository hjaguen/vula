# Snapshots, Backup & Recuperación Segura

Uno de los mayores temores de los desarrolladores al probar un nuevo entorno o sistema operativo es que sus configuraciones personales, atajos de GNOME o archivos `.bashrc` / `.zshrc` queden alterados sin posibilidad de retorno.

**Vula OS** incluye un motor nativo de **Snapshots y Rollback** (`vula backup`) diseñado con una filosofía de **cero riesgo**: en cualquier momento puedes congelar el estado de tu sistema o revertir cualquier cambio en un solo segundo.

---

## 🛡️ ¿Qué contiene un Snapshot de Vula?

Cada snapshot guarda de forma atómica:
1. **Volcado completo de GNOME Shell (`dconf`):** Todas las preferencias de ventanas, extensiones activas, temas y comportamiento del escritorio.
2. **Atajos de teclado personalizados (`custom-keybindings`):** Tu mapa exacto de combinaciones de teclas en GNOME.
3. **Dotfiles críticos del usuario:**
   - `~/.bashrc` y `~/.zshrc`
   - `~/.config/fish/config.fish`
   - `~/.tmux.conf`
   - `~/.config/starship.toml`
   - `~/.config/nvim/init.lua`
   - `~/.config/vula/config.yaml`
   - `~/.config/vula/shell.sh`
4. **Metadatos JSON:** Marca de tiempo, nombre, descripción y lista de archivos preservados.

---

## 📸 Creación de Snapshots (`vula backup create`)

Para congelar el estado actual de tu sistema:

```bash
# Snapshot manual rápido
vula backup create

# Snapshot con nombre descriptivo
vula backup create "antes-de-probar-tokyonight"

# Snapshot con nota explicativa
vula backup create "estable-dev" -d "Configuración perfecta antes de actualizar paquetes"
```

Los snapshots se almacenan organizados en `~/.config/vula/backups/<timestamp>_<nombre>/`.

---

## 📋 Listar Snapshots Disponibles (`vula backup list`)

Puedes ver el historial de copias de seguridad de tu máquina:

```bash
vula backup list
```

Salida típica:

```text
  VULA // VULA SYSTEM SNAPSHOTS  
Found 2 snapshots

  • 20261117-043836_pre-vula-install pre-vula-install (GNOME dconf, 7 dotfiles) - 2026-11-17 04:38
    Nota: Initial system snapshot prior to Vula OS configuration
  • 20261117-051200_estable-dev      estable-dev      (GNOME dconf, 7 dotfiles) - 2026-11-17 05:12
    Nota: Configuración perfecta antes de actualizar paquetes
```

---

## ⏪ Restauración Inmediata (`vula backup restore`)

Si deseas deshacer cambios de apariencia, recuperar tus atajos previos o restaurar tus dotfiles:

```bash
vula backup restore pre-vula-install
```

Puedes buscar por el nombre que le diste (`pre-vula-install`) o por el identificador de fecha (`20261117-043836`).

En un instante:
- Todos los dotfiles se copian de vuelta a `$HOME`.
- La base de datos `dconf` de GNOME se restaura automáticamente con `dconf load`.
- El entorno vuelve exactamente a como estaba.

---

## 🗑️ Eliminar Snapshots Antiguos (`vula backup delete`)

Para liberar espacio o limpiar snapshots temporales:

```bash
vula backup delete pre-vula-install
```

---

## ⚡ Protección Automática Pre-Instalación

No necesitas acordarte de hacer un snapshot antes de instalar:
- Al correr el instalador oficial (`curl -fsSL .../install.sh | bash`), Vula ejecuta automáticamente:
  ```bash
  vula backup create pre-vula-install
  ```
- Si ejecutas el instalador interactivo (`vula install`), también se toma un snapshot automático denominado `pre-installer-setup`.

---

## 🖥️ Acceso desde el HUD (`Super + Space`)

También puedes crear snapshots sin abrir la terminal:
1. Presiona `Super + Space`.
2. Escribe o selecciona **"Create System Snapshot"** (`💾`).
3. El HUD congelará el estado del sistema en segundo plano.
