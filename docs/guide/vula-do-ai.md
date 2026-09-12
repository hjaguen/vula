# Control del Sistema Operativo por IA (`vula do`)

El comando `vula do` es el motor de control del sistema en lenguaje natural de Vula. Convierte frases habladas o escritas en español en planes de ejecución estructurados y seguros.

---

## ⚡ Uso Básico

```bash
vula do "<instrucción en lenguaje natural>"
```

### Ejemplos de uso:

```bash
# Control de pantalla y sonido
vula do "ajusta el brillo al 40%"
vula do "sube el volumen al 80%"
vula do "bloquea la pantalla"

# Lanzar aplicaciones y crear WebApps
vula do "abre Visual Studio Code"
vula do "crea una app web para ChatGPT en https://chatgpt.com"

# Tiling y Personalización
vula do "cambia el tema a tokyonight"
vula do "pon el espacio entre ventanas a 10px"

# Procesamiento de Medios
vula do "convierte la última grabación de pantalla a mp4"
```

---

## 🛡️ Clasificación de Riesgo & Confirmación

Para garantizar la máxima seguridad en tu equipo, cada acción planificada por Vula es evaluada y clasificada en uno de 3 niveles de riesgo antes de su ejecución:

| Nivel de Riesgo | Tipo de Operación | Comportamiento |
| :--- | :--- | :--- |
| **`SAFE`** | Lanzar apps, ajustar volumen, cambiar tema, ajustar brillo. | Se ejecuta automáticamente sin interrupciones. |
| **`SENSITIVE`** | Mover/copiar archivos, terminar procesos (`pkill`), crear directorios. | Solicita confirmación interactiva mediante un formulario Charm Huh. |
| **`HIGH`** | Eliminar archivos (`rm`), apagar/reiniciar sistema, modificar archivos raíz (`/etc`, `/boot`). | Despliega advertencia de seguridad estricta y valida rutas restringidas. |

---

## 📜 Historial de Auditoría (`vula actions history`)

Todas las acciones ejecutadas o canceladas se persisten con registro de fecha, prompt original, nivel de riesgo y estado en `~/.config/vula/audit.log`.

Para revisar el historial en cualquier momento:

```bash
vula actions history
```
