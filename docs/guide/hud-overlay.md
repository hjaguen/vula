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

### 1. Modo Lanzador de Comandos & Búsqueda Fuzzy (`>_ COMANDOS`)
- **Descripción:** Permite buscar y ejecutar instantáneamente cualquier aplicación instalada en Ubuntu, perfiles de trabajo, herramientas del sistema o comandos.
- **Uso:** Escribe el nombre del comando o filtro para ubicarlo rápidamente.

<p align="center">
  <img src="/vula_hud_commands.png" alt="Modo Comandos Vula HUD" style="border-radius: 10px; box-shadow: 0 4px 16px rgba(0,0,0,0.3); margin: 1.5rem 0;" />
</p>

### 2. Modo Acción OS (`⚡ ACCIÓN OS`)
- **Descripción:** Ejecuta instrucciones directas sobre el sistema operativo (ajustar brillo, volumen, bloquear pantalla, etc.).
- **Uso:** Escribe `do <instrucción>` o selecciona este modo.

<p align="center">
  <img src="/vula_hud_action.png" alt="Modo Acción OS Vula HUD" style="border-radius: 10px; box-shadow: 0 4px 16px rgba(0,0,0,0.3); margin: 1.5rem 0;" />
</p>

### 3. Modo Chat con IA Assist (`🤖 IA CHAT`)
- **Descripción:** Consulta a la IA nativa de Vula directamente desde el panel flotante con contexto del sistema.
- **Uso:** Escribe tu consulta en lenguaje natural y presiona `Enter`.

<p align="center">
  <img src="/vula_hud_ai.png" alt="Modo Chat IA Vula HUD" style="border-radius: 10px; box-shadow: 0 4px 16px rgba(0,0,0,0.3); margin: 1.5rem 0;" />
</p>

### 4. Modo Asistente y Dictado por Voz (`🎙 VOZ`)
- **Descripción:** Graba audio de voz, lo transcribe localmente con Whisper STT y ejecuta la instrucción o responde por TTS.
- **Uso:** Presiona `Enter` para grabar audio.

<p align="center">
  <img src="/vula_hud_voice.png" alt="Modo Voz Vula HUD" style="border-radius: 10px; box-shadow: 0 4px 16px rgba(0,0,0,0.3); margin: 1.5rem 0;" />
</p>

---

## ⌨️ Atajos del HUD

| Tecla | Acción |
| :--- | :--- |
| `Super + Space` | Abrir / Alternar HUD |
| `Tab` | Cambiar entre modos (`Lanzador` → `Temas` → `Chat IA`) |
| `↑` / `↓` | Navegar elementos de la lista |
| `Enter` | Ejecutar acción seleccionada / Enviar mensaje |
| `Esc` | Cerrar HUD |
