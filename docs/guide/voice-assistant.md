# Asistente y Dictado por Voz

Vula integra un subsistema completo de voz con procesamiento local y baja latencia, combinando **Whisper STT** para reconocimiento de voz y **Piper TTS** para la síntesis de habla humana.

---

## 🎙️ Dictado de Voz a Texto (`Super + Alt + V`)

Transcribe tu voz en tiempo real e inyecta el texto transcrito directamente en la ventana o editor que tengas enfocado (VS Code, terminal, navegador, etc.).

- **Atajo global:** `Super + Alt + V`
- **Comando CLI:**
  ```bash
  vula voice record
  ```

---

## 🤖 Asistente Conversacional Activo (`Super + Alt + A`)

Escucha tu voz durante unos segundos, envía el audio a Whisper, procesa la consulta con la IA de Vula y te responde tanto en texto como con voz humana sintetizada mediante Piper TTS.

- **Atajo global:** `Super + Alt + A`
- **Comando CLI:**
  ```bash
  vula listen [segundos_de_escucha]
  ```

---

## 🔊 Síntesis de Texto a Voz (`vula voice speak`)

Para sintetizar cualquier texto localmente en voz alta:

```bash
vula voice speak "Hola Mauricio, el sistema está funcionando perfectamente."
```

---

## 🧪 Diagnóstico de Micrófono (`vula voice test`)

Si experimentas problemas con la captura de audio:

```bash
vula voice test
```

Este comando grabará durante 3 segundos y analizará el nivel máximo de dB de tu entrada de audio para confirmar que el micrófono interno esté capturando señal limpia.
