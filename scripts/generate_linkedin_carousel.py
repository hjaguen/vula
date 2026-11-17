#!/usr/bin/env python3
"""
Generate a professional, high-converting LinkedIn Carousel PDF for Vula OS.
Format: 16:9 Landscape Slide Deck (1920x1080 standard layout).
"""

import os
from reportlab.lib import colors
from reportlab.lib.colors import HexColor
from reportlab.pdfgen import canvas
from PIL import Image

WIDTH = 1920
HEIGHT = 1080
OUTPUT_PDF = "/home/mauricio/repos/vula/vula_linkedin_carousel.pdf"
ASSETS_DIR = "/home/mauricio/repos/vula/docs/public"

# Brand Palette
BG_DARK = HexColor("#0A0F1D")
BG_CARD = HexColor("#111827")
TEXT_WHITE = HexColor("#F8FAFC")
TEXT_MUTED = HexColor("#94A3B8")
CYAN_ACCENT = HexColor("#06B6D4")
PURPLE_ACCENT = HexColor("#8B5CF6")
BORDER_COLOR = HexColor("#1E293B")


def draw_header_footer(c, slide_num, total_slides, category="VULA OS // DEVELOPER ENVIRONMENT"):
    # Header Accent bar
    c.setStrokeColor(CYAN_ACCENT)
    c.setLineWidth(4)
    c.line(100, HEIGHT - 70, 240, HEIGHT - 70)
    
    c.setFont("Helvetica-Bold", 18)
    c.setFillColor(CYAN_ACCENT)
    c.drawString(260, HEIGHT - 76, category)
    
    # Footer
    c.setFont("Helvetica", 18)
    c.setFillColor(TEXT_MUTED)
    c.drawString(100, 50, "github.com/hjaguen/vula  •  Ubuntu 24.04 LTS  •  100% Go + Charm TUI")
    
    # Page indicator
    c.setFont("Helvetica-Bold", 20)
    c.setFillColor(TEXT_WHITE)
    c.drawRightString(WIDTH - 100, 50, f"{slide_num} / {total_slides}  →")


def create_slide_1(c):
    """Cover Slide"""
    # Background
    c.setFillColor(BG_DARK)
    c.rect(0, 0, WIDTH, HEIGHT, fill=1, stroke=0)
    
    # Top badge
    c.setFillColor(HexColor("#1E293B"))
    c.roundRect(100, HEIGHT - 150, 420, 44, 22, fill=1, stroke=0)
    c.setFont("Helvetica-Bold", 18)
    c.setFillColor(CYAN_ACCENT)
    c.drawString(130, HEIGHT - 122, "⚡ NEXT-GEN DEVELOPER ENVIRONMENT")
    
    # Title
    c.setFont("Helvetica-Bold", 72)
    c.setFillColor(TEXT_WHITE)
    c.drawString(100, HEIGHT - 240, "¿Por qué elegir entre la")
    c.setFillColor(CYAN_ACCENT)
    c.drawString(100, HEIGHT - 325, "estabilidad de Ubuntu LTS")
    c.setFillColor(TEXT_WHITE)
    c.drawString(100, HEIGHT - 410, "y la experiencia moderna de IA?")

    # Subtitle
    c.setFont("Helvetica", 28)
    c.setFillColor(TEXT_MUTED)
    c.drawString(100, HEIGHT - 480, "Te presento Vula OS: un orquestador 100% Go + Charm TUI diseñado para")
    c.drawString(100, HEIGHT - 520, "desarrolladores exigentes que aman el flujo keyboard-first.")
    
    # Embedded Hero Image
    img_path = os.path.join(ASSETS_DIR, "vula_brand_showcase.png")
    if os.path.exists(img_path):
        # Target box: x=1000, y=140, w=820, h=458
        c.drawImage(img_path, 980, 150, width=840, height=469, preserveAspectRatio=True, mask='auto')
        c.setStrokeColor(CYAN_ACCENT)
        c.setLineWidth(3)
        c.roundRect(980, 150, 840, 469, 12, fill=0, stroke=1)
    
    # 3 Pill Badges
    pills = [
        ("🛡️ Rock-Solid LTS", 100),
        ("🚀 100% Go & TUI", 340),
        ("🤖 Hybrid AI & Voice", 560),
    ]
    for text, x in pills:
        c.setFillColor(BG_CARD)
        c.setStrokeColor(BORDER_COLOR)
        c.roundRect(x, 150, 200, 50, 10, fill=1, stroke=1)
        c.setFont("Helvetica-Bold", 18)
        c.setFillColor(TEXT_WHITE)
        c.drawString(x + 20, 168, text)
        
    c.setFont("Helvetica-Bold", 24)
    c.setFillColor(CYAN_ACCENT)
    c.drawString(100, 80, "Desliza para conocer cómo funciona  →")
    c.showPage()


