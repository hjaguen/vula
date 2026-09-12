# Tiling Híbrido & Atajos de Teclado

Vula ofrece un sistema de mosaico de ventanas (*Tiling*) **Híbrido**, combinando la matriz interactiva de **Tactile** con el encaje automático y bordes resaltados de **Tiling Assistant**.

---

## 🧩 Modos de Tiling

### 1. Malla Interactiva Tactile (`Super + T`)
Presionar **`Super + T`** abre una cuadrícula visual sobre el escritorio con letras asignadas a cada zona (`Q`, `W`, `E`, `A`, `S`, `D`). Te permite ubicar o expandir cualquier ventana en zonas fijas usando exclusivamente el teclado.

### 2. Snap por Medios & Cuartos de Pantalla
- **Mitad izquierda / derecha:** `Super + Left` / `Super + Right`
- **Maximizar / Restaurar:** `Super + Up` / `Super + Down`
- **Cuartos de pantalla:**
  - Arriba izquierda: `Super + Alt + U`
  - Arriba derecha: `Super + Alt + I`
  - Abajo izquierda: `Super + Alt + J`
  - Abajo derecha: `Super + Alt + K`

### 3. Resaltado Dinámico de Borde Activo
La ventana enfocada resalta automáticamente con un borde de 2px en el **color de acento del tema activo de Vula**, haciendo la navegación visual mucho más clara.

---

## 🖱️ Arrastre Universal (`Super + Click Izquierdo`)

Puedes mantener presionada la tecla **`Super` y hacer clic izquierdo en cualquier parte dentro de una ventana** (incluso en aplicaciones sin barra de título como Alacritty, Ghostty o VS Code) para moverla libremente.

---

## 🚀 Navegación por Espacios de Trabajo & Dock

- **Saltar a Espacio de Trabajo 1 al 6:** `Super + 1` .. `Super + 6`
- **Mover Ventana Activa a Espacio 1 al 6:** `Shift + Super + 1` .. `Shift + Super + 6`
- **Saltar a App del Dock 1 al 9:** `Alt + 1` .. `Alt + 9`

---

## 🎹 Editor Interactivo de Atajos (`vula keys`)

```bash
# Ver tabla de atajos configurados
vula keys list

# Abrir el editor interactivo TUI de atajos
vula keys

# Cambiar un atajo desde CLI
vula keys set hud "<Super>space"
```
