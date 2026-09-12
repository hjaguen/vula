# HUD Overlay Interactivo (`Super + Space`)

El **Vula HUD** es la pieza central de interacción del escritorio. Se trata de una interfaz TUI flotante (estilo Raycast / Spotlight) construida con la suite de Charm (`Bubbletea`, `Bubbles`, `Lipgloss`) de ultra bajo consumo de memoria y ejecución instantánea.

![Vula HUD Showcase](/vula_hud_showcase.png)

---

## 🚀 Cómo invocar el HUD

Puedes abrir el HUD en cualquier momento presionando la combinación global:

```text
Super + Space
```

O ejecutando en la terminal:

```bash
vula hud
```

---

## 🎛️ Modos de Operación del HUD

El HUD cuenta con tres modos principales navegables con la tecla **`Tab`**:

### 1. Modo Lanzador de Aplicaciones & Búsqueda Fuzzy (`App Launcher`)
- **Descripción:** Permite buscar y ejecutar instantáneamente cualquier aplicación instalada en Ubuntu, aplicaciones web creadas con `vula webapp` o comandos del sistema.
- **Uso:** Simplemente empieza a escribir el nombre de la app (ejemplo: `code`, `brave`, `chatgpt`) y presiona `Enter` para lanzarla.

### 2. Modo Selector Interactivo de Temas (`Themes`)
- **Descripción:** Muestra la colección completa de **19 temas curados** en una lista paginada y adaptable.
- **Uso:** Navega con las flechas `↑` y `↓`. Al presionar `Enter`, Vula aplica instantáneamente el tema seleccionado, cambiando la paleta de colores, el fondo de pantalla en alta resolución y el modo claro/oscuro de GNOME.

### 3. Modo Chat con IA Assist (`AI Chat`)
- **Descripción:** Consulta a la IA nativa de Vula directamente desde el panel flotante. La IA recibe el contexto activo del sistema.
- **Uso:** Escribe tu consulta o comando (ejemplo: `"¿Cómo encuentro archivos de más de 100MB?"`) y recibe respuestas estructuradas en streaming.

---

## ⌨️ Atajos del HUD

| Tecla | Acción |
| :--- | :--- |
| `Super + Space` | Abrir / Alternar HUD |
| `Tab` | Cambiar entre modos (`Lanzador` → `Temas` → `Chat IA`) |
| `↑` / `↓` | Navegar elementos de la lista |
| `Enter` | Ejecutar acción seleccionada / Enviar mensaje |
| `Esc` | Cerrar HUD |