def create_slide_2(c):
    """Floating HUD Slide"""
    c.setFillColor(BG_DARK)
    c.rect(0, 0, WIDTH, HEIGHT, fill=1, stroke=0)
    draw_header_footer(c, 2, 5, "01 // FLOATING HUD LAUNCHER")
    
    # Left Content
    c.setFont("Helvetica-Bold", 56)
    c.setFillColor(TEXT_WHITE)
    c.drawString(100, HEIGHT - 180, "Vula Floating HUD")
    c.setFont("Helvetica-Bold", 32)
    c.setFillColor(CYAN_ACCENT)
    c.drawString(100, HEIGHT - 230, "Super + Space  •  Inspirado en Raycast / Spotlight")
    
    c.setFont("Helvetica", 24)
    c.setFillColor(TEXT_MUTED)
    points = [
        "• Carrusel de 3 pestañas: Acciones, IA conversacional y Diagnóstico.",
        "• Cero lentitud: Renderizado instantáneo con Bubbletea en terminal flotante.",
        "• Acceso rápido a perfiles de trabajo, selector de temas y tienda TUI.",
        "• Búsqueda difusa (fuzzy search) integrada para comandos del sistema."
    ]
    y = HEIGHT - 310
    for pt in points:
        c.drawString(100, y, pt)
        y -= 50
        
    # Image on right
    img_path = os.path.join(ASSETS_DIR, "vula_hud_showcase.png")
    if os.path.exists(img_path):
        c.drawImage(img_path, 880, 160, width=940, height=477, preserveAspectRatio=True, mask='auto')
        c.setStrokeColor(BORDER_COLOR)
        c.setLineWidth(2)
        c.roundRect(880, 160, 940, 477, 12, fill=0, stroke=1)
        
    c.showPage()


def create_slide_3(c):
    """OS Control AI Slide"""
    c.setFillColor(BG_DARK)
    c.rect(0, 0, WIDTH, HEIGHT, fill=1, stroke=0)
    draw_header_footer(c, 3, 5, "02 // NATURAL LANGUAGE OS CONTROL")
    
    # Left Content
    c.setFont("Helvetica-Bold", 56)
    c.setFillColor(TEXT_WHITE)
    c.drawString(100, HEIGHT - 180, "Control del OS con IA: vula do")
    c.setFont("Helvetica-Bold", 32)
    c.setFillColor(CYAN_ACCENT)
    c.drawString(100, HEIGHT - 230, "Ejecuta acciones en español sin recordar sintaxis compleja")
    
    c.setFont("Helvetica", 24)
    c.setFillColor(TEXT_MUTED)
    points = [
        "• 'vula do ajusta el brillo al 40%': Ajuste directo por D-Bus/hardware.",
        "• 'vula do crea una app web para Notion': Genera lanzador .desktop y favicon.",
        "• Guardarraíles de seguridad en 3 niveles (Safe, Sensitive, High-Risk).",
        "• Registro de auditoría persistente en ~/.config/vula/audit.log."
    ]
    y = HEIGHT - 310
    for pt in points:
        c.drawString(100, y, pt)
        y -= 50
        
    # Image on right
    img_path = os.path.join(ASSETS_DIR, "vula_hud_action.png")
    if os.path.exists(img_path):
        c.drawImage(img_path, 880, 160, width=940, height=477, preserveAspectRatio=True, mask='auto')
        c.setStrokeColor(BORDER_COLOR)
        c.setLineWidth(2)
        c.roundRect(880, 160, 940, 477, 12, fill=0, stroke=1)
        
    c.showPage()


def create_slide_4(c):
    """Workstation Profiles Slide"""
    c.setFillColor(BG_DARK)
    c.rect(0, 0, WIDTH, HEIGHT, fill=1, stroke=0)
    draw_header_footer(c, 4, 5, "03 // WORKSTATION PROFILES & THEMES")
    
    # Left Content
    c.setFont("Helvetica-Bold", 56)
    c.setFillColor(TEXT_WHITE)
    c.drawString(100, HEIGHT - 180, "Perfiles Listos en 1 Click")
    c.setFont("Helvetica-Bold", 32)
    c.setFillColor(CYAN_ACCENT)
    c.drawString(100, HEIGHT - 230, "Adapta tu entorno según tu rol diario")
    
    c.setFont("Helvetica", 24)
    c.setFillColor(TEXT_MUTED)
    points = [
        "• Full-Stack Dev: Node, Go, Python, Docker, Lazygit y JetBrains Mono.",
        "• DevOps & Cloud: kubectl, Helm, k9s, Terraform, Tmux y Hack Font.",
        "• 19 Temas Curados: Tokyo Night, Catppuccin, Nord, Everforest y más.",
        "• Tiling Inteligente con resaltado dinámico de bordes en color de acento."
    ]
    y = HEIGHT - 310
    for pt in points:
        c.drawString(100, y, pt)
        y -= 50
        
    # Image on right
    img_path = os.path.join(ASSETS_DIR, "vula_profiles_showcase.png")
    if os.path.exists(img_path):
        c.drawImage(img_path, 880, 160, width=940, height=524, preserveAspectRatio=True, mask='auto')
        c.setStrokeColor(BORDER_COLOR)
        c.setLineWidth(2)
        c.roundRect(880, 160, 940, 524, 12, fill=0, stroke=1)
        
    c.showPage()


def create_slide_5(c):
    """Final CTA Slide"""
    c.setFillColor(BG_DARK)
    c.rect(0, 0, WIDTH, HEIGHT, fill=1, stroke=0)
    draw_header_footer(c, 5, 5, "04 // OPEN SOURCE & COMMUNITY")
    
    # Centered Title
    c.setFont("Helvetica-Bold", 64)
    c.setFillColor(TEXT_WHITE)
    c.drawCentredString(WIDTH / 2, HEIGHT - 220, "Únete a la evolución de Vula OS")
    
    c.setFont("Helvetica", 28)
    c.setFillColor(CYAN_ACCENT)
    c.drawCentredString(WIDTH / 2, HEIGHT - 290, "100% Open Source (MIT)  •  Hecho con pasión para la comunidad Linux")
    
    # 2 Feature boxes in center
    box_w = 780
    box_h = 320
    box_y = 280
    
    # Box 1: Quick Install
    c.setFillColor(BG_CARD)
    c.setStrokeColor(BORDER_COLOR)
    c.roundRect(140, box_y, box_w, box_h, 16, fill=1, stroke=1)
    
    c.setFont("Helvetica-Bold", 26)
    c.setFillColor(CYAN_ACCENT)
    c.drawString(180, box_y + box_h - 60, "📦 Instalación Rápida")
    
    c.setFont("Courier-Bold", 22)
    c.setFillColor(TEXT_WHITE)
    c.drawString(180, box_y + box_h - 120, "git clone https://github.com/hjaguen/vula.git")
    c.drawString(180, box_y + box_h - 160, "cd vula && make install")
    c.drawString(180, box_y + box_h - 200, "vula desktop setup")
    
    c.setFont("Helvetica", 20)
    c.setFillColor(TEXT_MUTED)
    c.drawString(180, box_y + 40, "Requiere Ubuntu 24.04+ LTS y GNOME 46+")

    # Box 2: Links & Collaboration
    c.setFillColor(BG_CARD)
    c.setStrokeColor(BORDER_COLOR)
    c.roundRect(1000, box_y, box_w, box_h, 16, fill=1, stroke=1)
    
    c.setFont("Helvetica-Bold", 26)
    c.setFillColor(PURPLE_ACCENT)
    c.drawString(1040, box_y + box_h - 60, "⭐ Apoya el Proyecto")
    
    c.setFont("Helvetica-Bold", 24)
    c.setFillColor(TEXT_WHITE)
    c.drawString(1040, box_y + box_h - 120, "🔗 Repositorio: github.com/hjaguen/vula")
    c.drawString(1040, box_y + box_h - 170, "📖 Documentación: hjaguen.github.io/vula")
    
    c.setFont("Helvetica", 22)
    c.setFillColor(TEXT_MUTED)
    c.drawString(1040, box_y + 80, "¿Te gustaría colaborar o sugerir nuevas funciones?")
    c.drawString(1040, box_y + 45, "¡Déjame tu opinión en los comentarios del post!")
    
    # Bottom callout
    c.setFont("Helvetica-Bold", 28)
    c.setFillColor(TEXT_WHITE)
    c.drawCentredString(WIDTH / 2, 140, "Dale una estrella ⭐️ al repositorio en GitHub si te gusta la iniciativa.")
    
    c.showPage()


def main():
    c = canvas.Canvas(OUTPUT_PDF, pagesize=(WIDTH, HEIGHT))
    create_slide_1(c)
    create_slide_2(c)
    create_slide_3(c)
    create_slide_4(c)
    create_slide_5(c)
    c.save()
    print(f"PDF Carousel successfully generated at: {OUTPUT_PDF}")


if __name__ == "__main__":
    main()
